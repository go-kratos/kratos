package model

import "go-common/library/time"

// DBApp mysql app DB.
type DBApp struct {
	ID     int64     `json:"id" gorm:"primary_key"`
	Name   string    `json:"name"`
	Token  string    `json:"token"`
	Env    string    `json:"env"`
	Zone   string    `json:"zone"`
	TreeID int       `json:"tree_id"`
	Ctime  time.Time `json:"ctime"`
	Mtime  time.Time `json:"mtime"`
}

// TableName app
func (DBApp) TableName() string {
	return "app"
}

// App app local cache.
type App struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Token  string `json:"token"`
	Env    string `json:"env"`
	Zone   string `json:"zone"`
	TreeID int    `json:"tree_id"`
}
