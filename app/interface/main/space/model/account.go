package model

import (
	accmdl "go-common/app/service/main/account/api"
)

// NavNum nav num struct.
type NavNum struct {
	Video     int64 `json:"video"`
	Bangumi   int   `json:"bangumi"`
	Channel   *Num  `json:"channel"`
	Favourite *Num  `json:"favourite"`
	Tag       int   `json:"tag"`
	Article   int   `json:"article"`
	Playlist  int   `json:"playlist"`
	Album     int64 `json:"album"`
	Audio     int   `json:"audio"`
}

// Num num struct.
type Num struct {
	Master int `json:"master"`
	Guest  int `json:"guest"`
}

// UpStat up stat struct.
type UpStat struct {
	Archive struct {
		View int64 `json:"view"`
	} `json:"archive"`
	Article struct {
		View int64 `json:"view"`
	} `json:"article"`
}

// AccInfo account info.
type AccInfo struct {
	Mid       int64               `json:"mid"`
	Name      string              `json:"name"`
	Sex       string              `json:"sex"`
	Face      string              `json:"face"`
	Sign      string              `json:"sign"`
	Rank      int32               `json:"rank"`
	Level     int32               `json:"level"`
	JoinTime  int32               `json:"jointime"`
	Moral     int32               `json:"moral"`
	Silence   int32               `json:"silence"`
	Birthday  string              `json:"birthday"`
	Coins     float64             `json:"coins"`
	FansBadge bool                `json:"fans_badge"`
	Official  accmdl.OfficialInfo `json:"official"`
	Vip       struct {
		Type   int32 `json:"type"`
		Status int32 `json:"status"`
	} `json:"vip"`
	IsFollowed bool        `json:"is_followed"`
	TopPhoto   string      `json:"top_photo"`
	Theme      interface{} `json:"theme"`
}

// AccBlock acc block
type AccBlock struct {
	Status     int `json:"status"`
	IsDue      int `json:"is_due"`
	IsAnswered int `json:"is_answered"`
}

// TopPhoto top photo struct.
type TopPhoto struct {
	SImg string `json:"s_img"`
	LImg string `json:"l_img"`
}

// Relation .
type Relation struct {
	Relation   interface{} `json:"relation"`
	BeRelation interface{} `json:"be_relation"`
}

// FromCard from account card.
func (ai *AccInfo) FromCard(c *accmdl.ProfileStatReply) {
	ai.Mid = c.Profile.Mid
	ai.Name = c.Profile.Name
	ai.Rank = c.Profile.Rank
	ai.Face = c.Profile.Face
	ai.Sex = c.Profile.Sex
	ai.JoinTime = c.Profile.JoinTime
	ai.Silence = c.Profile.Silence
	ai.Birthday = c.Profile.Birthday.Time().Format("01-02")
	ai.Sign = c.Profile.Sign
	ai.Level = c.Profile.Level
	ai.Official = c.Profile.Official
	ai.Vip.Type = c.Profile.Vip.Type
	ai.Vip.Status = c.Profile.Vip.Status
	ai.Coins = c.Coins
}

var (
	// DefaultProfileStat .
	DefaultProfileStat = &accmdl.ProfileStatReply{
		Profile:   DefaultProfile,
		LevelInfo: accmdl.LevelInfo{},
	}
	// DefaultProfile .
	DefaultProfile = &accmdl.Profile{
		Name: "bilibili",
		Sex:  "保密",
		Face: "https://static.hdslb.com/images/member/noface.gif",
		Sign: "哔哩哔哩 (゜-゜)つロ 干杯~-bilibili",
		Rank: 5000,
	}
)
