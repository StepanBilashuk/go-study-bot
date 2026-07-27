package tgbot

import (
	"context"
	"sort"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"prepbot/internal/definitions"
)

// handleArchs lists the real-world architecture write-ups as tappable buttons.
func (b *Bot) handleArchs(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	chatID := update.Message.Chat.ID
	defs := b.getDefs()
	if len(defs.Arch) == 0 {
		b.reply(ctx, chatID, "No architecture write-ups loaded yet.")
		return
	}
	b.replyInline(ctx, chatID, "🏛 Real-world architectures (how they actually built it) — tap one:", archKeyboard(defs))
}

// handleArch shows one write-up for a typed /arch <slug>.
func (b *Bot) handleArch(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	chatID := update.Message.Chat.ID
	slug := commandArg(update.Message.Text, "/arch")
	if slug == "" {
		b.reply(ctx, chatID, "Usage: /arch <slug> — or /archs to browse.")
		return
	}
	b.sendArchCard(ctx, chatID, slug)
}

func (b *Bot) handleArchCallback(ctx context.Context, api *bot.Bot, update *models.Update) {
	cq := update.CallbackQuery
	if cq == nil {
		return
	}
	if _, err := api.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: cq.ID}); err != nil {
		_ = err
	}
	b.sendArchCard(ctx, cq.Message.Message.Chat.ID, strings.TrimPrefix(cq.Data, "arch:"))
}

func (b *Bot) sendArchCard(ctx context.Context, chatID int64, slug string) {
	defs := b.getDefs()
	content, ok := defs.Arch[slug]
	if !ok {
		b.reply(ctx, chatID, "Unknown architecture \""+slug+"\". Try /archs.")
		return
	}
	b.replyLong(ctx, chatID, content)
}

func archKeyboard(defs *definitions.Definitions) models.ReplyMarkup {
	slugs := make([]string, 0, len(defs.Arch))
	for s := range defs.Arch {
		slugs = append(slugs, s)
	}
	sort.Strings(slugs)

	var rows [][]models.InlineKeyboardButton
	for _, s := range slugs {
		title := strings.TrimSpace(strings.TrimPrefix(firstLine(defs.Arch[s]), "🏛"))
		if title == "" {
			title = s
		}
		rows = append(rows, []models.InlineKeyboardButton{
			{Text: "🏛 " + title, CallbackData: "arch:" + s},
		})
	}
	return &models.InlineKeyboardMarkup{InlineKeyboard: rows}
}

// firstLine returns the first non-empty line of s.
func firstLine(s string) string {
	for _, line := range strings.Split(s, "\n") {
		if l := strings.TrimSpace(line); l != "" {
			return l
		}
	}
	return ""
}
