package model

import (
	"time"
)

type Item struct {
	UserID       string        `json:"user_id,omitempty"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	URL          string        `json:"url"`
	Selector     string        `json:"selector"`
	Notification *Notification `json:"notification"`
	Timestamp    time.Time     `json:"time,omitempty"`
}
