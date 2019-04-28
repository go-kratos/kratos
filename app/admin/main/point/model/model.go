package model

// ArgPointHistory .
type ArgPointHistory struct {
	Mid             int64  `form:"mid"`
	ChangeType      int64  `form:"change_type"`
	StartChangeTime int64  `form:"begin_time"`
	EndChangeTime   int64  `form:"end_time"`
	BatchID         string `form:"batch_id"`
	RelationID      string `form:"relation_id"`
	PN              int64  `form:"pn" default:"1"`
	PS              int64  `form:"ps" default:"50"`
}

// ArgID .
type ArgID struct {
	ID int64 `form:"id" validate:"required,min=1,gte=1"`
}

// PageInfo common page info.
type PageInfo struct {
	Count       int         `json:"count"`
	CurrentPage int         `json:"currentPage,omitempty"`
	Item        interface{} `json:"item"`
}

// point add suc.
const (
	PointAddSuc = 1
)

// point change type
const (
	ExchangeVip            = iota + 1
	Charge                 //充电
	Contract               //承包
	PointSystem            //系统发放
	FYMReward              //分院帽奖励
	ExchangePendant        //兑换挂件
	MJActive               //萌节活动
	ReAcquirePointDedution //重复领取
)

// ArgPoint .
type ArgPoint struct {
	Mid    int64  `form:"mid" validate:"required,min=1,gte=1"`
	Point  int64  `form:"point"`
	Remark string `form:"remark"`
}
