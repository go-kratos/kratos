package show

import (
	"strconv"
	"time"

	"go-common/app/interface/main/app-show/model"
	"go-common/app/interface/main/app-show/model/activity"
	"go-common/app/interface/main/app-show/model/bangumi"
	"go-common/app/interface/main/app-show/model/banner"
	"go-common/app/interface/main/app-show/model/live"
	"go-common/app/interface/main/app-show/model/recommend"
	"go-common/app/service/main/archive/api"
	seasongrpc "go-common/app/service/openplatform/pgc-season/api/grpc/season/v1"
)

const (
	_activityForm = "2006-01-02 15:04:05"
)

// Show is module.
type Show struct {
	*Head
	Body   []*Item                     `json:"body"`
	Banner map[string][]*banner.Banner `json:"banner,omitempty"`
	Ext    *Ext                        `json:"ext,omitempty"`
}

// Slice is for sort.
type Slice []*Show

func (ss Slice) Len() int           { return len(ss) }
func (ss Slice) Less(i, j int) bool { return ss[i].Rank > ss[j].Rank }
func (ss Slice) Swap(i, j int)      { ss[i], ss[j] = ss[j], ss[i] }

// Head is show head.
type Head struct {
	ID        int    `json:"-"`
	CardID    int    `json:"card_id,omitempty"`
	Plat      int8   `json:"-"`
	Param     string `json:"param"`
	Type      string `json:"type"`
	Style     string `json:"style"`
	Title     string `json:"title"`
	Rank      int    `json:"-"`
	Build     int    `json:"-"`
	Condition string `json:"-"`
	Language  string `json:"-"`
	Date      int64  `json:"date,omitempty"`
	Cover     string `json:"cover,omitempty"`
	URI       string `json:"uri,omitempty"`
	Goto      string `json:"goto,omitempty"`
}

// Item is show item, contains av, bangumi, live, banner, feed...
type Item struct {
	Sid    int    `json:"-"`
	Title  string `json:"title"`
	Cover  string `json:"cover"`
	URI    string `json:"uri"`
	NewURI string `json:"-"`
	Param  string `json:"param"`
	Goto   string `json:"goto"`
	Random int    `json:"-"`
	// av
	Play    int    `json:"play,omitempty"`
	Danmaku int    `json:"danmaku,omitempty"`
	Area    string `json:"area,omitempty"`
	AreaID  int    `json:"area_id,omitempty"`
	Rname   string `json:"rname,omitempty"`
	// av stat
	Duration int64 `json:"duration,omitempty"`
	// live and feed
	Name string `json:"name,omitempty"`
	Face string `json:"face,omitempty"`
	// only live
	Online int `json:"online,omitempty"`
	// only feed
	CTime int64 `json:"ctime,omitempty"`
	// bangumi
	Finish     int8   `json:"finish,omitempty"`
	Index      string `json:"index,omitempty"`
	TotalCount string `json:"total_count,omitempty"`
	MTime      string `json:"mtime,omitempty"`
	// movie and bangumi badge
	Status    int8   `json:"status,omitempty"`
	CoverMark string `json:"cover_mark,omitempty"`
	Fav       int    `json:"favourite,omitempty"`
	// rank
	Rank int `json:"-"`
	// cpm
	RequestId   string `json:"request_id,omitempty"`
	CreativeId  int    `json:"creative_id,omitempty"`
	SrcId       int    `json:"src_id,omitempty"`
	IsAd        bool   `json:"is_ad"`
	IsAdReplace bool   `json:"-"`
	IsAdLoc     bool   `json:"is_ad_loc,omitempty"`
	CmMark      int    `json:"cm_mark"`
	AdCb        string `json:"ad_cb,omitempty"`
	ShowUrl     string `json:"show_url,omitempty"`
	ClickUrl    string `json:"click_url,omitempty"`
	ClientIp    string `json:"client_ip,omitempty"`
	// article
	CateID     int      `json:"cate_id,omitempty"`
	CateName   string   `json:"cate_name,omitempty"`
	Summary    string   `json:"summary,omitempty"`
	Covers     []string `json:"covers,omitempty"`
	Reply      int      `json:"reply,omitempty"`
	TemplateID int      `json:"template_id,omitempty"`
	BannerURL  string   `json:"banner_url,omitempty"`
	//new manager
	Desc  string `json:"desc,omitempty"`
	Stime string `json:"stime,omitempty"`
	Etime string `json:"etime,omitempty"`
	// rank
	Pts      int64   `json:"pts,omitempty"`
	Children []*Item `json:"children,omitempty"`
	Like     int     `json:"like,omitempty"`
	// region
	Rid  int `json:"rid,omitempty"`
	Reid int `json:"reid,omitempty"`
}

type Ext struct {
	LiveCnt int `json:"live_count,omitempty"`
}

// IsRandom check item whether or not random.
func (i *Item) IsRandom() bool {
	return i.Random == 1
}

// FromArc from recommend arc.
func (i *Item) FromArc(r *recommend.Arc) {
	i.Title = r.Title
	i.Cover = r.Pic
	i.Goto = model.GotoAv
	switch aid := r.Aid.(type) {
	case string:
		i.Param = aid
	case float64:
		i.Param = strconv.FormatInt(int64(aid), 10)
	}
	i.URI = model.FillURI(model.GotoAv, i.Param, nil)
	i.Danmaku = r.Danmaku
	v, ok := r.Views.(float64)
	if ok {
		i.Play = int(v)
	}
}

