package model

import (
	"go-common/library/time"
)

// Content content def.
type Content struct {
	ID         int64     `json:"id"`
	Title      string    `json:"title"`
	Subtitle   string    `json:"subtitle"`
	Desc       string    `json:"desc"`
	Cover      string    `json:"cover"`
	SeasonID   int       `json:"season_id"`
	CID        int       `json:"cid" gorm:"column:cid"`
	EPID       int       `json:"epid" gorm:"column:epid"`
	MenuID     int       `json:"menu_id"`
	State      int       `json:"state"`
	Valid      int       `json:"valid"`
	PayStatus  int       `json:"pay_status"`
	IsDeleted  int       `json:"is_deleted"`
	AuditTime  int       `json:"audit_time"`
	Ctime      time.Time `json:"ctime"`
	Mtime      time.Time `json:"mtime"`
	InjectTime time.Time `json:"inject_time"`
	Reason     string    `json:"reason"`
}

// ContentDetail def.
type ContentDetail struct {
	ID         int64     `json:"id"`
	Title      string    `json:"title"`
	Subtitle   string    `json:"subtitle"`
	Desc       string    `json:"desc"`
	Cover      string    `json:"cover"`
	SeasonID   int       `json:"season_id"`
	CID        int       `json:"cid" gorm:"column:cid"`
	EPID       int       `json:"epid" gorm:"column:epid"`
	MenuID     int       `json:"menu_id"`
	State      int       `json:"state"`
	Valid      int       `json:"valid"`
	PayStatus  int       `json:"pay_status"`
	IsDeleted  int       `json:"is_deleted"`
	AuditTime  int       `json:"audit_time"`
	Ctime      time.Time `json:"ctime"`
	Mtime      time.Time `json:"mtime"`
	InjectTime time.Time `json:"inject_time"`
	Reason     string    `json:"reason"`
	Order      int       `json:"order"`
}

// TableName tv_content
func (c Content) TableName() string {
	return "tv_content"
}

// TableName tv_content
func (*ContentDetail) TableName() string {
	return "tv_content"
}
