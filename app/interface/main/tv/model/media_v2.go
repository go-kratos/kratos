package model

import "go-common/library/ecode"

// SnDetailCore is the common part of pgc media v1 and v2
type SnDetailCore struct {
	Cover        string        `json:"cover"`
	Evaluate     string        `json:"evaluate"`
	Link         string        `json:"link"`
	MediaID      int           `json:"media_id"`
	Mode         int           `json:"mode"`
	Paster       *Paster       `json:"paster"`
	Publish      *Publish      `json:"publish"`
	Rating       *Rating       `json:"rating"`
	SeasonID     int64         `json:"season_id"`
	SeasonStatus int           `json:"season_status"`
	SeasonTitle  string        `json:"season_title"`
	SeasonType   int           `json:"season_type"`
	ShareURL     string        `json:"share_url"`
	SquareCover  string        `json:"square_cover"`
	Title        string        `json:"title"`
	TotalEp      int           `json:"total_ep"`
	Rights       *Rights       `json:"rights"`
	StyleLabel   []*ParamStyle `json:"style_label"`
}

// SnDetailV2 def
type SnDetailV2 struct {
	Episodes   []*EpisodeV2  `json:"episodes"`
	NewestEP   *NewEPV2      `json:"new_ep"`
	Stat       *StatV2       `json:"stat"`
	UserStatus *UserStatusV2 `json:"user_status"`
	Seasons    []*SeasonV2   `json:"seasons"`
	Section    []*Section    `json:"section"`
	Type       int           `json:"type"`
	SnDetailCore
}

// TypeTrans def.
func (v *SnDetailV2) TypeTrans() {
	v.SeasonType = v.Type
}

// Section def.
type Section struct {
	Episodes []*EpisodeV2 `json:"episodes"`
}

// EpisodeV2 def.
type EpisodeV2 struct {
	AID        int64       `json:"aid"`
	Badge      string      `json:"badge"`
	BadgeType  int         `json:"badge_type"`
	CID        int64       `json:"cid"`
	Cover      string      `json:"cover"`
	From       string      `json:"from"`
	ID         int64       `json:"id"`
	LongTitle  string      `json:"long_title"`
	ShareURL   string      `json:"share_url"`
	Status     int         `json:"status"`
	Title      string      `json:"title"`
	VID        string      `json:"vid"`
	WaterMark  bool        `json:"hidemark"` // true means in the whitelist
	CornerMark *CornerMark `json:"cornermark"`
}

// CornerMark def.
type CornerMark struct {
	Title string `json:"title"`
	Cover string `json:"cover"`
}

// SnVipCorner def.
type SnVipCorner struct {
	Title string `json:"title"`
	Cover string `json:"cover"`
}

// CmsInterv def.
func (v *EpisodeV2) CmsInterv(epCMS *EpCMS) {
	if epCMS.Cover != "" {
		v.Cover = epCMS.Cover
	}
	if epCMS.Title != "" {
		v.LongTitle = epCMS.Title
	}
}

// NewEPV2 def.
type NewEPV2 struct {
	Desc  string `json:"desc"`
	ID    int64  `json:"id"`
	IsNew int    `json:"is_new"`
	Title string `json:"title"`
}

// StatV2 def. 3 new fields
type StatV2 struct {
	Coin  int `json:"coin"`
	Reply int `json:"reply"`
	Share int `json:"share"`
	Stat
}

// UserStatusV2 def.
type UserStatusV2 struct {
	Follow   int            `json:"follow"`
	Pay      int            `json:"pay"`
	Progress *WatchProgress `json:"watch_progress"`
	Review   *ReviewV2      `json:"review"`
	Sponsor  int            `json:"sponsor"`
}

// ReviewV2 def.
type ReviewV2 struct {
	IsOpen int `json:"is_open"`
}

// Response standard structure
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// CodeErr generates the code error
func (r *Response) CodeErr() (err error) {
	if r.Code != ecode.OK.Code() {
		err = ecode.Int(r.Code)
	}
	return
}

// MediaRespV2 is the structure of PGC display api response
type MediaRespV2 struct {
	Response
	Result *SnDetailV2 `json:"result"`
}
