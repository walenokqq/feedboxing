package models

import "time"

type Project struct {
	ID          uint
	Title       string
	Description string
	Created_at  time.Time
}
