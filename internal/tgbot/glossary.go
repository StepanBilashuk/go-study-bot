package tgbot

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// handleGlossary reviews the English terms the user could not produce under
// pressure (collected by /debrief). Directly serves the B1.4→B2 goal (spec §2).
func (b *Bot) handleGlossary(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	chatID := update.Message.Chat.ID
	entries, err := b.db.ListGlossary(ctx, chatID, 40)
	if err != nil {
		b.reply(ctx, chatID, "DB error: "+err.Error())
		return
	}
	if len(entries) == 0 {
		b.reply(ctx, chatID, "📕 Glossary is empty. It fills up from /debrief with English terms you reach for but can't find.")
		return
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "📕 Glossary — %d term(s) to review:\n", len(entries))
	for _, e := range entries {
		if e.Context != "" {
			fmt.Fprintf(&sb, "• %s — %s\n", e.Term, e.Context)
		} else {
			fmt.Fprintf(&sb, "• %s\n", e.Term)
		}
	}
	b.reply(ctx, chatID, strings.TrimRight(sb.String(), "\n"))
}
