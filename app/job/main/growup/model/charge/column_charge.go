package charge

import (
	"go-common/library/time"
)

// Column column charge
type Column struct {
	ID           int64
	AID          int64
	Title        string
	MID          int64
	TagID        int64
	Words        int64
	IncCharge    int64
	IncViewCount int64
	Date         time.Time
	UploadTime   int64
	DBState      int
}

// ColumnStatis column charge
type ColumnStatis struct {
	ID          int64
	AID         int64
	Title       string
	MID         int64
	TagID       int64
	UploadTime  int64
	TotalCharge int64
	DBState     int
}
