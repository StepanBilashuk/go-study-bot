package tgbot

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"prepbot/internal/claude"
	"prepbot/internal/db"
)

// competencies is the fixed STAR competency set (spec §4).
var competencies = []string{"ownership", "mentoring", "conflict", "failure", "tradeoff"}

// handleStory (Phase 3) asks a story-mining question targeting the user's
// weakest-covered competency, then waits for the answer.
func (b *Bot) handleStory(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	chatID := update.Message.Chat.ID

	counts, err := b.db.CompetencyCounts(ctx, chatID)
	if err != nil {
		b.reply(ctx, chatID, "DB error: "+err.Error())
		return
	}
	target := weakestCompetency(counts)

	system, err := b.getPrompts().Render("story-mining", map[string]string{"competency": target})
	if err != nil {
		b.reply(ctx, chatID, "Prompt error: "+err.Error())
		return
	}
	question, err := b.claude.Complete(ctx, system, "Ask the question now.")
	if err != nil {
		b.reply(ctx, chatID, "Claude error: "+err.Error())
		return
	}

	b.setConv(chatID, &conversation{mode: modeStoryMining, storyComp: target})
	b.reply(ctx, chatID, fmt.Sprintf("📖 Story mining — %s\n\n%s\n\nReply with what happened (include a number — every story needs a metric).", target, question))
}

// handleStoryAnswer extracts a STAR story and saves it, rejecting metric-free
// stories (spec §2).
func (b *Bot) handleStoryAnswer(ctx context.Context, chatID int64, conv *conversation, answer string) {
	system, err := b.getPrompts().Render("story-extract", nil)
	if err != nil {
		b.reply(ctx, chatID, "Prompt error: "+err.Error())
		return
	}

	var story claude.Story
	err = b.claude.CompleteJSON(ctx, system, answer, func(raw []byte) error {
		s, e := claude.ParseStory(raw)
		story = s
		return e
	})
	if err != nil {
		// A metric-free story fails ParseStory — surface that specific ask.
		b.reply(ctx, chatID, "Not saved: "+err.Error()+"\nAdd a measurable result and send it again, or /story to pick a new prompt.")
		return
	}

	// Make sure the targeted competency is tagged.
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

// handleStories shows the competency matrix — which competencies still have no
// story (spec §7).
func (b *Bot) handleStories(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	chatID := update.Message.Chat.ID
	counts, err := b.db.CompetencyCounts(ctx, chatID)
	if err != nil {
		b.reply(ctx, chatID, "DB error: "+err.Error())
		return
	}

	var sb strings.Builder
	sb.WriteString("🗂 Competency matrix\n")
	missing := 0
	for _, c := range competencies {
		n := counts[c]
		mark := "✅"
		if n == 0 {
			mark = "❌"
			missing++
		}
		fmt.Fprintf(&sb, "%s %s: %d\n", mark, c, n)
	}
	if missing > 0 {
		fmt.Fprintf(&sb, "\n%d competencies still have no story — run /story.", missing)
	} else {
		sb.WriteString("\nEvery competency is covered. 💪")
	}
	b.reply(ctx, chatID, strings.TrimRight(sb.String(), "\n"))
}

// weakestCompetency returns the competency with the fewest stories (ties broken
// by the spec's ordering, which front-loads the gaps he lacks).
func weakestCompetency(counts map[string]int) string {
	best, bestN := competencies[0], counts[competencies[0]]
	for _, c := range competencies {
		if counts[c] < bestN {
			best, bestN = c, counts[c]
		}
	}
	return best
}

func contains(ss []string, s string) bool {
	for _, x := range ss {
		if x == s {
			return true
		}
	}
	return false
}
