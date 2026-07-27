package tgbot

import (
	"sort"
	"time"

	"prepbot/internal/db"
	"prepbot/internal/definitions"
)

// prereqConfidence is the confidence a topic's dependency must reach before the
// dependent topic can be advanced past stage 0. Uses calibration data so a
// topic he already knows doesn't block its successors.
const prereqConfidence = 3

// masteredStage is the top of the ladder; topics here are in spaced repetition.
const masteredStage = 4

// staleWindow bounds how recently a topic must have been touched to count
// toward a boss gate (spec §7: "none older than 10 days").
const staleWindow = 10 * 24 * time.Hour

// topicsByTrack returns a track's topics sorted by priority (NeetCode order).
func topicsByTrack(defs *definitions.Definitions, track definitions.Track) []definitions.Topic {
	var ts []definitions.Topic
	for _, t := range defs.Topics {
		if t.Track == track {
			ts = append(ts, t)
		}
	}
	sort.Slice(ts, func(i, j int) bool { return ts[i].Priority < ts[j].Priority })
	return ts
}

// prereqsMet reports whether every dependency of t is calibrated to at least
// prereqConfidence.
func prereqsMet(t definitions.Topic, prog map[string]db.Progress) bool {
	for _, dep := range t.DependsOn {
		if prog[dep].Confidence < prereqConfidence {
			return false
		}
	}
	return true
}

// pickTopic returns the next topic to work in a track: the highest-priority
// (earliest) topic not yet mastered whose prerequisites are met. When focus is
// non-empty (an upcoming interview), it prefers topics in that set, falling back
// to normal order if none is eligible. Returns (zero, false) if nothing fits.
func pickTopic(defs *definitions.Definitions, prog map[string]db.Progress, track definitions.Track, focus map[string]bool) (definitions.Topic, bool) {
	eligible := func(t definitions.Topic) bool {
		return prog[t.Slug].Stage < masteredStage && prereqsMet(t, prog)
	}
	if len(focus) > 0 {
		for _, t := range topicsByTrack(defs, track) {
			if focus[t.Slug] && eligible(t) {
				return t, true
			}
		}
	}
	for _, t := range topicsByTrack(defs, track) {
		if eligible(t) {
			return t, true
		}
	}
	return definitions.Topic{}, false
}

// pickReview returns the most-overdue mastered topic due for spaced repetition,
// or (zero, false) if none is due.
func pickReview(defs *definitions.Definitions, prog map[string]db.Progress, now time.Time) (definitions.Topic, bool) {
	var best definitions.Topic
	var bestDue time.Time
	found := false
	for slug, p := range prog {
		if p.Stage < masteredStage || p.NextDue == nil || p.NextDue.After(now) {
			continue
		}
		t, ok := defs.Topics[slug]
		if !ok {
			continue
		}
		if !found || p.NextDue.Before(bestDue) {
			best, bestDue, found = t, *p.NextDue, true
		}
	}
	return best, found
}

// resourcesFor returns up to 2 resources for a topic at its current stage,
// falling back to any-stage resources for that topic (spec §11: never more
// than 2).
func resourcesFor(defs *definitions.Definitions, slug string, stage int) []definitions.Resource {
	var atStage, anyStage []definitions.Resource
	for _, r := range defs.Resources {
		if r.Topic != slug {
			continue
		}
		anyStage = append(anyStage, r)
		if r.Stage == stage {
			atStage = append(atStage, r)
		}
	}
	pick := atStage
	if len(pick) == 0 {
		pick = anyStage
	}
	if len(pick) > 2 {
		pick = pick[:2]
	}
	return pick
}

// drillWeaknessRank orders the drill kinds by the candidate's diagnosed
// weakness, used to break ties (spec §10: estimation is #1, contradiction #2).
var drillWeaknessRank = map[string]int{
	string(definitions.DrillEstimation):    0,
	string(definitions.DrillContradiction): 1,
	string(definitions.DrillClarify):       2,
	string(definitions.DrillNextStep):      3,
}

// pickDrill returns the drill to serve: the kind with the weakest recent
// scores. Kinds never scored sort first; ties break by weakness rank, then by
// least-recently served (spec §7 /drill).
func pickDrill(defs *definitions.Definitions, stats map[string]db.DrillStat) (definitions.Drill, bool) {
	drills := make([]definitions.Drill, 0, len(defs.Drills))
	for _, d := range defs.Drills {
		drills = append(drills, d)
	}
	if len(drills) == 0 {
		return definitions.Drill{}, false
	}

	// score: lower is weaker (served first). No data → -1 (weakest).
	avg := func(kind string) float64 {
		s, ok := stats[kind]
		if !ok || s.AvgScore == nil {
			return -1
		}
		return *s.AvgScore
	}
	lastServed := func(kind string) time.Time {
		if s, ok := stats[kind]; ok && s.LastServed != nil {
			return *s.LastServed
		}
		return time.Time{} // zero = never served, sorts first
	}

	sort.Slice(drills, func(i, j int) bool {
		ki, kj := string(drills[i].Kind), string(drills[j].Kind)
		if ai, aj := avg(ki), avg(kj); ai != aj {
			return ai < aj
		}
		if ri, rj := drillWeaknessRank[ki], drillWeaknessRank[kj]; ri != rj {
			return ri < rj
		}
		return lastServed(ki).Before(lastServed(kj))
	})
	return drills[0], true
}

// stageName maps a ladder stage to its label (spec §6).
func stageName(stage int) string {
	names := []string{"Learn", "Guided", "Quiz", "Solo", "Review"}
	if stage < 0 || stage >= len(names) {
		return "?"
	}
	return names[stage]
}
