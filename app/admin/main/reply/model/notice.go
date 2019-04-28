package model

import (
	xtime "go-common/library/time"
)

// NoticeCondition NoticeCondition
type NoticeCondition uint8

// NoticeStatus NoticeStatus
type NoticeStatus uint8

// NoticePlat NoticePlat
type NoticePlat uint8

const (
	// PlatUnknow PlatUnknow
	PlatUnknow NoticePlat = 0
	// PlatWeb PlatUnknow
	PlatWeb NoticePlat = 1
	// PlatAndroid PlatAndroid
	PlatAndroid NoticePlat = 2
	// PlatIPhone PlatIPhone
	PlatIPhone NoticePlat = 3
	// PlatWpM wp mobile
	PlatWpM NoticePlat = 4
	// PlatIPad PlatIPad
	PlatIPad NoticePlat = 5
	// PlatPadHd ipad hd
	PlatPadHd NoticePlat = 6
	// PlatWpPc win10
	PlatWpPc NoticePlat = 7
)

const (
	// StatusOffline StatusOffline
	StatusOffline NoticeStatus = 0
	// StatusOnline StatusOnline
	StatusOnline NoticeStatus = 1
)

const (
	// ConditionEQ ConditionEQ
	ConditionEQ NoticeCondition = 0 // condition equal
	// ConditionGT ConditionGT
	ConditionGT NoticeCondition = 1 // greater
	// ConditionLT ConditionLT
	ConditionLT NoticeCondition = 2 // less
)

// Notice reply's public notice
type Notice struct {
	ID         uint32          `json:"id"`
	Plat       NoticePlat      `json:"plat"`
	Version    string          `json:"version"`
	Condition  NoticeCondition `json:"condi"`
	Build      uint32          `json:"build"`
	Title      string          `json:"title"`
	Content    string          `json:"content"`
	Link       string          `json:"link"`
	StartTime  xtime.Time      `json:"stime"`
	EndTime    xtime.Time      `json:"etime"`
	Status     NoticeStatus    `json:"status"`
	CreateTime xtime.Time      `json:"ctime"`
	ModifyTime xtime.Time      `json:"mtime"`
	//client's program type
	ClientType string `json:"client_type"`
}
