package pgService

import (
	"context"
	"database/sql"
)

type PgService struct {
	db *sql.DB
}

func NewPgService(db *sql.DB) *PgService {
	return &PgService{
		db: db,
	}
}

func (service *PgService) Insert(ctx context.Context, url string, intervalSec int64) error {
}
