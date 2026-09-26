package main

type Transport struct {
	tr *ReduceService
}

func NewTransport(tr *ReduceService) *Transport {
	return &Transport{
		tr: tr,
	}
}
