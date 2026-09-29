package worker

import (
	"context"
	"net/http"
	"time"
)

type Checker struct {
	client *http.Client
}

func NewChecker() *Checker {
	return &Checker{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (checker *Checker) Check(ctx context.Context, url string) (*int, *int, error) {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := checker.client.Do(req)
	latency := int(time.Since(start).Milliseconds())
	if err != nil {
		return nil, nil, nil
	}

	defer resp.Body.Close()

	statusCode := resp.StatusCode
	return &statusCode, &latency, nil
}
