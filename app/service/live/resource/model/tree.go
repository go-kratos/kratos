package model

import (
	xtime "go-common/library/time"
	"time"
)

//App app.
type App struct {
	ID     int64      `json:"id" gorm:"primary_key"`
	Name   string     `json:"name"`
	TreeID int64      `json:"tree_id"`
	Env    string     `json:"env"`
	Zone   string     `json:"zone"`
	Token  string     `json:"token"`
	Status int8       `json:"status"`
	Ctime  xtime.Time `json:"ctime"`
	Mtime  xtime.Time `json:"mtime"`
}

// TableName app
func (App) TableName() string {
	return "app"
}

// Node node.
type Node struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	TreeID int64  `json:"tree_id"`
}

// TreeNode TreeNode.
type TreeNode struct {
	Alias     string      `json:"alias"`
	CreatedAt string      `json:"created_at"`
	Name      string      `json:"name"`
	Path      string      `json:"path"`
	Tags      interface{} `json:"tags"`
	Type      int         `json:"type"`
}

// Res res.
type Res struct {
	Count   int         `json:"count"`
	Data    []*TreeNode `json:"data"`
	Page    int         `json:"page"`
	Results int         `json:"results"`
}

// AppPager app pager
type AppPager struct {
	Total int64  `json:"total"`
	Pn    int64  `json:"pn"`
	Ps    int64  `json:"ps"`
	Items []*App `json:"items"`
}

// Resp tree resp
type Resp struct {
	Data map[string]*Tree `json:"data"`
}

// Tree node.
type Tree struct {
	Name     string           `json:"name"`
	Type     int              `json:"type"`
	Path     string           `json:"path"`
	Tags     *TreeTag         `json:"tags"`
	Children map[string]*Tree `json:"children"`
}

//TreeTag tree tag.
type TreeTag struct {
	Ops string `json:"ops"`
	Rds string `json:"rds"`
}

//Env env.
type Env struct {
	Name     string `json:"name"`
	NikeName string `json:"nike_name"`
	Token    string `json:"token"`
}

//RoleNode roleNode .
type RoleNode struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
	Type int8   `json:"type"`
	Role int8   `json:"role"`
}

//UpdateTokenReq ...
type UpdateTokenReq struct {
	AppName string `form:"app_name" validate:"required"`
	Env     string `form:"env" validate:"required"`
	Zone    string `form:"zone" validate:"required"`
	TreeID  int64  `form:"tree_id" validate:"required"`
}

//CreateReq ...
type CreateReq struct {
	AppName string `form:"app_name" validate:"required"`
	TreeID  int64  `form:"tree_id" validate:"required"`
}

//AppListReq ...
type AppListReq struct {
	AppName string `form:"app_name"`
	Bu      string `form:"bu"`
	Team    string `form:"team"`
	Pn      int64  `form:"pn" default:"1" validate:"min=1"`
	Ps      int64  `form:"ps" default:"20" validate:"min=1"`
	Status  int8   `form:"status"`
}

//EnvsByTeamReq ...
type EnvsByTeamReq struct {
	AppName string `form:"app_name"`
	Zone    string `form:"zone"`
	Team    string `form:"team"`
}

//EnvsReq ...
type EnvsReq struct {
	AppName string `form:"app_name" validate:"required"`
	TreeID  int64  `form:"tree_id" validate:"required"`
	Zone    string `form:"zone" validate:"required"`
}

//NodeTreeReq ...
type NodeTreeReq struct {
	Node string `form:"node"`
	Team string `form:"team"`
}

//ZoneCopyReq ...
type ZoneCopyReq struct {
	AppName string `form:"app_name" validate:"required"`
	From    string `form:"from_zone" validate:"required"`
	To      string `form:"to_zone" validate:"required"`
	TreeID  int64  `form:"tree_id" validate:"required"`
}

//CasterEnvsReq ...
type CasterEnvsReq struct {
	TreeID int64  `form:"tree_id" validate:"required"`
	Zone   string `form:"zone" validate:"required"`
	Auth   string `form:"auth" validate:"required"`
}

//CacheData ...
type CacheData struct {
	Data  map[int64]*RoleNode `json:"data"`
	CTime time.Time           `json:"ctime"`
}

//AppStatusReq ...
type AppStatusReq struct {
	TreeID int64 `form:"tree_id" validate:"required"`
	Status int8  `form:"status" default:"1"`
}
