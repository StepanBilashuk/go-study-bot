package tgbot

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"prepbot/internal/definitions"
)

// handleLearn shows a topic's learning card for a typed /learn <slug>.
func (b *Bot) handleLearn(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	chatID := update.Message.Chat.ID
	slug := commandArg(update.Message.Text, "/learn")
	if slug == "" {
		b.reply(ctx, chatID, "Usage: /learn <topic-slug> — or tap a topic under /today.")
		return
	}
	b.sendLearnCard(ctx, chatID, slug)
}

// handleLearnCallback shows a topic's learning card when its /today button is
// tapped (callback data "learn:<slug>").
func (b *Bot) handleLearnCallback(ctx context.Context, api *bot.Bot, update *models.Update) {
	cq := update.CallbackQuery
	if cq == nil {
		return
	}
	// Stop the button's loading spinner.
	if _, err := api.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: cq.ID}); err != nil {
		// non-fatal
		_ = err
	}
	chatID := cq.Message.Message.Chat.ID
	slug := strings.TrimPrefix(cq.Data, "learn:")
	b.sendLearnCard(ctx, chatID, slug)
}

// sendLearnCard shows a topic's seeded theory (read in-bot), followed by its
// gate, prerequisites and resources. Going to Claude is reserved for the boss
// mock — theory lives here. Topics without seeded theory fall back to an
// emitted "teach me" prompt.
func (b *Bot) sendLearnCard(ctx context.Context, chatID int64, slug string) {
	defs := b.getDefs()
	t, ok := defs.Topics[slug]
	if !ok {
		b.reply(ctx, chatID, "Unknown topic \""+slug+"\".")
		return
	}

	var sb strings.Builder
	if theory := strings.TrimSpace(defs.Theory[slug]); theory != "" {
		sb.WriteString(theory)
		sb.WriteString("\n")
	} else {
		// Fallback: no seeded theory yet → emit a teach-me prompt for Claude.
		fmt.Fprintf(&sb, "📖 %s  (%s)\n\n", t.Name, t.Track)
		if prompt, err := b.getPrompts().Render("learn", map[string]string{
			"topic": t.Name, "track": string(t.Track), "gate": t.Gate,
		}); err == nil {
			sb.WriteString("Theory not seeded yet — run this in your Claude:\n\n" + prompt + "\n")
		}
	}

	fmt.Fprintf(&sb, "\n🎯 Gate: %s\n", t.Gate)
	if len(t.DependsOn) > 0 {
		fmt.Fprintf(&sb, "Depends on: %s\n", strings.Join(t.DependsOn, ", "))
	}
	if res := allResources(defs, slug); len(res) > 0 {
		sb.WriteString("Resources:\n")
		for _, line := range res {
			fmt.Fprintf(&sb, "  • %s\n", line)
		}
	}
	fmt.Fprintf(&sb, "\nReady? /quiz %s to self-test · /done %s when you pass the gate.", slug, slug)

	b.replyLong(ctx, chatID, sb.String())
}

// allResources returns every resource line for a topic (books, videos,
// articles) — the learn card shows all of them, not the /today cap of 2.
func allResources(defs *definitions.Definitions, slug string) []string {
	var out []string
	for _, r := range defs.Resources {
		if r.Topic != slug {
			continue
		}
		var line string
		if book, ok := defs.Books[r.Source]; ok && r.Type == "book" {
			if r.Chapter > 0 {
				line = fmt.Sprintf("%s ch.%d — %s", book.Title, r.Chapter, r.Section)
			} else {
				line = fmt.Sprintf("%s — %s", book.Title, r.Section)
			}
		} else if r.Section != "" {
			line = fmt.Sprintf("[%s] %s", r.Type, r.Section)
		} else {
			line = fmt.Sprintf("[%s] %s", r.Type, r.Source)
		}
		if r.Pages != nil && *r.Pages != "" {
			line += fmt.Sprintf(" (pp. %s)", *r.Pages)
		}
		if r.EstMin > 0 {
			line += fmt.Sprintf(" ~%dm", r.EstMin)
		}
		out = append(out, line)
	}
	return out
}

// learnKeyboard builds the inline "Learn" buttons under a /today message, one
// per shown topic.
func learnKeyboard(defs *definitions.Definitions, slugs []string) models.ReplyMarkup {
	var rows [][]models.InlineKeyboardButton
	for _, s := range slugs {
		name := s
		if t, ok := defs.Topics[s]; ok {
			name = t.Name
		}
		rows = append(rows, []models.InlineKeyboardButton{
			{Text: "📖 " + name, CallbackData: "learn:" + s},
		})
	}
	if len(rows) == 0 {
		return nil
	}
	return &models.InlineKeyboardMarkup{InlineKeyboard: rows}
}
