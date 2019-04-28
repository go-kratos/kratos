package model

import "go-common/library/time"

// Resource reprensents the resource table
type Resource struct {
	ID      int64     `json:"id" params:"id"`
	Name    string    `json:"name" params:"name"`
	Version int64     `json:"version" params:"version"`
	PoolID  int64     `json:"pool_id" params:"pool_id"`
	Ctime   time.Time `json:"ctime" params:"ctime"`
	Mtime   time.Time `json:"mtime" params:"mtime"`
}

// TableName gives the table name of the model
func (*Resource) TableName() string {
	return "resource"
}
