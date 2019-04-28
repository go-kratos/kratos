package model

import (
	"go-common/library/time"
)

// PointInfo .
type PointInfo struct {
	Mid          int64 `json:"mid"`
	PointBalance int64 `json:"pointBalance"`
	Ver          int64 `json:"ver"`
}

//PointHistory  point history
type PointHistory struct {
	ID           int64     `json:"id"`
	Mid          int64     `json:"mid"`
	Point        int64     `json:"point"`
	OrderID      string    `json:"orderID"`
	ChangeType   int       `json:"changeType"`
	ChangeTime   time.Time `json:"changeTime"`
	RelationID   string    `json:"relationID"`
	PointBalance int64     `json:"pointBalance"`
	Remark       string    `json:"remark"`
	Operator     string    `json:"operator"`
}

//OldPointHistory  point history
type OldPointHistory struct {
	ID           int64  `json:"id"`
	Mid          int64  `json:"mid"`
	Point        int64  `json:"point"`
	OrderID      string `json:"orderID"`
	ChangeType   int    `json:"changeType"`
	ChangeTime   int64  `json:"changeTime"`
	RelationID   string `json:"relationID"`
	PointBalance int64  `json:"pointBalance"`
	Remark       string `json:"remark"`
	Operator     string `json:"operator"`
}

//PointExchangePrice .
type PointExchangePrice struct {
	ID             int    `json:"id"`
	OriginPoint    int    `json:"originPoint"`
	CurrentPoint   int    `json:"currentPoint"`
	Month          int    `json:"month"`
	PromotionTip   string `json:"promotionTip"`
	PromotionColor string `json:"promotionColor"`
	OperatorID     string `json:"operatorId"`
}

//HandlerVip vip handler
type HandlerVip struct {
	Days   int
	Months int
	Mid    int
	Type   int
}

// VipPointConf vip point conf.
type VipPointConf struct {
	ID       int64     `json:"id"`
	AppID    int64     `json:"app_id"`
	Point    int64     `json:"point"`
	Operator string    `json:"operator"`
	Ctime    time.Time `json:"ctime"`
	Mtime    time.Time `json:"mtime"`
}
