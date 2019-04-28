package model

import (
	"fmt"
	"net/url"

	"go-common/library/time"
)

const (
	_epFree = 2
)

// SeasonCMS defines the elements could be changed from TV CMS side
type SeasonCMS struct {
	SeasonID    int64
	Cover       string
	Desc        string
	Title       string
	UpInfo      string
	Category    int8   // - cn, jp, movie, tv, documentary
	Area        string // - cn, jp, others
	Playtime    time.Time
	Role        string
	Staff       string
	NewestOrder int   // the newest passed ep's order
	NewestEPID  int64 // the newest passed ep's ID
	NewestNb    int   // the newest ep's number ( after keyword filter )
	TotalNum    int
	Style       string
	OriginName  string
	Alias       string
	PayStatus   int
}

// NeedVip returns whether the season need vip to watch
func (s *SeasonCMS) NeedVip() bool {
	return s.PayStatus == 1
}

// IdxSn is the structure of season in the index page
func (s *SeasonCMS) IdxSn() (idx *IdxSeason) {
	return &IdxSeason{
		SeasonID: s.SeasonID,
		Title:    s.Title,
		Cover:    s.Cover,
		Upinfo:   s.UpInfo,
	}
}

// EpCMS defines the elements could be changed from TV CMS side
type EpCMS struct {
	EPID      int64  `json:"epid"`
	Cover     string `json:"cover"`
	Title     string `json:"title"`
	Subtitle  string `json:"subtitle"`
	PayStatus int    `json:"pay_status"`
}

// IsFree def.
func (v *EpCMS) IsFree() bool {
	return v.PayStatus == _epFree
}

// EpDecor is used to decorate ep
type EpDecor struct {
	*EpCMS
	Watermark bool `json:"watermark"`
}

// ArcCMS reprensents the archive data structure in MC
type ArcCMS struct {
	Title   string
	AID     int64
	Content string
	Cover   string
	TypeID  int
	Pubtime time.Time
	Videos  int
	Valid   int
	Deleted int
	Result  int
}

// NotDeleted def.
func (s *ArcCMS) NotDeleted() bool {
	return s.Deleted == 0
}

// CanPlay returns whether the arc can play or not
func (s *ArcCMS) CanPlay() bool {
	return s.Valid == 1 && s.Result == 1 && s.Deleted == 0
}

// VideoCMS def.
type VideoCMS struct {
	// Media Info
	CID        int64
	Title      string
	AID        int64
	IndexOrder int
	// Auth Info
	Valid   int
	Deleted int
	Result  int
}

// CanPlay returns whether the arc can play or not
func (s *VideoCMS) CanPlay() bool {
	return s.Valid == 1 && s.Result == 1 && s.Deleted == 0
}

// NotDeleted def.
func (s *VideoCMS) NotDeleted() bool {
	return s.Deleted == 0
}

// Auditing returns whether the video is begin audited by the license owner
func (s *VideoCMS) Auditing() bool {
	return s.Result == 0 && s.Deleted == 0
}

// MediaParam def.
type MediaParam struct {
	SeasonID  int64  `form:"season_id"`
	EpID      int64  `form:"ep_id"`
	TrackPath string `form:"track_path" validate:"required"`
	AccessKey string `form:"access_key"`
	MobiAPP   string `form:"mobi_app" validate:"required"`
	Platform  string `form:"platform"`
	Build     int64  `form:"build"`
}

// GenerateUrl generates url.Values from tv media param struct
func (v *MediaParam) GenerateUrl() (params url.Values) {
	params = url.Values{}
	params.Set("build", fmt.Sprintf("%d", v.Build))
	params.Set("mobi_app", v.MobiAPP)
	params.Set("platform", v.Platform)
	params.Set("access_key", v.AccessKey)
	params.Set("track_path", v.TrackPath)
	params.Set("season_id", fmt.Sprintf("%d", v.SeasonID))
	return
}

// MediaResp is the structure of PGC display api response
type MediaResp struct {
	Response
	Result *SeasonDetail `json:"result"`
}

// SeasonDetail def
type SeasonDetail struct {
	Episodes     []*Episode  `json:"episodes"`
	IsNewDanmaku int         `json:"is_new_danmaku"`
	NewestEP     *NewestEP   `json:"newest_ep"`
	Stat         *Stat       `json:"stat"`
	UserStatus   *UserStatus `json:"user_status"`
	Sponsor      *Sponsor    `json:"sponsor"`
	SeriesID     int         `json:"series_id"`
	SnDetailCore
}

