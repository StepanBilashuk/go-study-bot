// Package tgbot owns the Telegram long-polling loop, command handlers, the
// per-chat conversation state used by the emit/import flows, and the daily push
// scheduler.
//
// The bot never calls the Anthropic API. Every AI command emits a prompt for
// the user to run in their own Claude; commands that update state then import
// the JSON the user pastes back (like /newcompany → /importcompany).
package tgbot

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"prepbot/internal/config"
	"prepbot/internal/db"
	"prepbot/internal/definitions"
	"prepbot/internal/prompts"
)

// convMode is what a bare (non-command) text message means for a given chat.
type convMode int

const (
	modeNone          convMode = iota
	modeAwaitImport            // an AI command emitted a prompt; awaiting the pasted JSON
	modeImportCompany          // /importcompany awaits the pasted JSON
)

// importKind identifies which emit/import flow is pending, so the pasted JSON is
// parsed and applied correctly.
type importKind string

const (
	importCalibration importKind = "calibration"
	importDebrief     importKind = "debrief"
	importDrill       importKind = "drill"
	importStory       importKind = "story"
)

// conversation is the in-memory state for one chat, lost on restart —
// acceptable for these short emit→paste-back flows.
type conversation struct {
	mode convMode
	kind importKind // which flow is awaiting a paste (mode == modeAwaitImport)

	debriefText string // the raw debrief text, to store on import
	drillKind   string // the drill kind, to log on import
	storyComp   string // the competency being mined, to tag on import
}

// Bot bundles the Telegram client with everything the handlers need.
type Bot struct {
	api *bot.Bot
	cfg config.Config
	db  *db.DB

	// defs and prompts are swapped atomically under mu on /reload.
	mu      sync.RWMutex
	defs    *definitions.Definitions
	prompts *prompts.Set

	convMu sync.Mutex
	conv   map[int64]*conversation
}

// New constructs the bot and registers the handlers.
func New(cfg config.Config, database *db.DB, defs *definitions.Definitions, ps *prompts.Set) (*Bot, error) {
	b := &Bot{
		cfg:     cfg,
		db:      database,
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
	api.RegisterHandler(bot.HandlerTypeMessageText, "/cancel", bot.MatchTypeExact, b.handleCancel)
	api.RegisterHandler(bot.HandlerTypeMessageText, "/menu", bot.MatchTypeExact, b.handleMenu)
	api.RegisterHandler(bot.HandlerTypeMessageText, "/learn", bot.MatchTypePrefix, b.handleLearn)
	api.RegisterHandler(bot.HandlerTypeMessageText, "/designs", bot.MatchTypeExact, b.handleDesigns)
	api.RegisterHandler(bot.HandlerTypeMessageText, "/design", bot.MatchTypePrefix, b.handleDesign)
	api.RegisterHandler(bot.HandlerTypeMessageText, "/prep", bot.MatchTypePrefix, b.handlePrep)
	// Inline buttons: topic theory + design case studies.
	api.RegisterHandler(bot.HandlerTypeCallbackQueryData, "learn:", bot.MatchTypePrefix, b.handleLearnCallback)
	api.RegisterHandler(bot.HandlerTypeCallbackQueryData, "design:", bot.MatchTypePrefix, b.handleDesignCallback)

	return b, nil
}

// Start runs long polling plus the daily push scheduler until ctx is cancelled.
func (b *Bot) Start(ctx context.Context) {
	// Register the native command menu (the ☰ button + "/" autocomplete).
	if _, err := b.api.SetMyCommands(ctx, &bot.SetMyCommandsParams{Commands: botCommands()}); err != nil {
		slog.Error("set my commands", "err", err)
	}
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

// replyLong sends text, splitting it into multiple messages under Telegram's
// 4096-character limit, preferring to break on line boundaries.
func (b *Bot) replyLong(ctx context.Context, chatID int64, text string) {
	const max = 3900
	for len(text) > max {
		cut := strings.LastIndex(text[:max], "\n")
		if cut <= 0 {
			cut = max
		}
		b.reply(ctx, chatID, strings.TrimRight(text[:cut], "\n"))
		text = strings.TrimLeft(text[cut:], "\n")
	}
	if strings.TrimSpace(text) != "" {
		b.reply(ctx, chatID, text)
	}
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

// handleDefault handles quick-keyboard button taps, routes bare text by
// conversation mode, and nudges on unknown commands.
func (b *Bot) handleDefault(ctx context.Context, api *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	chatID := update.Message.Chat.ID
	text := strings.TrimSpace(update.Message.Text)

	// Quick-keyboard button? Run the mapped command.
	if cmd, ok := quickButtons[text]; ok {
		b.runByCommand(ctx, api, update, cmd)
		return
	}

	if strings.HasPrefix(text, "/") {
		b.reply(ctx, chatID, "Unknown command. Send /help for the full list.")
		return
	}

	conv := b.getConv(chatID)
	switch conv.mode {
	case modeAwaitImport:
		b.handleImport(ctx, chatID, conv, text)
	case modeImportCompany:
		b.handleImportCompanyJSON(ctx, chatID, text)
	default:
		b.reply(ctx, chatID, "Send a command — /help for the list. (Waiting to paste JSON back? Run the emitting command first.)")
	}
}

// handleCancel drops any pending emit/import flow.
func (b *Bot) handleCancel(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	chatID := update.Message.Chat.ID
	b.clearConv(chatID)
	b.reply(ctx, chatID, "Cancelled.")
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
