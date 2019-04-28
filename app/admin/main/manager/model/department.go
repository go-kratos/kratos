package model

import (
	"go-common/library/time"
)

// DepartCustom struct info of table user_department
type DepartCustom struct {
	ID     int64     `json:"id" gorm:"column:id"`
	Name   string    `json:"name" gorm:"column:name"`
	Status int       `json:"-" gorm:"column:status"`
	Ctime  time.Time `json:"-" gorm:"-"`
	Mtime  time.Time `json:"-"  gorm:"-"`
}

// RoleCustom .
type RoleCustom struct {
	ID          int64     `json:"id" gorm:"column:id"`
	Name        string    `json:"name" gorm:"column:name"`
	Type        int64     `json:"-" gorm:"column:type"`
	Description string    `json:"-" gorm:"column:description"`
	RuleID      int64     `json:"-" gorm:"column:rule_id"`
	Data        string    `json:"-" gorm:"column:data"`
	Ctime       time.Time `json:"-" gorm:"-"`
	Mtime       time.Time `json:"-" gorm:"-"`
}

// UserCustom .
type UserCustom struct {
	ID           int64     `json:"id" gorm:"column:id"`
	Username     string    `json:"username" gorm:"column:username"`
	Nickname     string    `json:"nickname" gorm:"column:nickname"`
	Email        string    `json:"-" gorm:"column:email"`
	Phone        string    `json:"-" gorm:"column:phone"`
	DepartmentID int       `json:"-" gorm:"column:department_id"`
	State        int       `json:"-" gorm:"column:state"`
	Ctime        time.Time `json:"-" gorm:"-"`
	Mtime        time.Time `json:"-" gorm:"-"`
}

// Department struct info of table user_department
type Department struct {
	ID     int64     `json:"id"`
	Name   string    `json:"name"`
	Status int       `json:"status"`
	Ctime  time.Time `json:"ctime"`
	Mtime  time.Time `json:"mtime"`
}

// TableName return table name
func (a Department) TableName() string {
	return "user_department"
}
