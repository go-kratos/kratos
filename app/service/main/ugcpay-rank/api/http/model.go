package http

import (
	"go-common/app/service/main/ugcpay-rank/internal/model"
)

// Common .
type Common struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
}

// RetRankElecAllAV .
type RetRankElecAllAV struct {
	Common
	Data *RespRankElecAllAV `json:"data,omitempty"`
}

// RetRankElecMonth .
type RetRankElecMonth struct {
	Common
	Data *RespRankElecMonth `json:"data,omitempty"`
}

// RetRankElecMonthUP .
type RetRankElecMonthUP struct {
	Common
	Data *RespRankElecMonthUP `json:"data,omitempty"`
}

// ArgRankElecMonth .
type ArgRankElecMonth struct {
	UPMID    int64 `form:"up_mid" validate:"required"` // up主MID
	AVID     int64 `form:"av_id" validate:"required"`
	RankSize int   `form:"rank_size"` // 榜单大小
}

// ArgRankElecMonthUP .
type ArgRankElecMonthUP struct {
	UPMID    int64 `form:"up_mid" validate:"required"` // up主MID
	RankSize int   `form:"rank_size"`                  // 榜单大小
}

// RespRankElecElement 榜单元素信息
type RespRankElecElement struct {
	UName     string       `json:"uname"`      // 用户名
	Avatar    string       `json:"avatar"`     // 头像
	MID       int64        `json:"mid"`        // up主
	PayMID    int64        `json:"pay_mid"`    // 金主爸爸
	Rank      int          `json:"rank"`       // 排名
	TrendType uint8        `json:"trend_type"` // 上升趋势
	VIP       *RespVIPInfo `json:"vip_info"`   // VIP 信息
}

// Parse .
func (r *RespRankElecElement) Parse(ele *model.RankElecElementProto, upMID int64) {
	r.UName = ele.Nickname
	r.Avatar = ele.Avatar
	r.MID = upMID
	r.PayMID = ele.MID
	r.Rank = ele.Rank
	r.TrendType = ele.TrendType
	r.VIP = &RespVIPInfo{}
	if ele.VIP != nil {
		r.VIP.VIPDueMsec = ele.VIP.DueDate
		r.VIP.VIPStatus = ele.VIP.Status
		r.VIP.VIPType = ele.VIP.Type
	}
}

// RespRankElecElementDetail 榜单元素详情
type RespRankElecElementDetail struct {
	RespRankElecElement
	Message        string `json:"message"`        // 留言
	MessasgeHidden int    `json:"message_hidden"` //
}

// Parse .
func (r *RespRankElecElementDetail) Parse(ele *model.RankElecElementProto, upMID int64) {
	r.RespRankElecElement.Parse(ele, upMID)
	if ele.Message != nil {
		r.Message = ele.Message.Message
		if ele.Message.Hidden {
			r.MessasgeHidden = 1
		} else {
			r.MessasgeHidden = 0
		}
	}
}

// RespRankElecMonth .
type RespRankElecMonth struct {
	AVCount    int64                        `json:"av_count"`
	AVList     []*RespRankElecElementDetail `json:"av_list"`
	UPCount    int64                        `json:"up_count"`
	UPList     []*RespRankElecElementDetail `json:"up_list"`
	ShowInfo   *RespShowInfo                `json:"show_info"`
	TotalCount int64                        `json:"total_count"`
}

// Parse .
func (r *RespRankElecMonth) Parse(avRank *model.RankElecAVProto, upRank *model.RankElecUPProto) {
	r.AVCount = avRank.Count
	r.AVList = make([]*RespRankElecElementDetail, 0)
	for _, ele := range avRank.List {
		data := &RespRankElecElementDetail{}
		data.Parse(ele, avRank.UPMID)
		r.AVList = append(r.AVList, data)
	}

	r.TotalCount = upRank.CountUPTotalElec
	r.UPCount = upRank.Count
	r.UPList = make([]*RespRankElecElementDetail, 0)
	for _, ele := range upRank.List {
		data := &RespRankElecElementDetail{}
		data.Parse(ele, upRank.UPMID)
		r.UPList = append(r.UPList, data)
	}
	r.ShowInfo = &RespShowInfo{
		Show:  true,
		State: 0,
	}
}

// RespRankElecMonthUP .
type RespRankElecMonthUP struct {
	Count      int64                        `json:"count"` // UP主维度月充电数量
	List       []*RespRankElecElementDetail `json:"list"`
	TotalCount int64                        `json:"total_count"`
}

// Parse .
func (r *RespRankElecMonthUP) Parse(monthlyRank *model.RankElecUPProto) {
	r.Count = monthlyRank.Count
	r.List = make([]*RespRankElecElementDetail, 0)
	for _, ele := range monthlyRank.List {
		data := &RespRankElecElementDetail{}
		data.Parse(ele, monthlyRank.UPMID)
		r.List = append(r.List, data)
	}
	r.TotalCount = monthlyRank.CountUPTotalElec
}

// RespRankElecAllAV .
type RespRankElecAllAV struct {
	TotalCount int64                        `json:"total_count"`
	List       []*RespRankElecElementDetail `json:"list"`
}

// Parse .
func (r *RespRankElecAllAV) Parse(rank *model.RankElecAVProto) {
	r.TotalCount = rank.Count
	r.List = make([]*RespRankElecElementDetail, 0)
	for _, ele := range rank.List {
		data := &RespRankElecElementDetail{}
		data.Parse(ele, rank.UPMID)
		r.List = append(r.List, data)
	}
}

// RespElecReply .
type RespElecReply struct {
	ReplyMID  int64  `json:"reply_mid"`
	ReplyMSG  string `json:"reply_msg"`
	ReplyName string `json:"reply_name"`
	ReplyTime int64  `json:"reply_time"`
}

// RespVIPInfo .
type RespVIPInfo struct {
	VIPDueMsec int64 `json:"vipDueMsec"`
	VIPStatus  int32 `json:"vipStatus"`
	VIPType    int32 `json:"vipType"`
}

// RespShowInfo .
type RespShowInfo struct {
	Show  bool  `json:"show"`
	State int64 `json:"state"`
}
