// Package tgbot owns the Telegram long-polling loop, command handlers, the
// per-chat conversation state used by the multi-turn flows (/start, /drill),
// and the daily push scheduler.
package tgbot

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"prepbot/internal/claude"
	"prepbot/internal/config"
	"prepbot/internal/db"
	"prepbot/internal/definitions"
	"prepbot/internal/prompts"
)

// convMode is what a bare (non-command) text message means for a given chat.
type convMode int

const (
	modeNone         convMode = iota
	modeCalibrating           // /start is walking topics one at a time
	modeDrillPending          // /drill sent a challenge and awaits the answer
	modeQuizzing              // /quiz is walking items one at a time
	modeStoryMining           // /story asked a question and awaits the answer
	modeImportCompany         // /importcompany awaits the pasted JSON
)

// conversation is the in-memory state for one chat, lost on restart —
// acceptable for these interactive flows.
type conversation struct {
	mode        convMode
	calibTopics []string // ordered topic slugs still to calibrate
	calibIdx    int
	drillKind   string // kind awaiting an answer
	drillText   string // the challenge that was posed

	quizTopic string             // /quiz: topic under test
	quizItems []claude.QuizItem  // /quiz: the 10 items
	quizIdx   int                // /quiz: current item
	quizScore int                // /quiz: correct so far

	storyComp string // /story: competency being mined
}

// Bot bundles the Telegram client with everything the handlers need.
type Bot struct {
	api    *bot.Bot
	cfg    config.Config
	db     *db.DB
	claude *claude.Client

	// defs and prompts are swapped atomically under mu on /reload.
	mu      sync.RWMutex
	defs    *definitions.Definitions
	prompts *prompts.Set

	convMu sync.Mutex
	conv   map[int64]*conversation
}

// New constructs the bot and registers the Phase 1 handlers.
func New(cfg config.Config, database *db.DB, defs *definitions.Definitions, ps *prompts.Set) (*Bot, error) {
	b := &Bot{
		cfg:     cfg,
		db:      database,
		claude:  claude.New(cfg.AnthropicAPIKey, cfg.AnthropicModel),
		defs:    defs,
		prompts: ps,
		conv:    make(map[int64]*conversation),
	}

	api, err := bot.New(cfg.TelegramToken,
		bot.WithDefaultHandler(b.handleDefault),
		bot.WithMiddlewares(b.userMiddleware),
	)
	if err != nil {
		return nil, fmt.Errorf("create telegram bot: %w", err)
	}
	b.api = api

	// No-argument commands: exact match.
	api.RegisterHandler(bot.HandlerTypeMessageText, "/ping", bot.MatchTypeExact, b.handlePing)
	api.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, b.handleStart)
	api.RegisterHandler(bot.HandlerTypeMessageText, "/today", bot.MatchTypeExact, b.handleToday)
	api.RegisterHandler(bot.HandlerTypeMessageText, "/drill", bot.MatchTypeExact, b.handleDrill)
	api.RegisterHandler(bot.HandlerTypeMessageText, "/reload", bot.MatchTypeExact, b.handleReload)
	// Commands that take arguments: prefix match.
	api.RegisterHandler(bot.HandlerTypeMessageText, "/done", bot.MatchTypePrefix, b.handleDone)
	api.RegisterHandler(bot.HandlerTypeMessageText, "/debrief", bot.MatchTypePrefix, b.handleDebrief)
	api.RegisterHandler(bot.HandlerTypeMessageText, "/boss", bot.MatchTypePrefix, b.handleBoss)
	api.RegisterHandler(bot.HandlerTypeMessageText, "/push", bot.MatchTypePrefix, b.handlePush)
	// Phase 2-5.
	api.RegisterHandler(bot.HandlerTypeMessageText, "/quiz", bot.MatchTypePrefix, b.handleQuiz)
	api.RegisterHandler(bot.HandlerTypeMessageText, "/stories", bot.MatchTypeExact, b.handleStories)
	api.RegisterHandler(bot.HandlerTypeMessageText, "/story", bot.MatchTypeExact, b.handleStory)
	api.RegisterHandler(bot.HandlerTypeMessageText, "/newcompany", bot.MatchTypePrefix, b.handleNewCompany)
	api.RegisterHandler(bot.HandlerTypeMessageText, "/importcompany", bot.MatchTypeExact, b.handleImportCompany)
	api.RegisterHandler(bot.HandlerTypeMessageText, "/ready", bot.MatchTypeExact, b.handleReady)
	api.RegisterHandler(bot.HandlerTypeMessageText, "/interview", bot.MatchTypePrefix, b.handleInterview)
	api.RegisterHandler(bot.HandlerTypeMessageText, "/stats", bot.MatchTypeExact, b.handleStats)
	api.RegisterHandler(bot.HandlerTypeMessageText, "/glossary", bot.MatchTypeExact, b.handleGlossary)
	api.RegisterHandler(bot.HandlerTypeMessageText, "/progress", bot.MatchTypeExact, b.handleProgress)
	api.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, b.handleHelp)
	api.RegisterHandler(bot.HandlerTypeMessageText, "/whoami", bot.MatchTypeExact, b.handleWhoami)

	return b, nil
}

