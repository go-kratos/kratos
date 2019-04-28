package daily

import (
	xtime "go-common/library/time"
	"strconv"

	"go-common/app/interface/main/app-show/model"
	"go-common/app/service/main/archive/api"
)

// Show is module.
type Show struct {
	*Head
	Body []*Item `json:"body"`
}

// Head is show head.
type Head struct {
	ID        int    `json:"-"`
	ColumnID  int    `json:"column_id,omitempty"`
	Plat      int8   `json:"-"`
	Title     string `json:"title"`
	Desc      string `json:"desc"`
	Rank      int    `json:"-"`
	Build     int    `json:"-"`
	Condition string `json:"-"`
	Date      int64  `json:"date,omitempty"`
	Cover     string `json:"cover,omitempty"`
	Type      string `json:"type,omitempty"`
	Goto      string `json:"goto,omitempty"`
	Param     string `json:"param,omitempty"`
	URI       string `json:"uri,omitempty"`
}

type Item struct {
	Title string `json:"title"`
	Cover string `json:"cover"`
	URI   string `json:"uri"`
	Param string `json:"param"`
	Goto  string `json:"goto"`
	// up
	Name string `json:"name,omitempty"`
	// stat
	Play    int `json:"play,omitempty"`
	Danmaku int `json:"danmaku,omitempty"`
	Reply   int `json:"reply,omitempty"`
	Fav     int `json:"favourite,omitempty"`
	// movie and bangumi badge
	Status    int8   `json:"status,omitempty"`
	CoverMark string `json:"cover_mark,omitempty"`
	// ranking
	Pts int64 `json:"pts,omitempty"`
	// av
	PubDate xtime.Time `json:"pubdate"`
	// av stat
	Duration int64 `json:"duration,omitempty"`
	// region
	Rid   int    `json:"rid,omitempty"`
	Rname string `json:"rname,omitempty"`
	// tag
	TagID   int64  `json:"tag_id,omitempty"`
	TagName string `json:"tag_name,omitempty"`
}

// ColumnList
type ColumnList struct {
	Cid      int           `json:"cid,omitempty"`
	Ceid     int           `json:"ceid,omitempty"`
	Name     string        `json:"name,omitempty"`
	Cname    string        `json:"-"`
	Children []*ColumnList `json:"children,omitempty"`
}

// FromArchivePB from archive.
func (i *Item) FromArchivePB(a *api.Arc) {
	i.Title = a.Title
	i.Cover = a.Pic
	i.Param = strconv.FormatInt(a.Aid, 10)
	i.URI = model.FillURI(model.GotoAv, i.Param, nil)
	i.Goto = model.GotoAv
	i.Play = int(a.Stat.View)
	i.Danmaku = int(a.Stat.Danmaku)
	i.Name = a.Author.Name
	i.Reply = int(a.Stat.Reply)
	i.Fav = int(a.Stat.Fav)
	i.PubDate = a.PubDate
	i.Rid = int(a.TypeID)
	i.Rname = a.TypeName
	i.Duration = a.Duration
	if a.Access > 0 {
		i.Play = 0
	}
}
