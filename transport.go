package main

import (
	"link_watcher/middleware"
	"net/http"
)

type Transport struct {
	tr           *ReduceService
	CreateTargets http.Handler
}

func NewTransport(tr *ReduceService) *Transport {
	return &Transport{
		tr: tr,
	}
}

func (t *Transport) Handler() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/targets", middleware.LoggerMiddleware(t.CreateTarget))

	return mux
}
