package banner

import (
	"encoding/json"

	resource "go-common/app/service/main/resource/model"
)

// Banner struct
type Banner struct {
	ID         int64           `json:"id"`
	ParentID   int64           `json:"-"`
	Plat       int8            `json:"-"`
	Title      string          `json:"title"`
	Image      string          `json:"image"`
	Hash       string          `json:"hash"`
	URI        string          `json:"uri"`
	Value      string          `json:"-"`
	Channel    string          `json:"-"`
	Build      int             `json:"-"`
	Condition  string          `json:"-"`
	Area       string          `json:"-"`
	Rank       int64           `json:"-"`
	Rule       string          `json:"-"`
	Type       int8            `json:"-"`
	RequestID  string          `json:"request_id,omitempty"`
	CreativeID int             `json:"creative_id,omitempty"`
	SrcID      int             `json:"src_id,omitempty"`
	IsAd       bool            `json:"is_ad,omitempty"`
	IsAdLoc    bool            `json:"is_ad_loc,omitempty"`
	AdCb       string          `json:"ad_cb,omitempty"`
	ShowURL    string          `json:"show_url,omitempty"`
	ClickURL   string          `json:"click_url,omitempty"`
	ClientIP   string          `json:"client_ip,omitempty"`
	ServerType int             `json:"server_type"`
	ResourceID int             `json:"resource_id,omitempty"`
	Index      int             `json:"index,omitempty"`
	CmMark     int             `json:"cm_mark"`
	Extra      json.RawMessage `json:"extra,omitempty"`
}

func (b *Banner) ChangeBanner(banner *resource.Banner) {
	b.ID = int64(banner.ID)
	b.Title = banner.Title
	b.Image = banner.Image
	b.Hash = banner.Hash
	b.URI = banner.URI
	b.ResourceID = banner.ResourceID
	b.RequestID = banner.RequestId
	b.CreativeID = banner.CreativeId
	b.SrcID = banner.SrcId
	b.IsAd = banner.IsAd
	b.IsAdLoc = banner.IsAdLoc
	b.CmMark = banner.CmMark
	b.AdCb = banner.AdCb
	b.ShowURL = banner.ShowUrl
	b.ClickURL = banner.ClickUrl
	b.ClientIP = banner.ClientIp
	b.Index = banner.Index
	b.ServerType = banner.ServerType
	b.Extra = banner.Extra
}