// Start runs long polling plus the daily push scheduler until ctx is cancelled.
func (b *Bot) Start(ctx context.Context) {
	go b.runScheduler(ctx)
	b.api.Start(ctx)
}

// userMiddleware upserts the sender into the users table on every message, so
// per-user writes satisfy their foreign key and last_seen stays current.
func (b *Bot) userMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, api *bot.Bot, update *models.Update) {
		if update.Message != nil {
			var username string
			if update.Message.From != nil {
				username = update.Message.From.Username
			}
			if err := b.db.UpsertUser(ctx, update.Message.Chat.ID, username); err != nil {
				slog.Error("upsert user", "chat_id", update.Message.Chat.ID, "err", err)
			}
		}
		next(ctx, api, update)
	}
}

// --- shared helpers ---

func (b *Bot) getDefs() *definitions.Definitions {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.defs
}

func (b *Bot) getPrompts() *prompts.Set {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.prompts
}

// reply sends plain text to the originating chat, logging send failures.
func (b *Bot) reply(ctx context.Context, chatID int64, text string) {
	if _, err := b.api.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: text}); err != nil {
		slog.Error("send message", "chat_id", chatID, "err", err)
	}
}

// send pushes text to an arbitrary chat (used by the scheduler).
func (b *Bot) send(ctx context.Context, chatID int64, text string) {
	b.reply(ctx, chatID, text)
}

func (b *Bot) getConv(chatID int64) *conversation {
	b.convMu.Lock()
	defer b.convMu.Unlock()
	c, ok := b.conv[chatID]
	if !ok {
		c = &conversation{mode: modeNone}
		b.conv[chatID] = c
	}
	return c
}

func (b *Bot) setConv(chatID int64, c *conversation) {
	b.convMu.Lock()
	defer b.convMu.Unlock()
	b.conv[chatID] = c
}

func (b *Bot) clearConv(chatID int64) {
	b.convMu.Lock()
	defer b.convMu.Unlock()
	b.conv[chatID] = &conversation{mode: modeNone}
}

// commandArg returns the text after a "/command" prefix, trimmed.
func commandArg(text, command string) string {
	return strings.TrimSpace(strings.TrimPrefix(text, command))
}

// handlePing is the Step 1 smoke test, kept for liveness checks.
func (b *Bot) handlePing(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	b.reply(ctx, update.Message.Chat.ID, "pong")
}

// handleDefault routes bare text by conversation mode, and nudges on unknown
// commands.
func (b *Bot) handleDefault(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	chatID := update.Message.Chat.ID
	text := strings.TrimSpace(update.Message.Text)

	if strings.HasPrefix(text, "/") {
		b.reply(ctx, chatID, "Unknown command. Send /help for the full list.")
		return
	}

	conv := b.getConv(chatID)
	switch conv.mode {
	case modeCalibrating:
		b.handleCalibrationAnswer(ctx, chatID, conv, text)
	case modeDrillPending:
		b.handleDrillAnswer(ctx, chatID, conv, text)
	case modeQuizzing:
		b.handleQuizAnswer(ctx, chatID, conv, text)
	case modeStoryMining:
		b.handleStoryAnswer(ctx, chatID, conv, text)
	case modeImportCompany:
		b.handleImportCompanyJSON(ctx, chatID, text)
	default:
		b.reply(ctx, chatID, "Send a command like /today. Free-text goes to Claude only via /debrief <text>.")
	}
}

// awardXP records XP for an activity when the gamification flag is on (Phase 5).
// Failures are logged, never surfaced — XP is a side effect, not the task.
func (b *Bot) awardXP(ctx context.Context, userID int64, kind string, points int) {
	if !b.cfg.Gamification {
		return
	}
	if err := b.db.AwardXP(ctx, userID, kind, points); err != nil {
		slog.Error("award xp", "user", userID, "kind", kind, "err", err)
	}
}
