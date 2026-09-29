package http

import (
	"encoding/json"
	"link_watcher/errorResponse"
	"link_watcher/serviceErrors"
	"link_watcher/validation"
	"net/http"
	"strconv"
)

func (t *Transport) CreateTarget(w http.ResponseWriter, r *http.Request) {

	var req validation.TargetRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrInvalidJSON)
		return
	}

	err := validation.ValidateTargetRequest(&req)
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

	var req validation.TargetRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrInvalidJSON)
		return
	}

	err = validation.ValidateTargetRequest(&req)
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
