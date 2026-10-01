package checker

import (
	"context"
	"link_watcher/logger"
	"sync"
)

type RunRepository interface {
	CheckSave(ctx context.Context, targetId int64, url string) error
}

func RunChecker(ctx context.Context, r RunRepository, tasks <-chan Task, n int) {
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(checkers int) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			case task, ok := <-tasks:
				if !ok {
					return
				}
				if err := r.CheckSave(ctx, task.Id, task.Url); err != nil {
					logger.Error(err.Error())
				}
			}
		}(i)
	}
	wg.Wait()
}
