package banner

import (
	"encoding/json"
	xtime "go-common/library/time"
	"strconv"

	"go-common/app/interface/main/app-show/model"
	"go-common/app/interface/main/app-show/model/bangumi"
	resource "go-common/app/service/main/resource/model"
)

// Banner struct
type Banner struct {
	ID          int             `json:"id"`
	ParentID    int             `json:"-"`
	Plat        int8            `json:"-"`
	Module      string          `json:"-"`
	Position    string          `json:"-"`
	Title       string          `json:"title"`
	Image       string          `json:"image"`
	Hash        string          `json:"hash"`
	URI         string          `json:"uri"`
	Goto        string          `json:"-"`
	Value       string          `json:"-"`
	Param       string          `json:"-"`
	Channel     string          `json:"-"`
	Build       int             `json:"-"`
	Condition   string          `json:"-"`
	Area        string          `json:"-"`
	Rank        int             `json:"-"`
	Rule        string          `json:"-"`
	Type        int8            `json:"-"`
	Start       xtime.Time      `json:"-"`
	End         xtime.Time      `json:"-"`
	MTime       xtime.Time      `json:"-"`
	ResourceID  int             `json:"resource_id"`
	RequestId   string          `json:"request_id,omitempty"`
	CreativeId  int             `json:"creative_id,omitempty"`
	SrcId       int             `json:"src_id,omitempty"`
	IsAd        bool            `json:"is_ad"`
	IsAdReplace bool            `json:"-"`
	IsAdLoc     bool            `json:"is_ad_loc,omitempty"`
	CmMark      int             `json:"cm_mark"`
	AdCb        string          `json:"ad_cb,omitempty"`
	ShowUrl     string          `json:"show_url,omitempty"`
	ClickUrl    string          `json:"click_url,omitempty"`
	ClientIp    string          `json:"client_ip,omitempty"`
	Index       int             `json:"index"`
	ServerType  int             `json:"server_type"`
	Extra       json.RawMessage `json:"extra,omitempty"`
}

type JsonBanner struct {
	Area      string `json:"area"`
	Hash      string `json:"hash"`
	Build     int    `json:"build"`
	Condition string `json:"conditions"`
	Channel   string `json:"channel"`
}

// Banner limit
type Limit struct {
	Rule string `json:"-"`
}

// Json limit
type JsonLimit struct {
	Limit int      `json:"limit"`
	Resrc []string `json:"resrc"`
}

// PlatChange
func (b *Banner) BannerChange() {
	var tmp *JsonBanner
	if err := json.Unmarshal([]byte(b.Rule), &tmp); err == nil {
		b.Area = tmp.Area
		b.Build = tmp.Build
		b.Condition = tmp.Condition
		if tmp.Channel == "" {
			b.Channel = "*"
		} else {
			b.Channel = tmp.Channel
		}
		b.Hash = tmp.Hash
	}
	switch b.Plat {
	case 1: // resource iphone
		b.Plat = model.PlatIPhone
	case 2: // resource android
		b.Plat = model.PlatAndroid
	case 3: // resource pad
		b.Plat = model.PlatIPad
	case 4: // resource iphoneg
		b.Plat = model.PlatIPhoneI
	case 5: // resource androidg
		b.Plat = model.PlatAndroidG
	case 6: // resource padg
		b.Plat = model.PlatIPadI
	case 8: // resource androidi
		b.Plat = model.PlatAndroidI
	}
	if b.Value == "" {
		return
	}
	switch b.Type {
	case 6:
		//GotoAv
		b.URI = "bilibili://video/" + b.Value
	case 4:
		//GotoLive
		if b.Plat == model.PlatIPad {
			b.URI = "bilibili://player/live/" + b.Value
		} else {
			b.URI = "bilibili://live/" + b.Value
		}
	case 3:
		//GotoBangumi
		b.URI = "bilibili://bangumi/season/" + b.Value
	case 5:
		//GotoGame
		b.URI = "bilibili://game/" + b.Value
	case 2:
		//GotoWeb
		b.URI = b.Value
	}
}

// LimitChange
func (l *Limit) LimitChange() (data map[int]int) {
	data = map[int]int{}
	var (
		tmp   = &JsonLimit{}
		err   error
		resid int
	)
	if err = json.Unmarshal([]byte(l.Rule), tmp); err != nil {
		return
	}
	l.Rule = ""
	for _, residstr := range tmp.Resrc {
		resid, err = strconv.Atoi(residstr)
		if err != nil {
			return
		}
		data[resid] = tmp.Limit
	}
	return
}

// ResChangeBanner
func (b *Banner) ResChangeBanner(resb *resource.Banner) {
	b.ID = resb.ID
	b.Title = resb.Title
	b.Image = resb.Image
	b.Hash = resb.Hash
	b.URI = resb.URI
	b.ResourceID = resb.ResourceID
	b.RequestId = resb.RequestId
	b.CreativeId = resb.CreativeId
	b.SrcId = resb.SrcId
	b.IsAd = resb.IsAd
	b.IsAdLoc = resb.IsAdLoc
	b.CmMark = resb.CmMark
	b.AdCb = resb.AdCb
	b.ShowUrl = resb.ShowUrl
	b.ClickUrl = resb.ClickUrl
	b.ClientIp = resb.ClientIp
	b.Index = resb.Index
	b.ServerType = resb.ServerType
	b.Extra = resb.Extra
}

// BgmChangeBanner bangumiBanner change banner
func (b *Banner) BgmChangeBanner(bgmb *bangumi.Banner) {
	b.Title = bgmb.Title
	b.Image = bgmb.Image
	b.URI = bgmb.URI
}
