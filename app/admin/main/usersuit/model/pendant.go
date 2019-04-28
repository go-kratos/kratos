package model

import xtime "go-common/library/time"

// const .
const (
	// pendant type
	PendantCoinPrice     = int8(0)
	PendantBCoinPrice    = int8(1)
	PendantIntegralPrice = int8(2)
	// order plat
	PendantOrderPlatDefault = int8(-1)
	PendantOrderPlatPCAndH5 = int8(0)
	PendantOrderPlatPhone   = int8(1)
	// pkg status
	PendantPKGInvalid = int8(0)
	PendantPKGValid   = int8(1)
	PendantPKGOnEquip = int8(2)
	// pendant add style
	PendantAddStyleDay  = int8(1)
	PendantAddStyleDate = int8(2)

	// sourceType
	PendantSourceTypeAdmin = int8(1)
	PendantSourceTypePGC   = int8(2)
)

// var .
var (
	PriceTypes = []int8{PendantCoinPrice, PendantBCoinPrice, PendantIntegralPrice}
)

// PendantGroup .
type PendantGroup struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Rank   int16  `json:"rank"`
	Status int8   `json:"status"`
}

// PendantPrice .
type PendantPrice struct {
	PID   int64 `json:"pid"`
	TP    int8  `json:"type"`
	Price int   `json:"price"`
}

// PendantInfo .
type PendantInfo struct {
	ID         int64           `json:"id"`
	Name       string          `json:"name"`
	Image      string          `json:"image"`
	ImageModel string          `json:"image_model"`
	Status     int8            `json:"status"`
	Rank       int16           `json:"rank"`
	GID        int64           `json:"gid"`
	GroupName  string          `json:"group_name"`
	GroupRank  int16           `json:"group_rank"`
	Prices     []*PendantPrice `json:"prices"`
}

// PendantGroupRef .
type PendantGroupRef struct {
	GID int64 `json:"gid"`
	PID int64 `json:"pid"`
}

// PendantOrder .
type PendantOrder struct {
	BuyTime    int64  `json:"buy_time"`
	OrderID    string `json:"order_id"`
	PayID      string `json:"pay_id"`
	UID        int64  `json:"uid"`
	PID        int64  `json:"-"`
	PName      string `json:"pendant_name"`
	TimeLength int64  `json:"time_length"`
	Cost       string `json:"cost"`
	PayType    int8   `json:"pay_type"`
	Status     int8   `json:"status"`
	AppID      int8   `json:"appid"`
	Platform   string `json:"platform"`
}

// CoverToPlatform .
func (p *PendantOrder) CoverToPlatform() {
	switch p.AppID {
	case PendantOrderPlatDefault:
		p.Platform = "默认"
	case PendantOrderPlatPCAndH5:
		p.Platform = "PC/H5"
	case PendantOrderPlatPhone:
		p.Platform = "手机客户端"
	}
}

// PendantPKG .
type PendantPKG struct {
	ID      int64 `json:"id"`
	UID     int64 `json:"uid"`
	PID     int64 `json:"pid"`
	Expires int64 `json:"expires"`
	TP      int8  `json:"type"`
	Status  int8  `json:"status"`
	IsVip   int8  `json:"is_vip"`
}

// PendantOperLog .
type PendantOperLog struct {
	OID        int64      `json:"oper_id"`
	Action     string     `json:"action"`
	CTime      xtime.Time `json:"ctime"`
	MTime      xtime.Time `json:"mtime"`
	OperName   string     `json:"oper_name"`
	UID        int64      `json:"uid"`
	PID        int64      `json:"pid"`
	SourceType int8       `json:"source_type"`
}

// BulidPendantPrice .
func (pp *PendantPrice) BulidPendantPrice(arg *ArgPendantInfo, tp int8) {
	switch tp {
	case PendantCoinPrice:
		pp.Price = arg.CoinPrice
	case PendantBCoinPrice:
		pp.Price = arg.BcoinPrice
	case PendantIntegralPrice:
		pp.Price = arg.IntegralPrice
	}
	pp.TP = tp
	pp.PID = arg.PID
}
