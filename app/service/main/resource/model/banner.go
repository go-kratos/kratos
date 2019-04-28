package model

import (
	"encoding/json"
	"strconv"

	xtime "go-common/library/time"
)

// Banners struct
type Banners struct {
	Banner  map[int][]*Banner
	Version string
}

// Banner struct
type Banner struct {
	ID           int             `json:"id"`
	ParentID     int             `json:"-"`
	Plat         int8            `json:"-"`
	Module       string          `json:"-"`
	Position     string          `json:"-"`
	Title        string          `json:"title"`
	Image        string          `json:"image"`
	Hash         string          `json:"hash"`
	URI          string          `json:"uri"`
	Goto         string          `json:"-"`
	Value        string          `json:"-"`
	Param        string          `json:"-"`
	Channel      string          `json:"-"`
	Build        int             `json:"-"`
	Condition    string          `json:"-"`
	Area         string          `json:"-"`
	Rank         int             `json:"-"`
	Rule         string          `json:"-"`
	Type         int8            `json:"-"`
	Start        xtime.Time      `json:"stime"`
	End          xtime.Time      `json:"-"`
	MTime        xtime.Time      `json:"-"`
	ResourceID   int             `json:"resource_id"`
	RequestId    string          `json:"request_id,omitempty"`
	CreativeId   int             `json:"creative_id,omitempty"`
	SrcId        int             `json:"src_id,omitempty"`
	IsAd         bool            `json:"is_ad"`
	IsAdReplace  bool            `json:"-"`
	IsAdLoc      bool            `json:"is_ad_loc,omitempty"`
	CmMark       int             `json:"cm_mark"`
	AdCb         string          `json:"ad_cb,omitempty"`
	ShowUrl      string          `json:"show_url,omitempty"`
	ClickUrl     string          `json:"click_url,omitempty"`
	ClientIp     string          `json:"client_ip,omitempty"`
	Index        int             `json:"index"`
	ServerType   int             `json:"server_type"`
	Extra        json.RawMessage `json:"extra"`
	CreativeType int             `json:"creative_type"`
}

// JSONBanner bilibili_assignment rule
type JSONBanner struct {
	Area         string `json:"area"`
	Hash         string `json:"hash"`
	Build        int    `json:"build"`
	Condition    string `json:"cond"`
	Channel      string `json:"channel"`
	CreativeType int    `json:"creative_type"`
}

// Limit limit
type Limit struct {
	Rule string `json:"-"`
}

// JSONLimit limit
type JSONLimit struct {
	Limit int      `json:"limit"`
	Resrc []string `json:"resrc"`
}

// BannerChange change banner
func (b *Banner) BannerChange() {
	var tmp *JSONBanner
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
		b.CreativeType = tmp.CreativeType
	}
	switch b.Plat {
	case 1: // resource iphone
		b.Plat = PlatIPhone
	case 2: // resource android
		b.Plat = PlatAndroid
	case 3: // resource pad
		b.Plat = PlatIPad
	case 4: // resource iphoneg
		b.Plat = PlatIPhoneI
	case 5: // resource androidg
		b.Plat = PlatAndroidG
	case 6: // resource padg
		b.Plat = PlatIPadI
	case 8: // resource androidi
		b.Plat = PlatAndroidI
	}
	if b.Value == "" {
		return
	}
	switch b.Type {
	case 7:
		if b.Plat == PlatIPhone || b.Plat == PlatAndroid || b.Plat == PlatIPad || b.Plat == PlatIPhoneI || b.Plat == PlatAndroidG || b.Plat == PlatIPadI || b.Plat == PlatAndroidI {
			b.URI = "bilibili://pegasus/channel/" + b.Value + "/"
		} else {
			b.URI = "http://www.bilibili.com/tag/" + b.Value
		}
	case 6:
		//GotoAv
		b.URI = "bilibili://video/" + b.Value
	case 4:
		//GotoLive
		if b.Plat == PlatIPad {
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

// LimitChange change limit
func (l *Limit) LimitChange() (data map[int]int) {
	data = map[int]int{}
	var (
		tmp   = &JSONLimit{}
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
