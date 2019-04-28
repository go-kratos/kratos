package net

import (
	"time"
)

const (
	// TableNet .
	TableNet = "net"
)

// Recovered .
var Recovered = time.Time{}

// Net net.
type Net struct {
	ID          int64     `gorm:"primary_key" json:"id" form:"id" validate:"omitempty,gt=0"`
	BusinessID  int64     `gorm:"column:business_id" json:"business_id" form:"business_id" validate:"omitempty,gt=0"`
	ChName      string    `gorm:"column:ch_name" json:"ch_name" form:"ch_name" validate:"required,max=32"`
	Description string    `gorm:"column:description" json:"description" form:"description" validate:"max=60"`
	StartFlowID int64     `gorm:"column:start_flow_id" json:"start_flow_id" form:"start_flow_id"`
	PID         int64     `gorm:"column:pid" json:"pid"`
	UID         int64     `gorm:"column:uid" json:"uid"`
	DisableTime time.Time `gorm:"column:disable_time" json:"disable_time"`
	Ctime       time.Time `gorm:"column:ctime" json:"ctime"`
	Mtime       time.Time `gorm:"column:mtime" json:"mtime"`
}

//TableName table name
func (n *Net) TableName() string {
	return TableNet
}

// IsAvailable .
func (n *Net) IsAvailable() bool {
	return n.DisableTime.IsZero()
}
