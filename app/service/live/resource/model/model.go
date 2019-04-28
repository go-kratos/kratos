package model

import "go-common/library/time"

// TableName is used to identify table name in gorm
func (cf *SundryConfig) TableName() string {
	return "ap_sundry_config"
}

// TableName is used to identify table name in gorm
func (cf *SundyConfigObject) TableName() string {
	return "ap_sundry_config"
}

// InsertMaps is used to insertDb format
type InsertMaps struct {
	Team    int64  `protobuf:"varint,2,opt,name=team,proto3" json:"team"`
	Keyword string `protobuf:"bytes,3,opt,name=keyword,proto3" json:"keyword"`
	Name    string `protobuf:"bytes,4,opt,name=name,proto3" json:"name"`
	Value   string `protobuf:"bytes,5,opt,name=value,proto3" json:"value"`
	Status  int64  `protobuf:"bytes,5,opt,name=value,proto3" json:"status"`
}

// TableName is used to identify table name in gorm
func (cf *InsertMaps) TableName() string {
	return "ap_sundry_config"
}

// SundyConfigObject is used to format select
type SundyConfigObject struct {
	Id      int64     `json:"id" gorm:"comumn:id"`
	Team    int64     `json:"team" gorm:"comumn:team"`
	Keyword string    `json:"index" gorm:"comumn:keyword"`
	Name    string    `json:"name" gorm:"comumn:name"`
	Value   string    `json:"value" gorm:"comumn:value"`
	Ctime   time.Time `json:"ctime" gorm:"comumn:ctime"`
	Mtime   time.Time `json:"mtime" gorm:"comumn:mtime"`
	Status  int64     `json:"status" gorm:"comumn:status"`
}

// TableName is used to identify table name in gorm
func (cf *ServiceConfigObject) TableName() string {
	return "ap_services_config"
}

// TableName is used to identify table name in gorm
func (cf *InsertServiceConfig) TableName() string {
	return "ap_services_config"
}

// TableName is used to identify table name in gorm
func (cf *UpdateServiceConfig) TableName() string {
	return "ap_services_config"
}

// ServiceConfigObject is used to format select
type ServiceConfigObject struct {
	Id       int64     `json:"id" gorm:"comumn:id"`
	TreeName string    `protobuf:"bytes,2,opt,name=tree_name,proto3" json:"tree_name"`
	TreePath string    `protobuf:"bytes,3,opt,name=tree_path,proto3" json:"tree_path"`
	TreeId   int64     `protobuf:"varint,4,opt,name=tree_id,proto3" json:"tree_id"`
	Service  string    `protobuf:"bytes,5,opt,name=service,proto3" json:"service"`
	Keyword  string    `protobuf:"bytes,6,opt,name=keyword,proto3" json:"keyword"`
	Template int64     `protobuf:"varint,7,opt,name=template,proto3" json:"template"`
	Value    string    `protobuf:"bytes,8,opt,name=value,proto3" json:"value"`
	Name     string    `protobuf:"bytes,9,opt,name=name,proto3" json:"name"`
	Status   int64     `protobuf:"varint,10,opt,name=status,proto3" json:"status"`
	Ctime    time.Time `protobuf:"bytes,7,opt,name=ctime,proto3" json:"ctime"`
	Mtime    time.Time `protobuf:"bytes,8,opt,name=mtime,proto3" json:"mtime"`
}

// InsertServiceConfig is used to insertDb format
type InsertServiceConfig struct {
	TreeName string `protobuf:"bytes,3,opt,name=tree_name,proto3" json:"tree_name"`
	TreePath string `protobuf:"bytes,3,opt,name=tree_path,proto3" json:"tree_path"`
	TreeId   int64  `protobuf:"bytes,3,opt,name=tree_id,proto3" json:"tree_id"`
	Service  string `protobuf:"bytes,2,opt,name=service,json=servuce,proto3" `
	Keyword  string `protobuf:"bytes,3,opt,name=template,proto3" json:"keyword"`
	Template int64  `protobuf:"bytes,3,opt,name=template,proto3" json:"template"`
	Value    string `protobuf:"bytes,4,opt,name=value,proto3" json:"value"`
	Name     string `protobuf:"bytes,5,opt,name=name,proto3" json:"name"`
	Status   int64  `protobuf:"varint,6,opt,name=status,proto3" json:"status"`
}

// UpdateServiceConfig is used to insertDb format
type UpdateServiceConfig struct {
	Service  string `protobuf:"bytes,3,opt,name=service,proto3" json:"service"`
	Keyword  string `protobuf:"bytes,3,opt,name=template,proto3" json:"keyword"`
	Template int64  `protobuf:"bytes,3,opt,name=template,proto3" json:"template"`
	Value    string `protobuf:"bytes,4,opt,name=value,proto3" json:"value"`
	Name     string `protobuf:"bytes,5,opt,name=name,proto3" json:"name"`
	Status   int64  `protobuf:"varint,6,opt,name=status,proto3" json:"status"`
}
