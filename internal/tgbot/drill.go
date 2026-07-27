package tgbot

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"prepbot/internal/claude"
)

// handleDrill serves one process drill of the weakest kind (spec §7). It asks
// Claude to generate a fresh challenge, sends it, and parks the chat waiting for
// the candidate's answer (a bare text message).
func (b *Bot) handleDrill(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	chatID := update.Message.Chat.ID

	defs := b.getDefs()
	stats, err := b.db.DrillStats(ctx, chatID)
	if err != nil {
		b.reply(ctx, chatID, "DB error: "+err.Error())
		return
	}
	drill, ok := pickDrill(defs, stats)
	if !ok {
		b.reply(ctx, chatID, "No drills are defined. Add data/drills/process-drills.yaml.")
		return
	}

	system, err := b.getPrompts().Render("drill-"+string(drill.Kind), nil)
	if err != nil {
		b.reply(ctx, chatID, "Prompt error: "+err.Error())
		return
	}
	challenge, err := b.claude.Complete(ctx, system, "Generate one drill now.")
	if err != nil {
		b.reply(ctx, chatID, "Claude error: "+err.Error())
		return
	}

	b.setConv(chatID, &conversation{
		mode:      modeDrillPending,
		drillKind: string(drill.Kind),
		drillText: challenge,
	})
	b.reply(ctx, chatID, fmt.Sprintf("🎯 %s drill (%d min):\n\n%s\n\nReply with your answer.",
		drill.Name, drill.DurationMin, challenge))
}

// handleDrillAnswer scores a candidate's reply to a pending drill, logs it, and
// clears the pending state.
func (b *Bot) handleDrillAnswer(ctx context.Context, chatID int64, conv *conversation, answer string) {
	system, err := b.getPrompts().Render("drill-score", map[string]string{"kind": conv.drillKind})
	if err != nil {
		b.reply(ctx, chatID, "Prompt error: "+err.Error())
		return
	}
	user := fmt.Sprintf("Drill:\n%s\n\nCandidate answer:\n%s\n\nScore it.", conv.drillText, answer)

	var score claude.DrillScore
	err = b.claude.CompleteJSON(ctx, system, user, func(raw []byte) error {
		s, e := claude.ParseDrillScore(raw)
		score = s
		return e
	})
	if err != nil {
		b.reply(ctx, chatID, "Could not score the drill: "+err.Error())
		return
	}

	if err := b.db.InsertDrill(ctx, chatID, conv.drillKind, score.Outcome, &score.Score); err != nil {
		b.reply(ctx, chatID, "DB error: "+err.Error())
		return
	}
	b.awardXP(ctx, chatID, "drill", 10)
	b.clearConv(chatID)
	b.reply(ctx, chatID, fmt.Sprintf("Scored %s: %d/5 — %s", conv.drillKind, score.Score, score.Outcome))
}
