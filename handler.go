package main

import (
	"encoding/json"
	"link_watcher/errorResponse"
	"link_watcher/serviceErrors"
	"net/http"
	"net/url"
	"strconv"
)

type TargetRequest struct {
	Url         string `json:"url"`
	IntervalSec int64  `json:"intervalSec"`
}

func validateTargetRequest(req *TargetRequest) error {
	if req.Url == "" {
		return serviceErrors.ErrEmptyURL
	}

	if len(req.Url) > 2048 {
		return serviceErrors.ErrURLTooLarge
	}

	parsed, errParse := url.Parse(req.Url)
	if errParse != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return serviceErrors.ErrInvalidURL
	}

	if req.IntervalSec == 0 {
		req.IntervalSec = 60
	}
	if req.IntervalSec < 1 || req.IntervalSec > 60 {
		return serviceErrors.ErrInvalidInterval
	}
	return nil
}

func (t *Transport) CreateTarget(w http.ResponseWriter, r *http.Request) {

	var req TargetRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrInvalidJSON)
		return
	}

	err := validateTargetRequest(&req)
	if err != nil {
		errorResponse.ErrorResponseJSON(w, err)
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
	if err != nil || id <= 0 {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrInvalidId)
		return
	}

	var req TargetRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrInvalidJSON)
		return
	}

	err = validateTargetRequest(&req)
	if err != nil {
		errorResponse.ErrorResponseJSON(w, err)
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
	if err != nil || id <= 0 {
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
	if err != nil || id <= 0 {
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
