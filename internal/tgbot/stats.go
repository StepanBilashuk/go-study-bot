package tgbot

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"prepbot/internal/db"
	"prepbot/internal/definitions"
)

// handleStats (Phase 5, behind the GAMIFICATION flag) shows XP, level, streak,
// per-track mastery bars and book coverage (spec §7).
func (b *Bot) handleStats(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	chatID := update.Message.Chat.ID
	if !b.cfg.Gamification {
		b.reply(ctx, chatID, "Stats are off. Set GAMIFICATION=true to enable XP, levels and streaks.")
		return
	}

	defs := b.getDefs()
	prog, err := b.db.ListProgress(ctx, chatID)
	if err != nil {
		b.reply(ctx, chatID, "DB error: "+err.Error())
		return
	}
	xp, err := b.db.TotalXP(ctx, chatID)
	if err != nil {
		b.reply(ctx, chatID, "DB error: "+err.Error())
		return
	}
	streak, frozen, freezeBudget, err := b.db.Streak(ctx, chatID)
	if err != nil {
		b.reply(ctx, chatID, "DB error: "+err.Error())
		return
	}

	var sb strings.Builder
	sb.WriteString("🏆 Stats\n")
	fmt.Fprintf(&sb, "XP: %d (level %d)\n", xp, 1+xp/500)
	fmt.Fprintf(&sb, "🔥 Streak: %d day(s)  ❄️ freezes %d/%d used\n", streak, frozen, freezeBudget)

	for _, tr := range []struct {
		track definitions.Track
		label string
	}{
		{definitions.TrackAlgorithms, "Algorithms"},
		{definitions.TrackSystemDesign, "System design"},
	} {
		mastered, total := masteryCount(defs, prog, tr.track)
		fmt.Fprintf(&sb, "%s: %s %d/%d mastered\n", tr.label, bar(mastered, total), mastered, total)
	}

	fmt.Fprintf(&sb, "Books: %s", bookCoverage(defs, prog))
	b.reply(ctx, chatID, strings.TrimRight(sb.String(), "\n"))
}

// masteryCount returns (stage-4 topics, total topics) for a track.
func masteryCount(defs *definitions.Definitions, prog map[string]db.Progress, track definitions.Track) (int, int) {
	mastered, total := 0, 0
	for slug, t := range defs.Topics {
		if t.Track != track {
			continue
		}
		total++
		if prog[slug].Stage >= masteredStage {
			mastered++
		}
	}
	return mastered, total
}

// bar renders a 6-cell progress bar.
func bar(done, total int) string {
	const cells = 6
	filled := 0
	if total > 0 {
		filled = done * cells / total
	}
	return strings.Repeat("█", filled) + strings.Repeat("░", cells-filled)
}

// bookCoverage reports, per book, how many mapped chapters have been started
// (topic at stage ≥ 1).
func bookCoverage(defs *definitions.Definitions, prog map[string]db.Progress) string {
	type cov struct{ covered, total int }
	byBook := map[string]*cov{}
	seen := map[string]bool{} // book|chapter dedupe
	for _, r := range defs.Resources {
		if r.Type != "book" || r.Source == "" {
			continue
		}
		key := fmt.Sprintf("%s|%d", r.Source, r.Chapter)
		if seen[key] {
			continue
		}
		seen[key] = true
		c := byBook[r.Source]
		if c == nil {
			c = &cov{}
			byBook[r.Source] = c
		}
		c.total++
		if prog[r.Topic].Stage >= 1 {
			c.covered++
		}
	}
	if len(byBook) == 0 {
		return "—"
	}
	keys := make([]string, 0, len(byBook))
	for k := range byBook {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var parts []string
	for _, k := range keys {
		name := k
		if bk, ok := defs.Books[k]; ok {
			name = bk.Title
		}
		parts = append(parts, fmt.Sprintf("%s %d/%d ch", name, byBook[k].covered, byBook[k].total))
	}
	return strings.Join(parts, " · ")
}
