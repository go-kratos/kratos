package model

import "go-common/library/time"

//Team team.
type Team struct {
	ID    int64     `json:"id" gorm:"primary_key"`
	Name  string    `json:"name"`
	Env   string    `json:"env"`
	Zone  string    `json:"zone"`
	Ctime time.Time `json:"ctime"`
	Mtime time.Time `json:"mtime"`
}

// TableName team.
func (Team) TableName() string {
	return "project_team"
}
