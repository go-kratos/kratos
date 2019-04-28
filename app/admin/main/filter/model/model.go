package model

import (
	xtime "go-common/library/time"
)

// admin operation state
const (
	LogStateAdd  = 0
	LogStateEdit = 1
	LogStateDel  = 2
)

// Message origin message.
type Message struct {
	ID      int64  `json:"id"`
	Content string `json:"content"`
}

// Log present admin operation log
type Log struct {
	Key     string     `json:"key"`
	AdminID int64      `json:"adid"`
	Name    string     `json:"name"`
	Comment string     `json:"comment"`
	Ctime   xtime.Time `json:"ctime"`
	State   int8       `json:"state"`
}

// Page struct.
type Page struct {
	Num   int   `json:"num"`
	Size  int   `json:"size"`
	Total int64 `json:"total"`
}
