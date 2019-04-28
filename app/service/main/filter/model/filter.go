package model

import (
	xtime "go-common/library/time"
)

// FilterCacheRes struct .
type FilterCacheRes struct {
	Fmsg     string   `json:"fmsg"`
	Level    int8     `json:"level"`
	TpIDs    []int64  `json:"tpids"`
	HitRules []string `json:"hit_rules"`
	Limit    int      `json:"limit"`
	AI       *AiScore `json:"ai"`
}

// FilterAreaInfo .
type FilterAreaInfo struct {
	ID        int64  `json:"fid"`
	Mode      int8   `json:"mode"`
	Filter    string `json:"filter"`
	level     int8
	Source    int8    `json:"source"`
	Type      int8    `json:"type"`
	TpIDs     []int64 `json:"tpid"` //分区信息
	Area      string  `json:"area"`
	areaLevel int8
	Stime     xtime.Time `json:"stime"`
	Etime     xtime.Time `json:"etime"`
	Comment   string     `json:"comment"`
	State     int8       `json:"state"`
	CTime     xtime.Time `json:"ctime"` // 创建时间
}

// Level .
func (f *FilterAreaInfo) Level() int8 {
	if f.areaLevel > 0 && f.areaLevel != f.level {
		return f.areaLevel
	}
	return f.level
}

// SetLevel .
func (f *FilterAreaInfo) SetLevel(level, areaLevel int8) {
	f.level = level
	f.areaLevel = areaLevel
}
