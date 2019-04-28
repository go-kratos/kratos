package income

import (
	"go-common/library/time"
)

// ArchiveStatis archive statis
type ArchiveStatis struct {
	ID           int64
	Avs          int64
	MoneySection int
	MoneyTips    string
	Income       int64
	CategroyID   int64
	CDate        time.Time
}

// ArchiveIncome archive income
type ArchiveIncome struct {
	ID          int64     `json:"id"`
	AvID        int64     `json:"archive_id"`
	MID         int64     `json:"mid"`
	Type        int       `json:"type"`
	TagID       int64     `json:"category_id"`
	IsOriginal  int       `json:"is_original"`
	UploadTime  time.Time `json:"upload_time"`
	TotalIncome int64     `json:"total_income"`
	Income      int64     `json:"income"`
	TaxMoney    int64     `json:"tax_money"`
	Date        time.Time `json:"date"`
	DateFormat  string    `json:"date_format"`
	Nickname    string    `json:"nickname"`
	Avs         int       `json:"avs"`
}
