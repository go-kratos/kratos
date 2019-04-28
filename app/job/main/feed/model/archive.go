package model

import "go-common/app/service/main/archive/api"

// Archive model
type Archive struct {
	ID        int64  `json:"aid"`
	Mid       int64  `json:"mid"`
	TypeID    int16  `json:"typeid"`
	HumanRank int    `json:"humanrank"`
	Duration  int    `json:"duration"`
	Title     string `json:"title"`
	Cover     string `json:"cover"`
	Content   string `json:"content"`
	Tag       string `json:"tag"`
	Attribute int32  `json:"attribute"`
	Copyright int8   `json:"copyright"`
	AreaLimit int8   `json:"arealimit"`
	State     int    `json:"state"`
	Author    string `json:"author"`
	Access    int    `json:"access"`
	Forward   int    `json:"forward"`
	PubTime   string `json:"pubtime"`
	Round     int8   `json:"round"`
	CTime     string `json:"ctime"`
	MTime     string `json:"mtime"`
}

func (a *Archive) IsNormal() bool {
	arc := api.Arc{State: int32(a.State)}
	return arc.IsNormal()
}
