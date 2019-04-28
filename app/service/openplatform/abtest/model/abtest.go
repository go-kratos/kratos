package model

import (
	"go-common/library/time"
)

// AB AB测试实验
type AB struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Desc       string    `json:"desc"`
	Stra       Stra      `json:"stra"`
	Seed       int       `json:"seed"`
	Result     int       `json:"result"`
	Status     int       `json:"status"`
	Version    int       `json:"version"`
	Group      int       `json:"group"`
	Author     string    `json:"author"`
	Modifier   string    `json:"modifier"`
	CreateTime time.Time `json:"ctime"`
	ModifyTime time.Time `json:"mtime"`
}

// Stat .
type Stat struct {
	New map[int]map[int]int `json:"now"`
	Old map[int]map[int]int `json:"last"`
}

// Empty .
type Empty struct{}
