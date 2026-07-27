package db

import (
	"context"
	"fmt"
)

// UpsertUser records a user on first contact and refreshes username/last_seen
// on every interaction. Called before any per-user write so foreign keys hold.
func (db *DB) UpsertUser(ctx context.Context, userID int64, username string) error {
	_, err := db.Exec(ctx,
		`insert into users (user_id, username, last_seen)
		 values ($1, $2, now())
		 on conflict (user_id) do update set
		   username = excluded.username,
		   last_seen = now()`,
		userID, username)
	if err != nil {
		return fmt.Errorf("upsert user %d: %w", userID, err)
	}
	return nil
}

// PushUsers returns the ids of users who opted in to the daily push.
func (db *DB) PushUsers(ctx context.Context) ([]int64, error) {
	rows, err := db.Query(ctx, `select user_id from users where push_enabled = true`)
	if err != nil {
		return nil, fmt.Errorf("list push users: %w", err)
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan user id: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// SetPushEnabled toggles the daily push for a user.
func (db *DB) SetPushEnabled(ctx context.Context, userID int64, enabled bool) error {
	_, err := db.Exec(ctx, `update users set push_enabled = $2 where user_id = $1`, userID, enabled)
	if err != nil {
		return fmt.Errorf("set push_enabled %d: %w", userID, err)
	}
	return nil
}
