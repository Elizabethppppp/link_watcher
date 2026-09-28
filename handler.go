package main

import (
	"encoding/json"
	"link_watcher/errorResponse"
	"link_watcher/serviceErrors"
	"net/http"
	"net/url"
	"strconv"
)

type CreateTargetRequest struct {
	Url         string `json:"url"`
	IntervalSec int64  `json:"intervalSec"`
}

type UpdateTargetRequest struct {
	Url         string `json:"url"`
	IntervalSec int64  `json:"intervalSec"`
}

func (t *Transport) CreateTarget(w http.ResponseWriter, r *http.Request) {

	var req CreateTargetRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrInvalidJSON)
		return
	}

	if req.Url == "" {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrEmptyURL)
		return
	}

	if len(req.Url) > 2048 {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrURLTooLarge)
		return
	}

	parsed, errParse := url.Parse(req.Url)
	if errParse != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrInvalidURL)
		return
	}

	if req.IntervalSec == 0 {
		req.IntervalSec = 60
	}
	if req.IntervalSec < 1 || req.IntervalSec > 60 {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrInvalidInterval)
		return
	}

	target, err := t.red.Create(r.Context(), req.Url, req.IntervalSec)
	if err != nil {
		errorResponse.ErrorResponseJSON(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(target)
}

func (t *Transport) GetTarget(w http.ResponseWriter, r *http.Request) {
	target, err := t.red.GetTargets(r.Context())
	if err != nil {
		errorResponse.ErrorResponseJSON(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(target)
}

func (t *Transport) UpdateTarget(w http.ResponseWriter, r *http.Request) {

	idText := r.PathValue("id")
	id, err := strconv.ParseInt(idText, 10, 64)
	if err != nil {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrInvalidId)
		return
	}
	if id <= 0 {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrInvalidId)
		return
	}

	var req UpdateTargetRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrInvalidJSON)
		return
	}

	if req.Url == "" {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrEmptyURL)
		return
	}

	if len(req.Url) > 2048 {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrURLTooLarge)
		return
	}

	parsed, err := url.Parse(req.Url)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrInvalidURL)
		return
	}

	if req.IntervalSec == 0 {
		req.IntervalSec = 60
	}
	if req.IntervalSec < 1 || req.IntervalSec > 60 {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrInvalidInterval)
		return
	}

	target, err := t.red.UpdateTargetId(r.Context(), id, req.Url, req.IntervalSec)
	if err != nil {
		errorResponse.ErrorResponseJSON(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(target)

}

func (t *Transport) DeleteTarget(w http.ResponseWriter, r *http.Request) {
	idText := r.PathValue("id")
	id, err := strconv.ParseInt(idText, 10, 64)
	if err != nil {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrInvalidId)
		return
	}
	if id <= 0 {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrInvalidId)
		return
	}

	err = t.red.DeleteTargetId(r.Context(), id)
	if err != nil {
		errorResponse.ErrorResponseJSON(w, err)
		return

	}
	w.WriteHeader(http.StatusNoContent)
}

func (t *Transport) UpdateTracking(w http.ResponseWriter, r *http.Request) {
	idText := r.PathValue("id")
	id, err := strconv.ParseInt(idText, 10, 64)
	if err != nil {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrInvalidId)
		return
	}
	if id <= 0 {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrInvalidId)
		return
	}

	target, err := t.red.UpdateActiveTarget(r.Context(), id)
	if err != nil {
		errorResponse.ErrorResponseJSON(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(target)

}
