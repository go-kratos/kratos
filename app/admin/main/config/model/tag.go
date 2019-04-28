package model

import "go-common/library/time"

//Tag db table tag.
type Tag struct {
	ID        int64     `json:"id" gorm:"primary_key"`
	AppID     int64     `json:"app_id"`
	BuildID   int64     `json:"build_id"`
	ConfigIDs string    `json:"config_ids"`
	Mark      string    `json:"mark"`
	Operator  string    `json:"operator"`
	Ctime     time.Time `json:"ctime"`
	Mtime     time.Time `json:"mtime"`
	Force     int8      `json:"force"`
}

// TableName tag.
func (Tag) TableName() string {
	return "tag"
}

// TagPager tag pager
type TagPager struct {
	Total int64  `json:"total"`
	Pn    int64  `json:"pn"`
	Ps    int64  `json:"ps"`
	Items []*Tag `json:"items"`
}

//TagConfig tagConfig.
type TagConfig struct {
	*Tag
	Confs []*Config `json:"confs"`
}

// TagConfigPager tag configs pager.
type TagConfigPager struct {
	Total int64        `json:"total"`
	Pn    int64        `json:"pn"`
	Ps    int64        `json:"ps"`
	Items []*TagConfig `json:"items"`
}

//CreateTagReq ...
type CreateTagReq struct {
	AppName   string `form:"app_name" validate:"required"`
	Env       string `form:"env" validate:"required"`
	Zone      string `form:"zone" validate:"required"`
	ConfigIDs string `form:"config_ids" validate:"required"`
	Mark      string `form:"mark" validate:"required"`
	TreeID    int64  `form:"tree_id" validate:"required"`
}

//LastTagsReq ...
type LastTagsReq struct {
	AppName string `form:"app_name" validate:"required"`
	Zone    string `form:"zone" validate:"required"`
	Env     string `form:"env" validate:"required"`
	Build   string `form:"build" validate:"required"`
	TreeID  int64  `form:"tree_id" validate:"required"`
}

//TagsByBuildReq ...
type TagsByBuildReq struct {
	AppName string `form:"app_name" validate:"required"`
	Zone    string `form:"zone" validate:"required"`
	Env     string `form:"env" validate:"required"`
	Build   string `form:"build" validate:"required"`
	Pn      int64  `form:"pn" default:"1" validate:"min=1"`
	Ps      int64  `form:"ps" default:"20" validate:"min=1"`
	TreeID  int64  `form:"tree_id" validate:"required"`
}

//TagReq ...
type TagReq struct {
	TagID int64 `form:"tag_id" validate:"required"`
}

//UpdatetagReq ...
type UpdatetagReq struct {
	AppName   string `form:"app_name" validate:"required"`
	Env       string `form:"env" validate:"required"`
	Zone      string `form:"zone" validate:"required"`
	ConfigIDs string `form:"config_ids" validate:"required"`
	Mark      string `form:"mark" validate:"required"`
	Build     string `form:"build" validate:"required"`
	TreeID    int64  `form:"tree_id" validate:"required"`
	Force     int8   `form:"force"`
}

//UpdateTagIDReq ...
type UpdateTagIDReq struct {
	AppName string `form:"app_name" validate:"required"`
	Zone    string `form:"zone" validate:"required"`
	Env     string `form:"env" validate:"required"`
	Build   string `form:"build" validate:"required"`
	TreeID  int64  `form:"tree_id" validate:"required"`
	TagID   int64  `form:"tag_id" validate:"required"`
}

//TagConfigDiff ...
type TagConfigDiff struct {
	TagID   int64  `form:"tag_id" validate:"required"`
	Name    string `form:"name" validate:"required"`
	TreeID  int64  `form:"tree_id" validate:"required"`
	AppID   int64  `form:"app_id" validate:"required"`
	BuildID int64  `form:"build_id" validate:"required"`
}
