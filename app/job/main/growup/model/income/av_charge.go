package income

import (
	"go-common/library/time"
)

// AvCharge av daily charge
type AvCharge struct {
	ID             int64
	AvID           int64
	MID            int64
	TagID          int64
	IsOriginal     int
	DanmakuCount   int64
	CommentCount   int64
	CollectCount   int64
	CoinCount      int64
	ShareCount     int64
	ElecPayCount   int64
	TotalPlayCount int64
	WebPlayCount   int64
	AppPlayCount   int64
	H5PlayCount    int64
	LvUnknown      int64
	Lv0            int64
	Lv1            int64
	Lv2            int64
	Lv3            int64
	Lv4            int64
	Lv5            int64
	Lv6            int64
	VScore         int64
	IncCharge      int64
	TotalCharge    int64
	IsDeleted      int
	Date           time.Time
	UploadTime     time.Time
	CTime          time.Time
	MTime          time.Time
	DBState        int // 1-insert 2-update
}
