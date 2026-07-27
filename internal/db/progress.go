package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// Progress is a row of topic_progress for one user (spec §4).
type Progress struct {
	UserID      int64
	TopicSlug   string
	Stage       int
	Confidence  int
	Attempts    int
	LastTouched *time.Time
	NextDue     *time.Time
}

// GetProgress returns a user's progress on a topic, or (zero, false) if none.
func (db *DB) GetProgress(ctx context.Context, userID int64, slug string) (Progress, bool, error) {
	p := Progress{UserID: userID}
	err := db.QueryRow(ctx,
		`select stage, confidence, attempts, last_touched, next_due
		 from topic_progress where user_id = $1 and topic_slug = $2`, userID, slug,
	).Scan(&p.Stage, &p.Confidence, &p.Attempts, &p.LastTouched, &p.NextDue)
	if errors.Is(err, pgx.ErrNoRows) {
		return Progress{}, false, nil
	}
	if err != nil {
		return Progress{}, false, fmt.Errorf("get progress %d/%s: %w", userID, slug, err)
	}
	p.TopicSlug = slug
	return p, true, nil
}

// ListProgress returns all of a user's progress rows keyed by topic slug.
func (db *DB) ListProgress(ctx context.Context, userID int64) (map[string]Progress, error) {
	rows, err := db.Query(ctx,
		`select topic_slug, stage, confidence, attempts, last_touched, next_due
		 from topic_progress where user_id = $1`, userID)
	if err != nil {
		return nil, fmt.Errorf("list progress %d: %w", userID, err)
	}
	defer rows.Close()

	out := make(map[string]Progress)
	for rows.Next() {
		p := Progress{UserID: userID}
		if err := rows.Scan(&p.TopicSlug, &p.Stage, &p.Confidence, &p.Attempts, &p.LastTouched, &p.NextDue); err != nil {
			return nil, fmt.Errorf("scan progress: %w", err)
		}
		out[p.TopicSlug] = p
	}
	return out, rows.Err()
}

// UpsertProgress writes a full progress row (used by /done stage transitions).
func (db *DB) UpsertProgress(ctx context.Context, p Progress) error {
	_, err := db.Exec(ctx,
		`insert into topic_progress (user_id, topic_slug, stage, confidence, attempts, last_touched, next_due)
		 values ($1, $2, $3, $4, $5, $6, $7)
		 on conflict (user_id, topic_slug) do update set
		   stage = excluded.stage,
		   confidence = excluded.confidence,
		   attempts = excluded.attempts,
		   last_touched = excluded.last_touched,
		   next_due = excluded.next_due`,
		p.UserID, p.TopicSlug, p.Stage, p.Confidence, p.Attempts, p.LastTouched, p.NextDue)
	if err != nil {
		return fmt.Errorf("upsert progress %d/%s: %w", p.UserID, p.TopicSlug, err)
	}
	return nil
}

// SetConfidence sets a topic's confidence and last_touched for a user, creating
// the row at stage 0 if absent (used by /start calibration and /debrief).
func (db *DB) SetConfidence(ctx context.Context, userID int64, slug string, confidence int) error {
	_, err := db.Exec(ctx,
		`insert into topic_progress (user_id, topic_slug, stage, confidence, attempts, last_touched)
		 values ($1, $2, 0, $3, 0, now())
		 on conflict (user_id, topic_slug) do update set
		   confidence = excluded.confidence,
		   last_touched = now()`,
		userID, slug, confidence)
	if err != nil {
		return fmt.Errorf("set confidence %d/%s: %w", userID, slug, err)
	}
	return nil
}

// BumpReviewDue brings a user's mastered topic forward to now for review.
func (db *DB) BumpReviewDue(ctx context.Context, userID int64, slug string) error {
	_, err := db.Exec(ctx,
		`update topic_progress set next_due = now()
		 where user_id = $1 and topic_slug = $2 and stage = 4`, userID, slug)
	if err != nil {
		return fmt.Errorf("bump review due %d/%s: %w", userID, slug, err)
	}
	return nil
}
