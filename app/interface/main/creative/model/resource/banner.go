package source

import (
	"encoding/json"
	resource "go-common/app/service/main/resource/model"
	xtime "go-common/library/time"
)

// Banner str
type Banner struct {
	ID           int             `json:"id"`
	ParentID     int             `json:"-"`
	Plat         int8            `json:"-"`
	Module       string          `json:"-"`
	Position     string          `json:"-"`
	Title        string          `json:"title"`
	Content      string          `json:"content"`
	Image        string          `json:"image"`
	Pic          string          `json:"pic"`
	Hash         string          `json:"hash"`
	URI          string          `json:"uri"`
	Link         string          `json:"link"`
	Goto         string          `json:"-"`
	Value        string          `json:"-"`
	Param        string          `json:"-"`
	Channel      string          `json:"-"`
	Build        int             `json:"-"`
	Condition    string          `json:"-"`
	Area         string          `json:"-"`
	Rule         string          `json:"-"`
	Type         int8            `json:"-"`
	Start        xtime.Time      `json:"-"`
	End          xtime.Time      `json:"-"`
	MTime        xtime.Time      `json:"-"`
	ResourceID   int             `json:"resource_id"`
	RequestId    string          `json:"request_id"`
	CreativeId   int             `json:"creative_id"`
	SrcId        int             `json:"src_id"`
	IsAd         bool            `json:"is_ad"`
	IsAdReplace  bool            `json:"-"`
	IsAdLoc      bool            `json:"is_ad_loc"`
	CmMark       int             `json:"cm_mark"`
	AdCb         string          `json:"ad_cb"`
	ShowUrl      string          `json:"show_url"`
	ClickUrl     string          `json:"click_url"`
	ClientIp     string          `json:"client_ip"`
	Index        int             `json:"index"`
	Rank         int             `json:"rank"`
	ServerType   int             `json:"server_type"`
	Extra        json.RawMessage `json:"extra"`
	CreativeType int             `json:"creative_type"`
}

// ChangeBanner fn
func (b *Banner) ChangeBanner(banner *resource.Banner) {
	b.ID = banner.ID
	b.Rank = banner.Index
	b.Title = banner.Title
	b.Content = banner.Title
	b.Image = banner.Image
	b.Pic = banner.Image
	b.Hash = banner.Hash
	b.URI = banner.URI
	b.Link = banner.URI
	b.ResourceID = banner.ResourceID
	b.RequestId = banner.RequestId
	b.CreativeId = banner.CreativeId
	b.SrcId = banner.SrcId
	b.IsAd = banner.IsAd
	b.IsAdLoc = banner.IsAdLoc
	b.CmMark = banner.CmMark
	b.AdCb = banner.AdCb
	b.ShowUrl = banner.ShowUrl
	b.ClickUrl = banner.ClickUrl
	b.ClientIp = banner.ClientIp
	b.Index = banner.Index
	b.ServerType = banner.ServerType
	b.Extra = banner.Extra
	b.CreativeType = banner.CreativeType
	b.CreativeId = banner.CreativeId
}

// BannerList for operation list.
type BannerList struct {
	Banners []*Banner `json:"operations"`
	Pn      int       `json:"pn"`
	Ps      int       `json:"ps"`
	Total   int       `json:"total"`
}