// FromBangumi from bangumi.
func (i *Item) FromBangumi(b *bangumi.Bangumi) {
	i.Title = b.Title
	i.Cover = b.NewEp.Cover
	i.Goto = model.GotoBangumi
	i.Param = b.SeasonId
	i.URI = model.FillURI(model.GotoBangumiWeb, b.SeasonId, nil)
	i.Index = b.NewEp.Index
	i.TotalCount = b.TotalCount
	i.MTime = b.NewEp.UpTime
	i.Status = int8(b.SeasonStatus)
	// i.CoverMark = b.Badge
	i.Play, _ = strconv.Atoi(b.PlayCount)
	i.Fav, _ = strconv.Atoi(b.Favorites)
	if b.Finish == "1" {
		i.Finish = 1
	}
}

// FromLive from live.
func (i *Item) FromLive(r *live.Room) {
	i.Title = r.Title
	i.Cover = r.Cover.Src
	i.Name = r.Owner.Name
	i.Face = r.Owner.Face
	i.Goto = model.GotoLive
	i.Param = strconv.FormatInt(r.ID, 10)
	i.URI = model.FillURI(model.GotoLive, strconv.FormatInt(r.ID, 10), nil)
	i.Online = r.Online
	i.Area = r.Area
	i.AreaID = r.AreaID
}

// FromArchivePB from archive archive.
func (i *Item) FromArchivePB(a *api.Arc) {
	i.Title = a.Title
	i.Cover = a.Pic
	i.Goto = model.GotoAv
	i.Param = strconv.FormatInt(a.Aid, 10)
	i.URI = model.FillURI(model.GotoAv, i.Param, nil)
	i.Danmaku = int(a.Stat.Danmaku)
	i.Play = int(a.Stat.View)
	i.Like = int(a.Stat.Like)
	if a.Access > 0 {
		i.Play = 0
	}
}

// FromArchivePBBangumi from archive archive.
func (i *Item) FromArchivePBBangumi(a *api.Arc, season *seasongrpc.CardInfoProto, bangumiType int) {
	var (
		_bangumiSeasonID  = 1
		_bangumiEpisodeID = 2
	)
	i.Title = a.Title
	i.Cover = a.Pic
	i.Goto = model.GotoBangumi
	i.Param = strconv.Itoa(int(season.SeasonId))
	switch bangumiType {
	case _bangumiSeasonID:
		i.URI = model.FillURI(model.GotoBangumi, i.Param, nil)
	case _bangumiEpisodeID:
		if season.NewEp != nil && season.NewEp.Id > 0 {
			epid := strconv.Itoa(int(season.NewEp.Id))
			i.URI = model.FillURIBangumi(model.GotoBangumi, i.Param, epid, int(season.SeasonType))
		} else {
			i.URI = model.FillURI(model.GotoBangumi, i.Param, nil)
		}
	}
	i.Danmaku = int(a.Stat.Danmaku)
	i.Play = int(a.Stat.View)
	i.Like = int(a.Stat.Like)
	i.Fav = int(a.Stat.Fav)
	if a.Access > 0 {
		i.Play = 0
	}
}

// FromArchiveRank from archive archive.
func (i *Item) FromArchiveRank(a *api.Arc, scores map[int64]int64) {
	i.Title = a.Title
	i.Cover = a.Pic
	i.Goto = model.GotoAv
	i.Param = strconv.FormatInt(a.Aid, 10)
	i.URI = model.FillURI(model.GotoAv, i.Param, nil)
	i.Danmaku = int(a.Stat.Danmaku)
	i.Play = int(a.Stat.View)
	i.Title = a.Title
	i.Name = a.Author.Name
	i.Like = int(a.Stat.Like)
	if score, ok := scores[a.Aid]; ok {
		i.Pts = score
	}
	if a.Access > 0 {
		i.Play = 0
	}
}

// FromActivity
func (i *Item) FromActivity(a *activity.Activity, now time.Time) {
	stime, err := time.ParseInLocation(_activityForm, a.Stime, time.Local)
	if err != nil {
		return
	}
	etime, err := time.ParseInLocation(_activityForm, a.Etime, time.Local)
	if err != nil {
		return
	}
	if now.After(etime) {
		i.Status = 1
	} else if now.Before(stime) {
		i.Status = 2
	}
	i.Title = a.Name
	i.Cover = a.H5Cover
	i.Goto = model.GotoWeb
	i.Param = a.H5URL
	i.URI = model.FillURI(model.GotoWeb, i.Param, nil)
	i.Desc = a.Desc
	i.Stime = a.Stime
	i.Etime = a.Etime
}

// FromTopic
func (i *Item) FromTopic(a *activity.Activity) {
	i.Title = a.Name
	i.Cover = a.H5Cover
	i.Goto = model.GotoWeb
	i.Param = a.H5URL
	i.URI = model.FillURI(model.GotoWeb, i.Param, nil)
	i.Desc = a.Desc
}

func (h *Head) FillBuildURI(plat int8, build int) {
	switch h.Goto {
	case model.GotoDaily:
		if (plat == model.PlatIPhone && build > 6670) || (plat == model.PlatAndroid && build > 5250000) {
			h.URI = "bilibili://pegasus/list/daily/" + h.Param
		}
	}
}
