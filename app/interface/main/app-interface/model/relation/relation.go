package relation

import (
	accv1 "go-common/app/service/main/account/api"
	relation "go-common/app/service/main/relation/model"
)

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
	Uname          string             `json:"uname"`
	Face           string             `json:"face"`
	Sign           string             `json:"sign"`
	OfficialVerify accv1.OfficialInfo `json:"official_verify"`
	Vip            Vip                `json:"vip"`
	Live           int                `json:"live"`
}

type Tag struct {
	Mid            int64              `json:"mid"`
	Uname          string             `json:"uname"`
	Face           string             `json:"face"`
	Sign           string             `json:"sign"`
	OfficialVerify accv1.OfficialInfo `json:"official_verify"`
	Vip            Vip                `json:"vip"`
	Live           int                `json:"live"`
}

// ByMTime implements sort.Interface for []model.Following based on the MTime field.
type ByMTime []*relation.Following

func (mt ByMTime) Len() int           { return len(mt) }
func (mt ByMTime) Swap(i, j int)      { mt[i], mt[j] = mt[j], mt[i] }
func (mt ByMTime) Less(i, j int) bool { return mt[i].MTime < mt[j].MTime }
