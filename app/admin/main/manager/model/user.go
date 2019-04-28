package model

import (
	"go-common/library/time"
)

// const rbac const
const (
	// auth_item type
	TypePointer  = 1
	TypeCategory = 2
	TypeRole     = 3
	TypeGroup    = 4

	// Admin super admin
	Admin = 1
	// user state
	UserStateOn = 0
	UserDepOn   = 1
)

// Account dashboard user account
type Account struct {
	Username string `json:"username"`
}

// Auth .
type Auth struct {
	UID        int64    `json:"uid"`
	Username   string   `json:"username"`
	Nickname   string   `json:"nickname"`
	Perms      []string `json:"perms"`
	Admin      bool     `json:"admin"`
	Assignable bool     `json:"assignable"`
}

// Permissions .
type Permissions struct {
	UID   int64      `json:"uid"`
	Perms []string   `json:"perms"`
	Admin bool       `json:"admin"`
	Orgs  []*AuthOrg `json:"orgs"`
	Roles []*AuthOrg `json:"roles"`
}

// UserDept .
type UserDept struct {
	ID         int64  `json:"id" gorm:"column:id"`
	Department string `json:"department" gorm:"column:department"`
}

// user table

// User struct info of table user
type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Nickname     string    `json:"nickname"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	DepartmentID int       `json:"department_id"`
	State        int       `json:"state"`
	Ctime        time.Time `json:"ctime"`
	Mtime        time.Time `json:"mtime"`
}

// TableName return table name
func (a *User) TableName() string {
	return "user"
}

// UserPager def.
type UserPager struct {
	Pn    int     `json:"pn"`
	Ps    int     `json:"ps"`
	Items []*User `json:"items"`
}

// auth table

// AuthAssign struct info of table auth_assignment
type AuthAssign struct {
	ID     int64     `json:"id"`
	ItemID int64     `json:"item_id"`
	UserID int64     `json:"user_id"`
	Ctime  time.Time `json:"ctime"`
}

// TableName return table name
func (a AuthAssign) TableName() string {
	return "auth_assignment"
}

// AuthItem struct info of table auth_item
type AuthItem struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Type        int       `json:"type"`
	Description string    `json:"description"`
	RuleID      int64     `json:"rule_id"`
	Data        string    `json:"data"`
	Ctime       time.Time `json:"ctime"`
	Mtime       time.Time `json:"mtime"`
}

// TableName return table name
func (a AuthItem) TableName() string {
	return "auth_item"
}

// AuthItemChild stuct info of table auth_item_child
type AuthItemChild struct {
	ID     int64 `json:"id"`
	Parent int64 `json:"parent"`
	Child  int64 `json:"child"`
}

// TableName return table name
func (a AuthItemChild) TableName() string {
	return "auth_item_child"
}

// Role role info.
type Role struct {
	ID          int64    `json:"id"`
	OrgName     string   `json:"org_name"`
	RoleName    string   `json:"name"`
	Users       []string `json:"users"`
	Description string   `json:"description"`
}

// RoleAss .
type RoleAss struct {
	ID     int64     `json:"id" gorm:"column:id"`
	RoleID int64     `json:"role_id" gorm:"column:role_id"`
	UserID int64     `json:"user_id" gorm:"column:user_id"`
	Ctime  time.Time `json:"ctime" gorm:"-"`
}

// Org org info.
type Org struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Users       []string `json:"users"`
	Description string   `json:"description"`
}

// AuthOrg org info.
type AuthOrg struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type int    `json:"-"` // group or role, ignore it when transform to json
}

// TableName return table name
func (a AuthOrg) TableName() string {
	return "auth_item"
}

// RespPerm response Permission
type RespPerm struct {
	Res    []string
	Admin  bool
	Groups []*AuthOrg
	Roles  []*AuthOrg
}

// AssignRole auth assignment info.
type AssignRole struct {
	Parent     string      `json:"parent"`
	Name       string      `json:"name"`
	Assignable []*AuthItem `json:"children"`
	Items      []int64     `json:"items"`
}

// AssignOrg auth assignment info.
type AssignOrg struct {
	Name       string      `json:"name"`
	Assignable []*AuthItem `json:"children"`
	Items      []int64     `json:"items"`
}
