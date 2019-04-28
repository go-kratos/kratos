package model

import (
	"go-common/library/time"
)

// Setting def.
type Setting struct {
	ID       int64     `json:"id"`
	Property string    `json:"property"`
	Name     string    `json:"name"`
	Val      string    `json:"val"`
	CTime    time.Time `json:"ctime"`
	MTime    time.Time `json:"mtime"`
}
