package model

import xtime "go-common/library/time"

// Pager .
type Pager struct {
	Total int64  `json:"total"`
	PN    int    `json:"page"`
	PS    int    `json:"pagesize"`
	Order string `json:"order"`
	Sort  string `json:"sort"`
}

// ArgPendantGroupList .
type ArgPendantGroupList struct {
	GID int64 `form:"gid"`
	PN  int   `form:"pn"`
	PS  int   `form:"ps" validate:"max=100"`
}

// ArgPendantInfo .
type ArgPendantInfo struct {
	PID           int64  `form:"pid"`
	GID           int64  `form:"gid" validate:"required"`
	Name          string `form:"name" validate:"required"`
	Image         string `form:"image"`
	ImageModel    string `form:"image_model"`
	Rank          int16  `form:"rank"`
	Status        int8   `form:"status"`
	IntegralPrice int    `form:"integral_price"` // 积分
	BcoinPrice    int    `form:"bcoin_price"`    // B币
	CoinPrice     int    `form:"coin_price"`     // 硬币
}

// ArgPendantGroup .
type ArgPendantGroup struct {
	GID    int64  `form:"gid"`
	Name   string `form:"name" validate:"required"`
	Rank   int16  `form:"rank"`
	Status int8   `form:"status"`
}

// ArgPendantOrder .
type ArgPendantOrder struct {
	Start  xtime.Time `form:"start_time"`
	End    xtime.Time `form:"end_time"`
	Status int8       `form:"status"`
	PID    int64      `form:"pid"`
	PayID  string     `form:"pay_id"`
	UID    int64      `form:"uid"`
	PN     int        `form:"pn"`
	PS     int        `form:"ps" validate:"max=100"`
}

// ArgPendantPKG .
type ArgPendantPKG struct {
	UID     int64  `form:"uid"  validate:"required"`
	PID     int64  `form:"pid"  validate:"required"`
	Day     int64  `form:"day"  validate:"required"`
	Type    int8   `form:"type"`
	IsMsg   bool   `form:"is_msg"`
	Title   string `form:"title"`
	Content string `form:"content"`
	OID     int64  `form:"oper_id"  validate:"required"`
}

// ArgMedal medal struct .
type ArgMedal struct {
	GID         int64  `json:"gid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	ImageSmall  string `json:"image_small"`
	Condition   string `json:"condition"`
	Level       string `json:"level"`
	LevelRank   string `json:"level_rank"`
	Sort        int    `json:"sort"`
	IsOnline    int    `json:"is_online"`
}
