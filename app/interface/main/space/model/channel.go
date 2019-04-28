package model

import (
	v1 "go-common/app/service/main/archive/api"
	xtime "go-common/library/time"
)

// Channel channel struct
type Channel struct {
	Cid   int64      `json:"cid"`
	Mid   int64      `json:"mid"`
	Name  string     `json:"name"`
	Intro string     `json:"intro"`
	Mtime xtime.Time `json:"mtime"`
	Count int        `json:"count"`
	Cover string     `json:"cover"`
}

// ChannelExtra channel extra fields
type ChannelExtra struct {
	Aid   int64
	Cid   int64
	Count int
	Cover string
}

// ChannelDetail channel detail info
type ChannelDetail struct {
	*Channel
	Archives []*v1.Arc `json:"archives"`
}

// ChannelArc channel video struct
type ChannelArc struct {
	ID       int64      `json:"id"`
	Mid      int64      `json:"mid"`
	Cid      int64      `json:"cid"`
	Aid      int64      `json:"aid"`
	OrderNum int        `json:"order_num"`
	Mtime    xtime.Time `json:"mtime"`
}

// ChannelArcSort channel archive sort struct
type ChannelArcSort struct {
	Aid      int64 `json:"aid"`
	OrderNum int   `json:"order_num"`
}
