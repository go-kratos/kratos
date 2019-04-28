package region

import (
	"strconv"
	"time"

	"go-common/app/interface/main/app-show/model"
	"go-common/app/interface/main/app-show/model/activity"
	"go-common/app/interface/main/app-show/model/banner"
	"go-common/app/interface/main/app-show/model/recommend"
	"go-common/app/interface/main/app-show/model/tag"
	accv1 "go-common/app/service/main/account/api"
	"go-common/app/service/main/archive/api"
	seasongrpc "go-common/app/service/openplatform/pgc-season/api/grpc/season/v1"
	xtime "go-common/library/time"
)

const (
	_activityForm = "2006-01-02 15:04:05"
)

type Region struct {
	ID        int64     `json:"-"`
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
	Rtype     int8      `json:"type"`
	Entrance  int8      `json:"-"`
	IsBangumi int8      `json:"is_bangumi,omitempty"`
	Children  []*Region `json:"children,omitempty"`
	Config    []*Config `json:"config,omitempty"`
}

type Limit struct {
	ID        int64  `json:"-"`
	Rid       int64  `json:"-"`
	Build     int    `json:"-"`
	Condition string `json:"-"`
}

type Config struct {
	ID         int64  `json:"-"`
	Rid        int64  `json:"-"`
	ScenesID   int    `json:"-"`
	ScenesName string `json:"scenes_name,omitempty"`
	ScenesType string `json:"scenes_type,omitempty"`
}

type Show struct {
	Banner    map[string][]*banner.Banner `json:"banner,omitempty"`
	Card      []*Head                     `json:"card,omitempty"`
	Tag       *tag.Tag                    `json:"tag,omitempty"`
	TopTag    []*SimilarTag               `json:"top_tag,omitempty"`
	NewTag    *NewTag                     `json:"new_tag,omitempty"`
	Cbottom   xtime.Time                  `json:"cbottom,omitempty"`
	Ctop      xtime.Time                  `json:"ctop,omitempty"`
	Recommend []*ShowItem                 `json:"recommend,omitempty"`
	New       []*ShowItem                 `json:"new"`
	Dynamic   []*ShowItem                 `json:"dynamic,omitempty"`
}

type Head struct {
	CardID    int         `json:"card_id,omitempty"`
	Title     string      `json:"title,omitempty"`
	Cover     string      `json:"cover,omitempty"`
	Type      string      `json:"type,omitempty"`
	Date      int64       `json:"date,omitempty"`
	Plat      int8        `json:"-"`
	Build     int         `json:"-"`
	Condition string      `json:"-"`
	URI       string      `json:"uri,omitempty"`
	Goto      string      `json:"goto,omitempty"`
	Param     string      `json:"param,omitempty"`
	Body      []*ShowItem `json:"body,omitempty"`
}

type ShowItem struct {
	Title    string `json:"title"`
	Cover    string `json:"cover"`
	URI      string `json:"uri"`
	NewURI   string `json:"-"`
	Param    string `json:"param"`
	FirstCid int64  `json:"cid,omitempty"`
	Goto     string `json:"goto"`
	// up
	Mid            int64           `json:"mid,omitempty"`
	Name           string          `json:"name,omitempty"`
	Face           string          `json:"face,omitempty"`
	Follower       int             `json:"follower,omitempty"`
	Attribute      int             `json:"attribute,omitempty"`
	OfficialVerify *OfficialVerify `json:"official_verify,omitempty"`
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
	Desc        string `json:"desc,omitempty"`
	Stime       string `json:"stime,omitempty"`
	Etime       string `json:"etime,omitempty"`
	Like        int    `json:"like,omitempty"`
	RedirectURL string `json:"-"`
	UGCPay      int32  `json:"ugc_pay,omitempty"`
	Cooperation string `json:"cooperation,omitempty"`
}

type OfficialVerify struct {
	Type int    `json:"type"`
	Desc string `json:"desc"`
}

type SimilarTag struct {
	TagId   int    `json:"tid"`
	TagName string `json:"tname"`
	Rid     int    `json:"rid,omitempty"`
	Rname   string `json:"rname,omitempty"`
	Reid    int    `json:"reid,omitempty"`
	Rename  string `json:"rename,omitempty"`
}

type NewTag struct {
	Position int           `json:"pos"`
	Tag      []*SimilarTag `json:"tag"`
}

func (c *Config) ConfigChange() {
	switch c.ScenesID {
	case 0:
		c.ScenesName = "region"
		c.ScenesType = "bottom"
	case 1:
		c.ScenesName = "region"
		c.ScenesType = "top"
	case 2:
		c.ScenesName = "rank"
	case 3:
		c.ScenesName = "search"
	case 4:
		c.ScenesName = "tag"
	case 5:
		c.ScenesName = "attention"
	}
}

// FromArc from recommend arc.
func (i *ShowItem) FromArc(a *recommend.Arc) {
	i.fromArc(a)
	for _, as := range a.Others {
		child := &ShowItem{}
		child.fromArc(as)
		i.Children = append(i.Children, child)
	}
}

