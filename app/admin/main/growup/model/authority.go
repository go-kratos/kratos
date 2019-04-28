package model

import (
	"go-common/library/time"
)

// SUser simple user
type SUser struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Group simple task group
type Group struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Role simple task role
type Role struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// SPrivilege simple privilege
type SPrivilege struct {
	ID       int64         `json:"id" gorm:"column:id"`
	Title    string        `json:"title" gorm:"column:name"`
	Level    int64         `json:"level" gorm:"level"`
	IsRouter uint8         `json:"is_router" gorm:"is_router"`
	Children []*SPrivilege `json:"children"`
	Selected bool          `json:"selected"`
}

// User user info
type User struct {
	ID        int64     `json:"id" gorm:"column:id"`
	Username  string    `json:"username" gorm:"column:username"`
	Nickname  string    `json:"nickname" gorm:"column:nickname"`
	TaskGroup string    `json:"task_group" gorm:"column:task_group"`
	TaskRole  string    `json:"task_role" gorm:"column:task_role"`
	ATime     time.Time `json:"atime" gorm:"column:atime"`
	CTime     time.Time `json:"ctime" gorm:"column:ctime"`
	MTime     time.Time `json:"mtime" gorm:"column:mtime"`
	IsDeleted int       `json:"-"`
	Groups    []*Group  `json:"groups" gorm:"-"`
	Roles     []*Role   `json:"roles" gorm:"-"`
}

// TaskGroup task group
type TaskGroup struct {
	ID         int64     `json:"id" gorm:"column:id"`
	Name       string    `json:"name" gorm:"column:name"`
	Desc       string    `json:"desc" gorm:"column:desc"`
	Privileges string    `json:"privileges" gorm:"column:privileges"`
	ATime      time.Time `json:"atime" gorm:"column:atime"`
	CTime      time.Time `json:"ctime" gorm:"column:ctime"`
	MTime      time.Time `json:"mtime" gorm:"column:mtime"`
	IsDeleted  int       `json:"-"`
	Users      []*SUser  `json:"users" gorm:"-"`
}

// TaskRole task role
type TaskRole struct {
	ID         int64     `json:"id" gorm:"column:id"`
	Name       string    `json:"name" gorm:"column:name"`
	Desc       string    `json:"desc" gorm:"column:desc"`
	GroupID    int64     `json:"group_id" gorm:"column:group_id"`
	Privileges string    `json:"privileges" gorm:"column:privileges"`
	ATime      time.Time `json:"atime" gorm:"column:atime"`
	CTime      time.Time `json:"ctime" gorm:"column:ctime"`
	MTime      time.Time `json:"mtime" gorm:"column:mtime"`
	IsDeleted  int       `json:"-"`
	Users      []*SUser  `json:"users" gorm:"-"`
	GroupName  string    `json:"group_name" gorm:"-"`
}

// Privilege privilege
type Privilege struct {
	ID        int64     `json:"id" gorm:"column:id"`
	Name      string    `json:"name" gorm:"column:name"`
	Level     int64     `json:"level" gorm:"level"`
	FatherID  int64     `json:"father_id" gorm:"father_id"`
	IsRouter  uint8     `json:"is_router" gorm:"is_router"`
	CTime     time.Time `json:"ctime" gorm:"column:ctime"`
	MTime     time.Time `json:"mtime" gorm:"column:mtime"`
	IsDeleted int       `json:"-"`
}
