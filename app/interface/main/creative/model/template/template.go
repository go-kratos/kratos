package template

import "go-common/library/time"

const (
	// StateNormal 正常
	StateNormal = 0
	// StateDel 删除
	StateDel = 1
)

// Template archive template.
type Template struct {
	ID   int64  `json:"tid"`
	Name string `json:"name"`
	// Arctype   string `json:"-"`
	TypeID    int16     `json:"typeid"`
	Title     string    `json:"title"`
	Tag       string    `json:"tags"`
	Content   string    `json:"description"`
	Copyright int8      `json:"copyright"`
	State     int8      `json:"-"`
	CTime     time.Time `json:"-"`
	MTime     time.Time `json:"-"`
}

// Copyright get int8 val
func Copyright(cp string) int8 {
	if cp == "Original" {
		return 1
	} else if cp == "Copy" {
		return 2
	} else {
		return 0
	}
}
