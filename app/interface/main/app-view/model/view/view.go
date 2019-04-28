package view

import (
	"encoding/json"
	"strconv"

	"go-common/app/interface/main/app-view/model"
	"go-common/app/interface/main/app-view/model/act"
	"go-common/app/interface/main/app-view/model/ad"
	"go-common/app/interface/main/app-view/model/bangumi"
	"go-common/app/interface/main/app-view/model/creative"
	"go-common/app/interface/main/app-view/model/elec"
	"go-common/app/interface/main/app-view/model/game"
	"go-common/app/interface/main/app-view/model/live"
	"go-common/app/interface/main/app-view/model/manager"
	"go-common/app/interface/main/app-view/model/special"
	"go-common/app/interface/main/app-view/model/tag"
	dm2 "go-common/app/interface/main/dm2/model"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	resmdl "go-common/app/service/main/resource/model"
	"go-common/app/service/openplatform/pgc-season/api/grpc/season/v1"
	xtime "go-common/library/time"
)

// BnjView 2018
type BnjView struct {
	BeginTime int64 `json:"begin_time"`
	*api.Arc
	// owner_ext
	OwnerExt OwnerExt `json:"owner_ext"`
	// now user
	ReqUser  *ReqUser `json:"req_user,omitempty"`
	LiveRoom struct {
		ID       int64  `json:"id"`
		DmServer string `json:"dm_server"`
		DmPort   int64  `json:"dm_port"`
		Title    string `json:"title"`
		Uname    string `json:"uname"`
		Cover    string `json:"cover"`
	} `json:"live_room"`
	Pages   []*Page    `json:"pages,omitempty"`
	Elec    *elec.Info `json:"elec,omitempty"`
	Stat    *api.Stat  `json:"stat,omitempty"`
	Lottery struct {
		ActID   int64    `json:"act_id"`
		Times   int64    `json:"times"`
		Rule    string   `json:"rule"`
		List    []string `json:"list"`
		Winners []string `json:"winners"`
	} `json:"lottery"`
	Relates     []*act.Relate    `json:"relates"`
	PastReviews []*BnjPastReview `json:"past_review"`
}

// BnjPastReview struct
type BnjPastReview struct {
	AID   int64  `json:"aid"`
	Img   string `json:"img"`
	Title string `json:"title"`
}

// vip active subID.
const (
	VIPActiveView   = 1
	VIPActivePGC    = 2
	VIPActiveCenter = 3
)

// View struct
type View struct {
	*ViewStatic
	// owner_ext
	OwnerExt OwnerExt `json:"owner_ext"`
	// now user
	ReqUser *ReqUser `json:"req_user,omitempty"`
	// tag info
	Tag []*tag.Tag `json:"tag,omitempty"`
	// movie
	Movie *bangumi.Movie `json:"movie,omitempty"`
	// bangumi
	Season *bangumi.Season `json:"season,omitempty"`
	// bp
	Bp json.RawMessage `json:"bp,omitempty"`
	// elec
	Elec *elec.Info `json:"elec,omitempty"`
	// history
	History *History `json:"history,omitempty"`
	// audio
	Audio *Audio `json:"audio,omitempty"`
	// contribute data
	Contributions []*Contribution `json:"contributions,omitempty"`
	// relate data
	Relates     []*Relate `json:"relates,omitempty"`
	ReturnCode  string    `json:"-"`
	UserFeature string    `json:"-"`
	IsRec       int8      `json:"-"`
	// dislike reason
	Dislikes []*Dislike `json:"dislike_reasons,omitempty"`
	// dm
	DMSeg int `json:"dm_seg,omitempty"`
	// paster
	Paster *resmdl.Paster `json:"paster,omitempty"`
	// player_icon
	PlayerIcon *PlayerIcon `json:"player_icon,omitempty"`
	// vip_active
	VIPActive string `json:"vip_active,omitempty"`
	// cm
	CMs []*CM `json:"cms,omitempty"`
	// cm config
	CMConfig *CMConfig `json:"cm_config,omitempty"`
	// asset
	Asset       *Asset          `json:"asset,omitempty"`
	ActivityURL string          `json:"activity_url,omitempty"`
	Bgm         []*creative.Bgm `json:"bgm,omitempty"`
	Staff       []*Staff        `json:"staff,omitempty"`
}

// Staff from cooperation
type Staff struct {
	Mid            int64  `json:"mid,omitempty"`
	Title          string `json:"title,omitempty"`
	Face           string `json:"face,omitempty"`
	Name           string `json:"name,omitempty"`
	OfficialVerify struct {
		Type int    `json:"type"`
		Desc string `json:"desc"`
	} `json:"official_verify"`
	Vip struct {
		Type          int    `json:"vipType"`
		DueDate       int64  `json:"vipDueDate"`
		DueRemark     string `json:"dueRemark"`
		AccessStatus  int    `json:"accessStatus"`
		VipStatus     int    `json:"vipStatus"`
		VipStatusWarn string `json:"vipStatusWarn"`
	} `json:"vip"`
	Attention int `json:"attention"`
}

