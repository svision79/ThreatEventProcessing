package models

import "time"

type Event struct {
	ID          int       `json:"id"`
	Source      string    `json:"source"`
	ThreatType  string    `json:"threat_type"`
	DetectedAt  time.Time `json:"detected_at"`
	ProcessedAt time.Time `json:"processed_at"`
	Details     string    `json:"details"`
}
