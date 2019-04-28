package pprof

import (
	"go-common/library/time"
)

// TableName .
func (*Warn) TableName() string {
	return "pprof_warn"
}

// Warn .
type Warn struct {
	ID      int64     `gorm:"column:id" json:"id"`
	AppID   string    `gorm:"column:app_id" json:"app_id"`
	SvgName string    `gorm:"column:svg_name" json:"svg_name"`
	IP      string    `gorm:"column:ip" json:"ip"`
	Kind    int64     `gorm:"column:kind" json:"kind"`
	Ctime   time.Time `gorm:"column:ctime" json:"ctime"`
	Mtime   time.Time `gorm:"column:mtime" json:"mtime"`
	URL     string    `gorm:"-" json:"url"`
}

// Response .
type Response struct {
	Code int  `json:"code"`
	Data *Ins `json:"data"`
}

// Warning .
type Warning struct {
	Tags struct {
		App string `json:"app"`
	} `json:"tags"`
}

// Ins .
type Ins struct {
	Instances []struct {
		Treeid   int      `json:"treeid"`
		Hostname string   `json:"hostname"`
		Addrs    []string `json:"addrs"`
		Status   int      `json:"status"`
	} `json:"instances"`
}

// Params .
type Params struct {
	AppID     string    `form:"app_id" default:""`
	SvgName   string    `form:"svg_name" default:""`
	Kind      int64     `form:"kind" default:"0"`
	IP        string    `form:"ip" default:""`
	StartTime time.Time `form:"start_time" default:"0"`
	EndTime   time.Time `form:"end_time" default:"0"`
}
