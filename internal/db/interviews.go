package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// Interview is a scheduled interview date for a company (Phase 4).
type Interview struct {
	CompanySlug string
	Date        time.Time
}

// SetInterview records or updates a user's interview date for a company.
func (db *DB) SetInterview(ctx context.Context, userID int64, companySlug string, date time.Time) error {
	_, err := db.Exec(ctx,
		`insert into interviews (user_id, company_slug, date) values ($1, $2, $3)
		 on conflict (user_id, company_slug) do update set date = excluded.date`,
		userID, companySlug, date)
	if err != nil {
		return fmt.Errorf("set interview %d/%s: %w", userID, companySlug, err)
	}
	return nil
}

// ClearInterview removes a user's interview date for a company.
func (db *DB) ClearInterview(ctx context.Context, userID int64, companySlug string) error {
	_, err := db.Exec(ctx,
		`delete from interviews where user_id = $1 and company_slug = $2`, userID, companySlug)
	if err != nil {
		return fmt.Errorf("clear interview %d/%s: %w", userID, companySlug, err)
	}
	return nil
}

// NextInterview returns a user's nearest upcoming interview (date >= today), or
// (zero, false) if none.
func (db *DB) NextInterview(ctx context.Context, userID int64) (Interview, bool, error) {
	var iv Interview
	err := db.QueryRow(ctx,
		`select company_slug, date from interviews
		 where user_id = $1 and date >= current_date
		 order by date asc limit 1`, userID,
	).Scan(&iv.CompanySlug, &iv.Date)
	if errors.Is(err, pgx.ErrNoRows) {
		return Interview{}, false, nil
	}
	if err != nil {
		return Interview{}, false, fmt.Errorf("next interview %d: %w", userID, err)
	}
	return iv, true, nil
}
