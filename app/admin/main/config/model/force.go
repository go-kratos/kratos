package model

import "go-common/library/time"

//Force ...
type Force struct {
	ID       int64     `json:"id"`
	AppID    int64     `json:"app_id"`
	Hostname string    `json:"hostname"`
	IP       string    `json:"ip"`
	Version  int64     `json:"version"`
	Operator string    `json:"operator"`
	Ctime    time.Time `json:"ctime"`
	Mtime    time.Time `json:"mtime"`
}

// TableName force.
func (Force) TableName() string {
	return "force"
}

//CreateForceReq ...
type CreateForceReq struct {
	Env     string `form:"env" validate:"required"`
	Zone    string `form:"zone" validate:"required"`
	Build   string `form:"build" validate:"required"`
	TreeID  int64  `form:"tree_id" validate:"required"`
	Version int64  `form:"version"`
	Hosts   string `form:"hosts"`
}

//ClearForceReq ...
type ClearForceReq struct {
	Env    string `form:"env" validate:"required"`
	Zone   string `form:"zone" validate:"required"`
	Build  string `form:"build" validate:"required"`
	TreeID int64  `form:"tree_id" validate:"required"`
	Hosts  string `form:"hosts"`
}

//MapHosts ...
type MapHosts map[string]string
