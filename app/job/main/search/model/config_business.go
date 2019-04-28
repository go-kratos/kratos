package model

import (
	"go-common/library/database/sql"
)

// SQLBusiness single table offset
type SQLBusiness struct {
	Business string
	AppIds   string
	AssetDB  string
	AssetES  string
	AssetDtb string
}

// Bsn single table offset
type Bsn struct {
	Business string
	AppInfo  []BsnAppInfo
	AssetDB  map[string]*sql.Config
	AssetES  []string
	//AssetDtb []AssetDtb
}

// BsnAppInfo .
type BsnAppInfo struct {
	AppID       string `json:"appid"`
	IncrWay     string `json:"incr_way"`
	IncrOpen    bool   `json:"incr_open"`
	RecoverLock bool
}

// AssetDtb .
// type AssetDtb struct {
// 	dtb   map[string]*databus.Config
// 	size  int
// 	sleep int
// }
