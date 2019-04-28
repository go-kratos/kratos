package model

import (
	"time"
)

// Notice notice
type Notice struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Type      int       `json:"type"`
	Platform  int       `json:"-"`
	Link      string    `json:"link"`
	Status    int       `json:"-"`
	IsDeleted int       `json:"-"`
	CTime     time.Time `json:"ctime"`
	MTime     time.Time `json:"-"`
}
