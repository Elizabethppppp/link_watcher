package checker

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"
)

type CheckerRepository interface {
	InsertChecks(ctx context.Context, targetId int64, statusCode *int, latencyMs *int) error
}

type Checker struct {
	client *http.Client
	repo   CheckerRepository
}

func NewChecker(repo CheckerRepository, timeout time.Duration) *Checker {
	return &Checker{
		client: &http.Client{
			Timeout: timeout,
		},
		repo: repo,
	}
}

func (c *Checker) Check(ctx context.Context, urlStr string) (*int, *int, error) {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, http.NoBody)
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		var urlErr *url.Error
		switch {
		case errors.As(err, &urlErr) && urlErr.Timeout():
			log.Printf("timeout: %v", err)
		case errors.Is(err, context.Canceled):
			log.Printf("canceled: %v", err)
		default:
			log.Printf("error: %v", err)
		}
		return nil, nil, nil
	}
	defer resp.Body.Close()

	latency := int(time.Since(start).Milliseconds())
	return &resp.StatusCode, &latency, nil
}

func (c *Checker) CheckSave(ctx context.Context, targetId int64, url string) error {
	statusCode, latency, _ := c.Check(ctx, url)
	return c.repo.InsertChecks(ctx, targetId, statusCode, latency)
}
