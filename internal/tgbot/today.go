package tgbot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"prepbot/internal/db"
	"prepbot/internal/definitions"
)

func (b *Bot) handleToday(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	userID := update.Message.Chat.ID
	text, slugs, err := b.buildToday(ctx, userID)
	if err != nil {
		b.reply(ctx, userID, "Could not build today: "+err.Error())
		return
	}
	// Topics are tappable → their learning card. (The quick-bar reply keyboard
	// set earlier via /menu or /help persists underneath.)
	b.replyInline(ctx, userID, text, learnKeyboard(b.getDefs(), slugs))
}

// buildToday assembles one user's daily plan: at most 3 topics (algorithms,
// system design, review) + 1 drill + up to 2 resources per topic (spec §7). The
// 3-slot structure enforces the "never more than 3 topics" rule. It also returns
// the slugs of the topics shown, for the tappable Learn buttons.
func (b *Bot) buildToday(ctx context.Context, userID int64) (string, []string, error) {
	defs := b.getDefs()
	prog, err := b.db.ListProgress(ctx, userID)
	if err != nil {
		return "", nil, err
	}
	stats, err := b.db.DrillStats(ctx, userID)
	if err != nil {
		return "", nil, err
	}
	now := time.Now()

	// Focus set: if an interview is coming up, prefer that company's topics.
	var focus map[string]bool
	var focusNote string
	if iv, ok, err := b.db.NextInterview(ctx, userID); err == nil && ok {
		if c, ok := defs.Companies[iv.CompanySlug]; ok {
			focus = make(map[string]bool, len(c.RequiredTopics))
			for _, rt := range c.RequiredTopics {
				focus[rt] = true
			}
			days := int(time.Until(iv.Date).Hours()/24 + 0.5)
			focusNote = fmt.Sprintf(" (focus: %s in %d days)", c.Name, days)
		}
	}

	var sb strings.Builder
	sb.WriteString("📅 Today" + focusNote + "\n")

	var slugs []string
	if t, ok := pickTopic(defs, prog, definitions.TrackAlgorithms, focus); ok {
		writeTopic(&sb, defs, prog, t, "Algorithms")
		slugs = append(slugs, t.Slug)
	}
	if t, ok := pickTopic(defs, prog, definitions.TrackSystemDesign, focus); ok {
		writeTopic(&sb, defs, prog, t, "System design")
		slugs = append(slugs, t.Slug)
	}
	if t, ok := pickReview(defs, prog, now); ok {
		writeTopic(&sb, defs, prog, t, "Review")
		slugs = append(slugs, t.Slug)
	}

	if d, ok := pickDrill(defs, stats); ok {
		fmt.Fprintf(&sb, "\n🎯 Drill: %s (%s, %d min) — run /drill\n", d.Name, d.Kind, d.DurationMin)
	}
	sb.WriteString("\nTap a topic below for theory + best practices.")
	return strings.TrimRight(sb.String(), "\n"), slugs, nil
}

func writeTopic(sb *strings.Builder, defs *definitions.Definitions, prog map[string]db.Progress, t definitions.Topic, slot string) {
	p := prog[t.Slug]
	fmt.Fprintf(sb, "\n▸ %s: %s — stage %d (%s)\n", slot, t.Name, p.Stage, stageName(p.Stage))
	fmt.Fprintf(sb, "   gate: %s\n", t.Gate)
	for _, r := range resourcesFor(defs, t.Slug, p.Stage) {
		writeResource(sb, defs, r)
	}
}

func writeResource(sb *strings.Builder, defs *definitions.Definitions, r definitions.Resource) {
	label := r.Section
	if book, ok := defs.Books[r.Source]; ok && r.Type == "book" {
		if r.Chapter > 0 {
			fmt.Fprintf(sb, "   • %s ch.%d — %s", book.Title, r.Chapter, r.Section)
		} else {
			fmt.Fprintf(sb, "   • %s — %s", book.Title, r.Section)
		}
	} else {
		fmt.Fprintf(sb, "   • %s", label)
	}
	if r.Pages != nil && *r.Pages != "" {
		fmt.Fprintf(sb, " (pp. %s)", *r.Pages)
	}
	if r.EstMin > 0 {
		fmt.Fprintf(sb, " ~%dm", r.EstMin)
	}
	sb.WriteString("\n")
}
