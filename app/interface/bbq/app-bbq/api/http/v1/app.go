package v1

import (
	"go-common/app/interface/bbq/app-bbq/model"
	"go-common/library/time"
)

// AppSettingRequest .
type AppSettingRequest struct {
	Base
	VersionCode int `json:"version_code" form:"version_code" validate:"required"`
}

// AppUpdate .
type AppUpdate struct {
	NewVersion uint8             `json:"new_version"`
	Info       *model.AppVersion `json:"info,omitempty"`
}

// AppSettingResponse .
type AppSettingResponse struct {
	Public    map[string]interface{} `json:"public"`
	Update    *AppUpdate             `json:"update"`
	Resources []*model.AppResource   `json:"resources"`
}

// AppPackage .
type AppPackage struct {
	ID          int64     `json:"id"`
	Platform    uint8     `json:"platform"`
	VersionName string    `json:"version_name"`
	VersionCode uint32    `json:"version_code"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Download    string    `json:"download"`
	MD5         string    `json:"md5"`
	Size        int32     `json:"size"`
	Force       uint8     `json:"force"`
	Status      uint8     `json:"status"`
	CTime       time.Time `json:"ctime"`
}
