package tgbot

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"prepbot/internal/claude"
)

// handleDebrief is the core loop (spec §7, §13): free text in, Claude extracts
// gaps, confidence updates, and glossary terms; the bot applies them.
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

	system, err := b.getPrompts().Render("debrief", map[string]string{"topic_slugs": strings.Join(slugs, ", ")})
	if err != nil {
		b.reply(ctx, chatID, "Prompt error: "+err.Error())
		return
	}

	var result claude.DebriefResult
	err = b.claude.CompleteJSON(ctx, system, text, func(raw []byte) error {
		r, e := claude.ParseDebrief(raw)
		result = r
		return e
	})
	if err != nil {
		b.reply(ctx, chatID, "Could not process the debrief: "+err.Error())
		return
	}

	// Persist the raw debrief plus the extracted gaps.
	gapsJSON, _ := json.Marshal(result.Gaps)
	if err := b.db.InsertSession(ctx, chatID, text, gapsJSON, result.TopicsTouched); err != nil {
		b.reply(ctx, chatID, "DB error: "+err.Error())
		return
	}

	// Apply confidence updates (only for known slugs) and reschedule reviews
	// on any mastered topic whose confidence dropped.
	applied := 0
	for _, cu := range result.ConfidenceUpdates {
		if _, ok := defs.Topics[cu.Topic]; !ok {
			continue
		}
		if err := b.db.SetConfidence(ctx, chatID, cu.Topic, cu.NewConfidence); err != nil {
			b.reply(ctx, chatID, "DB error: "+err.Error())
			return
		}
		applied++
	}
	for _, g := range result.Gaps {
		if _, ok := defs.Topics[g.Topic]; !ok {
			continue
		}
		if err := b.db.BumpReviewDue(ctx, chatID, g.Topic); err != nil {
			b.reply(ctx, chatID, "DB error: "+err.Error())
			return
		}
	}
	for _, term := range result.NewGlossaryTerms {
		if err := b.db.InsertGlossary(ctx, chatID, term.Term, term.Context); err != nil {
			b.reply(ctx, chatID, "DB error: "+err.Error())
			return
		}
	}

	b.awardXP(ctx, chatID, "debrief", 10)

	var sb strings.Builder
	sb.WriteString("📝 Debrief processed.\n")
	if len(result.TopicsTouched) > 0 {
		fmt.Fprintf(&sb, "Topics: %s\n", strings.Join(result.TopicsTouched, ", "))
	}
	fmt.Fprintf(&sb, "Gaps found: %d\n", len(result.Gaps))
	for _, g := range result.Gaps {
		fmt.Fprintf(&sb, "  • [%s] %s: %s\n", g.Severity, g.Topic, g.Detail)
	}
	fmt.Fprintf(&sb, "Confidence updated: %d topic(s)\n", applied)
	if len(result.NewGlossaryTerms) > 0 {
		fmt.Fprintf(&sb, "Glossary added: %d term(s)", len(result.NewGlossaryTerms))
	}
	b.reply(ctx, chatID, strings.TrimRight(sb.String(), "\n"))
}
