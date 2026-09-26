package main

import (
	"context"
	"link_watcher/pgService"
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

func (service *ReduceService) Create(ctx context.Context, url string, intervalSec int64) (pgService.Target, error) {
	var tar pgService.Target
	tar, err := service.repo.Insert(ctx, url, intervalSec)
	if err != nil {
		return service.repo.Insert(ctx, url, intervalSec)
	}
	return nil
}
