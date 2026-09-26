package pgService

import (
	"context"
	"database/sql"
	"errors"
	"link_watcher/serviceErrors"

	"github.com/jackc/pgx/v5/pgconn"
)

type PgService struct {
	db *sql.DB
}

func NewPgService(db *sql.DB) *PgService {
	return &PgService{
		db: db,
	}
}

func (pg *PgService) Insert(ctx context.Context, url string, intervalSec int64) (Target, error) {
	query := `INSERT INTO target (url, interval_sec) VALUES ($1, $2) RETURNING id, url, is_tracking ,interval_sec, created_at, updated_at`

	var target Target
	err := pg.db.QueryRowContext(ctx, query, url, intervalSec).Scan(
		&target.Id,
		&target.Url,
		&target.IsTracking,
		&target.IntervalSec,
		&target.CreatedAt,
		&target.UpdatedAt)
	if err != nil {
		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) && pgerr.Code == "23505" {
			return Target{}, serviceErrors.ErrConflict
		}
	}
	return target, nil
}
