package model

import (
	"time"
)

type History struct {
	ID        string
	Status    Status
	Text      string
	Timestamp time.Time
}
