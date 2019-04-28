package model

import "go-common/library/time"

// InviteCode 邀请码表
type InviteCode struct {
	Code     int64     `json:"code"`
	DeviceID string    `json:"device_id"`
	Ctime    time.Time `json:"ctime"`
	Mtime    time.Time `json:"mtime"`
}
