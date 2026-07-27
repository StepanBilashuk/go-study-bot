package tgbot

import (
	"context"
	"log/slog"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// quickButtons maps a reply-keyboard button label to the command it runs. Taps
// arrive as ordinary text messages, so handleDefault translates the label back
// to a command (all mapped commands are no-argument).
var quickButtons = map[string]string{
	"📅 Today":    "/today",
	"🎲 Drill":    "/drill",
	"🎯 Boss":     "/boss",
	"📊 Progress": "/progress",
	"🏢 Ready":    "/ready",
	"☰ Commands":  "/help",
}

// quickKeyboard is the persistent bottom bar of the most-used daily actions.
func quickKeyboard() models.ReplyMarkup {
	return &models.ReplyKeyboardMarkup{
		Keyboard: [][]models.KeyboardButton{
			{{Text: "📅 Today"}, {Text: "🎲 Drill"}},
			{{Text: "🎯 Boss"}, {Text: "📊 Progress"}},
			{{Text: "🏢 Ready"}, {Text: "☰ Commands"}},
		},
		ResizeKeyboard: true,
		IsPersistent:   true,
	}
}

// replyKb sends text and installs the persistent quick-bar keyboard.
func (b *Bot) replyKb(ctx context.Context, chatID int64, text string) {
	if _, err := b.api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        text,
		ReplyMarkup: quickKeyboard(),
	}); err != nil {
		slog.Error("send message", "chat_id", chatID, "err", err)
	}
}

// replyInline sends text with an inline keyboard attached. The persistent
// quick-bar reply keyboard set earlier is unaffected.
func (b *Bot) replyInline(ctx context.Context, chatID int64, text string, markup models.ReplyMarkup) {
	if _, err := b.api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        text,
		ReplyMarkup: markup,
	}); err != nil {
		slog.Error("send message", "chat_id", chatID, "err", err)
	}
}

// handleMenu installs the quick-bar and points at the native command menu.
func (b *Bot) handleMenu(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	b.replyKb(ctx, update.Message.Chat.ID,
		"⌨️ Quick buttons are pinned below. Tap ☰ (left of the input) for the full command list, or send /help.")
}

// runByCommand sets the update's text to cmd and dispatches to the handler, so a
// quick-button tap runs the same code path as typing the command.
func (b *Bot) runByCommand(ctx context.Context, api *bot.Bot, update *models.Update, cmd string) {
	if update.Message == nil {
		return
	}
	update.Message.Text = cmd
	switch strings.Fields(cmd)[0] {
	case "/today":
		b.handleToday(ctx, api, update)
	case "/drill":
		b.handleDrill(ctx, api, update)
	case "/boss":
		b.handleBoss(ctx, api, update)
	case "/progress":
		b.handleProgress(ctx, api, update)
	case "/ready":
		b.handleReady(ctx, api, update)
	case "/help":
		b.handleHelp(ctx, api, update)
	default:
		b.reply(ctx, update.Message.Chat.ID, "Unknown menu action.")
	}
}

// botCommands is the native command menu (the ☰ button + "/" autocomplete),
// registered on startup via SetMyCommands. Daily-loop commands come first.
func botCommands() []models.BotCommand {
	return []models.BotCommand{
		{Command: "today", Description: "Daily plan — 3 topics + a drill"},
		{Command: "learn", Description: "Topic theory: /learn <slug> (or tap in /today)"},
		{Command: "done", Description: "Advance a topic: /done <slug>"},
		{Command: "drill", Description: "Process drill (paste the score back)"},
		{Command: "debrief", Description: "Log a debrief: /debrief <text>"},
		{Command: "boss", Description: "Mock brief (add 'behavioral')"},
		{Command: "quiz", Description: "Recognition quiz: /quiz <slug>"},
		{Command: "story", Description: "Mine a STAR story"},
		{Command: "stories", Description: "Competency matrix"},
		{Command: "glossary", Description: "English terms to review"},
		{Command: "ready", Description: "Readiness per company"},
		{Command: "newcompany", Description: "Research prompt: /newcompany <name>"},
		{Command: "importcompany", Description: "Import a researched company JSON"},
		{Command: "interview", Description: "Set a date: /interview <slug> <YYYY-MM-DD>"},
		{Command: "progress", Description: "Per-track progress summary"},
		{Command: "stats", Description: "XP, streak, mastery"},
		{Command: "start", Description: "Calibrate all topics"},
		{Command: "push", Description: "Daily push: /push on|off"},
		{Command: "menu", Description: "Show the quick-button keyboard"},
		{Command: "cancel", Description: "Cancel a pending paste-back"},
		{Command: "whoami", Description: "Show your chat id"},
		{Command: "reload", Description: "Reload YAML from disk"},
		{Command: "help", Description: "Full command list"},
		{Command: "ping", Description: "Liveness check"},
	}
}
