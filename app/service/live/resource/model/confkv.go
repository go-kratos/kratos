package model

import "time"

// Confkv def
type Confkv struct {
	ID    int64     `json:"id" gorm:"column:id"`
	Key   string    `json:"key" form:"key"`
	Value string    `json:"value" form:"value"`
	Ctime time.Time `json:"ctime"`
	Mtime time.Time `json:"mtime"`
}

// TableName resource
func (c Confkv) TableName() string {
	return "confkv"
}
