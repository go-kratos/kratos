package manager

import "go-common/library/time"

// User for manager.
type User struct {
	UserID   int64     `json:"user_id"`
	System   string    `json:"-"`
	UserName string    `json:"user_name"`
	RoleID   int64     `json:"role_id"`
	RoleName string    `json:"role_name"`
	CTime    time.Time `json:"-"`
	MTime    time.Time `json:"-"`
}

// Role for manager.
type Role struct {
	RoleID   int64             `json:"role_id"`
	System   string            `json:"-"`
	RoleName string            `json:"role_name"`
	Auths    map[string]*Auth  `json:"auths"`
	Models   map[string]*Model `json:"models"`
	CTime    time.Time         `json:"-"`
	MTime    time.Time         `json:"-"`
}

// Auth for manager.
type Auth struct {
	AuthID   int64     `json:"auth_id"`
	System   string    `json:"-"`
	AuthName string    `json:"auth_name"`
	AuthFlag string    `json:"auth_flag"`
	CTime    time.Time `json:"-"`
	MTime    time.Time `json:"-"`
}

// Users User sorted.
type Users []*User

func (u Users) Len() int           { return len(u) }
func (u Users) Less(i, j int) bool { return int64(u[i].UserID) < int64(u[j].UserID) }
func (u Users) Swap(i, j int)      { u[i], u[j] = u[j], u[i] }

// Roles Role sorted.
type Roles []*Role

func (r Roles) Len() int           { return len(r) }
func (r Roles) Less(i, j int) bool { return int64(r[i].RoleID) < int64(r[j].RoleID) }
func (r Roles) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

// Auths Auth sorted.
type Auths []*Auth

func (a Auths) Len() int           { return len(a) }
func (a Auths) Less(i, j int) bool { return int64(a[i].AuthID) < int64(a[j].AuthID) }
func (a Auths) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
