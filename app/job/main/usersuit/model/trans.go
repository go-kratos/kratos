package model

import (
	xtime "go-common/library/time"
)

// PendantGroupTrans pendant group pendant
type PendantGroupTrans struct {
	Gid        int64
	GroupName  string
	Rank       string
	GroupImage string
	IsOnline   int32
	ModifyTime string
}

// PendantTrans pendant trans
type PendantTrans struct {
	Pid           int64
	Name          string
	Image         string
	ImageModel    string
	DisplayExpire int64
	Gid           int64
	Rank          int32
	IsOnline      int32
	ModifyTime    string
}

// PendantPriceTrans pendantPrice
type PendantPriceTrans struct {
	ID     int64
	Pid    int64
	Ecode  string
	Price  int64
	Active int32
	Ctime  string
	Mtime  string
}

// PTransStatus pendant history
type PTransStatus struct {
	ID          int64  `json:"id"`
	Mid         int64  `json:"mid"`
	Pid         int64  `json:"pid"`
	Expire      int64  `json:"expire"`
	IsActivated int64  `json:"is_activated"`
	Mtime       string `json:"modify_time"`
}

// PTransHistory pendant trans pendant info
type PTransHistory struct {
	ID           int64  `json:"id"`
	Mid          int64  `json:"mid"`
	OrderID      string `json:"order_id"`
	PayID        string `json:"pay_id"`
	AppID        int64  `json:"appid"`
	Status       int32  `json:"status"`
	Pid          int64  `json:"pid"`
	TimeLength   int64  `json:"time_length"`
	Cost         string `json:"cost"`
	BuyTime      int64  `json:"buy_time"`
	IsCallback   int32  `json:"is_callback"`
	CallbackTime int64  `json:"callback_time"`
	Mtime        string `json:"modify_time"`
}

// PGTransHistory pendant grant history
type PGTransHistory struct {
	ID             int64
	Mid            int64
	OperatorName   string
	OperatorTime   int64
	OperatorAction string
	OperatorType   int32
	Mtime          string
}

// MedalGroup struct .
type MedalGroup struct {
	GID      int64
	Name     string `json:"group_name"`
	PID      int64  `json:"parent_gid"`
	Rank     int8   `json:"group_rank"`
	IsOnline int8   `json:"is_online"`
	IsDel    int8
	Ctime    xtime.Time
	Mtime    string `json:"modify_time"`
}

// MedalOwner struct db bus trans.
type MedalOwner struct {
	ID          int64
	MID         int64
	NID         int64
	IsActivated int8 `json:"is_activated"`
	Ctime       int64
	Mtime       string `json:"modify_time"`
}

// MedalInfo struct .
type MedalInfo struct {
	NID         int64
	GID         int64
	Name        string
	Description string
	Image       string `json:"image"`
	ImageSmall  string `json:"image_small"`
	Condition   string
	Level       int8
	LevelRank   string `json:"level_rank"`
	Sort        int8
	IsOnline    int8 `json:"is_online"`
	CTime       xtime.Time
	MTime       string `json:"modify_time"`
}
