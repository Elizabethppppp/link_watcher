package main

import (
	"encoding/json"
	"errors"
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
	if req.IntervalSec < 1 || req.IntervalSec > 60 {
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
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(target)
}

func (t *Transport) GetTarget(w http.ResponseWriter, r *http.Request) {
	target, err := t.red.GetTargets(r.Context())
	if errors.Is(err, serviceErrors.ErrInternal) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("id is invalid"))
		return
	}

	var req UpdateTargetRequest

	if errDecode := json.NewDecoder(r.Body).Decode(&req); errDecode != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errDecode.Error()))
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
	if req.IntervalSec < 1 || req.IntervalSec > 60 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("interval is not less than 0"))
		return
	}

	target, err := t.red.UpdateTargetId(r.Context(), id, req.Url, req.IntervalSec)
	if err != nil {
		if errors.Is(err, serviceErrors.ErrConflict) {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(err.Error()))
			return
		}
		if errors.Is(err, serviceErrors.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(target)

}

func (t *Transport) DeleteTarget(w http.ResponseWriter, r *http.Request) {
	idText := r.PathValue("id")
	if idText == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("id is empty"))
		return
	}
	id, err := strconv.ParseInt(idText, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("id is invalid"))
		return
	}

	err = t.red.DeleteTargetId(r.Context(), id)
	if err != nil {
		if errors.Is(err, serviceErrors.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return

	}
	w.WriteHeader(http.StatusNoContent)
}

func (t *Transport) UpdateTracking(w http.ResponseWriter, r *http.Request) {
	idText := r.PathValue("id")
	id, err := strconv.ParseInt(idText, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("id is invalid"))
		return
	}

	target, err := t.red.UpdateActiveTarget(r.Context(), id)
	if err != nil {
		if errors.Is(err, serviceErrors.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(target)

}
