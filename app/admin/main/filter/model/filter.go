package model

import (
	xtime "go-common/library/time"
)

const (
	FilterStateNormal  = 0
	FilterStateDeleted = 1
	FilterStateExpired = 2

	RegMode = 0
	StrMode = 1
)

// FilterForSearch 返回前端搜索列表时显示
type FilterForSearch struct {
	FilterInfo
	Areas     []string `json:"areas"`
	ShowLevel int8     `json:"level"`
}

func (f *FilterForSearch) LoadFromFilter(filter *FilterInfo) {
	f.FilterInfo = *filter
	if len(filter.Areas) == 1 && filter.Areas[0].Level != filter.Level {
		f.ShowLevel = filter.Areas[0].Level
	} else {
		f.ShowLevel = filter.Level
	}
	for _, fa := range filter.Areas {
		f.Areas = append(f.Areas, fa.Area)
	}
}

// FilterForGet 返回前端编辑敏感词时显示
type FilterForGet struct {
	FilterInfo
	Areas     []string  `json:"areas"`
	AreaLevel AreaLevel `json:"level"`
}

func (f *FilterForGet) LoadFromFilter(filter *FilterInfo) {
	f.FilterInfo = *filter
	f.AreaLevel.Level = filter.Level
	f.AreaLevel.Area = make(map[string]int8)
	for _, fa := range filter.Areas {
		f.Areas = append(f.Areas, fa.Area)
		if fa.Level != filter.Level {
			f.AreaLevel.Area[fa.Area] = fa.Level
		}
	}
}

// FilterInfo .
type FilterInfo struct {
	ID      int64         `json:"fid"`
	Mode    int8          `json:"mode"`
	Filter  string        `json:"filter"`
	Level   int8          `json:"-"`
	Source  int8          `json:"source"`
	Type    int8          `json:"type"`
	TpIDs   []int64       `json:"tpid"` //分区信息
	Areas   []*FilterArea `json:"-"`
	Stime   xtime.Time    `json:"stime"`
	Etime   xtime.Time    `json:"etime"`
	Comment string        `json:"comment"`
	State   int8          `json:"state"`
	CTime   xtime.Time    `json:"ctime"` // 创建时间
}

// FilterArea rule in area
type FilterArea struct {
	Area  string
	Level int8
}
