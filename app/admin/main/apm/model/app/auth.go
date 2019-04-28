package app

import (
	xtime "go-common/library/time"
)

// TableName case tablename
func (*Auth) TableName() string {
	return "app_auth"
}

// Auth app auth
type Auth struct {
	ID            int64      `gorm:"column:id" json:"id"`
	ServiceTreeID int64      `gorm:"column:service_tree_id" json:"service_tree_id"`
	AppTreeID     int64      `gorm:"column:app_tree_id" json:"app_tree_id"`
	ServiceID     string     `gorm:"column:service_id" json:"service_id"`
	AppID         string     `gorm:"column:app_id" json:"app_id"`
	RPCMethod     string     `gorm:"column:rpc_method" json:"rpc_method"`
	HTTPMethod    string     `gorm:"column:http_method" json:"http_method"`
	Quota         int64      `gorm:"column:quota" json:"quota"`
	CTime         xtime.Time `gorm:"column:ctime" json:"ctime"`
	MTime         xtime.Time `gorm:"column:mtime" json:"mtime"`
}
