package model

import (
	"strconv"

	accmdl "go-common/app/service/main/account/api"
	account "go-common/app/service/main/account/model"
)

// Space space top photo
type Space struct {
	SImg string `json:"s_img"`
	LImg string `json:"l_img"`
}

// Card  Card  and Space and Relation and Archive Count.
type Card struct {
	Card         *AccountCard `json:"card"`
	Space        *Space       `json:"space,omitempty"`
	Following    bool         `json:"following"`
	ArchiveCount int          `json:"archive_count"`
	ArticleCount int          `json:"article_count"`
	Follower     int64        `json:"follower"`
}

// AccountCard struct.
type AccountCard struct {
	Mid         string  `json:"mid"`
	Name        string  `json:"name"`
	Approve     bool    `json:"approve"`
	Sex         string  `json:"sex"`
	Rank        string  `json:"rank"`
	Face        string  `json:"face"`
	DisplayRank string  `json:"DisplayRank"`
	Regtime     int64   `json:"regtime"`
	Spacesta    int     `json:"spacesta"`
	Birthday    string  `json:"birthday"`
	Place       string  `json:"place"`
	Description string  `json:"description"`
	Article     int     `json:"article"`
	Attentions  []int64 `json:"attentions"`
	Fans        int     `json:"fans"`
	Friend      int     `json:"friend"`
	Attention   int     `json:"attention"`
	Sign        string  `json:"sign"`
	LevelInfo   struct {
		Cur     int         `json:"current_level"`
		Min     int         `json:"current_min"`
		NowExp  int         `json:"current_exp"`
		NextExp interface{} `json:"next_exp"`
	} `json:"level_info"`
	Pendant        account.PendantInfo   `json:"pendant"`
	Nameplate      account.NameplateInfo `json:"nameplate"`
	Official       accmdl.OfficialInfo
	OfficialVerify struct {
		Type int    `json:"type"`
		Desc string `json:"desc"`
	} `json:"official_verify"`
	Vip struct {
		Type          int    `json:"vipType"`
		DueRemark     string `json:"dueRemark"`
		AccessStatus  int    `json:"accessStatus"`
		VipStatus     int    `json:"vipStatus"`
		VipStatusWarn string `json:"vipStatusWarn"`
	} `json:"vip"`
}

// FromCard from account catd.
func (ac *AccountCard) FromCard(c *account.Card) {
	ac.Mid = strconv.FormatInt(c.Mid, 10)
	ac.Name = c.Name
	// ac.Approve =
	ac.Sex = c.Sex
	ac.Rank = strconv.FormatInt(int64(c.Rank), 10)
	ac.DisplayRank = "0"
	ac.Face = c.Face
	// ac.Regtime =
	if c.Silence == 1 {
		ac.Spacesta = -2
	}
	// ac.Birthday =
	// ac.Place =
	// ac.Description =
	// ac.Article =
	// ac.Attentions = []int64{}
	// ac.Fans =
	// ac.Friend
	// ac.Attention =
	ac.Sign = c.Sign
	ac.LevelInfo.Cur = int(c.Level)
	ac.LevelInfo.NextExp = 0
	// ac.LevelInfo.Min =
	ac.Pendant = c.Pendant
	ac.Nameplate = c.Nameplate
	if c.Official.Role == 0 {
		ac.OfficialVerify.Type = -1
	} else {
		if c.Official.Role <= 2 {
			ac.OfficialVerify.Type = 0
			ac.OfficialVerify.Desc = c.Official.Title
		} else {
			ac.OfficialVerify.Type = 1
			ac.OfficialVerify.Desc = c.Official.Title
		}
	}
	ac.Official = c.Official
	ac.Vip.Type = int(c.Vip.Type)
	ac.Vip.VipStatus = int(c.Vip.Status)
}

// DefaultProfile .
var DefaultProfile = &accmdl.ProfileStatReply{
	Profile: &account.Profile{
		Sex:  "保密",
		Rank: 10000,
		Face: "https://static.hdslb.com/images/member/noface.gif",
		Sign: "没签名",
	},
	LevelInfo: accmdl.LevelInfo{},
}
