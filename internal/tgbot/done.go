package tgbot

import (
	"context"
	"fmt"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"prepbot/internal/db"
	"prepbot/internal/scheduler"
)

// handleDone advances a topic one stage if its prerequisites are met, otherwise
// names what is missing (spec §7). The topic-specific gate is a human judgment
// the bot can't verify, so /done trusts the user cleared it; the automatable
// check is the dependency prerequisite.
func (b *Bot) handleDone(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	chatID := update.Message.Chat.ID
	slug := commandArg(update.Message.Text, "/done")
	if slug == "" {
		b.reply(ctx, chatID, "Usage: /done <topic-slug>")
		return
	}

	defs := b.getDefs()
	topic, ok := defs.Topics[slug]
	if !ok {
		b.reply(ctx, chatID, fmt.Sprintf("Unknown topic %q. Check the slug.", slug))
		return
	}

	prog, err := b.db.ListProgress(ctx, chatID)
	if err != nil {
		b.reply(ctx, chatID, "DB error: "+err.Error())
		return
	}

	// Prerequisite gate: to advance past stage 0, dependencies must be known.
	cur := prog[slug]
	if cur.Stage == 0 {
		for _, dep := range topic.DependsOn {
			if prog[dep].Confidence < prereqConfidence {
				depName := dep
				if dt, ok := defs.Topics[dep]; ok {
					depName = dt.Name
				}
				b.reply(ctx, chatID, fmt.Sprintf(
					"⛔ %s is blocked: prerequisite %s needs confidence ≥%d (currently %d). Work it up first.",
					topic.Name, depName, prereqConfidence, prog[dep].Confidence))
				return
			}
		}
	}

	now := time.Now()
	next := db.Progress{
		UserID:      chatID,
		TopicSlug:   slug,
		Confidence:  cur.Confidence,
		LastTouched: &now,
	}
	if next.Confidence == 0 {
		next.Confidence = 1 // sensible floor for a topic touched for the first time
	}

	var msg string
	if cur.Stage < masteredStage {
		next.Stage = cur.Stage + 1
		if next.Stage == masteredStage {
			// First time reaching Review → enter spaced repetition.
			next.Attempts = 0
			due := scheduler.NextReview(now, 0)
			next.NextDue = &due
			msg = fmt.Sprintf("🎯 %s mastered — enters spaced repetition, next review in %d day(s).",
				topic.Name, scheduler.IntervalDays()[0])
		} else {
			next.Attempts = cur.Attempts + 1
			next.NextDue = &now // stays in the active rotation
			msg = fmt.Sprintf("✅ %s → stage %d (%s).\nGate for the next stage: %s",
				topic.Name, next.Stage, stageName(next.Stage), topic.Gate)
		}
	} else {
		// A spaced-repetition review was completed.
		next.Stage = masteredStage
		next.Attempts = cur.Attempts + 1
		due := scheduler.NextReview(now, next.Attempts)
		next.NextDue = &due
		days := int(due.Sub(now).Hours()/24 + 0.5)
		msg = fmt.Sprintf("♻️ %s review logged — next review in %d day(s).", topic.Name, days)
	}

	if err := b.db.UpsertProgress(ctx, next); err != nil {
		b.reply(ctx, chatID, "DB error: "+err.Error())
		return
	}
	if next.Stage == masteredStage && cur.Stage < masteredStage {
		b.awardXP(ctx, chatID, "topic_closed", 100) // reached Review the first time
	} else {
		b.awardXP(ctx, chatID, "topic_closed", 20) // a stage advance / review
	}
	b.reply(ctx, chatID, msg)
}
