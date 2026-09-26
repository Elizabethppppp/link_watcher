package model

import "time"

type Target struct {
	Id          string    `json:"id"`
	Url         string    `json:"url"`
	IsTracking  bool      `json:"is_tracking"`
	IntervalSec int64     `json:"interva_sec"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
