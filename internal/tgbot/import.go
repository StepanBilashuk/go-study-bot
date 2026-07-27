package tgbot

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"prepbot/internal/claude"
	"prepbot/internal/db"
)

// emitHeader is the standard framing above an emitted prompt: run it in your own
// Claude, then paste the JSON reply back.
func emitHeader(name string) string {
	return "👇 " + name + " — copy the prompt below into your Claude, then paste the JSON reply back here (or /cancel):\n\n"
}

// handleImport dispatches a pasted JSON reply to the right apply function based
// on which emit/import flow is pending.
func (b *Bot) handleImport(ctx context.Context, chatID int64, conv *conversation, text string) {
	raw := []byte(claude.ExtractJSON(text))
	switch conv.kind {
	case importCalibration:
		b.applyCalibration(ctx, chatID, raw)
	case importDebrief:
		b.applyDebrief(ctx, chatID, conv, raw)
	case importDrill:
		b.applyDrill(ctx, chatID, conv, raw)
	case importStory:
		b.applyStory(ctx, chatID, conv, raw)
	default:
		b.clearConv(chatID)
		b.reply(ctx, chatID, "Nothing to import. Run the command again.")
	}
}

// badJSON reports a parse failure without dropping the pending flow, so the user
// can just paste corrected JSON.
func (b *Bot) badJSON(ctx context.Context, chatID int64, err error) {
	b.reply(ctx, chatID, "Couldn't read that JSON: "+err.Error()+"\nPaste the JSON again, or /cancel.")
}

func (b *Bot) applyCalibration(ctx context.Context, chatID int64, raw []byte) {
	res, err := claude.ParseCalibration(raw)
	if err != nil {
		b.badJSON(ctx, chatID, err)
		return
	}
	defs := b.getDefs()
	applied, unknown := 0, 0
	for _, s := range res.Scores {
		if _, ok := defs.Topics[s.Topic]; !ok {
			unknown++
			continue
		}
		if err := b.db.SetConfidence(ctx, chatID, s.Topic, s.Confidence); err != nil {
			b.reply(ctx, chatID, "DB error: "+err.Error())
			return
		}
		applied++
	}
	b.clearConv(chatID)
	msg := fmt.Sprintf("✅ Calibrated %d topic(s). Run /today.", applied)
	if unknown > 0 {
		msg += fmt.Sprintf(" (%d unknown slugs ignored)", unknown)
	}
	b.reply(ctx, chatID, msg)
}

func (b *Bot) applyDebrief(ctx context.Context, chatID int64, conv *conversation, raw []byte) {
	res, err := claude.ParseDebrief(raw)
	if err != nil {
		b.badJSON(ctx, chatID, err)
		return
	}
	defs := b.getDefs()

	gapsJSON, _ := json.Marshal(res.Gaps)
	if err := b.db.InsertSession(ctx, chatID, conv.debriefText, gapsJSON, res.TopicsTouched); err != nil {
		b.reply(ctx, chatID, "DB error: "+err.Error())
		return
	}
	applied := 0
	for _, cu := range res.ConfidenceUpdates {
		if _, ok := defs.Topics[cu.Topic]; !ok {
			continue
		}
		if err := b.db.SetConfidence(ctx, chatID, cu.Topic, cu.NewConfidence); err != nil {
			b.reply(ctx, chatID, "DB error: "+err.Error())
			return
		}
		applied++
	}
	for _, g := range res.Gaps {
		if _, ok := defs.Topics[g.Topic]; !ok {
			continue
		}
		if err := b.db.BumpReviewDue(ctx, chatID, g.Topic); err != nil {
			b.reply(ctx, chatID, "DB error: "+err.Error())
			return
		}
	}
	for _, term := range res.NewGlossaryTerms {
		if err := b.db.InsertGlossary(ctx, chatID, term.Term, term.Context); err != nil {
			b.reply(ctx, chatID, "DB error: "+err.Error())
			return
		}
	}
	b.awardXP(ctx, chatID, "debrief", 10)
	b.clearConv(chatID)

	var sb strings.Builder
	sb.WriteString("📝 Debrief imported.\n")
	if len(res.TopicsTouched) > 0 {
		fmt.Fprintf(&sb, "Topics: %s\n", strings.Join(res.TopicsTouched, ", "))
	}
	fmt.Fprintf(&sb, "Gaps: %d\n", len(res.Gaps))
	for _, g := range res.Gaps {
		fmt.Fprintf(&sb, "  • [%s] %s: %s\n", g.Severity, g.Topic, g.Detail)
	}
	fmt.Fprintf(&sb, "Confidence updated: %d · glossary added: %d", applied, len(res.NewGlossaryTerms))
	b.reply(ctx, chatID, sb.String())
}

func (b *Bot) applyDrill(ctx context.Context, chatID int64, conv *conversation, raw []byte) {
	score, err := claude.ParseDrillScore(raw)
	if err != nil {
		b.badJSON(ctx, chatID, err)
		return
	}
	if err := b.db.InsertDrill(ctx, chatID, conv.drillKind, score.Outcome, &score.Score); err != nil {
		b.reply(ctx, chatID, "DB error: "+err.Error())
		return
	}
	b.awardXP(ctx, chatID, "drill", 10)
	b.clearConv(chatID)
	b.reply(ctx, chatID, fmt.Sprintf("Logged %s drill: %d/5 — %s", conv.drillKind, score.Score, score.Outcome))
}

func (b *Bot) applyStory(ctx context.Context, chatID int64, conv *conversation, raw []byte) {
	story, err := claude.ParseStory(raw)
	if err != nil {
		b.badJSON(ctx, chatID, err)
		return
	}
	if !contains(story.Competencies, conv.storyComp) {
		story.Competencies = append(story.Competencies, conv.storyComp)
	}
	if err := b.db.InsertStory(ctx, db.Story{
		UserID:       chatID,
		Title:        story.Title,
		Situation:    story.Situation,
		Task:         story.Task,
		Action:       story.Action,
		Result:       story.Result,
		Metrics:      story.Metrics,
		TechTags:     story.TechTags,
		Competencies: story.Competencies,
	}); err != nil {
		b.reply(ctx, chatID, "DB error: "+err.Error())
		return
	}
	b.clearConv(chatID)
	b.reply(ctx, chatID, fmt.Sprintf("✅ Saved: %s\nMetric: %s\nCompetencies: %s",
		story.Title, strings.Join(story.Metrics, "; "), strings.Join(story.Competencies, ", ")))
}
