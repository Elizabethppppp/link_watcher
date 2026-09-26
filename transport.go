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

	mux.Handle("/targets", middleware.LoggerMiddleware(http.HandlerFunc(t.CreateTarget)))

	return mux
}
