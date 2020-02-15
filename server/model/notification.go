package model

import "time"

type Notification struct {
	PageURL    string    `json:"page_url"`
	Selector   string    `json:"selector"`
	Text       string    `json:"text"`
	WebHookURL string    `json:"webhook_url"`
	Timestamp  time.Time `json:"time,omitempty"`
}