// CmsInterv def.
func (v *SnDetailCore) CmsInterv(snCMS *SeasonCMS) {
	if snCMS.Title != "" {
		v.Title = snCMS.Title
	}
	if snCMS.Cover != "" {
		v.Cover = snCMS.Cover
	}
	if snCMS.Desc != "" {
		v.Evaluate = snCMS.Desc
	}
}

// UserStatus def
type UserStatus struct {
	Follow        int            `json:"follow"`
	IsVip         int            `json:"is_vip"`
	Pay           int            `json:"pay"`
	PayPackPaid   int            `json:"pay_pack_paid"`
	Sponsor       int            `json:"sponsor"`
	WatchProgress *WatchProgress `json:"watch_progress"`
}

// WatchProgress def.
type WatchProgress struct {
	LastEpID    int    `json:"last_ep_id"`
	LastEPIndex string `json:"last_ep_index"`
	LastTime    int64  `json:"last_time"`
}

// Stat def
type Stat struct {
	Danmakus  int `json:"danmakus"`
	Favorites int `json:"favorites"`
	Views     int `json:"views"`
}

// List def
type List struct {
	Face  string `json:"face"`
	UID   int    `json:"uid"`
	Uname string `json:"uname"`
}

// Sponsor def
type Sponsor struct {
	List          []*List        `json:"list"`
	PointActivity *PointActivity `json:"point_activity"`
	TotalBpCount  int            `json:"total_bp_count"`
	WeekBpCount   int            `json:"week_bp_count"`
}

// PointActivity def
type PointActivity struct {
	Content string `json:"content"`
	Link    string `json:"link"`
	Tip     string `json:"tip"`
}

// Season def
type Season struct {
	SeasonV2
	Title string `json:"title"`
}

// SeasonV2 def
type SeasonV2 struct {
	IsNew       int    `json:"is_new"`
	SeasonID    int    `json:"season_id"`
	SeasonTitle string `json:"season_title"`
}

// Rights def
type Rights struct {
	AllowBp       int    `json:"allow_bp"`
	AllowDownload int    `json:"allow_download"`
	AllowReview   int    `json:"allow_review"`
	AreaLimit     int    `json:"area_limit"`
	BanAreaShow   int    `json:"ban_area_show"`
	Copyright     string `json:"copyright"`
	IsPreview     int    `json:"is_preview"`
}

// Rating def
type Rating struct {
	Count int     `json:"count"`
	Score float64 `json:"score"`
}

// Publish def
type Publish struct {
	IsFinish    int    `json:"is_finish"`
	IsStarted   int    `json:"is_started"`
	PubTime     string `json:"pub_time"`
	PubTimeShow string `json:"pub_time_show"`
	Weekday     int    `json:"weekday"`
}

// Paster def
type Paster struct {
	AID       int    `json:"aid"`
	CID       int    `json:"cid"`
	AllowJump int    `json:"allow_jump"`
	Duration  int    `json:"duration"`
	Type      int    `json:"type"`
	URL       string `json:"url"`
}

// NewestEP def
type NewestEP struct {
	Desc  string `json:"desc"`
	ID    int    `json:"id"`
	Index string `json:"index"`
	IsNew int    `jsontt:"is_new"`
}

// Episode def
type Episode struct {
	AID           int    `json:"aid"`
	CID           int    `json:"cid"`
	Cover         string `json:"cover"`
	EPID          int64  `json:"ep_id"`
	EpisodeStatus int    `json:"episode_status"`
	From          string `json:"from"`
	Index         string `json:"index"`
	IndexTitle    string `json:"index_title"`
	MID           int    `json:"mid"`
	Page          int    `json:"page"`
	ShareURL      string `json:"share_url"`
	VID           string `json:"vid"`
	WaterMark     bool   `json:"hidemark"` // true means in the whitelist
}

// CmsInterv def.
func (v *Episode) CmsInterv(epCMS *EpCMS) {
	if epCMS.Cover != "" {
		v.Cover = epCMS.Cover
	}
	if epCMS.Title != "" {
		v.IndexTitle = epCMS.Title
	}
}

// ParamStyle .
type ParamStyle struct {
	Name    string `json:"name"`
	StyleID int    `json:"style_id"`
}
