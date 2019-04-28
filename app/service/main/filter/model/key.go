package model

import (
	xtime "go-common/library/time"
)

// KeyAreaInfo rule in key
type KeyAreaInfo struct {
	ID      int64      `json:"fid"`
	FKID    int64      `json:"fkid"` // area 表id
	Key     string     `json:"key"`
	Mode    int8       `json:"mode"`
	Filter  string     `json:"filter"`
	Level   int8       `json:"level"`
	TpIDs   []int64    `json:"tpid"`
	Area    string     `json:"area"`
	State   int8       `json:"state"`
	Comment string     `json:"comment"`
	STime   xtime.Time `json:"stime"`
	ETime   xtime.Time `json:"etime"`
	CTime   xtime.Time `json:"ctime"` // 创建时间
}

// KeyTestResult struct .
type KeyTestResult struct {
	*KeyAreaInfo
	Shelve bool `json:"shelve"` // 是否上架
}
