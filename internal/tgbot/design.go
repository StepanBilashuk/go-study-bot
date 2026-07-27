package tgbot

import (
	"context"
	"sort"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"prepbot/internal/definitions"
)

// handleDesigns lists the seeded "design X" case studies as tappable buttons.
func (b *Bot) handleDesigns(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	chatID := update.Message.Chat.ID
	defs := b.getDefs()
	if len(defs.Designs) == 0 {
		b.reply(ctx, chatID, "No design case studies loaded yet.")
		return
	}
	b.replyInline(ctx, chatID, "🏗 System-design case studies — tap one for the full walkthrough:", designKeyboard(defs))
}

// handleDesign shows one case study for a typed /design <slug>.
func (b *Bot) handleDesign(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	chatID := update.Message.Chat.ID
	slug := commandArg(update.Message.Text, "/design")
	if slug == "" {
		b.reply(ctx, chatID, "Usage: /design <slug> — or /designs to browse.")
		return
	}
	b.sendDesignCard(ctx, chatID, slug)
}

// handleDesignCallback shows a case study when its button is tapped.
func (b *Bot) handleDesignCallback(ctx context.Context, api *bot.Bot, update *models.Update) {
	cq := update.CallbackQuery
	if cq == nil {
		return
	}
	if _, err := api.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: cq.ID}); err != nil {
		_ = err
	}
	b.sendDesignCard(ctx, cq.Message.Message.Chat.ID, strings.TrimPrefix(cq.Data, "design:"))
}

func (b *Bot) sendDesignCard(ctx context.Context, chatID int64, slug string) {
	defs := b.getDefs()
	content, ok := defs.Designs[slug]
	if !ok {
		b.reply(ctx, chatID, "Unknown design \""+slug+"\". Try /designs.")
		return
	}
	content = strings.TrimRight(content, "\n") +
		"\n\n🎯 Practice: run /boss and use this as the design task."
	b.replyLong(ctx, chatID, content)
}

func designKeyboard(defs *definitions.Definitions) models.ReplyMarkup {
	slugs := make([]string, 0, len(defs.Designs))
	for s := range defs.Designs {
		slugs = append(slugs, s)
	}
	sort.Slice(slugs, func(i, j int) bool {
		// The RESHADED framework card always comes first.
		if (slugs[i] == "framework") != (slugs[j] == "framework") {
			return slugs[i] == "framework"
		}
		return slugs[i] < slugs[j]
	})

	var rows [][]models.InlineKeyboardButton
	for _, s := range slugs {
		rows = append(rows, []models.InlineKeyboardButton{
			{Text: "🏗 " + designTitle(defs.Designs[s], s), CallbackData: "design:" + s},
		})
	}
	return &models.InlineKeyboardMarkup{InlineKeyboard: rows}
}

// designTitle uses the first non-empty line of the walkthrough as the button
// label, stripping a leading emoji, falling back to the slug.
func designTitle(content, slug string) string {
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		line = strings.TrimSpace(strings.TrimPrefix(line, "🏗"))
		return line
	}
	return slug
}
