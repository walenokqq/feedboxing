package models

import (
	"encoding/json"
	"time"
)

type Form struct {
	ID          uint
	ProjectID   uint
	Description string
	Schema      json.RawMessage
	Created_at  time.Time
}
