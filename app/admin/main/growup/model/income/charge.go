package income

import "go-common/library/time"

// ArchiveCharge av charge
type ArchiveCharge struct {
	AID         int64     `json:"aid"`
	AvID        int64     `json:"-"`
	MID         int64     `json:"mid"`
	Nickname    string    `json:"nickname"`
	CategoryID  int64     `json:"category_id"`
	Charge      int64     `json:"charge"`
	TotalCharge int64     `json:"total_charge"`
	UploadTime  time.Time `json:"upload_time"`
	Date        time.Time `json:"date"`
}

// ArchiveChargeStatis av charge statis
type ArchiveChargeStatis struct {
	ID           int64
	Avs          int64
	MoneySection int
	MoneyTips    string
	Charge       int64
	CategroyID   int64
	CDate        time.Time
}
