package model

import "go-common/library/time"

//Build build.
type Build struct {
	ID       int64     `json:"id"`
	AppID    int64     `json:"app_id"`
	Name     string    `json:"name"`
	TagID    int64     `json:"tag_id"`
	Mark     string    `json:"mark"`
	Operator string    `json:"operator"`
	Ctime    time.Time `json:"ctime"`
	Mtime    time.Time `json:"mtime"`
}

// TableName build.
func (Build) TableName() string {
	return "build"
}

//CreateBuildReq ...
type CreateBuildReq struct {
	AppName string `form:"app_name" validate:"required"`
	Env     string `form:"env" validate:"required"`
	Zone    string `form:"zone" validate:"required"`
	Name    string `form:"name" validate:"required"`
	TagID   int64  `form:"tag_id" validate:"required"`
	TreeID  int64  `form:"tree_id" validate:"required"`
}

//BuildsReq ...
type BuildsReq struct {
	AppName string `form:"app_name" validate:"required"`
	Env     string `form:"env" validate:"required"`
	Zone    string `form:"zone" validate:"required"`
	TreeID  int64  `form:"tree_id" validate:"required"`
}

//BuildReq ...
type BuildReq struct {
	BuildID int64 `form:"build_id" validate:"required"`
}
