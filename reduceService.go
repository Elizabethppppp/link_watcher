package main

import (
	"context"
	"link_watcher/model"
)

type TargetRepository interface {
	Insert(ctx context.Context, url string, intervalSec int64) (model.Target, error)
	GetAllTargets(ctx context.Context) ([]model.Target, error)
	Update(ctx context.Context, id int64, url string, intervalSec int64) (model.Target, error)
	Delete(ctx context.Context, id int64) error
	UpdateActive(ctx context.Context, id int64) (model.Target, error)
}

type ReduceService struct {
	repo TargetRepository
}

func NewReduceService(repo TargetRepository) *ReduceService {
	return &ReduceService{
		repo: repo,
	}
}

func (service *ReduceService) Create(ctx context.Context, url string, intervalSec int64) (model.Target, error) {
	return service.repo.Insert(ctx, url, intervalSec)
}

func (service *ReduceService) GetTargets(ctx context.Context) ([]model.Target, error) {
	return service.repo.GetAllTargets(ctx)
}

func (service *ReduceService) UpdateTargetId(ctx context.Context, id int64, url string, intervalSec int64) (model.Target, error) {
	return service.repo.Update(ctx, id, url, intervalSec)
}

func (service *ReduceService) DeleteTargetId(ctx context.Context, id int64) error {
	return service.repo.Delete(ctx, id)
}

func (service *ReduceService) UpdateActiveTarget(ctx context.Context, id int64) (model.Target, error) {
	return service.repo.UpdateActive(ctx, id)
}
