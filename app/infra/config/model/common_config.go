package model

import "go-common/library/time"

// CommonConf common config.
type CommonConf struct {
	ID       int64     `json:"id" gorm:"primary_key"`
	TeamID   int64     `json:"team_id"`
	Name     string    `json:"name"`
	Comment  string    `json:"comment"`
	State    int8      `json:"state"`
	Mark     string    `json:"mark"`
	Operator string    `json:"operator"`
	Ctime    time.Time `json:"ctime"`
	Mtime    time.Time `json:"mtime"`
}

// TableName commonConfig.
func (CommonConf) TableName() string {
	return "common_config"
}
