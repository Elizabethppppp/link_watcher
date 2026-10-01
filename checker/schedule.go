package checker

import (
	"context"
	"link_watcher/model"
	"log"
	"time"

	"errors"
)

type ScheduleRepository interface {
	GetTargetsIsTrackingNow(ctx context.Context) ([]model.Target, error)
}

type Schedule struct {
	repo  ScheduleRepository
	tasks chan<- Task
	tick  time.Duration
}

type Task struct {
	Id  int64
	Url string
}

func NewSchedule(repo ScheduleRepository, tasks chan<- Task, tick time.Duration) *Schedule {
	return &Schedule{repo: repo, tasks: tasks, tick: tick}
}

func (s *Schedule) Start(ctx context.Context) {
	ticker := time.NewTicker(s.tick)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.queue(ctx)
		}
	}

}

func (s *Schedule) queue(ctx context.Context) {
	targets, err := s.repo.GetTargetsIsTrackingNow(ctx)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		log.Printf("schedule: get due targets: %v", err)
		return
	}
	for _, target := range targets {
		select {
		case s.tasks <- Task{Id: target.Id, Url: target.Url}:
		case <-ctx.Done():
			return
		default:
			log.Printf("schedule: queue full, skip target %d", target.Id)
		}
	}
}
