package tgbot

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"prepbot/internal/claude"
)

// quizGate is the stage-2 threshold: 8/10 correct (spec §6).
const quizGate = 8

// handleQuiz (Phase 2) generates a 10-item pattern-recognition quiz for a topic
// and walks it one item at a time.
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
		b.reply(ctx, chatID, fmt.Sprintf("Unknown topic %q.", slug))
		return
	}

	system, err := b.getPrompts().Render("quiz", map[string]string{
		"topic": topic.Name,
		"track": string(topic.Track),
	})
	if err != nil {
		b.reply(ctx, chatID, "Prompt error: "+err.Error())
		return
	}

	var quiz claude.Quiz
	err = b.claude.CompleteJSON(ctx, system, "Generate the quiz now.", func(raw []byte) error {
		q, e := claude.ParseQuiz(raw)
		quiz = q
		return e
	})
	if err != nil {
		b.reply(ctx, chatID, "Could not generate the quiz: "+err.Error())
		return
	}

	b.setConv(chatID, &conversation{
		mode:      modeQuizzing,
		quizTopic: slug,
		quizItems: quiz.Items,
	})
	b.reply(ctx, chatID, fmt.Sprintf("🧠 Quiz: %s — %d items, ~60s each. Reply A/B/C/D. No solving, just recognition.\n\n%s",
		topic.Name, len(quiz.Items), formatQuizItem(quiz.Items[0], 0, len(quiz.Items))))
}

// handleQuizAnswer grades one answer and advances, or finishes with a score.
func (b *Bot) handleQuizAnswer(ctx context.Context, chatID int64, conv *conversation, answer string) {
	if conv.quizIdx >= len(conv.quizItems) {
		b.clearConv(chatID)
		return
	}
	item := conv.quizItems[conv.quizIdx]
	choice, ok := parseChoice(answer, len(item.Options))
	if !ok {
		b.reply(ctx, chatID, fmt.Sprintf("Reply with a letter A-%c.", 'A'+rune(len(item.Options)-1)))
		return
	}

	var fb string
	if choice == item.Answer {
		conv.quizScore++
		fb = "✅ Correct."
	} else {
		fb = fmt.Sprintf("❌ Answer: %c) %s", 'A'+rune(item.Answer), item.Options[item.Answer])
	}
	fb += " " + item.Explanation

	conv.quizIdx++
	if conv.quizIdx >= len(conv.quizItems) {
		defs := b.getDefs()
		name := conv.quizTopic
		if t, ok := defs.Topics[conv.quizTopic]; ok {
			name = t.Name
		}
		total := len(conv.quizItems)
		msg := fmt.Sprintf("%s\n\n🏁 %s quiz: %d/%d.", fb, name, conv.quizScore, total)
		if conv.quizScore >= quizGate {
			msg += fmt.Sprintf("\nQuiz gate met (≥%d/10) — run /done %s to advance.", quizGate, conv.quizTopic)
		} else {
			msg += "\nBelow the 8/10 gate — worth another pass."
		}
		b.clearConv(chatID)
		b.reply(ctx, chatID, msg)
		return
	}

	b.setConv(chatID, conv)
	b.reply(ctx, chatID, fmt.Sprintf("%s\n\n%s", fb, formatQuizItem(conv.quizItems[conv.quizIdx], conv.quizIdx, len(conv.quizItems))))
}

func formatQuizItem(item claude.QuizItem, idx, total int) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Q%d/%d: %s\n", idx+1, total, item.Question)
	for i, opt := range item.Options {
		fmt.Fprintf(&sb, "%c) %s\n", 'A'+rune(i), opt)
	}
	return strings.TrimRight(sb.String(), "\n")
}

// parseChoice maps a reply like "b", "B", "2", or "B) ..." to a 0-based index.
func parseChoice(s string, n int) (int, bool) {
	s = strings.TrimSpace(strings.ToUpper(s))
	if s == "" {
		return 0, false
	}
	c := s[0]
	if c >= 'A' && c < 'A'+byte(n) {
		return int(c - 'A'), true
	}
	if c >= '1' && c <= '9' {
		i := int(c - '1')
		if i < n {
			return i, true
		}
	}
	return 0, false
}
