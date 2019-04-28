package model

import "time"

// Module etc.
const (
	UpdateFactor int8 = iota + 1 // 修改因子操作
	UpdateSetting
	UpdateStat
)

// stat check text.
const (
	WaiteCheck = "待确认"
	HadCheck   = "已确认"

	MaxRemarkLen = 100
)

// Log def.
type Log struct {
	ID        int64     `json:"id"`
	RefID     int64     `json:"ref_id"`     // 关联ID
	Name      string    `json:"name"`       //操作用户
	Module    int8      `json:"module"`     //操作名称
	Context   string    `json:"context"`    //操作内容
	Ctime     time.Time `json:"ctime"`      //创建时间
	CtimeUnix int64     `json:"ctime_unix"` //创建时间
}
