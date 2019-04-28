package model

import (
	xtime "go-common/library/time"
)

type AreaGroup struct {
	ID    int        `json:"id"`
	Name  string     `json:"name"`
	Ctime xtime.Time `json:"ctime"`
	Mtime xtime.Time `json:"mtime"`
}

type Area struct {
	ID         int        `json:"id"`
	GroupID    int        `json:"group_id"`
	Name       string     `json:"name"`
	ShowName   string     `json:"show_name"`
	CommonFlag bool       `json:"common_flag"`
	Ctime      xtime.Time `json:"ctime"`
	Mtime      xtime.Time `json:"mtime"`
}

type AreaLog struct {
	ID      int        `json:"id"`
	AdID    int        `json:"adid"`
	AdName  string     `json:"ad_name"`
	Comment string     `json:"comment"`
	State   int        `json:"state"`
	Ctime   xtime.Time `json:"ctime"`
}

type AreaGroupLog struct {
	ID      int        `json:"id"`
	AdID    int        `json:"adid"`
	AdName  string     `json:"ad_name"`
	Comment string     `json:"comment"`
	State   int        `json:"state"`
	Ctime   xtime.Time `json:"ctime"`
}

type AreaLevel struct {
	Level int8            `json:"level"` // 全站默认
	Area  map[string]int8 `json:"area"`
}
