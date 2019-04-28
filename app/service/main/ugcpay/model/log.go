package model

import (
	"time"
)

// LogAccount .
type LogAccount struct {
	ID        int64
	AccountID int64
	From      int64
	To        int64
	Ver       int64
	State     string
	CTime     time.Time
	MTime     time.Time
}
