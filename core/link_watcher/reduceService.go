package link_watcher

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
	GetCurrentSummary(ctx context.Context, targetId int64) (model.Summary, error)
}

type ReduceService struct {
	repo TargetRepository
}

func NewReduceService(repo TargetRepository) *ReduceService {
	return &ReduceService{
		repo: repo,
	}
}

func (s *ReduceService) Create(ctx context.Context, url string, intervalSec int64) (model.Target, error) {
	return s.repo.Insert(ctx, url, intervalSec)
}

func (s *ReduceService) GetTargets(ctx context.Context) ([]model.Target, error) {
	return s.repo.GetAllTargets(ctx)
}

func (s *ReduceService) UpdateTargetId(ctx context.Context, id int64, url string, intervalSec int64) (model.Target, error) {
	return s.repo.Update(ctx, id, url, intervalSec)
}

func (s *ReduceService) DeleteTargetId(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *ReduceService) UpdateActiveTarget(ctx context.Context, id int64) (model.Target, error) {
	return s.repo.UpdateActive(ctx, id)
}

func (s *ReduceService) GetOneSummary(ctx context.Context, targetId int64) (model.Summary, error) {
	return s.repo.GetCurrentSummary(ctx, targetId)
}
