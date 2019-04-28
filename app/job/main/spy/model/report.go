package model

import "time"

// Report etc.
type Report struct {
	ID          int64
	Name        string
	DateVersion string
	Val         int64
	Ctime       time.Time
}

const (
	//BlockCount block count
	BlockCount = "封禁总数"
	// SecurityLoginCount security login count
	SecurityLoginCount = "二次验证总数"
)
