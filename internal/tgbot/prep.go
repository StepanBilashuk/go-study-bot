package tgbot

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// handlePrep shows a company's interview-prep card (real question themes mapped
// to /learn and /design), or lists which companies have one.
func (b *Bot) handlePrep(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	chatID := update.Message.Chat.ID
	defs := b.getDefs()
	slug := strings.ToLower(commandArg(update.Message.Text, "/prep"))

	if slug == "" {
		if len(defs.Prep) == 0 {
			b.reply(ctx, chatID, "No prep cards yet.")
			return
		}
		slugs := make([]string, 0, len(defs.Prep))
		for s := range defs.Prep {
			slugs = append(slugs, s)
		}
		sort.Strings(slugs)
		b.reply(ctx, chatID, "Interview prep cards: "+strings.Join(prefixed(slugs, "/prep "), " · "))
		return
	}

	card, ok := defs.Prep[slug]
	if !ok {
		b.reply(ctx, chatID, fmt.Sprintf("No prep card for %q. Try /prep to list.", slug))
		return
	}
	b.replyLong(ctx, chatID, card)
}

func prefixed(ss []string, p string) []string {
	out := make([]string, len(ss))
	for i, s := range ss {
		out[i] = p + s
	}
	return out
}
