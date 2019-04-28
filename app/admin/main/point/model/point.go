package model

import "time"

// PointConf .
type PointConf struct {
	ID         int64     `json:"id" form:"id"`
	AppID      int64     `json:"app_id" form:"app_id"`
	Point      int64     `json:"point" form:"point"`
	Operator   string    `json:"operator" form:"operator"`
	ChangeType int64     `json:"change_type" form:"change_type"`
	Name       string    `json:"business_name" form:"name"`
	Ctime      time.Time `json:"-"`
	Mtime      time.Time `json:"mtime"`
}

// PointHistory .
type PointHistory struct {
	ID           int64   `json:"id"`
	Mid          int64   `json:"mid"`
	OrderID      string  `json:"order_id"`
	RelationID   string  `json:"relation_id"`
	PointBalance float64 `json:"point_balance"`
	ChangeTime   string  `json:"change_time"`
	ChangeType   int8    `json:"change_type"`
	Remark       string  `json:"remark"`
	Operator     string  `json:"operator"`
}

// AppInfo .
type AppInfo struct {
	ID       int64     `json:"_"`
	Name     string    `json:"name"`
	AppKey   string    `json:"app_key"`
	PurgeURL string    `json:"purge_url"`
	Ctime    time.Time `json:"-"`
	Mtime    time.Time `json:"mtime"`
}
