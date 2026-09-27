package main

import (
	"context"
	"link_watcher/model"
)

type TargetRepository interface {
	Insert(ctx context.Context, url string, intervalSec int64) (model.Target, error)
	GetAllTargets(ctx context.Context) ([]model.Target, error)
	Update(ctx context.Context, id, url string, intervalSec int64) (model.Target, error)
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

	tar, err := service.repo.Insert(ctx, url, intervalSec)
	if err != nil {
		return model.Target{}, err
	}
	return tar, nil
}

func (service *ReduceService) GetTargets(ctx context.Context) ([]model.Target, error) {
	targets, err := service.repo.GetAllTargets(ctx)
	if err != nil {
		return nil, err
	}
	return targets, nil
}

func (service *ReduceService) UpdateTargetId(ctx context.Context, id, url string, intervalSec int64) (model.Target, error) {
	target, err := service.repo.Update(ctx, id, url, intervalSec)
	if err != nil {
		return model.Target{}, err
	}
	return target, nil
}
