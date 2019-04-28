package model

import (
	v1 "go-common/app/service/main/archive/api"
	"go-common/library/time"
)

// WxArchive .
type WxArchive struct {
	Aid      int64       `json:"aid"`
	TypeID   int32       `json:"type_id"`
	TypeName string      `json:"tname"`
	Pic      string      `json:"pic"`
	Title    string      `json:"title"`
	PubDate  time.Time   `json:"pubdate"`
	Ctime    time.Time   `json:"ctime"`
	Tags     []*WxArcTag `json:"tags"`
	Duration int64       `json:"duration"`
	Author   v1.Author   `json:"author"`
	Stat     v1.Stat     `json:"stat"`
}

// WxArcTag .
type WxArcTag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// FromArchive .
func (w *WxArchive) FromArchive(a *v1.Arc) {
	w.Aid = a.Aid
	w.TypeID = a.TypeID
	w.TypeName = a.TypeName
	w.Pic = a.Pic
	w.Title = a.Title
	w.PubDate = a.PubDate
	w.Ctime = a.Ctime
	w.Duration = a.Duration
	w.Author = a.Author
	w.Stat = a.Stat
}
