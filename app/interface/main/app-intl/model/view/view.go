package view

import (
	"encoding/json"
	"strconv"

	"go-common/app/interface/main/app-intl/model"
	"go-common/app/interface/main/app-intl/model/bangumi"
	"go-common/app/interface/main/app-intl/model/manager"
	"go-common/app/interface/main/app-intl/model/tag"
	dm2 "go-common/app/interface/main/dm2/model"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	xtime "go-common/library/time"
)

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
	// player_icon
	PlayerIcon *PlayerIcon `json:"player_icon,omitempty"`
	// vip_active
	VIPActive string `json:"vip_active,omitempty"`
	// cm config
	CMConfig *CMConfig `json:"cm_config,omitempty"`
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
	Vip struct {
		Type          int    `json:"vipType"`
		DueDate       int64  `json:"vipDueDate"`
		DueRemark     string `json:"dueRemark"`
		AccessStatus  int    `json:"accessStatus"`
		VipStatus     int    `json:"vipStatus"`
		VipStatusWarn string `json:"vipStatusWarn"`
	} `json:"vip"`
	Assists  []int64 `json:"assists"`
	Fans     int     `json:"fans"`
	Archives int     `json:"archives"`
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
	TrackID   string          `json:"trackid"`
	Oid       int64           `json:"id"`
	Source    string          `json:"source"`
	AvFeature json.RawMessage `json:"av_feature"`
	Goto      string          `json:"goto"`
}

// FromAv func
func (r *Relate) FromAv(a *api.Arc, from, trackid string, ap *archive.PlayerInfo) {
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
}

// FromOperate func
func (r *Relate) FromOperate(i *manager.Relate, a *api.Arc, from string) {
	switch i.Goto {
	case model.GotoAv:
		r.FromAv(a, from, "", nil)
	}
	if r.Title == "" {
		r.Title = i.Title
	}
	if r.RcmdReason == "" {
		r.RcmdReason = i.RecReason
	}
}
