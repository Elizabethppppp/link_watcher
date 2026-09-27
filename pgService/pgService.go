package pgService

import (
	"context"
	"database/sql"
	"errors"
	"link_watcher/model"
	"link_watcher/serviceErrors"

	"github.com/jackc/pgx/v5"
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

func (pg *PgService) Insert(ctx context.Context, url string, intervalSec int64) (model.Target, error) {
	query := `INSERT INTO target (url, interval_sec) VALUES ($1, $2) RETURNING id, url, is_tracking ,interval_sec, created_at, updated_at`

	var target model.Target
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
			return model.Target{}, serviceErrors.ErrConflict
		}
		return model.Target{}, serviceErrors.ErrInternal
	}
	return target, nil
}

func (pg *PgService) GetAllTargets(ctx context.Context) ([]model.Target, error) {
	query := `SELECT * FROM target`
	rows, err := pg.db.QueryContext(ctx, query)
	if err != nil {
		return nil, serviceErrors.ErrInternal
	}
	defer rows.Close()
	targets := make([]model.Target, 0)
	for rows.Next() {
		var target model.Target
		err := rows.Scan(
			&target.Id,
			&target.Url,
			&target.IsTracking,
			&target.IntervalSec,
			&target.CreatedAt,
			&target.UpdatedAt)
		if err != nil {
			return nil, serviceErrors.ErrInternal
		}
		targets = append(targets, target)
	}
	if err := rows.Err(); err != nil {
		return nil, serviceErrors.ErrInternal
	}
	return targets, nil
}

func (pg *PgService) Update(ctx context.Context, id, url string, intervalSec int64) (model.Target, error) {
	query := `UPDATE target SET url = $1, interval_sec = $2,updated_at = NOW() WHERE id = $3 RETURNING id, url, is_tracking, interval_sec, created_at, updated_at`
	var target model.Target
	err := pg.db.QueryRowContext(ctx, query, url, intervalSec, id).Scan(
		&target.Id,
		&target.Url,
		&target.IsTracking,
		&target.IntervalSec,
		&target.CreatedAt,
		&target.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Target{}, serviceErrors.ErrNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.Target{}, serviceErrors.ErrConflict
		}
		return model.Target{}, serviceErrors.ErrInternal
	}
	return target, nil
}
