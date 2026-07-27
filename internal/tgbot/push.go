package tgbot

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// handlePush toggles the daily push for the sender: /push on | /push off.
func (b *Bot) handlePush(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	chatID := update.Message.Chat.ID
	arg := strings.ToLower(commandArg(update.Message.Text, "/push"))

	var enabled bool
	switch arg {
	case "on", "":
		enabled = true
	case "off":
		enabled = false
	default:
		b.reply(ctx, chatID, "Usage: /push on | /push off")
		return
	}
	if err := b.db.SetPushEnabled(ctx, chatID, enabled); err != nil {
		b.reply(ctx, chatID, "DB error: "+err.Error())
		return
	}
	if enabled {
		b.reply(ctx, chatID, "🔔 Daily plan push is ON.")
	} else {
		b.reply(ctx, chatID, "🔕 Daily plan push is OFF.")
	}
}

// runScheduler is the Phase 1 scheduler goroutine (spec §3): once a day at the
// configured local hour it pushes /today to every opted-in user. It checks
// every minute and fires at most once per day.
func (b *Bot) runScheduler(ctx context.Context) {
	slog.Info("daily push scheduler started", "hour", b.cfg.PushHour)

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	var lastPush string // YYYY-MM-DD of the last push run, to fire once per day
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			today := now.Format("2006-01-02")
			if now.Hour() != b.cfg.PushHour || lastPush == today {
				continue
			}
			lastPush = today
			b.pushDailyToAll(ctx)
		}
	}
}

func (b *Bot) pushDailyToAll(ctx context.Context) {
	users, err := b.db.PushUsers(ctx)
	if err != nil {
		slog.Error("daily push: list users", "err", err)
		return
	}
	for _, userID := range users {
		text, err := b.buildToday(ctx, userID)
		if err != nil {
			slog.Error("daily push: build today", "user", userID, "err", err)
			continue
		}
		b.send(ctx, userID, "⏰ Daily plan\n\n"+text)
	}
	slog.Info("daily push sent", "users", len(users))
}