// ViewStatic struct
type ViewStatic struct {
	*archive.Archive3
	Pages []*Page `json:"pages,omitempty"`
}

// ReqUser struct
type ReqUser struct {
	Attention int  `json:"attention"`
	Favorite  int8 `json:"favorite"`
	Like      int8 `json:"like"`
	Dislike   int8 `json:"dislike"`
	Coin      int8 `json:"coin"`
}

// Page struct
type Page struct {
	*archive.Page3
	Metas  []*Meta          `json:"metas"`
	DMLink string           `json:"dmlink"`
	Audio  *Audio           `json:"audio,omitempty"`
	DM     *dm2.SubjectInfo `json:"dm,omitempty"`
}

// Meta struct
type Meta struct {
	Quality int    `json:"quality"`
	Format  string `json:"format"`
	Size    int64  `json:"size"`
}

// History struct
type History struct {
	Cid      int64 `json:"cid"`
	Progress int64 `json:"progress"`
}

// CM struct
type CM struct {
	RequestID string     `json:"request_id,omitempty"`
	RscID     int64      `json:"rsc_id,omitempty"`
	SrcID     int64      `json:"src_id,omitempty"`
	IsAdLoc   bool       `json:"is_ad_loc,omitempty"`
	IsAd      bool       `json:"is_ad,omitempty"`
	CmMark    int        `json:"cm_mark,omitempty"`
	ClientIP  string     `json:"client_ip,omitempty"`
	Index     int        `json:"index,omitempty"`
	AdInfo    *ad.AdInfo `json:"ad_info,omitempty"`
}

// CMConfig struct
type CMConfig struct {
	AdsControl  json.RawMessage `json:"ads_control,omitempty"`
	MonitorInfo json.RawMessage `json:"monitor_info,omitempty"`
}

// Dislike struct
type Dislike struct {
	ID   int    `json:"reason_id"`
	Name string `json:"reason_name"`
}

// OwnerExt struct
type OwnerExt struct {
	OfficialVerify struct {
		Type int    `json:"type"`
		Desc string `json:"desc"`
	} `json:"official_verify,omitempty"`
	Live *live.Live `json:"live,omitempty"`
	Vip  struct {
		Type          int    `json:"vipType"`
		DueDate       int64  `json:"vipDueDate"`
		DueRemark     string `json:"dueRemark"`
		AccessStatus  int    `json:"accessStatus"`
		VipStatus     int    `json:"vipStatus"`
		VipStatusWarn string `json:"vipStatusWarn"`
	} `json:"vip"`
	Assists  []int64 `json:"assists"`
	Fans     int     `json:"fans"`
	Archives int     `json:"-"`
}

// Relate struct
type Relate struct {
	Aid         int64       `json:"aid,omitempty"`
	Pic         string      `json:"pic,omitempty"`
	Title       string      `json:"title,omitempty"`
	Author      *api.Author `json:"owner,omitempty"`
	Stat        api.Stat    `json:"stat,omitempty"`
	Duration    int64       `json:"duration,omitempty"`
	Goto        string      `json:"goto,omitempty"`
	Param       string      `json:"param,omitempty"`
	URI         string      `json:"uri,omitempty"`
	Rating      float64     `json:"rating,omitempty"`
	Reserve     string      `json:"reserve,omitempty"`
	From        string      `json:"from,omitempty"`
	Desc        string      `json:"desc,omitempty"`
	RcmdReason  string      `json:"rcmd_reason,omitempty"`
	Badge       string      `json:"badge,omitempty"`
	Cid         int64       `json:"cid,omitempty"`
	SeasonType  int32       `json:"season_type,omitempty"`
	RatingCount int32       `json:"rating_count,omitempty"`
	// cm ad
	AdIndex      int             `json:"ad_index,omitempty"`
	CmMark       int             `json:"cm_mark,omitempty"`
	SrcID        int64           `json:"src_id,omitempty"`
	RequestID    string          `json:"request_id,omitempty"`
	CreativeID   int64           `json:"creative_id,omitempty"`
	CreativeType int64           `json:"creative_type,omitempty"`
	Type         int             `json:"type,omitempty"`
	Cover        string          `json:"cover,omitempty"`
	ButtonTitle  string          `json:"button_title,omitempty"`
	View         int             `json:"view,omitempty"`
	Danmaku      int             `json:"danmaku,omitempty"`
	IsAd         bool            `json:"is_ad,omitempty"`
	IsAdLoc      bool            `json:"is_ad_loc,omitempty"`
	AdCb         string          `json:"ad_cb,omitempty"`
	ShowURL      string          `json:"show_url,omitempty"`
	ClickURL     string          `json:"click_url,omitempty"`
	ClientIP     string          `json:"client_ip,omitempty"`
	Extra        json.RawMessage `json:"extra,omitempty"`
	Button       *Button         `json:"button,omitempty"`
	CardIndex    int             `json:"card_index,omitempty"`
	Source       string          `json:"-"`
	AvFeature    json.RawMessage `json:"-"`
	TrackID      string          `json:"trackid"`
}

