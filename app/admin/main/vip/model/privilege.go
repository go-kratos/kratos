package model

import (
	"go-common/library/time"
)

// Privilege info.
type Privilege struct {
	ID          int64     `gorm:"column:id" json:"id"`
	Name        string    `gorm:"column:privileges_name" json:"name"`
	Title       string    `gorm:"column:title" json:"title"`
	Explain     string    `gorm:"column:explains" json:"explain"`
	Type        int8      `gorm:"column:privileges_type" json:"type"`
	Operator    string    `gorm:"column:operator" json:"operator"`
	State       int8      `gorm:"column:state" json:"state"`
	Deleted     int8      `gorm:"column:deleted" json:"deleted"`
	IconURL     string    `gorm:"column:icon_url" json:"icon_url"`
	IconGrayURL string    `gorm:"column:icon_gray_url" json:"icon_gray_url"`
	Order       int64     `gorm:"column:order_num" json:"order"`
	LangType    int8      `gorm:"column:lang_type" json:"lang_type"`
	Ctime       time.Time `gorm:"column:ctime" json:"ctime"`
	Mtime       time.Time `gorm:"column:mtime" json:"mtime"`
}

// TableName for grom.
func (s *Privilege) TableName() string {
	return "vip_privileges"
}

// PrivilegeResources privilege resources.
type PrivilegeResources struct {
	ID       int64     `gorm:"column:id" json:"id"`
	PID      int64     `gorm:"column:pid" json:"pid"`
	Link     string    `gorm:"column:link" json:"link"`
	ImageURL string    `gorm:"column:image_url" json:"image_url"`
	Type     int8      `gorm:"column:resources_type" json:"type"`
	Ctime    time.Time `gorm:"column:ctime" json:"ctime"`
	Mtime    time.Time `gorm:"column:mtime" json:"mtime"`
}

// TableName for grom.
func (s *PrivilegeResources) TableName() string {
	return "vip_privileges_resources"
}

// PrivilegeResp  resp.
type PrivilegeResp struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Explain     string `json:"explain"`
	Type        int8   `json:"type"`
	Operator    string `json:"operator"`
	State       int8   `json:"state"`
	IconURL     string `json:"icon_url"`
	IconGrayURL string `json:"icon_gray_url"`
	Order       int64  `json:"order"`
	WebLink     string `json:"web_link"`
	WebImageURL string `json:"web_image_url"`
	AppLink     string `json:"app_link"`
	AppImageURL string `json:"app_image_url"`
	LangType    int8   `json:"lang_type"`
}
