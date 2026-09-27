package main

import (
	"link_watcher/middleware"
	"net/http"
)

type Transport struct {
	red           *ReduceService
	CreateTargets http.Handler
}

func NewTransport(red *ReduceService) *Transport {
	return &Transport{
		red: red,
	}
}

func (t *Transport) Handler() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("POST /targets", middleware.LoggerMiddleware(http.HandlerFunc(t.CreateTarget)))
	mux.Handle("GET /targets", middleware.LoggerMiddleware(http.HandlerFunc(t.GetTarget)))
	mux.Handle("PUT /targets/{id}", middleware.LoggerMiddleware(http.HandlerFunc(t.UpdateTarget)))
	mux.Handle("DELETE /targets/{id}", middleware.LoggerMiddleware(http.HandlerFunc(t.DeleteTarget)))
	mux.Handle("PATCH /targets/{id}/tracking", middleware.LoggerMiddleware(http.HandlerFunc(t.UpdateTracking)))

	return mux
}
