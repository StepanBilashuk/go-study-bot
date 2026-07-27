package tgbot

import (
	"context"
	"fmt"
	"sort"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"prepbot/internal/claude"
	"prepbot/internal/definitions"
)

// handleStart begins one-time calibration: it walks the topics one at a time,
// Claude scores each two-sentence answer, and the bot sets initial confidence
// (spec §7). State is held in memory for the chat.
func (b *Bot) handleStart(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	chatID := update.Message.Chat.ID
	defs := b.getDefs()

	topics := make([]definitions.Topic, 0, len(defs.Topics))
	for _, t := range defs.Topics {
		topics = append(topics, t)
	}
	if len(topics) == 0 {
		b.reply(ctx, chatID, "No topics loaded — nothing to calibrate.")
		return
	}
	// Stable order: track then priority.
	sort.Slice(topics, func(i, j int) bool {
		if topics[i].Track != topics[j].Track {
			return topics[i].Track < topics[j].Track
		}
		return topics[i].Priority < topics[j].Priority
	})
	slugs := make([]string, len(topics))
	for i, t := range topics {
		slugs[i] = t.Slug
	}

	b.setConv(chatID, &conversation{mode: modeCalibrating, calibTopics: slugs, calibIdx: 0})
	b.reply(ctx, chatID, fmt.Sprintf(
		"🧭 Calibration: %d topics. For each, explain it in two sentences (or say when you'd apply it). I'll score it 1-5 and set your starting confidence.\n\n%s",
		len(slugs), b.calibQuestion(defs, slugs, 0)))
}

func (b *Bot) calibQuestion(defs *definitions.Definitions, slugs []string, idx int) string {
	t := defs.Topics[slugs[idx]]
	return fmt.Sprintf("Topic %d/%d — %s:\nExplain it in two sentences, or say when you'd apply it.", idx+1, len(slugs), t.Name)
}

// handleCalibrationAnswer scores one topic and advances to the next.
func (b *Bot) handleCalibrationAnswer(ctx context.Context, chatID int64, conv *conversation, answer string) {
	defs := b.getDefs()
	if conv.calibIdx >= len(conv.calibTopics) {
		b.clearConv(chatID)
		return
	}
	slug := conv.calibTopics[conv.calibIdx]
	topic := defs.Topics[slug]

	system, err := b.getPrompts().Render("calibration", nil)
	if err != nil {
		b.reply(ctx, chatID, "Prompt error: "+err.Error())
		return
	}
	user := fmt.Sprintf("Topic slug: %s\nTopic name: %s\nCandidate's answer: %s", slug, topic.Name, answer)

	var result claude.CalibrationResult
	err = b.claude.CompleteJSON(ctx, system, user, func(raw []byte) error {
		r, e := claude.ParseCalibration(raw)
		result = r
		return e
	})
	if err != nil {
		b.reply(ctx, chatID, "Could not score that answer: "+err.Error()+"\nTry rephrasing, or /today to stop calibration.")
		return
	}

	conf, note := result.Scores[0].Confidence, result.Scores[0].Note
	if err := b.db.SetConfidence(ctx, chatID, slug, conf); err != nil {
		b.reply(ctx, chatID, "DB error: "+err.Error())
		return
	}

	conv.calibIdx++
	b.setConv(chatID, conv)

	if conv.calibIdx >= len(conv.calibTopics) {
		b.clearConv(chatID)
		b.reply(ctx, chatID, fmt.Sprintf("%s: %d/5 (%s)\n\n✅ Calibration complete — %d topics scored. Run /today.",
			topic.Name, conf, note, len(conv.calibTopics)))
		return
	}
	b.reply(ctx, chatID, fmt.Sprintf("%s: %d/5 (%s)\n\n%s",
		topic.Name, conf, note, b.calibQuestion(defs, conv.calibTopics, conv.calibIdx)))
}
