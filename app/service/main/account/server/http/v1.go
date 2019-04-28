package http

import (
	"strconv"

	v1 "go-common/app/service/main/account/api"
	"go-common/app/service/main/account/model"
	bm "go-common/library/net/http/blademaster"
)

// v1Info
func v1Info(c *bm.Context) {
	p := new(model.ParamMid)
	if err := c.Bind(p); err != nil {
		return
	}
	card, err := accSvc.Card(c, p.Mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	i := &V1Info{}
	i.FromCard(card)
	c.JSON(i, nil)
}

// v1Infos
func v1Infos(c *bm.Context) {
	p := new(model.ParamMids)
	if err := c.Bind(p); err != nil {
		return
	}
	cards, err := accSvc.Cards(c, p.Mids)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	im := make(map[int64]*V1Info, len(cards))
	for _, card := range cards {
		i := &V1Info{}
		i.FromCard(card)
		im[card.Mid] = i
	}
	c.JSON(im, nil)
}

// card
func v1Card(c *bm.Context) {
	p := new(model.ParamMid)
	if err := c.Bind(p); err != nil {
		return
	}
	ps, err := accSvc.ProfileWithStat(c, p.Mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	card := &V1Card{}
	card.FromProfile(ps)
	c.JSON(card, nil)
}

// vip
func v1Vip(c *bm.Context) {
	p := new(model.ParamMid)
	if err := c.Bind(p); err != nil {
		return
	}
	vi, err := accSvc.Vip(c, p.Mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	v := &V1Vip{}
	v.FromVip(vi)
	c.JSON(v, nil)
}

// V1Info info.
type V1Info struct {
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
	Pendant        v1.PendantInfo    `json:"pendant"`
	Nameplate      v1.NameplateInfo  `json:"nameplate"`
	OfficialVerify model.OldOfficial `json:"official_verify"`
	Vip            struct {
		Type          int    `json:"vipType"`
		DueDate       int64  `json:"vipDueDate"`
		DueRemark     string `json:"dueRemark"`
		AccessStatus  int    `json:"accessStatus"`
		VipStatus     int    `json:"vipStatus"`
		VipStatusWarn string `json:"vipStatusWarn"`
	} `json:"vip"`
}

// FromCard from card.
func (i *V1Info) FromCard(c *v1.Card) {
	i.Mid = strconv.FormatInt(c.Mid, 10)
	i.Name = c.Name
	i.Sex = c.Sex
	i.Sign = c.Sign
	i.Avatar = c.Face
	i.Rank = strconv.FormatInt(int64(c.Rank), 10)
	i.DisplayRank = "0"
	i.LevelInfo.Cur = int(c.Level)
	i.LevelInfo.Min = 0
	i.LevelInfo.NowExp = 0
	i.LevelInfo.NextExp = 0
	i.Pendant = c.Pendant
	i.Nameplate = c.Nameplate
	i.OfficialVerify = model.CvtOfficial(c.Official)
	i.Vip.Type = int(c.Vip.Type)
	i.Vip.VipStatus = int(c.Vip.Status)
	i.Vip.DueDate = c.Vip.DueDate
}

// V1Card card
type V1Card struct {
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
	Pendant        v1.PendantInfo    `json:"pendant"`
	Nameplate      v1.NameplateInfo  `json:"nameplate"`
	OfficialVerify model.OldOfficial `json:"official_verify"`
	Vip            struct {
		Type          int    `json:"vipType"`
		DueDate       int64  `json:"vipDueDate"`
		DueRemark     string `json:"dueRemark"`
		AccessStatus  int    `json:"accessStatus"`
		VipStatus     int    `json:"vipStatus"`
		VipStatusWarn string `json:"vipStatusWarn"`
	} `json:"vip"`
}

// FromProfile from profile.
func (i *V1Card) FromProfile(c *model.ProfileStat) {
	i.Mid = strconv.FormatInt(c.Mid, 10)
	i.Name = c.Name
	i.Sex = c.Sex
	i.Sign = c.Sign
	i.Face = c.Face
	i.Rank = strconv.FormatInt(int64(c.Rank), 10)
	i.DisplayRank = "0"
	i.Regtime = int64(c.JoinTime)
	if c.Silence == 1 {
		i.Spacesta = -2
	}
	i.Attentions = []int64{}
	i.Fans = int(c.Follower)
	i.Attention = int(c.Following)
	i.LevelInfo.Cur = int(c.Level)
	i.LevelInfo.Min = int(c.LevelExp.Min)
	i.LevelInfo.NowExp = int(c.LevelExp.NowExp)
	i.LevelInfo.NextExp = c.LevelExp.NextExp
	if c.LevelExp.NowExp == -1 {
		i.LevelInfo.NextExp = "--"
	}
	i.Pendant = c.Pendant
	i.Nameplate = c.Nameplate
	i.OfficialVerify = model.CvtOfficial(c.Official)
	i.Vip.Type = int(c.Vip.Type)
	i.Vip.VipStatus = int(c.Vip.Status)
	i.Vip.DueDate = c.Vip.DueDate
}

// V1Vip vip
type V1Vip struct {
	Type          int    `json:"vipType"`
	DueDate       int64  `json:"vipDueDate"`
	DueRemark     string `json:"dueRemark"`
	AccessStatus  int    `json:"accessStatus"`
	VipStatus     int    `json:"vipStatus"`
	VipStatusWarn string `json:"vipStatusWarn"`
}

// FromVip from vip.
func (v *V1Vip) FromVip(vi *v1.VipInfo) {
	v.Type = int(vi.Type)
	v.VipStatus = int(vi.Status)
	v.DueDate = vi.DueDate
}
