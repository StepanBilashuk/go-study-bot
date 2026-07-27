package tgbot

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// handleDrill emits a drill of the weakest kind: the generation prompt plus an
// instruction to score the answer and return JSON. The user runs the whole
// drill in their Claude and pastes the {score,outcome} back; applyDrill logs it
// so the weakest-kind rotation keeps working (spec §7, §10).
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

	gen, err := b.getPrompts().Render("drill-"+string(drill.Kind), nil)
	if err != nil {
		b.reply(ctx, chatID, "Prompt error: "+err.Error())
		return
	}
	prompt := gen + "\n\nAfter I answer, score me 1-5 and end your reply with ONLY this JSON " +
		"(no prose, no fences):\n{\"score\": 1-5, \"outcome\": \"<max 15 words on what was strong or missing>\"}"

	b.setConv(chatID, &conversation{mode: modeAwaitImport, kind: importDrill, drillKind: string(drill.Kind)})
	b.reply(ctx, chatID, emitHeader(fmt.Sprintf("%s drill (%d min)", drill.Name, drill.DurationMin))+prompt)
}
