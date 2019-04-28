package income

import (
	"go-common/library/time"
)

// ColumnCharge column charge
type ColumnCharge struct {
	ID           int64
	ArticleID    int64
	Title        string
	MID          int64
	TagID        int64
	IncCharge    int64
	IncViewCount int64
	Date         time.Time
	UploadTime   time.Time
}
