package net

import (
	"time"
)

const (
	// TableFlow .
	TableFlow = "net_flow"
)

// Flow 节点
type Flow struct {
	ID          int64     `gorm:"primary_key" json:"id"`
	NetID       int64     `gorm:"column:net_id" json:"net_id"`
	Name        string    `gorm:"column:name" json:"name"`
	ChName      string    `gorm:"column:ch_name" json:"ch_name"`
	Description string    `gorm:"column:description" json:"description"`
	UID         int64     `gorm:"column:uid" json:"uid"`
	DisableTime time.Time `gorm:"column:disable_time" json:"disable_time"`
	Ctime       time.Time `gorm:"column:ctime" json:"ctime"`
	Mtime       time.Time `gorm:"column:mtime" json:"mtime"`
}

// TableName .
func (f *Flow) TableName() string {
	return TableFlow
}

// IsAvailable .
func (f *Flow) IsAvailable() bool {
	return f.DisableTime.IsZero()
}

// FlowArr .
type FlowArr []*Flow

func (a FlowArr) Len() int {
	return len(a)
}

func (a FlowArr) Less(i, j int) bool {
	return a[i].ID < a[j].ID
}

func (a FlowArr) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
