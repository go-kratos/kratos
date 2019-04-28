package model

import (
	"go-common/library/time"
)

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

// AvBreach av_breach_record
type AvBreach struct {
	ID     int64
	MID    int64
	AvID   int64
	Money  int64
	Reason string
	Date   time.Time
}

// BreachRecord breach record
type BreachRecord struct {
	ID        int64     `json:"id" gorm:"column:id"`
	AvID      int64     `json:"av_id" gorm:"column:av_id"`
	CType     int       `json:"ctype" gorm:"column:ctype"`
	CTime     time.Time `json:"ctime" gorm:"column:ctime"`
	MTime     time.Time `json:"mtime" gorm:"column:mtime"`
	IsDeleted int       `json:"-"`
}

// PorderRes porder response
type PorderRes struct {
	Code    int       `json:"code"`
	Data    []*Porder `json:"data"`
	Message string    `json:"message"`
	TTL     int       `json:"ttl"`
}

// Porder porder
type Porder struct {
	AID        int64     `json:"aid"`
	IndustryID int64     `json:"industry_id"`
	BrandID    int64     `json:"brand_id"`
	BrandName  string    `json:"brand_name"`
	Official   int64     `json:"official"`
	ShowType   string    `json:"show_type"`
	Advertiser string    `json:"advertiser"`
	Agent      string    `json:"agent"`
	State      int64     `json:"state"`
	ShowFront  int64     `json:"show_front"`
	CTime      time.Time `json:"ctime"`
	MTime      time.Time `json:"mtime"`
}

// ExecuteOrder execute order
type ExecuteOrder struct {
	AvID  int64     `json:"av_id"`
	MID   int64     `json:"mid"`
	CTime time.Time `json:"ctime"`
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
	Owner *Owner `json:"owner"`
}

// Owner archive owner
type Owner struct {
	MID  int64  `json:"mid"`
	Name string `json:"name"`
}
