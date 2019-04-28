package space

import (
	accv1 "go-common/app/service/main/account/api"
	accmdl "go-common/app/service/main/account/model"
)

// Mine my center struct
type Mine struct {
	Mid          int64   `json:"mid"`
	Name         string  `json:"name"`
	Face         string  `json:"face"`
	Coin         float64 `json:"coin"`
	BCoin        float64 `json:"bcoin"`
	Sex          int32   `json:"sex"`
	Rank         int32   `json:"rank"`
	Silence      int32   `json:"silence"`
	EndTime      int64   `json:"end_time,omitempty"`
	ShowVideoup  int     `json:"show_videoup"`
	ShowCreative int     `json:"show_creative"`
	Level        int32   `json:"level"`
	VipType      int32   `json:"vip_type"`
	AudioType    int     `json:"audio_type"`
	Dynamic      int64   `json:"dynamic"`
	Following    int64   `json:"following"`
	Follower     int64   `json:"follower"`
	NewFollowers int64   `json:"new_followers"`
	Official     struct {
		Type int8   `json:"type"`
		Desc string `json:"desc"`
	} `json:"official_verify"`
	Pendant           *Pendant       `json:"pendant,omitempty"`
	Sections          []*Section     `json:"sections,omitempty"`
	IpadSections      []*SectionItem `json:"ipad_sections,omitempty"`
	IpadUpperSections []*SectionItem `json:"ipad_upper_sections,omitempty"`
}

// Section for mine page, like 【个人中心】【我的服务】
type Section struct {
	Title string         `json:"title"`
	Items []*SectionItem `json:"items"`
}

// SectionItem like 【离线缓存】 【历史记录】,a part of section
type SectionItem struct {
	Title     string `json:"title"`
	URI       string `json:"uri"`
	Icon      string `json:"icon"`
	NeedLogin int8   `json:"need_login,omitempty"`
	RedDot    int8   `json:"red_dot,omitempty"`
}

// Myinfo myinfo
type Myinfo struct {
	Mid            int64              `json:"mid"`
	Name           string             `json:"name"`
	Sign           string             `json:"sign"`
	Coins          float64            `json:"coins"`
	Birthday       string             `json:"birthday"`
	Face           string             `json:"face"`
	Sex            int                `json:"sex"`
	Level          int32              `json:"level"`
	Rank           int32              `json:"rank"`
	Silence        int32              `json:"silence"`
	EndTime        int64              `json:"end_time,omitempty"`
	Vip            accmdl.VipInfo     `json:"vip"`
	EmailStatus    int32              `json:"email_status"`
	TelStatus      int32              `json:"tel_status"`
	Official       accv1.OfficialInfo `json:"official"`
	Identification int32              `json:"identification"`
	Pendant        *Pendant           `json:"pendant,omitempty"`
}

// MineParam struct
type MineParam struct {
	MobiApp  string `form:"mobi_app"`
	Device   string `form:"device"`
	Build    int    `form:"build"`
	Platform string `form:"platform"`
	Mid      int64  `form:"mid"`
	Filtered string `form:"filtered"`
}

// Pendant struct
type Pendant struct {
	Image string `json:"image"`
}
