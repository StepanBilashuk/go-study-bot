package tgbot

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// competencies is the fixed STAR competency set (spec §4).
var competencies = []string{"ownership", "mentoring", "conflict", "failure", "tradeoff"}

// handleStory emits a story-mining prompt targeting the user's weakest-covered
// competency. The user runs the interview in their Claude, which extracts a STAR
// story to JSON; applyStory saves it, rejecting metric-free stories (spec §2,§3).
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

	prompt, err := b.getPrompts().Render("story", map[string]string{"competency": target})
	if err != nil {
		b.reply(ctx, chatID, "Prompt error: "+err.Error())
		return
	}

	b.setConv(chatID, &conversation{mode: modeAwaitImport, kind: importStory, storyComp: target})
	b.reply(ctx, chatID, emitHeader("Story mining — "+target)+prompt)
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

// weakestCompetency returns the competency with the fewest stories.
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
