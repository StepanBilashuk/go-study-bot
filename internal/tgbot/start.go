package tgbot

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"prepbot/internal/definitions"
)

// handleStart emits a calibration prompt covering all topics. The user runs it
// in their own Claude, explains each topic, and pastes the scores JSON back;
// applyCalibration then sets initial confidence (spec §7).
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
	sort.Slice(topics, func(i, j int) bool {
		if topics[i].Track != topics[j].Track {
			return topics[i].Track < topics[j].Track
		}
		return topics[i].Priority < topics[j].Priority
	})

	var list strings.Builder
	for _, t := range topics {
		fmt.Fprintf(&list, "- %s — %s\n", t.Slug, t.Name)
	}
	prompt, err := b.getPrompts().Render("calibration", map[string]string{"topics": strings.TrimRight(list.String(), "\n")})
	if err != nil {
		b.reply(ctx, chatID, "Prompt error: "+err.Error())
		return
	}

	b.setConv(chatID, &conversation{mode: modeAwaitImport, kind: importCalibration})
	b.reply(ctx, chatID, emitHeader("Calibration")+prompt)
}
