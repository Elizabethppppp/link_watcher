package model

import "time"

type Target struct {
	Id          int64     `json:"id"`
	Url         string    `json:"url"`
	IsTracking  bool      `json:"isTracking"`
	IntervalSec int64     `json:"intervalSec"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
