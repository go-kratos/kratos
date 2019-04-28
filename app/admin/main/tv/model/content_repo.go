package model

import (
	"go-common/library/time"
)

// ContentRepo def.
type ContentRepo struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Subtitle    string    `json:"subtitle"`
	Desc        string    `json:"desc"`
	Cover       string    `json:"cover"`
	SeasonID    int       `json:"season_id"`
	CID         int       `json:"cid" gorm:"column:cid"`
	EPID        int       `json:"epid" gorm:"column:epid"`
	MenuID      int       `json:"menu_id"`
	State       int       `json:"state"`
	Valid       int       `json:"valid"`
	PayStatus   int       `json:"pay_status"`
	IsDeleted   int       `json:"is_deleted"`
	AuditTime   int       `json:"audit_time"`
	Ctime       time.Time `json:"ctime"`
	Mtime       time.Time `json:"mtime_nb,omitempty"`
	MtimeFormat string    `json:"mtime"`
	InjectTime  time.Time `json:"inject_time"`
	// InjectTimeFormat string    `json:"inject_time"`
	Reason      string `json:"reason"`
	SeasonTitle string `json:"season_title" gorm:"column:season_title"`
	Category    int8   `json:"category" gorm:"column:category"`
}

// TableName tv_content
func (*ContentRepo) TableName() string {
	return "tv_content"
}

// ContentRepoPager def.
type ContentRepoPager struct {
	TotalCount int64          `json:"total_count"`
	Pn         int            `json:"pn"`
	Ps         int            `json:"ps"`
	Items      []*ContentRepo `json:"items"`
}
