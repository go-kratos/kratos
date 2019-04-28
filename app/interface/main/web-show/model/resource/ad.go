package resource

import xtime "go-common/library/time"

//  StrategyOnly int8
const (
	StrategyOnly  = int8(0) // only
	StrategyShare = int8(1) // share
	StrategyRank  = int8(2) // rank
)

// VideoAD is Ads of videos
type VideoAD struct {
	ID       int        `json:"-"`
	Name     string     `json:"name"`
	AidS     string     `json:"-"`
	Aid      int64      `json:"aid"`
	Cid      int64      `json:"cid"`
	URL      string     `json:"url"`
	Skipable int8       `json:"skipable"`
	Strategy int8       `json:"strategy"`
	MTime    xtime.Time `json:"-"`
}

// Ad struct
type Ad struct {
	RequestID string                         `json:"request_id"`
	AdsInfo   map[string]map[string]*AdsInfo `json:"ads_info"`
}

// AdsInfo struct
type AdsInfo struct {
	Index  int64   `json:"index"`
	IsAd   bool    `json:"is_ad"`
	CmMark int8    `json:"cm_mark"`
	AdInfo *AdInfo `json:"ad_info"`
}

// CreativeImage type
const (
	CreativeImage = int8(0)
	CreativeVideo = int8(1)
)

// AdInfo struct
type AdInfo struct {
	CreativeID      int64 `json:"creative_id"`
	CreativeType    int8  `json:"creative_type"`
	CreativeContent struct {
		Title        string `json:"title"`
		Desc         string `json:"description"`
		VideoID      int64  `json:"video_id"`
		UserName     string `json:"username"`
		ImageURL     string `json:"image_url"`
		ImageMD5     string `json:"image_md5"`
		LogURL       string `json:"log_url"`
		LogMD5       string `json:"log_md5"`
		URL          string `json:"url"`
		ClickURL     string `json:"click_url"`
		ShowURL      string `json:"show_url"`
		ThumbnailURL string `json:"thumbnail_url"`
	} `json:"creative_content"`
	AdCb string `json:"ad_cb"`
}
