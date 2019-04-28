package model

import (
	"strconv"

	accountv1 "go-common/app/service/main/account/api"
	relation "go-common/app/service/main/relation/model"
	bm "go-common/library/net/http/blademaster"
)

// Vip .
type Vip struct {
	Type          int    `json:"vipType"`
	DueDate       int64  `json:"vipDueDate"`
	DueRemark     string `json:"dueRemark"`
	AccessStatus  int    `json:"accessStatus"`
	VipStatus     int    `json:"vipStatus"`
	VipStatusWarn string `json:"vipStatusWarn"`
}

// Following is user followinng info.
type Following struct {
	*relation.Following
	Uname string `json:"uname"`
	Face  string `json:"face"`
	Sign  string `json:"sign"`
	// OfficialVerify member.OfficialInfo `json:"official_verify"`
	OfficialVerify struct {
		Type int8   `json:"type"`
		Desc string `json:"desc"`
	} `json:"official_verify"`
	Vip Vip `json:"vip"`
}

// Tag is user info.
type Tag struct {
	Mid   int64  `json:"mid"`
	Uname string `json:"uname"`
	Face  string `json:"face"`
	Sign  string `json:"sign"`
	// OfficialVerify member.OfficialInfo `json:"official_verify"`
	OfficialVerify struct {
		Type int8   `json:"type"`
		Desc string `json:"desc"`
	} `json:"official_verify"`
	Vip Vip `json:"vip"`
}

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
	Pendant        accountv1.PendantInfo   `json:"pendant"`
	Nameplate      accountv1.NameplateInfo `json:"nameplate"`
	OfficialVerify accountv1.OfficialInfo  `json:"official_verify"`
	Vip            struct {
		Type          int    `json:"vipType"`
		DueDate       int64  `json:"vipDueDate"`
		DueRemark     string `json:"dueRemark"`
		AccessStatus  int    `json:"accessStatus"`
		VipStatus     int    `json:"vipStatus"`
		VipStatusWarn string `json:"vipStatusWarn"`
	} `json:"vip"`
}

// RecommendInfo is
type RecommendInfo struct {
	Info
	RecommendContent
	Fans           int64               `json:"fans"`
	TypeName       string              `json:"type_name"`
	SecondTypeName string              `json:"second_type_name"`
	TrackID        string              `json:"track_id"`
	Relation       *relation.Following `json:"relation"`
}

// TagSuggestRecommendInfo is
type TagSuggestRecommendInfo struct {
	TagName  string           `json:"tagname"`
	UpList   []*RecommendInfo `json:"up_list"`
	MatchCnt int64            `json:"match_cnt"`
}

/*
{
    "code": 0,
    "trackid": "123",
    "msg": "success",
    "data": [
        {
            "up_id": 123,
            "rec_reason": "游戏区热门up主",
            "rec_type": 1,
            "tid": 4,
            "second_tid": 173
        }
    ]
}
*/

// RecommendContent is
type RecommendContent struct {
	UpID      int64  `json:"up_id"`
	RecReason string `json:"rec_reason"`
	RecType   int64  `json:"rec_type"`
	Tid       int16  `json:"tid"`
	SecondTid int16  `json:"second_tid"`
}

// RecommendResponse is
type RecommendResponse struct {
	Code    int64               `json:"code"`
	TrackID string              `json:"trackid"`
	Msg     string              `json:"msg"`
	Data    []*RecommendContent `json:"data"`
}

// TagSuggestRecommendContent is
type TagSuggestRecommendContent struct {
	TagName  string              `json:"tagname"`
	UpList   []*RecommendContent `json:"up_list"`
	MatchCnt int64               `json:"match_cnt"`
}

// UpIDs is
func (tsrc *TagSuggestRecommendContent) UpIDs() []int64 {
	upIDs := make([]int64, 0, len(tsrc.UpList))
	for _, up := range tsrc.UpList {
		upIDs = append(upIDs, up.UpID)
	}
	return upIDs
}

// TagSuggestRecommendResponse is
type TagSuggestRecommendResponse struct {
	Code    int64                         `json:"code"`
	TrackID string                        `json:"trackid"`
	Msg     string                        `json:"msg"`
	Data    []*TagSuggestRecommendContent `json:"data"`
}

// FromCard from card.
func (i *Info) FromCard(c *accountv1.Card) {
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
	i.OfficialVerify = c.Official
	i.Vip.Type = int(c.Vip.Type)
	i.Vip.VipStatus = int(c.Vip.Status)
	i.Vip.DueDate = c.Vip.DueDate
}

// BatchModifyResult is
type BatchModifyResult struct {
	FailedFids []int64 `json:"failed_fids"`
}

// ArgRecommend is
type ArgRecommend struct {
	Mid      int64
	Device   *bm.Device
	RemoteIP string
	MainTids string `form:"main_tids"`
	SubTids  string `form:"sub_tids"`
	PageSize int64  `form:"pagesize" default:"10"`
}

// ArgTagSuggestRecommend is
type ArgTagSuggestRecommend struct {
	Mid       int64
	Device    *bm.Device
	RemoteIP  string
	TagName   string `form:"tagname"`
	ContextID string `form:"context_id" validate:"required"`
	PageSize  int64  `form:"pagesize" default:"10"`
}

// ArgAchieveGet is
type ArgAchieveGet struct {
	Mid   int64
	Award string `form:"award" validate:"required"`
}

// ArgAchieve is
type ArgAchieve struct {
	AwardToken string `form:"award_token" validate:"required"`
}

// AchieveReply is
type AchieveReply struct {
	relation.Achieve
	Metadata map[string]interface{} `json:"metadata"`
}

// ArgSameFollowing is
type ArgSameFollowing struct {
	Mid       int64  `form:"mid"`
	VMid      int64  `form:"vmid" validate:"required"`
	Order     string `form:"order"`
	PS        int64  `form:"ps"`
	PN        int64  `form:"pn"`
	ReVersion uint64 `form:"re_version"`
}
