package income

import (
	"go-common/library/time"
)

// AvBreach av breach record
type AvBreach struct {
	AvID       int64     `json:"archive_id"`
	MID        int64     `json:"mid"`
	CDate      time.Time `json:"cdate"`
	Money      int64     `json:"money"`
	CType      int       `json:"ctype"`
	Reason     string    `json:"reason"`
	UploadTime time.Time `json:"upload_time"`
	Nickname   string    `json:"nickname"`
}
