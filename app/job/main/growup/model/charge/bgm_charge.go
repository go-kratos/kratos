package charge

import (
	"go-common/library/time"
)

// BgmCharge bgm charge
type BgmCharge struct {
	ID        int64
	SID       int64
	AID       int64
	MID       int64
	CID       int64
	IncCharge int64
	Date      time.Time
	JoinAt    time.Time
	Title     string
	DBState   int
}

// BgmStatis bgm statis
type BgmStatis struct {
	ID          int64
	SID         int64
	AID         int64
	MID         int64
	CID         int64
	Title       string
	JoinAt      time.Time
	TotalCharge int64
	DBState     int
}

// Bgm background music
type Bgm struct {
	ID     int64
	MID    int64
	SID    int64
	AID    int64
	CID    int64
	JoinAt time.Time
	Title  string
}
