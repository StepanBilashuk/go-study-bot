package db

import (
	"context"
	"fmt"
	"time"
)

// AwardXP records an XP event (spec §4, Phase 5). kind is one of
// topic_closed | debrief | boss | drill.
func (db *DB) AwardXP(ctx context.Context, userID int64, kind string, points int) error {
	_, err := db.Exec(ctx,
		`insert into xp_events (user_id, kind, points) values ($1, $2, $3)`,
		userID, kind, points)
	if err != nil {
		return fmt.Errorf("award xp %d/%s: %w", userID, kind, err)
	}
	return nil
}

// TotalXP returns a user's summed XP.
func (db *DB) TotalXP(ctx context.Context, userID int64) (int, error) {
	var total int
	err := db.QueryRow(ctx,
		`select coalesce(sum(points), 0) from xp_events where user_id = $1`, userID,
	).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("total xp %d: %w", userID, err)
	}
	return total, nil
}

// Streak returns the number of active days in the user's current streak, the
// number of missed days bridged by freezes, and the user's freeze budget
// (spec §16). A day is "active" if the user earned XP on it. Walking back from
// today: today with no activity yet does not break the streak or spend a
// freeze; each earlier missed day spends one freeze, and the streak ends when
// the budget runs out.
func (db *DB) Streak(ctx context.Context, userID int64) (streak, freezesUsed, freezesTotal int, err error) {
	if err = db.QueryRow(ctx,
		`select coalesce(streak_freezes, 2) from users where user_id = $1`, userID,
	).Scan(&freezesTotal); err != nil {
		return 0, 0, 0, fmt.Errorf("streak freezes %d: %w", userID, err)
	}

	rows, err := db.Query(ctx,
		`select distinct (created_at at time zone 'UTC')::date
		 from xp_events where user_id = $1 order by 1 desc`, userID)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("streak %d: %w", userID, err)
	}
	defer rows.Close()

	active := make(map[string]bool)
	for rows.Next() {
		var d time.Time
		if err := rows.Scan(&d); err != nil {
			return 0, 0, 0, fmt.Errorf("scan streak day: %w", err)
		}
		active[d.Format("2006-01-02")] = true
	}
	if err := rows.Err(); err != nil {
		return 0, 0, 0, err
	}
	if len(active) == 0 {
		return 0, 0, freezesTotal, nil
	}

	const day = 24 * time.Hour
	freezesLeft := freezesTotal
	today := time.Now().UTC().Truncate(day)
	for i := 0; i < 1000; i++ {
		d := today.Add(-time.Duration(i) * day)
		if active[d.Format("2006-01-02")] {
			streak++
			continue
		}
		if i == 0 {
			continue // today not done yet — neither breaks nor spends a freeze
		}
		if freezesLeft > 0 {
			freezesLeft--
			freezesUsed++
			continue
		}
		break
	}
	return streak, freezesUsed, freezesTotal, nil
}
