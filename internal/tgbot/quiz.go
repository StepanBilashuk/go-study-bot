package tgbot

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// quizGate is the stage-2 threshold: 8/10 correct (spec §6).
const quizGate = 8

// handleQuiz emits a pattern-recognition quiz prompt for a topic. The user runs
// the whole quiz in their own Claude — there's no state to import, so this is
// emit-only. If they clear the 8/10 gate, they run /done to advance (spec §7).
func (b *Bot) handleQuiz(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	chatID := update.Message.Chat.ID
	slug := commandArg(update.Message.Text, "/quiz")
	if slug == "" {
		b.reply(ctx, chatID, "Usage: /quiz <topic-slug>")
		return
	}
	defs := b.getDefs()
	topic, ok := defs.Topics[slug]
	if !ok {
		b.reply(ctx, chatID, "Unknown topic \""+slug+"\".")
		return
	}

	prompt, err := b.getPrompts().Render("quiz", map[string]string{
		"topic": topic.Name,
		"track": string(topic.Track),
	})
	if err != nil {
		b.reply(ctx, chatID, "Prompt error: "+err.Error())
		return
	}
	b.reply(ctx, chatID, "🧠 Run this quiz in your Claude — answer the 10 items there, no need to paste anything back. "+
		"Clear 8/10 and you've met the stage-2 gate → /done "+slug+".\n\n"+prompt)
}
