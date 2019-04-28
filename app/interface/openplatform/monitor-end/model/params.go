package model

import "go-common/library/time"

// LogParams .
type LogParams struct {
	Source string `form:"source"`
	Log    string `form:"log"`
	IsAPP  int    `form:"is_app"`
}

// CollectParams .
type CollectParams struct {
	SubEvent  string `form:"sub_event" json:"sub_event"`
	Event     string `form:"event" json:"event"`
	Product   string `form:"product" json:"product"`
	Source    string `form:"source" json:"source"`
	Code      int    `form:"code" json:"code"`
	ExtJSON   string `form:"ext_json"`
	Mid       int64
	IP        string
	Buvid     string
	UserAgent string
}

// Group .
type Group struct {
	ID        int64     `form:"id" json:"id"`
	Name      string    `form:"name" json:"name"`
	Receivers string    `form:"receivers" json:"receivers"`
	Interval  int       `form:"interval" json:"interval"`
	Ctime     time.Time `json:"ctime"`
	Mtime     time.Time `json:"mtime"`
}

// Target .
type Target struct {
	ID         int64     `form:"id" json:"id"`
	SubEvent   string    `form:"sub_event" json:"sub_event"`
	Event      string    `form:"event" json:"event"`
	Product    string    `form:"product" json:"product"`
	Source     string    `form:"source" json:"source"`
	GroupIDs   string    `form:"gid" json:"-"`
	Groups     []*Group  `json:"groups"`
	States     string    `form:"states" json:"-"`
	State      int       `form:"state" json:"state"`
	Threshold  int       `form:"threshold" json:"threshold"`
	Duration   int       `form:"duration" json:"duration"`
	DeleteTime time.Time `json:"deleted_time"`
	Ctime      time.Time `json:"ctime"`
	Mtime      time.Time `json:"mtime"`
}

// Targets .
type Targets struct {
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"pagesize"`
	Draw     int       `form:"draw" json:"draw"`
	Targets  []*Target `json:"targets"`
}

// GroupListParams .
type GroupListParams struct {
	Pn   int    `form:"pn" json:"pn"`
	Ps   int    `form:"ps" json:"ps"`
	Name string `form:"name" json:"name"`
}

// Groups .
type Groups struct {
	Total    int      `json:"total"`
	Page     int      `json:"page"`
	PageSize int      `json:"pagesize"`
	Groups   []*Group `json:"groups"`
}

// Product .
type Product struct {
	ID       int64     `form:"id" json:"id"`
	Name     string    `form:"name" json:"name"`
	GroupIDs string    `form:"gid" json:"-"`
	Groups   []*Group  `json:"groups"`
	State    int       `form:"state" json:"state"`
	Ctime    time.Time `json:"ctime"`
	Mtime    time.Time `json:"mtime"`
}

// Products .
type Products struct {
	Total    int        `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"pagesize"`
	Products []*Product `json:"products"`
}
