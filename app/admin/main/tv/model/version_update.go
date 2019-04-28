package model

import (
	"go-common/library/time"
)

// VersionUpdate .
type VersionUpdate struct {
	ID         int64     `json:"id"`
	VID        int       `json:"vid" gorm:"column:vid"`
	Channel    string    `json:"channel"`
	Coverage   int32     `json:"coverage"`
	Size       int       `json:"size"`
	URL        string    `json:"url" gorm:"column:url"`
	Md5        string    `json:"md5"`
	State      int8      `json:"state"`
	Ctime      time.Time `json:"ctime"`
	Mtime      time.Time `json:"mtime"`
	Sdkint     int       `json:"sdkint"`
	Model      string    `json:"model"`
	Policy     int8      `json:"policy"`
	IsForce    int8      `json:"is_force"`
	PolicyName string    `json:"policy_name"`
	IsPush     int8      `json:"is_push"`
}

// VersionUpdateLimit .
type VersionUpdateLimit struct {
	ID    int64  `json:"id"`
	UPID  int32  `json:"up_id" gorm:"column:up_id"`
	Condi string `json:"condi"`
	Value int    `json:"value"`
}

// VersionUpdateDetail .
type VersionUpdateDetail struct {
	*VersionUpdate
	VerLimit []*VersionUpdateLimit `json:"ver_limit"`
}

// TableName version_update
func (v VersionUpdate) TableName() string {
	return "version_update"
}

// TableName version_update_limit
func (l VersionUpdateLimit) TableName() string {
	return "version_update_limit"
}

// VersionUpdatePager def.
type VersionUpdatePager struct {
	TotalCount int64                  `json:"total_count"`
	Pn         int                    `json:"pn"`
	Ps         int                    `json:"ps"`
	Items      map[string]interface{} `json:"items"`
}

// Version .
type Version struct {
	ID          int64     `json:"id"`
	Plat        int8      `json:"plat"`
	Description string    `json:"description"`
	Version     string    `json:"version"`
	Build       int       `json:"build"`
	State       int8      `json:"state"`
	Ptime       time.Time `json:"ptime"`
	Ctime       time.Time `json:"ctime"`
	Mtime       time.Time `json:"mtime"`
}

// TableName version
func (*Version) TableName() string {
	return "version"
}

// VersionPager def.
type VersionPager struct {
	TotalCount int64      `json:"total_count"`
	Pn         int        `json:"pn"`
	Ps         int        `json:"ps"`
	Items      []*Version `json:"items"`
}
