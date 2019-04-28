package model

import "go-common/library/time"

var (
	//ConfigIng config ing.
	ConfigIng = int8(1)
	//ConfigEnd config ing.
	ConfigEnd = int8(2)
)

// Config config.
type Config struct {
	ID       int64     `json:"id" gorm:"primary_key"`
	AppID    int64     `json:"app_id"`
	Name     string    `json:"name"`
	Comment  string    `json:"comment"`
	From     int64     `json:"from"`
	State    int8      `json:"state"`
	Mark     string    `json:"mark"`
	Operator string    `json:"operator"`
	Ctime    time.Time `json:"ctime"`
	Mtime    time.Time `json:"mtime"`
}

// TableName config
func (Config) TableName() string {
	return "config"
}
