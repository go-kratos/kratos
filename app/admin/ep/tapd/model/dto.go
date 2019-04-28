package model

import (
	"time"
)

// HookUrl Hook Url.
type HookUrl struct {
	ID          int64     `json:"id" gorm:"auto_increment;primary_key;column:id"`
	URL         string    `json:"url" gorm:"column:url"`
	WorkspaceID int       `json:"workspace_id" gorm:"column:workspace_id"`
	Status      int       `json:"status" gorm:"column:status"`
	UpdateBy    string    `json:"update_by" gorm:"column:update_by"`
	CTime       time.Time `json:"ctime" gorm:"column:ctime"`
	UTime       time.Time `json:"mtime" gorm:"column:mtime"`
}

// UrlEvent Url Event.
type UrlEvent struct {
	ID     int64     `json:"id" gorm:"auto_increment;primary_key;column:id"`
	Event  string    `json:"event" gorm:"column:event"`
	UrlID  int64     `json:"url_id" gorm:"column:url_id"`
	Status int       `json:"status" gorm:"column:status"`
	CTime  time.Time `json:"ctime" gorm:"column:ctime"`
	UTime  time.Time `json:"mtime" gorm:"column:mtime"`
}

// EventLog Event Log.
type EventLog struct {
	ID          int64     `json:"id" gorm:"auto_increment;primary_key;column:id"`
	Event       string    `json:"event" gorm:"column:event"`
	WorkspaceID int       `json:"workspace_id" gorm:"column:workspace_id"`
	EventID     int       `json:"event_id" gorm:"column:event_id"`
	CTime       time.Time `json:"ctime" gorm:"column:ctime"`
	UTime       time.Time `json:"mtime" gorm:"column:mtime"`
}
