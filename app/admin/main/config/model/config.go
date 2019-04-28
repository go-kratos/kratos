package model

import "go-common/library/time"

var (
	//ConfigIng config ing.
	ConfigIng = int8(1)
	//ConfigEnd config end.
	ConfigEnd = int8(2)
)

// Config config.
type Config struct {
	ID        int64     `json:"id" gorm:"primary_key"`
	AppID     int64     `json:"app_id"`
	Name      string    `json:"name"`
	Comment   string    `json:"comment"`
	From      int64     `json:"from"`
	State     int8      `json:"state"`
	Mark      string    `json:"mark"`
	Operator  string    `json:"operator"`
	IsDelete  int8      `json:"is_delete"`
	NewCommon int64     `gorm:"-" json:"new_common"`
	Ctime     time.Time `json:"ctime"`
	Mtime     time.Time `json:"mtime"`
}

// TableName config.
func (Config) TableName() string {
	return "config"
}

//BuildFile file.
type BuildFile struct {
	*Config
	LastConf *Config `json:"last_conf"`
}

//ConfigRes configRes.
type ConfigRes struct {
	Files        []*Config    `json:"files"`
	BuildFiles   []*BuildFile `json:"build_files"`
	BuildNewFile []*Config    `json:"build_new_file"`
}

//ConfigRefs ConfigRefs.
type ConfigRefs struct {
	Configs   []*ConfigRef `json:"configs"`
	Ref       *ConfigRef   `json:"ref"`
	Name      string       `json:"name"`
	DeleteMAX int64        `json:"delete_max"`
}

//ConfigRef configRef.
type ConfigRef struct {
	ID       int64  `json:"id"`
	Mark     string `json:"mark"`
	IsDelete int8   `json:"is_delete"`
}

//ConfigSearch config search resp.
type ConfigSearch struct {
	App        string   `json:"app"`
	TreeID     int64    `json:"tree_id"`
	Builds     []string `json:"build"`
	ConfID     int64    `json:"config_id"`
	Mark       string   `json:"mark"`
	ConfName   string   `json:"conf_name"`
	ConfValues []string `json:"conf_value"`
}

//CanalTagUpdateReq ...
type CanalTagUpdateReq struct {
	AppName   string `form:"app_name" validate:"required"`
	Env       string `form:"env" validate:"required"`
	Zone      string `form:"zone" validate:"required"`
	ConfigIDs string `form:"config_ids"`
	TreeID    int64  `form:"tree_id" validate:"required"`
	Token     string `form:"token" validate:"required"`
	User      string `form:"user" validate:"required"`
	Mark      string `form:"mark" default:"canal发版"`
	Build     string `form:"build" default:"docker-1"`
	Force     int8   `form:"force"`
}

//CanalNameConfigsReq ...
type CanalNameConfigsReq struct {
	Env     string `form:"env" validate:"required"`
	Zone    string `form:"zone" validate:"required"`
	AppName string `form:"app_name" validate:"required"`
	Name    string `form:"name" validate:"required"`
	TreeID  int64  `form:"tree_id" validate:"required"`
	Token   string `form:"token" validate:"required"`
}

//CanalConfigCreateReq ...
type CanalConfigCreateReq struct {
	AppName string `form:"app_name" validate:"required"`
	Env     string `form:"env" validate:"required"`
	Zone    string `form:"zone" validate:"required"`
	Name    string `form:"name" validate:"required"`
	State   int8   `form:"state" validate:"required"`
	From    int64  `form:"from" default:"0"`
	Comment string `form:"comment" validate:"required"`
	Mark    string `form:"mark" validate:"required"`
	TreeID  int64  `form:"tree_id" validate:"required"`
	User    string `form:"user" validate:"required"`
	Token   string `form:"token" valildate:"required"`
}

//CreateConfigReq ...
type CreateConfigReq struct {
	AppName  string `form:"app_name" validate:"required"`
	Env      string `form:"env" validate:"required"`
	Zone     string `form:"zone" validate:"required"`
	Name     string `form:"name" validate:"required"`
	State    int8   `form:"state" validate:"required"`
	From     int64  `form:"from" default:"0"`
	Comment  string `form:"comment" validate:"required"`
	Mark     string `form:"mark" validate:"required"`
	TreeID   int64  `form:"tree_id" validate:"required"`
	SkipLint bool   `form:"skiplint"`
}

//UpdateConfValueReq ...
type UpdateConfValueReq struct {
	Name      string `form:"name"`
	ID        int64  `form:"config_id" validate:"required"`
	Mtime     int64  `form:"mtime" validate:"required"`
	State     int8   `form:"state" validate:"required"`
	Comment   string `form:"comment" validate:"required"`
	Mark      string `form:"mark" validate:"required"`
	NewCommon int64  `form:"new_common"`
	Ignore    int8   `form:"ignore"`
	SkipLint  bool   `form:"skiplint"`
}

