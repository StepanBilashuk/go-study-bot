package tgbot

import (
	"context"
	"sort"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// handleDebrief emits a debrief-extraction prompt with the user's text embedded.
// The user runs it in their Claude and pastes the JSON back; applyDebrief then
// stores the session and applies gaps/confidence/glossary (spec §7, §13).
func (b *Bot) handleDebrief(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	chatID := update.Message.Chat.ID
	text := commandArg(update.Message.Text, "/debrief")
	if text == "" {
		b.reply(ctx, chatID, "Usage: /debrief <what actually happened, in plain text>")
		return
	}

	defs := b.getDefs()
	slugs := make([]string, 0, len(defs.Topics))
	for s := range defs.Topics {
		slugs = append(slugs, s)
	}
	sort.Strings(slugs)

	prompt, err := b.getPrompts().Render("debrief", map[string]string{
		"topic_slugs": strings.Join(slugs, ", "),
		"debrief":     text,
	})
	if err != nil {
		b.reply(ctx, chatID, "Prompt error: "+err.Error())
		return
	}

	b.setConv(chatID, &conversation{mode: modeAwaitImport, kind: importDebrief, debriefText: text})
	b.reply(ctx, chatID, emitHeader("Debrief")+prompt)
}
