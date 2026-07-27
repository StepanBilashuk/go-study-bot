package tgbot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"prepbot/internal/definitions"
)

// handleProgress shows a compact one-line-per-track summary — mastered / in
// progress / not started — plus review-due count. Deliberately no per-topic
// list, to avoid recreating the paralysis the bot exists to remove (spec §1).
func (b *Bot) handleProgress(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	chatID := update.Message.Chat.ID
	prog, err := b.db.ListProgress(ctx, chatID)
	if err != nil {
		b.reply(ctx, chatID, "DB error: "+err.Error())
		return
	}
	defs := b.getDefs()
	now := time.Now()

	var sb strings.Builder
	sb.WriteString("📊 Progress\n")
	for _, tr := range []struct {
		track definitions.Track
		label string
	}{
		{definitions.TrackAlgorithms, "Algorithms"},
		{definitions.TrackSystemDesign, "System design"},
	} {
		mastered, inProgress, notStarted := 0, 0, 0
		for slug, t := range defs.Topics {
			if t.Track != tr.track {
				continue
			}
			switch p := prog[slug]; {
			case p.Stage >= masteredStage:
				mastered++
			case p.Stage > 0 || p.Confidence > 0:
				inProgress++
			default:
				notStarted++
			}
		}
		total := mastered + inProgress + notStarted
		fmt.Fprintf(&sb, "%s: %d mastered · %d in progress · %d not started (of %d)\n",
			tr.label, mastered, inProgress, notStarted, total)
	}

	due := 0
	for _, p := range prog {
		if p.Stage >= masteredStage && p.NextDue != nil && !p.NextDue.After(now) {
			due++
		}
	}
	fmt.Fprintf(&sb, "Due for review: %d", due)
	b.reply(ctx, chatID, sb.String())
}
