package model

import (
	"time"
)

// Log is the universal tag model, contains any type of tags
// The Business field and the Round field will from any business definition
type Log struct {
	AdminID  int32     `gorm:"column:adminid"`
	Oid      int64     `gorm:"column:oid"`
	Business int8      `gorm:"column:business"`
	Target   int32     `gorm:"column:target"`
	Module   int8      `gorm:"column:module"`
	Remark   string    `gorm:"column:remark"`
	Note     string    `gorm:"column:note"`
	CTime    time.Time `gorm:"column:ctime"`
	MTime    time.Time `gorm:"column:mtime"`
}

// TableName Tag tablename
func (*Log) TableName() string {
	return "workflow_log"
}
