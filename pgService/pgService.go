package pgService

import (
	"context"
	"database/sql"
	"errors"
	"link_watcher/model"
	"link_watcher/serviceErrors"

	"github.com/jackc/pgerrcode"
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
	const query = `INSERT INTO target (url, interval_sec) VALUES ($1, $2) RETURNING id, url, is_tracking ,interval_sec, created_at, updated_at`

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
		if errors.As(err, &pgerr) && pgerr.Code == pgerrcode.UniqueViolation {
			return model.Target{}, serviceErrors.ErrConflict
		}
		return model.Target{}, serviceErrors.ErrInternal
	}
	return target, nil
}

func (pg *PgService) GetAllTargets(ctx context.Context) ([]model.Target, error) {
	const query = `SELECT id, url, is_tracking ,interval_sec, created_at, updated_at FROM target`
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

func (pg *PgService) Update(ctx context.Context, id int64, url string, intervalSec int64) (model.Target, error) {
	const query = `UPDATE target SET url = $1, interval_sec = $2,updated_at = NOW() WHERE id = $3 RETURNING id, url, is_tracking, interval_sec, created_at, updated_at`
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
		if errors.Is(err, sql.ErrNoRows) {
			return model.Target{}, serviceErrors.ErrNotFound
		}
		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) && pgerr.Code == pgerrcode.UniqueViolation {
			return model.Target{}, serviceErrors.ErrConflict
		}
		return model.Target{}, serviceErrors.ErrInternal
	}
	return target, nil
}

func (pg *PgService) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM target WHERE id = $1`
	result, err := pg.db.ExecContext(ctx, query, id)
	if err != nil {
		return serviceErrors.ErrInternal
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return serviceErrors.ErrInternal
	}
	if affected == 0 {
		return serviceErrors.ErrNotFound
	}
	return nil
}

func (pg *PgService) UpdateActive(ctx context.Context, id int64) (model.Target, error) {
	const query = `UPDATE target SET is_tracking = NOT is_tracking,updated_at=NOW() WHERE id = $1 RETURNING id, url, is_tracking, interval_sec, created_at, updated_at`
	var target model.Target
	err := pg.db.QueryRowContext(ctx, query, id).Scan(
		&target.Id,
		&target.Url,
		&target.IsTracking,
		&target.IntervalSec,
		&target.CreatedAt,
		&target.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Target{}, serviceErrors.ErrNotFound
		}
		return model.Target{}, serviceErrors.ErrInternal
	}
	return target, nil
}

func (pg *PgService) InsertChecks(ctx context.Context, targetId int64, statusCode *int, latencyMs *int) error {
	const query = `INSERT INTO checks (target_id, status_code, latency_ms) values ($1, $2, $3)`

	_, err := pg.db.ExecContext(ctx, query, targetId, statusCode, latencyMs)
	if err != nil {
		return serviceErrors.ErrInternal
	}
	return nil
}

func (pg *PgService) GetTargetsIsTrackingNow(ctx context.Context) ([]model.Target, error) {
	const query = `SELECT id, url ,interval_sec
                FROM target 
				LEFT JOIN LATERAL ( 
					    SELECT MAX(checked_at) AS last_checked_at
					    FROM checks
					    WHERE target_id = target.id
					 ) AS last_check ON TRUE
					where target.is_tracking=true
					AND (
					    last_check.last_checked_at IS NULL 
					    OR last_check.last_checked_at + make_interval(secs => target.interval_sec) <= NOW()
					)`
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
			&target.IntervalSec)
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

func (pg *PgService) GetCurrentSummary(ctx context.Context, targetId int64) (model.Summary, error) {
	const query = `SELECT COUNT(*) AS total_checks,
COUNT (*) FILTER (WHERE status_code BETWEEN 200 AND 299) AS success_checks,
ROUND(100.0 * COUNT(*) FILTER (WHERE status_code BETWEEN 200 AND 299) / NULLIF(COUNT(*), 0), 2) AS success_rate,
    AVG(latency_ms) FILTER (WHERE status_code BETWEEN 200 AND 299) AS avg_latency_ms,
    MAX(checked_at) AS last_checked_at
FROM checks
WHERE target_id = $1 AND checked_at >= CURRENT_DATE`

	var checks model.Summary
	checks.TargetId = targetId

	err := pg.db.QueryRowContext(ctx, query, targetId).Scan(
		&checks.TotalChecks,
		&checks.SuccessChecks,
		&checks.SuccessRate,
		&checks.AvgLatencyMs,
		&checks.LastCheckedAt,
	)
	if err != nil {
		return model.Summary{}, serviceErrors.ErrInternal
	}
	return checks, nil

}