// FromArcBangumi from recommend arc bangumi.
func (i *ShowItem) FromArcBangumi(a *recommend.Arc, sids map[int64]int64) {
	aidInt := fromAid(a.Aid)
	if sid, ok := sids[aidInt]; ok && sid != 0 {
		i.fromArcBangumi(a, sid)
	} else {
		i.fromArc(a)
	}
	for _, as := range a.Others {
		child := &ShowItem{}
		aidInt = fromAid(as.Aid)
		if sid, ok := sids[aidInt]; ok && sid != 0 {
			child.fromArcBangumi(as, sid)
		} else {
			child.fromArc(as)
		}
		i.Children = append(i.Children, child)
	}
}

// fromAid
func fromAid(aidInter interface{}) (aid int64) {
	switch aidType := aidInter.(type) {
	case string:
		if aidtmp, err := strconv.ParseInt(aidType, 10, 64); err == nil && aidtmp != 0 {
			aid = aidtmp
		}
	case float64:
		aid = int64(aidType)
	}
	return
}

func (i *ShowItem) fromArc(a *recommend.Arc) {
	i.Title = a.Title
	i.Cover = a.Pic
	switch aid := a.Aid.(type) {
	case string:
		i.Param = aid
	case float64:
		i.Param = strconv.FormatInt(int64(aid), 10)
	}
	i.URI = model.FillURI(model.GotoAv, i.Param, nil)
	i.Goto = model.GotoAv
	v, ok := a.Views.(float64)
	if ok {
		i.Play = int(v)
	}
	i.Danmaku = a.Danmaku
	i.Name = a.Author
	i.Reply = int(a.Comments)
	i.Fav = int(a.Favorites)
	i.Pts = a.Pts
}

func (i *ShowItem) fromArcBangumi(a *recommend.Arc, sid int64) {
	i.Title = a.Title
	i.Cover = a.Pic
	i.Param = strconv.FormatInt(sid, 10)
	i.URI = model.FillURI(model.GotoBangumi, i.Param, nil)
	i.Goto = model.GotoBangumi
	v, ok := a.Views.(float64)
	if ok {
		i.Play = int(v)
	}
	i.Danmaku = a.Danmaku
	i.Name = a.Author
	i.Reply = int(a.Comments)
	i.Fav = int(a.Favorites)
	i.Pts = a.Pts
}

// FromArchivePB from archive archive.
func (i *ShowItem) FromArchivePB(a *api.Arc) {
	i.Title = a.Title
	i.Cover = a.Pic
	i.Param = strconv.FormatInt(a.Aid, 10)
	i.URI = model.FillURI(model.GotoAv, i.Param, model.AvHandler(a))
	i.Goto = model.GotoAv
	i.Play = int(a.Stat.View)
	i.Danmaku = int(a.Stat.Danmaku)
	i.Name = a.Author.Name
	i.Face = a.Author.Face
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
	i.UGCPay = a.Rights.UGCPay
}

// FromBangumi from archive archive.
func (i *ShowItem) FromBangumiArchivePB(a *api.Arc, season *seasongrpc.CardInfoProto, bangumiType int) {
	var (
		_bangumiSeasonID  = 1
		_bangumiEpisodeID = 2
	)
	if season == nil {
		return
	}
	i.Title = a.Title
	i.Cover = a.Pic
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

// FromArchivePBRank from archive archive.
func (i *ShowItem) FromArchivePBRank(a *api.Arc, scores map[int64]int64) {
	i.Title = a.Title
	i.Cover = a.Pic
	i.Param = strconv.FormatInt(a.Aid, 10)
	i.URI = model.FillURI(model.GotoAv, i.Param, nil)
	i.RedirectURL = a.RedirectURL
	i.Goto = model.GotoAv
	i.Play = int(a.Stat.View)
	i.Danmaku = int(a.Stat.Danmaku)
	i.Mid = a.Author.Mid
	i.Name = a.Author.Name
	i.Face = a.Author.Face
	i.Reply = int(a.Stat.Reply)
	i.Fav = int(a.Stat.Fav)
	i.PubDate = a.PubDate
	i.Rid = int(a.TypeID)
	i.Rname = a.TypeName
	i.Duration = a.Duration
	i.Like = int(a.Stat.Like)
	i.FirstCid = a.FirstCid
	if score, ok := scores[a.Aid]; ok {
		i.Pts = score
	}
	if a.Access > 0 {
		i.Play = 0
	}
	if a.Rights.IsCooperation > 0 {
		i.Cooperation = "等联合创作"
	}
}

// FromTopic
func (i *ShowItem) FromTopic(a *activity.Activity) {
	i.Title = a.Name
	i.Cover = a.H5Cover
	i.Goto = model.GotoWeb
	i.Param = a.H5URL
	i.URI = model.FillURI(model.GotoWeb, i.Param, nil)
	i.Desc = a.Desc
}

// FromActivity
func (i *ShowItem) FromActivity(a *activity.Activity, now time.Time) {
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

// FromOfficialVerify from official
func (i *OfficialVerify) FromOfficialVerify(a accv1.OfficialInfo) {
	if a.Role == 0 {
		i.Type = -1
	} else {
		if a.Role <= 2 {
			i.Type = 0
		} else {
			i.Type = 1
		}
		i.Desc = a.Title
	}
}

func (h *Head) FillBuildURI(plat int8, build int) {
	switch h.Goto {
	case model.GotoDaily:
		if (plat == model.PlatIPhone && build > 6670) || (plat == model.PlatAndroid && build > 5250000) {
			h.URI = "bilibili://pegasus/list/daily/" + h.Param
		}
	}
}
