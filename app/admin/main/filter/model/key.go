package model

import (
	xtime "go-common/library/time"
)

// KeyInfo rule in key
type KeyInfo struct {
	ID      int64      `json:"fid"`
	Key     string     `json:"key"`
	Mode    int8       `json:"mode"`
	Filter  string     `json:"filter"`
	Level   int8       `json:"level"`
	TpIDs   []int64    `json:"tpid"`
	Areas   []string   `json:"areas"`
	State   int8       `json:"state"`
	Comment string     `json:"comment"`
	Stime   xtime.Time `json:"stime"`
	Etime   xtime.Time `json:"etime"`
	Shelve  bool       `json:"shelve"` // 是否上架
	CTime   xtime.Time `json:"ctime"`  // 创建时间
}
