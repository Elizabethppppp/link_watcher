package main

import (
	"encoding/json"
	"errors"
	"link_watcher/serviceErrors"
	"net/http"
	"net/url"
)

type CreateTargetRequest struct {
	Url         string `json:"url"`
	IntervalSec int64  `json:"interval"`
}

func (t *Transport) CreateTarget(w http.ResponseWriter, r *http.Request) {

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

	parsed, errParse := url.Parse(req.Url)
	if errParse != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("url is invalid"))
		return
	}

	if req.IntervalSec == 0 {
		req.IntervalSec = 60
	}
	if req.IntervalSec < 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("interval is not less than 0"))
		return
	}

	target, err := t.red.Create(r.Context(), req.Url, req.IntervalSec)
	if err != nil {
		if errors.Is(err, serviceErrors.ErrConflict) {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(err.Error()))
			return
		}
		if errors.Is(err, serviceErrors.ErrInternal) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(target)
}
