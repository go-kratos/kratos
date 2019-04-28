package model

import (
	"go-common/library/time"
)

// Bucket in accord with bucket table in database
type Bucket struct {
	ID           int                   `json:"id" gorm:"column:id"`
	Name         string                `json:"name" gorm:"column:bucket_name"`
	Property     int                   `json:"property" gorm:"column:property"`
	KeyID        string                `json:"key_id" gorm:"column:key_id"`
	KeySecret    string                `json:"key_secret" gorm:"column:key_secret"`
	PurgeCDN     bool                  `json:"purge_cdn" gorm:"column:purge_cdn"`
	CacheControl int                   `json:"cache_control" gorm:"column:purge_cdn"`
	CTime        time.Time             `json:"ctime" gorm:"column:ctime"`
	MTime        time.Time             `json:"mtime" gorm:"column:mtime"`
	DirLimit     map[string]*DirConfig `json:"dir_limit" gorm:"-"`
}

// DirLimit in accord with dir_limit table in database
type DirLimit struct {
	ID            int       `json:"id" gorm:"column:id"`
	BucketName    string    `json:"bucket_name" gorm:"column:bucket_name"`
	Dir           string    `json:"dir" gorm:"column:dir"`
	DirPicConfig  string    `json:"dir_pic_config" gorm:"column:config_pic"`
	DirRateConfig string    `json:"dir_rate_config" gorm:"column:config_rate"`
	CTime         time.Time `json:"ctime" gorm:"column:ctime"`
	MTime         time.Time `json:"mtime" gorm:"column:mtime"`
}

// TableName return table name.
func (b Bucket) TableName() string {
	return "bucket"
}

// TableName return table name.
func (l DirLimit) TableName() string {
	return "dir_limit"
}
