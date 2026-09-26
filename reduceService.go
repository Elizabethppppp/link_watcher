package main

import (
	"context"
)

type TargetRepository interface {
	Insert(ctx context.Context, url string, intervalSec int64) error
}

type ReduceService struct {
	repo TargetRepository
}

func NewReduceService(repo TargetRepository) *ReduceService {
	return &ReduceService{
		repo: repo,
	}
}

func (service *ReduceService) Create(ctx context.Context, url string, intervalSec int64) error {
	if err := service.repo.Insert(ctx, url, intervalSec); err != nil {
		return service.repo.Insert(ctx, url, intervalSec)
	}
	return nil
}
