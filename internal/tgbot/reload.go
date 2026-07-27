package tgbot

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"prepbot/internal/definitions"
	"prepbot/internal/prompts"
)

// handleReload re-reads all YAML (definitions + prompts) from disk and swaps it
// in atomically. A validation error leaves the running set untouched (spec §7).
func (b *Bot) handleReload(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	chatID := update.Message.Chat.ID

	newDefs, err := definitions.Load(b.cfg.DataDir)
	if err != nil {
		b.reply(ctx, chatID, "Reload aborted — data error:\n"+err.Error())
		return
	}
	newPrompts, err := prompts.Load(b.cfg.PromptsDir)
	if err != nil {
		b.reply(ctx, chatID, "Reload aborted — prompt error:\n"+err.Error())
		return
	}

	b.mu.Lock()
	b.defs = newDefs
	b.prompts = newPrompts
	b.mu.Unlock()

	b.reply(ctx, chatID, fmt.Sprintf(
		"🔄 Reloaded: %d topics, %d drills, %d companies, %d prompts.",
		len(newDefs.Topics), len(newDefs.Drills), len(newDefs.Companies), newPrompts.Count()))
}
