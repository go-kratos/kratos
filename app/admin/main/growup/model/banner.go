package model

import (
	"go-common/library/time"
)

// Banner banner
type Banner struct {
	ID      int64     `json:"id"`
	Image   string    `json:"image"`
	Link    string    `json:"link"`
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
	CTime   time.Time `json:"-"`
	MTime   time.Time `json:"-"`
}
