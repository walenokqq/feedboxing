package models

import "time"

type Feedback struct {
	ID         uint
	FormID     uint
	Data       JSONB
	Status     string
	Created_at time.Time
}
