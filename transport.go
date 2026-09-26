package main

import (
	"net/http"
)

type Transport struct {
	tr *ReduceService
}

func NewTransport(tr *ReduceService) *Transport {
	return &Transport{
		tr: tr,
	}
}

func (t *Transport) Handler() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/targets", t.CreateTarget)
}
