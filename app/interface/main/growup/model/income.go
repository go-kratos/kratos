package model

import (
	"go-common/library/time"
)

// ArchiveIncome av income
type ArchiveIncome struct {
	ID          int64     `json:"id"`
	ArchiveID   int64     `json:"archive_id"`
	Avs         []int64   `json:"avs,omitempty"`
	MID         int64     `json:"mid"`
	Income      int64     `json:"income"`
	MonthIncome int64     `json:"month_income"`
	TotalIncome int64     `json:"total_income"`
	Breach      *AvBreach `json:"breach"`
	Title       string    `json:"title"`
	Icon        string    `json:"icon"`
	Date        time.Time `json:"date"`
	CTime       time.Time `json:"ctime"`
	MTime       time.Time `json:"mtime"`
	IsDeleted   int       `json:"-"`
}

// Blacklist black list
type Blacklist struct {
	ID        int64     `json:"id" gorm:"column:id"`
	AvID      int64     `json:"av_id" gorm:"column:av_id"`
	MID       int64     `json:"mid" gorm:"column:mid"`
	Reason    int       `json:"reason" gorm:"column:reason"`
	CType     int       `json:"ctype" gorm:"column:ctype"`
	HasSigned int       `json:"has_signed" gorm:"column:has_signed"`
	Nickname  string    `json:"nickname" gorm:"column:nickname"`
	CTime     time.Time `json:"ctime" gorm:"column:ctime"`
	MTime     time.Time `json:"mtime" gorm:"column:mtime"`
	IsDeleted int       `json:"-"`
}

// ArchiveRes archive response
type ArchiveRes struct {
	Code    int                 `json:"code"`
	Data    map[string]*Archive `json:"data"`
	Message string              `json:"message"`
}

// Archive archive
type Archive struct {
	AID   int64  `json:"aid"`
	Title string `json:"title"`
}

// TagInfo tag_info
type TagInfo struct {
	ID    int64
	Radio int64
	Icon  string
}

// UpIncome up_income
type UpIncome struct {
	ID               int64
	MID              int64
	Income           int64
	AvIncome         int64
	ColumnIncome     int64
	BgmIncome        int64
	BaseIncome       int64
	AvBaseIncome     int64
	ColumnBaseIncome int64
	BgmBaseIncome    int64
	TotalIncome      int64
	Date             time.Time
}

// UpIncomeStat for up daily income analytics
type UpIncomeStat struct {
	MID         int64     `json:"-"`
	Income      int64     `json:"income"`
	BaseIncome  int64     `json:"base_income"`
	ExtraIncome int64     `json:"extra_income"`
	Breach      int64     `json:"breach"`
	Date        time.Time `json:"date"`
}