// Button struct
type Button struct {
	Title string `json:"title,omitempty"`
	URI   string `json:"uri,omitempty"`
}

// Contribution struct
type Contribution struct {
	Aid    int64      `json:"aid,omitempty"`
	Pic    string     `json:"pic,omitempty"`
	Title  string     `json:"title,omitempty"`
	Author api.Author `json:"owner,omitempty"`
	Stat   api.Stat   `json:"stat,omitempty"`
	CTime  xtime.Time `json:"ctime,omitempty"`
}

// Audio struct
type Audio struct {
	Title    string `json:"title"`
	Cover    string `json:"cover_url"`
	SongID   int    `json:"song_id"`
	Play     int    `json:"play_count"`
	Reply    int    `json:"reply_count"`
	UpperID  int    `json:"upper_id"`
	Entrance string `json:"entrance"`
	SongAttr int    `json:"song_attr"`
}

// PlayerIcon struct
type PlayerIcon struct {
	URL1  string     `json:"url1,omitempty"`
	Hash1 string     `json:"hash1,omitempty"`
	URL2  string     `json:"url2,omitempty"`
	Hash2 string     `json:"hash2,omitempty"`
	CTime xtime.Time `json:"ctime,omitempty"`
}

// VipPlayURL playurl token struct.
type VipPlayURL struct {
	From  string `json:"from"`
	Ts    int64  `json:"ts"`
	Aid   int64  `json:"aid"`
	Cid   int64  `json:"cid"`
	Mid   int64  `json:"mid"`
	VIP   int    `json:"vip"`
	SVIP  int    `json:"svip"`
	Owner int    `json:"owner"`
	Fcs   string `json:"fcs"`
}

// NewRelateRec struct
type NewRelateRec struct {
	TrackID    string          `json:"trackid"`
	Oid        int64           `json:"id"`
	Source     string          `json:"source"`
	AvFeature  json.RawMessage `json:"av_feature"`
	Goto       string          `json:"goto"`
	Title      string          `json:"title"`
	IsDalao    int8            `json:"is_dalao"`
	RcmdReason struct {
		Content string `json:"content"`
	} `json:"rcmd_reason"`
}

type Asset struct {
	Paid  int8  `json:"paid"`
	Price int64 `json:"price"`
	Msg   struct {
		Desc1 string `json:"desc1"`
		Desc2 string `json:"desc2"`
	} `json:"msg"`
}

// FromAv func
func (r *Relate) FromAv(a *api.Arc, from, trackid string, ap *archive.PlayerInfo, cooperation bool) {
	r.Aid = a.Aid
	r.Title = a.Title
	r.Pic = a.Pic
	r.Author = &a.Author
	r.Stat = a.Stat
	r.Duration = a.Duration
	r.Cid = a.FirstCid
	r.Goto = model.GotoAv
	r.Param = strconv.FormatInt(a.Aid, 10)
	r.URI = model.FillURI(r.Goto, r.Param, model.AvHandler(a, trackid, ap))
	r.From = from
	if a.AttrVal(archive.AttrBitIsCooperation) == archive.AttrYes && r.Author != nil && r.Author.Name != "" && cooperation {
		r.Author.Name = r.Author.Name + " 等联合创作"
	}
}

// FromGame func
func (r *Relate) FromGame(i *game.Info, from string) {
	if i.GameLink == "" {
		return
	}
	r.Title = "相关游戏：" + i.GameName
	r.Pic = i.GameIcon
	r.Rating = i.Grade
	if i.GameStatus == 1 || i.GameStatus == 2 {
		var reserve string
		if i.BookNum < 10000 {
			reserve = strconv.FormatInt(i.BookNum, 10) + "人预约"
		} else {
			reserve = strconv.FormatFloat(float64(i.BookNum)/10000, 'f', 1, 64) + "万人预约"
		}
		r.Reserve = reserve
	}
	r.Goto = model.GotoGame
	r.URI = model.FillURI(r.Goto, i.GameLink, nil)
	r.Param = strconv.FormatInt(i.GameBaseID, 10)
	r.Button = &Button{Title: "进入", URI: r.URI}
	r.From = from
}

