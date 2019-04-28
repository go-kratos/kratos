package model

import (
	xtime "go-common/library/time"
)

// WhiteInfo rule in white
type WhiteInfo struct {
	ID      int64      `json:"fid"`
	Content string     `json:"filter"`
	Mode    int8       `json:"mode"`
	TpIDs   []int64    `json:"tpid"`
	Areas   []string   `json:"areas"`
	Comment string     `json:"comment"`
	State   int8       `json:"state"`
	CTime   xtime.Time `json:"ctime"` // 创建时间
}
