package cm

import (
	"encoding/json"
)

type Ad struct {
	RequestID string                     `json:"request_id,omitempty"`
	AdsInfo   map[int64]map[int]*AdsInfo `json:"ads_info,omitempty"`
	ClientIP  string                     `json:"-"`
}

type AdsInfo struct {
	Index     int     `json:"index,omitempty"`
	IsAd      bool    `json:"is_ad,omitempty"`
	CmMark    int64   `json:"cm_mark,omitempty"`
	AdInfo    *AdInfo `json:"ad_info,omitempty"`
	CardIndex int     `json:"card_index,omitempty"`
}

type AdInfo struct {
	CreativeID      int64 `json:"creative_id,omitempty"`
	CreativeType    int   `json:"creative_type,omitempty"`
	CardType        int   `json:"card_type,omitempty"`
	CreativeContent *struct {
		Title    string `json:"title,omitempty"`
		Desc     string `json:"description,omitempty"`
		VideoID  int64  `json:"video_id,omitempty"`
		UserName string `json:"username,omitempty"`
		ImageURL string `json:"image_url,omitempty"`
		ImageMD5 string `json:"image_md5,omitempty"`
		LogURL   string `json:"log_url,omitempty"`
		LogMD5   string `json:"log_md5,omitempty"`
		URL      string `json:"url,omitempty"`
		ClickURL string `json:"click_url,omitempty"`
		ShowURL  string `json:"show_url,omitempty"`
	} `json:"creative_content,omitempty"`
	AdCb      string          `json:"ad_cb,omitempty"`
	Resource  int64           `json:"resource,omitempty"`
	Source    int             `json:"source,omitempty"`
	RequestID string          `json:"request_id,omitempty"`
	IsAd      bool            `json:"is_ad,omitempty"`
	CmMark    int64           `json:"cm_mark,omitempty"`
	Index     int             `json:"index,omitempty"`
	IsAdLoc   bool            `json:"is_ad_loc,omitempty"`
	CardIndex int             `json:"card_index,omitempty"`
	ClientIP  string          `json:"client_ip,omitempty"`
	Extra     json.RawMessage `json:"extra,omitempty"`
}
