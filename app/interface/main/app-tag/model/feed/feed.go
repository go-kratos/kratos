package feed

import (
	"strconv"

	"go-common/app/interface/main/app-tag/model"
	"go-common/app/service/main/archive/api"
	xtime "go-common/library/time"
)

// Item is feed item, contains av, bangumi, live, banner, feed...
type Item struct {
	Title   string `json:"title,omitempty"`
	Cover   string `json:"cover,omitempty"`
	URI     string `json:"uri,omitempty"`
	Param   string `json:"param,omitempty"`
	Goto    string `json:"goto,omitempty"`
	Rid     int16  `json:"tid,omitempty"`
	TName   string `json:"tname,omitempty"`
	Desc    string `json:"desc,omitempty"`
	Play    int    `json:"play,omitempty"`
	Danmaku int    `json:"danmaku,omitempty"`
	Reply   int    `json:"reply,omitempty"`
	Fav     int    `json:"favorite,omitempty"`
	Coin    int    `json:"coin,omitempty"`
	Share   int    `json:"share,omitempty"`
	Like    int    `json:"like,omitempty"`
	// av stat
	Duration int64 `json:"duration,omitempty"`
	// upper
	Mid   int64      `json:"mid,omitempty"`
	Name  string     `json:"name,omitempty"`
	Face  string     `json:"face,omitempty"`
	CTime xtime.Time `json:"ctime,omitempty"`
}

func (i *Item) FromArc(a *api.Arc) {
	if i.Title == "" {
		i.Title = a.Title
	}
	if i.Cover == "" {
		i.Cover = a.Pic
	}
	i.Param = strconv.FormatInt(a.Aid, 10)
	i.Goto = model.GotoAv
	i.URI = model.FillURI(i.Goto, i.Param)
	i.Rid = int16(a.TypeID)
	i.TName = a.TypeName
	i.Desc = a.Desc
	i.fillStat(a)
	i.Duration = a.Duration
	i.Mid = a.Author.Mid
	i.Name = a.Author.Name
	i.Face = a.Author.Face
	i.CTime = a.PubDate
}

func (i *Item) fillStat(a *api.Arc) {
	if a.Access == 0 {
		i.Play = int(a.Stat.View)
	}
	i.Danmaku = int(a.Stat.Danmaku)
	i.Reply = int(a.Stat.Reply)
	i.Fav = int(a.Stat.Fav)
	i.Coin = int(a.Stat.Coin)
	i.Share = int(a.Stat.Share)
	i.Like = int(a.Stat.Like)
}
