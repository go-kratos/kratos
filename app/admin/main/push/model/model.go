package model

import (
	"math/rand"
	"time"
)

const (
	// UploadTypeMid 上传文件内容为 mid
	UploadTypeMid = 1
	// UploadTypeToken 上传文件内容为 token
	UploadTypeToken = 2
)

// Page .
type Page struct {
	Num   int `json:"num"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

// Pager def.
type Pager struct {
	Total int `json:"total"`
	Pn    int `json:"page" form:"pn" validate:"min=1" default:"1"`
	Ps    int `json:"pagesize" form:"ps" validate:"min=1" default:"20"`
}

// App .
type App struct {
	ID            int64      `json:"id" form:"id"`
	Name          string     `json:"name" form:"name" validate:"required"`
	PushLimitUser int        `json:"push_limit_user" form:"push_limit_user"`
	Ctime         time.Time  `json:"ctime"`
	Mtime         time.Time  `json:"mtime"`
	Dtime         int64      `json:"dtime"`
	Business      []Business `json:"-"`
	Auths         []Auth     `json:"-"`
}

// Auth .
type Auth struct {
	ID         int64     `json:"id" form:"id"`
	AppID      int64     `json:"app_id" form:"app_id"`
	PlatformID int       `json:"platform_id" form:"platform_id"`
	Name       string    `json:"name" form:"name"`
	Key        string    `json:"key" form:"key"`
	Value      string    `json:"value" form:"value"`
	BundleID   string    `json:"bundle_id" form:"bundle_id"`
	Mtime      time.Time `json:"mtime"`
	Ctime      time.Time `json:"ctime"`
	Dtime      int       `json:"dtime"`
}

// Business .
type Business struct {
	ID            int64     `json:"id" form:"id"`
	Ctime         time.Time `json:"ctime"`
	Mtime         time.Time `json:"mtime"`
	Dtime         int       `json:"dtime"`
	AppID         int64     `json:"app_id" form:"app_id"`
	Name          string    `json:"name" form:"name"`
	Desc          string    `json:"desc" gorm:"column:description" form:"desc"`
	Token         string    `json:"token"`
	Sound         int       `json:"sound" form:"sound"`
	Vibration     int       `json:"vibration" form:"vibration"`
	ReceiveSwitch int       `json:"receive_switch" form:"receive_switch"`
	PushSwitch    int       `json:"push_switch" form:"push_switch"`
	AppName       string    `json:"app_name" gorm:"-"`
	SilentTime    string    `json:"silent_time" form:"silent_time"`
	PushLimitUser int       `json:"push_limit_user" form:"push_limit_user"`
	Whitelist     int       `json:"whitelist" form:"whitelist"`
}

// TableName .
func (b Business) TableName() string {
	return "push_business"
}

// Task .
type Task struct {
	ID             string    `json:"id" form:"id"`
	Job            string    `json:"job" form:"job"`
	Type           int       `json:"type" form:"type"`
	AppID          int64     `json:"app_id" form:"app_id"`
	PlatformID     int       `json:"platform_id"`
	BusinessID     int64     `json:"business_id" form:"business_id"`
	Platform       string    `json:"platform"`
	Title          string    `json:"title" form:"title"`
	Summary        string    `json:"summary" form:"summary"`
	LinkType       int       `json:"link_type" form:"link_type"`
	LinkValue      string    `json:"link_value" form:"link_value"`
	Build          string    `json:"build" form:"build"`
	Sound          int       `json:"sound" form:"sound"`
	Vibration      int       `json:"vibration" form:"vibration"`
	MidFile        string    `json:"mid_file" form:"mid_file"`
	Progress       string    `json:"progress"`
	PushTime       time.Time `json:"-"`
	ExpireTime     time.Time `json:"-"`
	PassThrough    int       `json:"pass_through" form:"pass_through"`
	PushTimeUnix   int64     `json:"push_time" form:"push_time" gorm:"-"`
	ExpireTimeUnix int64     `json:"expire_time" form:"expire_time" gorm:"-"`
	Status         int       `json:"status"`
	ImageURL       string    `json:"image_url" form:"image_url"`
	Group          string    `json:"group" form:"group"`
	Extra          string    `json:"extra"`
	Mtime          time.Time `json:"mtime"`
	Ctime          time.Time `json:"ctime"`
	Dtime          int       `json:"dtime"`
}

// RandomString gets random string by length.
func RandomString(l int) string {
	bs := []byte("0123456789abcdefghijklmnopqrstuvwxyz")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var res []byte
	for i := 0; i < l; i++ {
		res = append(res, bs[r.Intn(len(bs))])
	}
	return string(res)
}