//ValueReq ...
type ValueReq struct {
	ConfigID int64 `form:"config_id" validate:"required"`
}

//ConfigsByBuildIDReq ...
type ConfigsByBuildIDReq struct {
	BuildID int64 `form:"build_id" validate:"required"`
}

//ConfigsByTagIDReq ...
type ConfigsByTagIDReq struct {
	TagID int64 `form:"tag_id" validate:"required"`
}

//ConfigsByAppNameReq ...
type ConfigsByAppNameReq struct {
	Env     string `form:"env" validate:"required"`
	Zone    string `form:"zone" validate:"required"`
	AppName string `form:"app_name" validate:"required"`
	TreeID  int64  `form:"tree_id" validate:"required"`
}

//ConfigSearchAllReq ...
type ConfigSearchAllReq struct {
	Env  string `form:"env" validate:"required"`
	Zone string `form:"zone" validate:"required"`
	Like string `form:"like" validate:"required"`
}

//ConfigSearchAppReq ...
type ConfigSearchAppReq struct {
	Env     string `form:"env" validate:"required"`
	Zone    string `form:"zone" validate:"required"`
	AppName string `form:"app_name" validate:"required"`
	Like    string `form:"like" validate:"required"`
	BuildID int64  `form:"build_id" validate:"required"`
	TreeID  int64  `form:"tree_id" validate:"required"`
}

//ConfigsByNameReq ...
type ConfigsByNameReq struct {
	Env     string `form:"env" validate:"required"`
	Zone    string `form:"zone" validate:"required"`
	AppName string `form:"app_name" validate:"required"`
	Name    string `form:"name" validate:"required"`
	TreeID  int64  `form:"tree_id" validate:"required"`
}

//ConfigsReq ...
type ConfigsReq struct {
	Env     string `form:"env" validate:"required"`
	Zone    string `form:"zone" validate:"required"`
	AppName string `form:"app_name" validate:"required"`
	BuildID int64  `form:"build_id"`
	TreeID  int64  `form:"tree_id" validate:"required"`
}

//ConfigRefsReq ...
type ConfigRefsReq struct {
	Env     string `form:"env" validate:"required"`
	Zone    string `form:"zone" validate:"required"`
	AppName string `form:"app_name" validate:"required"`
	BuildID int64  `form:"build_id" validate:"required"`
	TreeID  int64  `form:"tree_id" validate:"required"`
}

//NamesByAppNameReq ...
type NamesByAppNameReq struct {
	Env     string `form:"env" validate:"required"`
	Zone    string `form:"zone" validate:"required"`
	AppName string `form:"app_name" validate:"required"`
	TreeID  int64  `form:"tree_id" validate:"required"`
}

//DiffReq ...
type DiffReq struct {
	ConfigID int64 `form:"config_id" validate:"required"`
	BuildID  int64 `form:"build_id"`
}

//ConfigDelReq ...
type ConfigDelReq struct {
	ConfigID int64 `form:"config_id" validate:"required"`
}

//ConfigBuildInfosReq ...
type ConfigBuildInfosReq struct {
	Env     string `form:"env" validate:"required"`
	Zone    string `form:"zone" validate:"required"`
	AppName string `form:"app_name" validate:"required"`
	BuildID int64  `form:"build_id"`
	TreeID  int64  `form:"tree_id" validate:"required"`
}

//ConfigUpdateReq ...
type ConfigUpdateReq struct {
	AppName string `form:"app_name" validate:"required"`
	Env     string `form:"env" validate:"required"`
	Zone    string `form:"zone" validate:"required"`
	TreeID  int64  `form:"tree_id" validate:"required"`
	Data    string `form:"data" validate:"required"`
	Token   string `form:"token" validate:"required"`
	User    string `form:"user" validate:"required"`
}

//TagUpdateReq ...
type TagUpdateReq struct {
	AppName   string `form:"app_name" validate:"required"`
	Env       string `form:"env" validate:"required"`
	Zone      string `form:"zone" validate:"required"`
	ConfigIDs string `form:"config_ids"`
	Mark      string `form:"mark" validate:"required"`
	Build     string `form:"build" validate:"required"`
	TreeID    int64  `form:"tree_id" validate:"required"`
	Token     string `form:"token" validate:"required"`
	User      string `form:"user" validate:"required"`
	Names     string `form:"names"`
	Increment int    `form:"increment"`
	Force     int8   `form:"force"`
}

//SetTokenReq ...
type SetTokenReq struct {
	Env    string `form:"env" validate:"required"`
	Zone   string `form:"zone" validate:"required"`
	App    string `form:"service" validate:"required"`
	Token  string `form:"token" validate:"required"`
	TreeID int64  `form:"tree_id" validate:"required"`
}

//HostsReq ...
type HostsReq struct {
	Env    string `form:"env" validate:"required"`
	Zone   string `form:"zone" validate:"required"`
	App    string `form:"service" validate:"required"`
	TreeID int64  `form:"tree_id" validate:"required"`
}
