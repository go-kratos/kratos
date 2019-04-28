package model

import "go-common/library/time"

// CommonConf commonConf.
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

// CommonConfPager app pager
type CommonConfPager struct {
	Total int64         `json:"total"`
	Pn    int64         `json:"pn"`
	Ps    int64         `json:"ps"`
	Items []*CommonConf `json:"items"`
}

// CommonName app pager
type CommonName struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
}

// CommonTemp app temp
type CommonTemp struct {
	ID int64 `json:"id"`
}

// CommonCounts app counts
type CommonCounts struct {
	Counts int64 `json:"counts"`
}

//CreateComConfigReq ...
type CreateComConfigReq struct {
	Team     string `form:"team" validate:"required"`
	Env      string `form:"env" validate:"required"`
	Zone     string `form:"zone" validate:"required"`
	Name     string `form:"name" validate:"required"`
	State    int8   `form:"state" validate:"required"`
	Comment  string `form:"comment" validate:"required"`
	Mark     string `form:"mark" validate:"required"`
	SkipLint bool   `form:"skiplint"`
}

//ComValueReq ...
type ComValueReq struct {
	ConfigID int64 `form:"config_id" validate:"required"`
}

//ConfigsByTeamReq ...
type ConfigsByTeamReq struct {
	Env  string `form:"env" validate:"required"`
	Zone string `form:"zone" validate:"required"`
	Team string `form:"team" validate:"required"`
	Pn   int64  `form:"pn" default:"1" validate:"min=1"`
	Ps   int64  `form:"ps" default:"20" validate:"min=1"`
}

//ComConfigsByNameReq ...
type ComConfigsByNameReq struct {
	Env  string `form:"env" validate:"required"`
	Zone string `form:"zone" validate:"required"`
	Team string `form:"team" validate:"required"`
	Name string `form:"name" validate:"required"`
}

//UpdateComConfValueReq ...
type UpdateComConfValueReq struct {
	ID       int64  `form:"config_id" validate:"required"`
	State    int8   `form:"state" validate:"required"`
	ConfigID int64  `form:"config_id" validate:"required"`
	Name     string `form:"name" validate:"required"`
	Comment  string `form:"comment" validate:"required"`
	Mark     string `form:"mark" validate:"required"`
	Mtime    int64  `form:"mtime" validate:"required"`
	SkipLint bool   `form:"skiplint"`
}

//NamesByTeamReq ...
type NamesByTeamReq struct {
	Env  string `form:"env" validate:"required"`
	Zone string `form:"zone" validate:"required"`
	Team string `form:"team" validate:"required"`
}

// TagMap ...
type TagMap struct {
	*Tag
	AppName   string `json:"app_name"`
	BuildName string `json:"build_name"`
	TreeID    int64  `json:"tree_id"`
}
