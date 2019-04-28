package region

import (
	"go-common/app/interface/main/app-tag/model"
	"go-common/app/interface/main/app-tag/model/bangumi"
	"go-common/app/interface/main/app-tag/model/tag"
	"go-common/app/service/main/archive/api"
	xtime "go-common/library/time"
	"strconv"
)

type Region struct {
	Rid       int       `json:"tid"`
	Reid      int       `json:"reid"`
	Name      string    `json:"name"`
	Logo      string    `json:"logo"`
	Goto      string    `json:"goto"`
	Param     string    `json:"param"`
	Rank      string    `json:"-"`
	Plat      int8      `json:"-"`
	Build     int       `json:"-"`
	Condition string    `json:"-"`
	Area      string    `json:"-"`
	Language  string    `json:"-"`
	URI       string    `json:"uri,omitempty"`
	Islogo    int8      `json:"-"`
	Children  []*Region `json:"children,omitempty"`
}

type Show struct {
	Tag       *tag.Tag          `json:"tag,omitempty"`
	TopTag    []*tag.SimilarTag `json:"top_tag,omitempty"`
	NewTag    *NewTag           `json:"new_tag,omitempty"`
	Cbottom   xtime.Time        `json:"cbottom,omitempty"`
	Ctop      xtime.Time        `json:"ctop,omitempty"`
	Recommend []*ShowItem       `json:"recommend,omitempty"`
	New       []*ShowItem       `json:"new"`
	Dynamic   []*ShowItem       `json:"dynamic,omitempty"`
}

type TagTab struct {
	Tag    *tag.Tag          `json:"tag,omitempty"`
	TopTag []*tag.SimilarTag `json:"top_tag,omitempty"`
}

type ShowItem struct {
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
	Pts      int64       `json:"pts,omitempty"`
	Children []*ShowItem `json:"children,omitempty"`
	// av
	PubDate xtime.Time `json:"pubdate"`
	// av stat
	Duration int64 `json:"duration,omitempty"`
	// region
	Rid   int    `json:"rid,omitempty"`
	Rname string `json:"rname,omitempty"`
	Reid  int    `json:"reid,omitempty"`
	//new manager
	Desc  string `json:"desc,omitempty"`
	Stime string `json:"stime,omitempty"`
	Etime string `json:"etime,omitempty"`
	Like  int    `json:"like,omitempty"`
}

type NewTag struct {
	Position int               `json:"pos"`
	Tag      []*tag.SimilarTag `json:"tag"`
}

// FromArchive from archive archive.
func (i *ShowItem) FromArchive(a *api.Arc) {
	i.Title = a.Title
	i.Cover = a.Pic
	i.Param = strconv.FormatInt(a.Aid, 10)
	i.URI = model.FillURI(model.GotoAv, i.Param)
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
	i.Like = int(a.Stat.Like)
	if a.Access > 0 {
		i.Play = 0
	}
}

// FromArchive from archive archive.
func (i *ShowItem) FromBangumi(a *api.Arc, season *bangumi.SeasonInfo, bangumiType int) {
	var (
		_bangumiSeasonID  = 1
		_bangumiEpisodeID = 2
	)
	if season == nil {
		return
	}
	i.Title = a.Title
	i.Cover = a.Pic
	i.Param = strconv.FormatInt(season.SeasonID, 10)
	switch bangumiType {
	case _bangumiSeasonID:
		i.URI = model.FillURI(model.GotoBangumi, i.Param)
	case _bangumiEpisodeID:
		epid := strconv.Itoa(season.EpisodeID)
		i.URI = model.FillURIBangumi(model.GotoBangumi, i.Param, epid, season.SeasonType)
	}
	i.Goto = model.GotoBangumi
	i.Play = int(a.Stat.View)
	i.Danmaku = int(a.Stat.Danmaku)
	i.Name = a.Author.Name
	i.Reply = int(a.Stat.Reply)
	i.Fav = int(a.Stat.Fav)
	i.PubDate = a.PubDate
	i.Rid = int(a.TypeID)
	i.Rname = a.TypeName
	i.Duration = a.Duration
	i.Like = int(a.Stat.Like)
	if a.Access > 0 {
		i.Play = 0
	}
}
