package scheduler

import (
	"testing"
	"time"
)

func TestNextReview(t *testing.T) {
	now := time.Date(2026, 7, 27, 8, 0, 0, 0, time.UTC)
	tests := []struct {
		completed int
		wantDays  int
	}{
		{0, 1},   // entering SR → first review in 1 day
		{1, 3},   // after 1 review → 3 days
		{2, 7},   // after 2 → 7 days
		{3, 21},  // after 3 → 21 days
		{4, 21},  // clamps to the last interval
		{99, 21}, // stays clamped
		{-5, 1},  // negatives treated as 0
	}
	for _, tt := range tests {
		got := NextReview(now, tt.completed)
		want := now.Add(time.Duration(tt.wantDays) * 24 * time.Hour)
		if !got.Equal(want) {
			t.Errorf("NextReview(now, %d) = %v, want %v (%d days)", tt.completed, got, want, tt.wantDays)
		}
	}
}

func TestIntervalDays(t *testing.T) {
	got := IntervalDays()
	want := []int{1, 3, 7, 21}
	if len(got) != len(want) {
		t.Fatalf("IntervalDays() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("IntervalDays() = %v, want %v", got, want)
		}
	}
}
