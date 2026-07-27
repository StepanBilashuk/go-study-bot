package tgbot

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"prepbot/internal/db"
	"prepbot/internal/definitions"
)

// bossGateSize is the number of confident, fresh topics needed to unlock a boss
// (spec §7: "at least 5 topics in the block at confidence ≥ 4, none older than
// 10 days").
const bossGateSize = 5
const bossGateConfidence = 4

// handleBoss checks readiness and, if the gate is met, emits a paste-ready mock
// brief. Default block is system design; "/boss behavioral" emits a behavioral
// brief (spec §7).
func (b *Bot) handleBoss(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	chatID := update.Message.Chat.ID
	arg := strings.ToLower(commandArg(update.Message.Text, "/boss"))

	prog, err := b.db.ListProgress(ctx, chatID)
	if err != nil {
		b.reply(ctx, chatID, "DB error: "+err.Error())
		return
	}
	defs := b.getDefs()

	if arg == "behavioral" {
		b.reply(ctx, chatID, b.behavioralBrief(ctx, chatID, defs))
		return
	}

	// System-design boss.
	track := definitions.TrackSystemDesign
	topics := topicsByTrack(defs, track)
	now := time.Now()

	var eligible, shortfall []string
	for _, t := range topics {
		p := prog[t.Slug]
		fresh := p.LastTouched != nil && now.Sub(*p.LastTouched) <= staleWindow
		if p.Confidence >= bossGateConfidence && fresh {
			eligible = append(eligible, t.Slug)
		} else {
			shortfall = append(shortfall, fmt.Sprintf("%s %d/%d", t.Slug, p.Confidence, bossGateConfidence))
		}
	}

	if len(eligible) < bossGateSize {
		var sb strings.Builder
		fmt.Fprintf(&sb, "🔒 System-design boss locked: %d/%d topics ready.\n", len(eligible), bossGateSize)
		sb.WriteString("Missing (confidence ≥4, touched within 10 days):\n")
		for _, s := range shortfall {
			fmt.Fprintf(&sb, "  • %s\n", s)
		}
		b.reply(ctx, chatID, strings.TrimRight(sb.String(), "\n"))
		return
	}

	brief, err := b.sdBrief(defs, prog, topics)
	if err != nil {
		b.reply(ctx, chatID, "Prompt error: "+err.Error())
		return
	}
	b.awardXP(ctx, chatID, "boss", 50)
	b.reply(ctx, chatID, "🟢 System-design boss unlocked. Paste this into Claude:\n\n"+brief)
}

// sdBrief renders the system-design mock, deep-diving the weakest topics.
func (b *Bot) sdBrief(defs *definitions.Definitions, prog map[string]db.Progress, topics []definitions.Topic) (string, error) {
	weak := weakestTopics(prog, topics, 3)
	company := pickCompanyForTrack(defs, definitions.TrackSystemDesign)
	task := "Design the backend for a ride-hailing platform: ingest location events from 200k active drivers and let operations query recent activity by city region in near-real-time."
	return b.getPrompts().Render("mock-sd", map[string]string{
		"company":     company,
		"task":        task,
		"weak_topics": strings.Join(weak, ", "),
	})
}

// behavioralBrief renders the behavioral mock, targeting the competencies the
// user has NO story for yet, so the mock probes exactly the STAR-bank gaps
// (spec §3). If every competency is covered, it probes the full set.
func (b *Bot) behavioralBrief(ctx context.Context, userID int64, defs *definitions.Definitions) string {
	company, values := "a European product company", "ownership, autonomy"
	// Prefer a target company that actually declares values.
	slugs := make([]string, 0, len(defs.Companies))
	for s := range defs.Companies {
		slugs = append(slugs, s)
	}
	sort.Strings(slugs)
	for _, s := range slugs {
		c := defs.Companies[s]
		if len(c.Values) > 0 {
			company, values = c.Name, strings.Join(c.Values, ", ")
			break
		}
	}

	// Target the uncovered competencies from the STAR bank.
	counts, err := b.db.CompetencyCounts(ctx, userID)
	if err != nil {
		return "DB error: " + err.Error()
	}
	var missing, covered []string
	for _, c := range competencies {
		if counts[c] == 0 {
			missing = append(missing, c)
		} else {
			covered = append(covered, c)
		}
	}
	probe := missing
	header := "🎭 Behavioral boss — probing your UNCOVERED competencies (no story yet)."
	if len(probe) == 0 {
		probe = competencies
		header = "🎭 Behavioral boss — every competency has a story; probing the full set for depth."
	}

	brief, err := b.getPrompts().Render("mock-behavioral", map[string]string{
		"company":      company,
		"values":       values,
		"competencies": strings.Join(probe, ", "),
	})
	if err != nil {
		return "Prompt error: " + err.Error()
	}

	var sb strings.Builder
	sb.WriteString(header + "\n")
	if len(covered) > 0 {
		fmt.Fprintf(&sb, "Already have stories for: %s\n", strings.Join(covered, ", "))
	}
	sb.WriteString("\nPaste this into Claude:\n\n" + brief)
	return sb.String()
}

// weakestTopics returns the n lowest-confidence topic slugs in a track.
func weakestTopics(prog map[string]db.Progress, topics []definitions.Topic, n int) []string {
	sorted := append([]definitions.Topic(nil), topics...)
	sort.SliceStable(sorted, func(i, j int) bool {
		return prog[sorted[i].Slug].Confidence < prog[sorted[j].Slug].Confidence
	})
	var out []string
	for i := 0; i < len(sorted) && i < n; i++ {
		out = append(out, sorted[i].Slug)
	}
	return out
}

// pickCompanyForTrack returns the name of the target company with the most
// required topics in the given track (best domain fit), or a generic label.
func pickCompanyForTrack(defs *definitions.Definitions, track definitions.Track) string {
	inTrack := map[string]bool{}
	for slug, t := range defs.Topics {
		if t.Track == track {
			inTrack[slug] = true
		}
	}
	best, bestCount := "a European product company", 0
	slugs := make([]string, 0, len(defs.Companies))
	for s := range defs.Companies {
		slugs = append(slugs, s)
	}
	sort.Strings(slugs) // deterministic tie-break
	for _, s := range slugs {
		c := defs.Companies[s]
		count := 0
		for _, rt := range c.RequiredTopics {
			if inTrack[rt] {
				count++
			}
		}
		if count > bestCount {
			best, bestCount = c.Name, count
		}
	}
	return best
}
