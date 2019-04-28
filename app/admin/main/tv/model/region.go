package model

import (
	"go-common/library/time"

	"github.com/siddontang/go-mysql/mysql"
)

// RegCore .
type RegCore struct {
	ID        int    `json:"id" form:"id"`
	PageID    int    `json:"page_id" form:"page_id"`
	Title     string `json:"title" form:"title"`
	Valid     int    `json:"valid" form:"valid"`
	IndexType int    `json:"index_type" form:"index_type"`
	IndexTid  int    `json:"index_tid" form:"index_tid"`
	Deleted   int    `json:"deleted" form:"deleted"`
	Rank      int    `json:"rank"`
}

// RegDB .
type RegDB struct {
	RegCore
	Mtime time.Time `json:"mtime" form:"mtime"`
}

// RegList .
type RegList struct {
	RegCore
	Mtime string `json:"mtime"`
}

// ToList ctime format .
func (v *RegDB) ToList() *RegList {
	return &RegList{
		RegCore: v.RegCore,
		Mtime:   v.Mtime.Time().Format(mysql.TimeFormat),
	}
}

// TableName return table name .
func (*RegDB) TableName() string {
	return "tv_pages"
}

// Param .
type Param struct {
	Title  string `form:"title"`
	PageID string `form:"page_id"`
	State  string `form:"state"`
}
