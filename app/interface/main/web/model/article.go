package model

import (
	artmdl "go-common/app/interface/openplatform/article/model"
	accmdl "go-common/app/service/main/account/model"
	"go-common/library/time"
	"strconv"
)

// Info struct.
type Info struct {
	Mid         string `json:"mid"`
	Name        string `json:"uname"`
	Sex         string `json:"sex"`
	Sign        string `json:"sign"`
	Avatar      string `json:"avatar"`
	Rank        string `json:"rank"`
	DisplayRank string `json:"DisplayRank"`
	LevelInfo   struct {
		Cur     int         `json:"current_level"`
		Min     int         `json:"current_min"`
		NowExp  int         `json:"current_exp"`
		NextExp interface{} `json:"next_exp"`
	} `json:"level_info"`
	Pendant        accmdl.PendantInfo   `json:"pendant"`
	Nameplate      accmdl.NameplateInfo `json:"nameplate"`
	OfficialVerify struct {
		Type int    `json:"type"`
		Desc string `json:"desc"`
	} `json:"official_verify"`
	Vip struct {
		Type          int    `json:"vipType"`
		DueDate       int64  `json:"vipDueDate"`
		DueRemark     string `json:"dueRemark"`
		AccessStatus  int    `json:"accessStatus"`
		VipStatus     int    `json:"vipStatus"`
		VipStatusWarn string `json:"vipStatusWarn"`
	} `json:"vip"`
	// article
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	PublishTime time.Time `json:"publish_time"`
	Following   bool      `json:"following"`
}

// FromCard from card.
func (i *Info) FromCard(c *accmdl.Card) {
	i.Mid = strconv.FormatInt(c.Mid, 10)
	i.Name = c.Name
	i.Sex = c.Sex
	i.Sign = c.Sign
	i.Avatar = c.Face
	i.Rank = strconv.FormatInt(int64(c.Rank), 10)
	i.DisplayRank = "0"
	i.LevelInfo.Cur = int(c.Level)
	i.LevelInfo.NextExp = 0
	// i.LevelInfo.Min =
	i.Pendant = c.Pendant
	i.Nameplate = c.Nameplate
	if c.Official.Role == 0 {
		i.OfficialVerify.Type = -1
	} else {
		if c.Official.Role <= 2 {
			i.OfficialVerify.Type = 0
			i.OfficialVerify.Desc = c.Official.Title
		} else {
			i.OfficialVerify.Type = 1
			i.OfficialVerify.Desc = c.Official.Title
		}
	}
	i.Vip.Type = int(c.Vip.Type)
	i.Vip.VipStatus = int(c.Vip.Status)
	i.Vip.DueDate = c.Vip.DueDate
}

// Meta struct.
type Meta struct {
	*artmdl.Meta
	Like int `json:"like"`
}

// ArticleUpInfo struct.
type ArticleUpInfo struct {
	ArtCount    int   `json:"art_count"`
	Follower    int64 `json:"follower"`
	IsFollowing bool  `json:"is_following"`
}
