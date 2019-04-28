package model

import "go-common/library/time"

// DBTag tag table in mysql.
type DBTag struct {
	ID        int64     `json:"id" gorm:"primary_key"`
	AppID     int64     `json:"app_id"`
	ConfigIDs string    `json:"config_ids"`
	Mark      string    `json:"mark"`
	Force     int8      `json:"force"`
	Operator  string    `json:"operator"`
	Ctime     time.Time `json:"ctime"`
	Mtime     time.Time `json:"mtime"`
}

// TableName tag.
func (DBTag) TableName() string {
	return "tag"
}
