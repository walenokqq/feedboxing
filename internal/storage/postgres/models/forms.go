package models

import "time"

type Form struct {
	ID          uint
	ProjectID   uint
	Description string
	// Schema JSONB
	Created_at time.Time
}
