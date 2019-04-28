package model

import (
	"go-common/app/admin/main/vip/model"
	"go-common/library/time"
)

// Privilege info.
type Privilege struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	Explain     string    `json:"explain"`
	Type        int8      `json:"type"`
	Operator    string    `json:"operator"`
	State       int8      `json:"state"`
	Deleted     int8      `json:"deleted"`
	IconURL     string    `json:"icon_url"`
	IconGrayURL string    `json:"icon_gray_url"`
	Order       int64     `json:"order"`
	LangType    int64     `json:"-"`
	Ctime       time.Time `json:"ctime"`
	Mtime       time.Time `json:"mtime"`
}

// PrivilegeResources privilege resources.
type PrivilegeResources struct {
	ID       int64     `json:"id"`
	PID      int64     `json:"pid"`
	Link     string    `json:"link"`
	ImageURL string    `json:"image_url"`
	Type     int8      `json:"type"`
	Ctime    time.Time `json:"ctime"`
	Mtime    time.Time `json:"mtime"`
}

// PrivilegeDetailResp privilege detail resp.
type PrivilegeDetailResp struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Title    string `json:"title"`
	Explain  string `json:"explain"`
	IconURL  string `json:"icon_url"`
	Type     int8   `json:"type"`
	Link     string `json:"link"`
	ImageURL string `json:"image_url"`
}

// PrivilegeResp privilege resp.
type PrivilegeResp struct {
	Name    string `json:"name"`
	IconURL string `json:"icon_url"`
	Type    int8   `json:"type"`
}

// PrivilegesResp privileges resp.
type PrivilegesResp struct {
	Title string           `json:"title"`
	List  []*PrivilegeResp `json:"list"`
}

// ResourcesType get type by platform.
func ResourcesType(p string) int8 {
	if p == "pc" {
		return model.WebResources
	}
	return model.AppResources
}
