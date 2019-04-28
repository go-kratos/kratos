package model

import (
	xtime "go-common/library/time"
)

const (
	// PrivateReadBit 私有读位
	PrivateReadBit = 0
	// PrivateWriteBit 私有写位
	PrivateWriteBit = 1
	//status

	// Public = 0
	Public = int(0)
	// PrivateRead = 1
	PrivateRead = int(1 << PrivateReadBit)
	// PrivateWrite = 2
	PrivateWrite = int(1 << PrivateWriteBit)
	// PrivateReadWrite = 3
	PrivateReadWrite = int(PrivateRead | PrivateWrite)
)

// Bucket bucekt table orm
type Bucket struct {
	ID           int         `json:"id" gorm:"column:id"`
	BucketName   string      `json:"bucket_name" gorm:"column:bucket_name"`
	Property     int         `json:"property" gorm:"column:property"`
	KeyID        string      `json:"key_id" gorm:"column:key_id"`
	KeySecret    string      `json:"key_secret" gorm:"column:key_secret"`
	PurgeCDN     bool        `json:"purge_cdn" gorm:"column:purge_cdn"`
	CacheControl int         `json:"cache_control" gorm:"column:cache_control"`
	Domain       string      `json:"domain" gorm:"column:domain"`
	CTime        xtime.Time  `json:"ctime" gorm:"column:ctime"`
	MTime        xtime.Time  `json:"mtime" gorm:"column:mtime"`
	DirLimit     []*DirLimit `json:"dir_limit" gorm:"-"`
}

// TableName bucket
func (b Bucket) TableName() string {
	return "bucket"
}

// Page common page response
type Page struct {
	PS    int `json:"ps"`
	PN    int `json:"pn"`
	Total int `json:"total"`
}

// BucketListPage bucket/list result
type BucketListPage struct {
	Items []*Bucket `json:"items"`
	Page  *Page     `json:"page"`
}
