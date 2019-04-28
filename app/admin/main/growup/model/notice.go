package model

import (
	"go-common/library/time"
)

// Notice notice
type Notice struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Type      int       `json:"type"`
	Platform  int       `json:"platform"`
	Link      string    `json:"link"`
	Status    int       `json:"status"`
	IsDeleted int       `json:"-"`
	CTime     time.Time `json:"ctime"`
	MTime     time.Time `json:"mtime"`
}
