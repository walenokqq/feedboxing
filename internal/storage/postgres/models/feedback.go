package models

import (
	"encoding/json"
	"time"
)

type Feedback struct {
	ID         uint
	FormID     uint
	Data       json.RawMessage
	Status     string
	Created_at time.Time
}
