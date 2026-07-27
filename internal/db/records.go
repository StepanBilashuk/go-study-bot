package db

import (
	"context"
	"fmt"
)

// InsertSession stores a user's raw /debrief and Claude's extracted gaps.
func (db *DB) InsertSession(ctx context.Context, userID int64, rawText string, extractedGaps []byte, topicsTouched []string) error {
	_, err := db.Exec(ctx,
		`insert into sessions (user_id, raw_text, extracted_gaps, topics_touched) values ($1, $2, $3, $4)`,
		userID, rawText, extractedGaps, topicsTouched)
	if err != nil {
		return fmt.Errorf("insert session %d: %w", userID, err)
	}
	return nil
}

// InsertGlossary stores an English term a user could not produce.
func (db *DB) InsertGlossary(ctx context.Context, userID int64, term, context string) error {
	_, err := db.Exec(ctx,
		`insert into glossary (user_id, term, context) values ($1, $2, $3)`, userID, term, context)
	if err != nil {
		return fmt.Errorf("insert glossary %d/%q: %w", userID, term, err)
	}
	return nil
}

// GlossaryEntry is one collected English term.
type GlossaryEntry struct {
	Term    string
	Context string
}

// ListGlossary returns a user's glossary terms, newest first, up to limit.
func (db *DB) ListGlossary(ctx context.Context, userID int64, limit int) ([]GlossaryEntry, error) {
	rows, err := db.Query(ctx,
		`select term, context from glossary where user_id = $1 order by added_at desc limit $2`,
		userID, limit)
	if err != nil {
		return nil, fmt.Errorf("list glossary %d: %w", userID, err)
	}
	defer rows.Close()

	var out []GlossaryEntry
	for rows.Next() {
		var e GlossaryEntry
		if err := rows.Scan(&e.Term, &e.Context); err != nil {
			return nil, fmt.Errorf("scan glossary: %w", err)
		}
		out = append(out, e)
	}
	return out, rows.Err()
}
