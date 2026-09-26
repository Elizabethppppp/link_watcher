package pgService

import (
	"context"
	"database/sql"
	"link_watcher/serviceErrors"
)

type PgService struct {
	db *sql.DB
}

func NewPgService(db *sql.DB) *PgService {
	return &PgService{
		db: db,
	}
}

func (pg *PgService) Insert(ctx context.Context, url string, intervalSec int64) error {
	_, err := pg.db.ExecContext(ctx, "INSERT INTO target (url, interval_sec) VALUES ($1, $2) RETURNING id, url, is_tracking ,interval_sec, created_at, updated_at")
	if err != nil {
		return serviceErrors.ErrInternal
	}
	return nil
}
