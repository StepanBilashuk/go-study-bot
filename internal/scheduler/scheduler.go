// Package scheduler holds the Phase 1 spaced-repetition schedule. Intervals are
// hardcoded at 1, 3, 7, 21 days (spec §6, §14). This is one of only two things
// the project unit-tests, so it is kept pure and dependency-free.
package scheduler

import "time"

const day = 24 * time.Hour

// intervals is the hardcoded Phase 1 ladder. A topic that reaches stage 4
// enters spaced repetition and returns at these spacings.
var intervals = []time.Duration{1 * day, 3 * day, 7 * day, 21 * day}

// IntervalDays returns the hardcoded interval ladder in days.
func IntervalDays() []int {
	out := make([]int, len(intervals))
	for i, d := range intervals {
		out[i] = int(d / day)
	}
	return out
}

// NextReview returns when a spaced-repetition topic should next surface, given
// how many reviews it has already completed. completed==0 schedules the first
// review 1 day out; the sequence is 1, 3, 7, 21 days and then clamps to 21.
// Negative inputs are treated as 0.
func NextReview(now time.Time, completed int) time.Time {
	if completed < 0 {
		completed = 0
	}
	i := completed
	if i >= len(intervals) {
		i = len(intervals) - 1
	}
	return now.Add(intervals[i])
}
