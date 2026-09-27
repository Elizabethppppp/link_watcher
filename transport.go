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

	return mux
}
