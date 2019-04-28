package model

import (
	xtime "go-common/library/time"
)

// TableName case tablename.
func (*ServiceName) TableName() string {
	return "service_name"
}

// ServiceName service name.
type ServiceName struct {
	ID            int        `gorm:"column:id" json:"id"`
	Name          string     `gorm:"column:name" json:"name"`
	Remark        string     `gorm:"column:remark" json:"remark"`
	Token         string     `gorm:"column:token" json:"token"`
	ConfigID      string     `gorm:"column:config_id" json:"config_id"`
	ProjectTeamID string     `gorm:"column:project_team_id" json:"project_team_id"`
	Environment   int8       `gorm:"column:environment" json:"environment"`
	Public        string     `gorm:"column:public" json:"public"`
	CTime         xtime.Time `gorm:"column:ctime" json:"ctime"`
	MTime         xtime.Time `gorm:"column:mtime" json:"mtime"`
}

// TableName case tablename.
func (*ServiceConfig) TableName() string {
	return "service_config"
}

// ServiceConfig service config.
type ServiceConfig struct {
	ID        int        `gorm:"column:id" json:"id"`
	ServiceID int        `gorm:"column:service_id" json:"service_id"`
	Suffix    string     `gorm:"column:suffix" json:"suffix"`
	Config    string     `gorm:"column:config" json:"config"`
	State     int8       `gorm:"column:state" json:"state"`
	Operator  string     `gorm:"column:operator" json:"operator"`
	Remark    string     `gorm:"column:remark" json:"remark"`
	CTime     xtime.Time `gorm:"column:ctime" json:"ctime"`
	MTime     xtime.Time `gorm:"column:mtime" json:"mtime"`
}

// TableName case tablename.
func (*ServiceConfigValue) TableName() string {
	return "service_config_value"
}

// ServiceConfigValue service config value.
type ServiceConfigValue struct {
	ID          int        `gorm:"column:id" json:"id"`
	ConfigID    int        `gorm:"column:config_id" json:"config_id"`
	Name        string     `gorm:"column:name" json:"name"`
	Config      string     `gorm:"column:config" json:"config"`
	State       int8       `gorm:"column:state" json:"state"`
	Operator    string     `gorm:"column:operator" json:"operator"`
	NamespaceID int        `gorm:"column:namespace_id" json:"namespace_id"`
	CTime       xtime.Time `gorm:"column:ctime" json:"ctime"`
	MTime       xtime.Time `gorm:"column:mtime" json:"mtime"`
}

// TableName case tablename.
func (*BuildVersion) TableName() string {
	return "build_version"
}

// BuildVersion build version.
type BuildVersion struct {
	ID        int        `gorm:"column:id" json:"id"`
	ServiceID int        `gorm:"column:service_id" json:"service_id"`
	Version   string     `gorm:"column:version" json:"version"`
	Remark    string     `gorm:"column:remark" json:"remark"`
	State     int8       `gorm:"column:state" json:"state"`
	ConfigID  int        `gorm:"column:config_id" json:"config_id"`
	CTime     xtime.Time `gorm:"column:ctime" json:"ctime"`
	MTime     xtime.Time `gorm:"column:mtime" json:"mtime"`
}

//ApmCopyReq ...
type ApmCopyReq struct {
	Name    string `form:"name" validate:"required"`
	TreeID  int64  `form:"tree_id" validate:"required"`
	ApmName string `form:"apmname" validate:"required"`
}
