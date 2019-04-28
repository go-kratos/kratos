package model

import (
	"time"
)

// Bill .
type Bill struct {
	ID       int64
	MID      int64
	Biz      string
	Currency string
	In       int64
	Out      int64
	Ver      int64
	Version  int64
	CTime    time.Time
	MTime    time.Time
}
