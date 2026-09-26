package main

import (
	"context"
	"encoding/json"
	"link_watcher/logger"
	"net/http"
	"strings"
)

type CreateTargetRequest struct {
	Url      string `json:"url"`
	Interval int    `json:"interval"`
}

type Target struct {
	Id         string `json:"id"`
	Url        string `json:"url"`
	IsTracking bool   `json:"is_tracking"`
	IntervaSec int    `json:"interva_sec"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

func (l *Link) CreateTarget(w http.ResponseWriter, r *http.Request) {

	var req CreateTargetRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if req.Url == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("url is empty"))
		return
	}

	if len(req.Url) > 2048 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("url is too large"))
		return
	}

	if req.Interval == 0 {
		req.Interval = 60
	}
	if req.Interval < 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("interval is not less than 0"))
		return
	}

	query := `INSERT INTO target (url, interval_sec) VALUES ($1, $2) RETURNING id, url, is_tracking ,interval_sec, created_at, updated_at`

	var tar Target
	ctx := context.Background()

	err := l.DB.QueryRowContext(ctx, query, req.Url, req.Interval).Scan(
		&tar.Id,
		&tar.Url,
		&tar.IsTracking,
		&tar.IntervaSec,
		&tar.CreatedAt,
		&tar.UpdatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "23505") || strings.Contains(err.Error(), "unique") {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("Target with this URL already exists"))
			return
		}
		logger.Error("DB insert error", "error", err.Error(), "url", req.Url, "interval", req.Interval)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tar)
}