// FromSpecial func
func (r *Relate) FromSpecial(sp *special.Card, from string) {
	r.Title = sp.Title
	r.Pic = sp.Cover
	r.Goto = model.GotoSpecial
	// TODO FUCK game
	r.URI = model.FillURI(model.OperateType[sp.ReType], sp.ReValue, nil)
	r.Desc = sp.Desc
	r.Param = strconv.FormatInt(sp.ID, 10)
	r.RcmdReason = sp.Badge
	r.From = from
}

// FromOperate func
func (r *Relate) FromOperate(i *NewRelateRec, a *api.Arc, info *game.Info, sp *special.Card, from string, cooperation bool) {
	switch i.Goto {
	case model.GotoAv:
		r.FromAv(a, from, "", nil, cooperation)
	case model.GotoGame:
		r.FromGame(info, from)
	case model.GotoSpecial:
		r.FromSpecial(sp, from)
	}
	if r.Title == "" {
		r.Title = i.Title
	}
	if i.RcmdReason.Content != "" {
		r.RcmdReason = i.RcmdReason.Content
	}
}

// FromOperate func
func (r *Relate) FromOperateOld(i *manager.Relate, a *api.Arc, info *game.Info, sp *special.Card, from string, cooperation bool) {
	switch i.Goto {
	case model.GotoAv:
		r.FromAv(a, from, "", nil, cooperation)
	case model.GotoGame:
		r.FromGame(info, from)
	case model.GotoSpecial:
		r.FromSpecial(sp, from)
	}
	if r.Title == "" {
		r.Title = i.Title
	}
	if r.RcmdReason == "" {
		r.RcmdReason = i.RecReason
	}
}

// FromCM func
func (r *Relate) FromCM(ad *ad.AdInfo) {
	r.AdIndex = ad.Index
	r.CmMark = ad.CmMark
	r.SrcID = ad.Source
	r.RequestID = ad.RequestID
	r.CreativeID = ad.CreativeID
	r.CreativeType = ad.CreativeType
	r.Type = ad.CardType
	r.URI = ad.URI
	r.Param = ad.Param
	r.Goto = model.GotoCm
	r.View = ad.View
	r.Danmaku = ad.Danmaku
	r.IsAd = ad.IsAd
	r.IsAdLoc = ad.IsAdLoc
	r.AdCb = ad.AdCb
	r.ClientIP = ad.ClientIP
	r.Extra = ad.Extra
	r.CardIndex = ad.CardIndex
	if ad.CreativeContent != nil {
		r.Aid = ad.CreativeContent.VideoID
		r.Cover = ad.CreativeContent.ImageURL
		r.Title = ad.CreativeContent.Title
		r.ButtonTitle = ad.CreativeContent.ButtonTitle
		r.Desc = ad.CreativeContent.Desc
		r.ShowURL = ad.CreativeContent.ShowURL
		r.ClickURL = ad.CreativeContent.ClickURL
	}
}

// FromCM func
func (c *CM) FromCM(ad *ad.AdInfo) {
	c.RequestID = ad.RequestID
	c.RscID = ad.Resource
	c.SrcID = ad.Source
	c.IsAd = ad.IsAd
	c.IsAdLoc = ad.IsAdLoc
	c.Index = ad.Index
	c.CmMark = ad.CmMark
	c.ClientIP = ad.ClientIP
	c.AdInfo = ad
}

// FromBangumi func
func (r *Relate) FromBangumi(ban *v1.CardInfoProto) {
	r.Title = ban.Title
	r.Pic = ban.NewEp.Cover
	r.Stat = api.Stat{
		Danmaku: int32(ban.Stat.Danmaku),
		View:    int32(ban.Stat.View),
		Fav:     int32(ban.Stat.Follow),
	}
	r.Goto = model.GotoBangumi
	r.Param = strconv.FormatInt(int64(ban.SeasonId), 10)
	r.URI = model.FillURI(r.Goto, r.Param, nil)
	r.SeasonType = ban.SeasonType
	r.Badge = ban.SeasonTypeName
	r.Desc = ban.NewEp.IndexShow
	if ban.Rating != nil {
		r.Rating = float64(ban.Rating.Score)
		r.RatingCount = ban.Rating.Count
	}
}

// TripleParam struct
type TripleParam struct {
	MobiApp string `form:"mobi_app"`
	Build   string `form:"build"`
	AID     int64  `form:"aid"`
	Ak      string `form:"access_key"`
	From    string `form:"from"`
}

// TripleRes struct
type TripleRes struct {
	Like      bool  `json:"like"`
	Coin      bool  `json:"coin"`
	Fav       bool  `json:"fav"`
	Prompt    bool  `json:"prompt"`
	Multiply  int64 `json:"multiply"`
	UpID      int64 `json:"-"`
	Anticheat bool  `json:"-"`
}

// Videoshot videoshot
type Videoshot struct {
	*archive.Videoshot
	Points []*creative.Points `json:"points,omitempty"`
}
