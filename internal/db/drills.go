package db

import (
	"context"
	"fmt"
	"time"
)

// DrillStat aggregates one user's drill_log for one drill kind.
type DrillStat struct {
	Kind       string
	AvgScore   *float64 // nil when no scored attempts yet
	Count      int
	LastServed *time.Time
}

// InsertDrill records that a user was served/scored a drill of the given kind.
// score may be nil when a drill is served but not yet scored.
func (db *DB) InsertDrill(ctx context.Context, userID int64, kind, outcome string, score *int) error {
	_, err := db.Exec(ctx,
		`insert into drill_log (user_id, drill_slug, outcome, score) values ($1, $2, $3, $4)`,
		userID, kind, outcome, score)
	if err != nil {
		return fmt.Errorf("insert drill %d/%s: %w", userID, kind, err)
	}
	return nil
}

// DrillStats returns a user's per-kind averages and recency, keyed by kind.
// Kinds with no rows are absent from the map.
func (db *DB) DrillStats(ctx context.Context, userID int64) (map[string]DrillStat, error) {
	rows, err := db.Query(ctx,
		`select drill_slug, avg(score)::float8, count(*), max(date)
		 from drill_log where user_id = $1 group by drill_slug`, userID)
	if err != nil {
		return nil, fmt.Errorf("drill stats %d: %w", userID, err)
	}
	defer rows.Close()

	out := make(map[string]DrillStat)
	for rows.Next() {
		var s DrillStat
		if err := rows.Scan(&s.Kind, &s.AvgScore, &s.Count, &s.LastServed); err != nil {
			return nil, fmt.Errorf("scan drill stat: %w", err)
		}
		out[s.Kind] = s
	}
	return out, rows.Err()
}
