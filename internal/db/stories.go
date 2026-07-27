package db

import (
	"context"
	"fmt"
	"time"
)

// Story is a STAR story row (spec §4, Phase 3).
type Story struct {
	ID           int64
	UserID       int64
	Title        string
	Situation    string
	Task         string
	Action       string
	Result       string
	Metrics      []string
	TechTags     []string
	Competencies []string
	Strength     int
	LastUsed     *time.Time
}

// InsertStory saves a mined STAR story. Callers must have already enforced the
// "every story needs a metric" rule (spec §2).
func (db *DB) InsertStory(ctx context.Context, s Story) error {
	_, err := db.Exec(ctx,
		`insert into stories
		   (user_id, title, situation, task, action, result, metrics, tech_tags, competencies, strength)
		 values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		s.UserID, s.Title, s.Situation, s.Task, s.Action, s.Result,
		s.Metrics, s.TechTags, s.Competencies, 3)
	if err != nil {
		return fmt.Errorf("insert story %d: %w", s.UserID, err)
	}
	return nil
}

// ListStories returns a user's stories.
func (db *DB) ListStories(ctx context.Context, userID int64) ([]Story, error) {
	rows, err := db.Query(ctx,
		`select id, title, situation, task, action, result, metrics, tech_tags, competencies, strength, last_used
		 from stories where user_id = $1 order by id`, userID)
	if err != nil {
		return nil, fmt.Errorf("list stories %d: %w", userID, err)
	}
	defer rows.Close()

	var out []Story
	for rows.Next() {
		s := Story{UserID: userID}
		if err := rows.Scan(&s.ID, &s.Title, &s.Situation, &s.Task, &s.Action, &s.Result,
			&s.Metrics, &s.TechTags, &s.Competencies, &s.Strength, &s.LastUsed); err != nil {
			return nil, fmt.Errorf("scan story: %w", err)
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// CompetencyCounts returns how many of a user's stories cover each competency.
func (db *DB) CompetencyCounts(ctx context.Context, userID int64) (map[string]int, error) {
	stories, err := db.ListStories(ctx, userID)
	if err != nil {
		return nil, err
	}
	counts := make(map[string]int)
	for _, s := range stories {
		for _, c := range s.Competencies {
			counts[c]++
		}
	}
	return counts, nil
}
