package checker

import (
	"context"
	"link_watcher/logger"
	"sync"
)

type RunRepository interface {
	CheckSave(ctx context.Context, targetId int64, url string) error
}

type Run struct {
	saver   RunRepository
	workers int
	tasks   <-chan Task
}

func NewRun(saver RunRepository, tasks <-chan Task, workers int) *Run {
	return &Run{
		saver:   saver,
		workers: workers,
		tasks:   tasks,
	}
}

func (r *Run) RunChecker(ctx context.Context) {
	var wg sync.WaitGroup

	wg.Add(r.workers)
	for range r.workers {

		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case task, ok := <-r.tasks:
					if !ok {
						return
					}
					if ctx.Err() != nil {
						return
					}
					if err := r.saver.CheckSave(ctx, task.Id, task.Url); err != nil {
						logger.Error("check save failed", "error", err.Error(), "targetId", task.Id, "url", task.Url)
					}
				}
			}
		}()
	}
	wg.Wait()
}
