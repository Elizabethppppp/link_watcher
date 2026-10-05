package model

import "time"

type Summary struct {
	TargetID      int64      `json:"targetId"`
	TotalChecks   int64      `json:"totalChecks"`
	SuccessChecks int64      `json:"successChecks"`
	SuccessRate   float64    `json:"successRate"`
	AvgLatencyMs  *float64   `json:"avgLatencyMs"`
	LastCheckedAt *time.Time `json:"lastCheckedAt"`
}
