package tgbot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// handleInterview sets or clears an interview date (spec §12: "When an interview
// date is set, /today reorders around that company's profile").
//
//	/interview <company-slug> <YYYY-MM-DD>   set a date
//	/interview <company-slug> clear          remove it
func (b *Bot) handleInterview(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	chatID := update.Message.Chat.ID
	arg := commandArg(update.Message.Text, "/interview")
	fields := strings.Fields(arg)
	if len(fields) != 2 {
		b.reply(ctx, chatID, "Usage: /interview <company-slug> <YYYY-MM-DD>   (or: /interview <slug> clear)")
		return
	}
	slug, when := fields[0], fields[1]

	defs := b.getDefs()
	if _, ok := defs.Companies[slug]; !ok {
		b.reply(ctx, chatID, fmt.Sprintf("Unknown company %q. See /ready for slugs.", slug))
		return
	}

	if when == "clear" {
		if err := b.db.ClearInterview(ctx, chatID, slug); err != nil {
			b.reply(ctx, chatID, "DB error: "+err.Error())
			return
		}
		b.reply(ctx, chatID, "🗓 Cleared the interview for "+slug+".")
		return
	}

	date, err := time.Parse("2006-01-02", when)
	if err != nil {
		b.reply(ctx, chatID, "Date must be YYYY-MM-DD, got "+when)
		return
	}
	if err := b.db.SetInterview(ctx, chatID, slug, date); err != nil {
		b.reply(ctx, chatID, "DB error: "+err.Error())
		return
	}
	days := int(time.Until(date).Hours()/24 + 0.5)
	b.reply(ctx, chatID, fmt.Sprintf("🗓 %s interview set for %s (%d days). /today will now prioritise its required topics.",
		defs.Companies[slug].Name, when, days))
}
