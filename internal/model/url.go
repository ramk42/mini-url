// Package model contains the URL model
package model

import (
	"time"
)

const Sluglength = 6

type URL struct {
	ID        string
	Original  string
	Slug      string
	CreatedAt time.Time
	ExpiresAt *time.Time
	Clicks    int64
}
