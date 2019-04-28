package model

import (
	xtime "go-common/library/time"
)

// CreditAnswerMSG param struct
type CreditAnswerMSG struct {
	MID   int64      `json:"mid"`
	MTime xtime.Time `json:"mtime"`
}
