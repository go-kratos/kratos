package model

import (
	xtime "go-common/library/time"
)

// AppInfo App info.
type AppInfo struct {
	AppTreeID int64      `json:"app_tree_id"`
	AppID     string     `json:"app_id"`
	Limit     int32      `json:"limit"`
	MTime     xtime.Time `json:"mtime"`
}

// AppAuth AppAuth.
type AppAuth struct {
	ServiceTreeID int64      `json:"service_tree_id"`
	AppTreeID     int64      `json:"app_tree_id"`
	RPCMethod     string     `json:"rpc_method"`
	HTTPMethod    string     `json:"http_method"`
	Quota         int32      `json:"quota"`
	MTime         xtime.Time `json:"mtime"`
}

// Scope Scope.
type Scope struct {
	AppTreeID   int64    `json:"app_tree_id"`
	RPCMethods  []string `json:"rpc_methods"`
	HTTPMethods []string `json:"http_methods"`
	Quota       int32    `json:"quota"`
	Sign        string   `json:"sign"`
}

// AppToken AppToken.
type AppToken struct {
	AppTreeID int64  `json:"app_tree_id"`
	AppID     string `json:"app_id"`
	AppAuth   string `json:"app_auth"`
}
