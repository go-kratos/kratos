package model

import (
	"time"
)

// RealnameStatus is.
type RealnameStatus int8

const (
	// RealnameStatusFalse is.
	RealnameStatusFalse RealnameStatus = 0
	// RealnameStatusTrue is.
	RealnameStatusTrue RealnameStatus = 1
)

// RealnameApplyStatus is.
type RealnameApplyStatus int8

const (
	// RealnameApplyStatusPending is.
	RealnameApplyStatusPending RealnameApplyStatus = iota
	// RealnameApplyStatusPass is.
	RealnameApplyStatusPass
	// RealnameApplyStatusBack is.
	RealnameApplyStatusBack
	// RealnameApplyStatusNone is.
	RealnameApplyStatusNone
)

// IsPass return is apply passed
func (r RealnameApplyStatus) IsPass() bool {
	switch r {
	case RealnameApplyStatusPass:
		return true
	default:
		return false
	}
}

// RealnameChannel is
type RealnameChannel int8

// RealnameChannel enum
const (
	RealnameChannelMain RealnameChannel = iota
	RealnameChannelAlipay
)

// RealnameApplyStatusInfo is.
type RealnameApplyStatusInfo struct {
	Status   RealnameApplyStatus `json:"status"`
	Remark   string              `json:"remark"`
	Realname string              `json:"realname"`
	Card     string              `json:"card"`
}

// RealnameCacheInfo model in cache
type RealnameCacheInfo struct {
	*RealnameInfo
	RealCard string `json:"real_card"`
}

// RealnameBrief is.
type RealnameBrief struct {
	Realname string         `json:"realname"`
	Card     string         `json:"card"`
	CardType int            `json:"card_type"`
	Status   RealnameStatus `json:"status"`
}

// RealnameInfo is.
type RealnameInfo struct {
	ID       int64               `json:"id"`
	MID      int64               `json:"mid"`
	Channel  RealnameChannel     `json:"channel"`
	Realname string              `json:"realname"`
	Country  int                 `json:"country"`
	CardType int                 `json:"card_type"`
	Card     string              `json:"card"`
	CardMD5  string              `json:"card_md5"`
	Status   RealnameApplyStatus `json:"status"`
	Reason   string              `json:"reason"`
	CTime    time.Time           `json:"ctime"`
	MTime    time.Time           `json:"mtime"`
}

// RealnameDetail is.
type RealnameDetail struct {
	*RealnameBrief
	Gender  string `json:"gender"`
	HandIMG string `json:"hand_img"`
}

// RealnameApply is.
type RealnameApply struct {
	ID           int64               `json:"id"`
	MID          int64               `json:"mid"`
	Realname     string              `json:"realname"`
	Country      int16               `json:"country"`
	CardType     int8                `json:"card_type"`
	CardNum      string              `json:"card_num"`
	CardMD5      string              `json:"card_md5"`
	HandIMG      int                 `json:"hand_img"`
	FrontIMG     int                 `json:"front_img"`
	BackIMG      int                 `json:"back_img"`
	Status       RealnameApplyStatus `json:"status"`
	Operator     string              `json:"operator"`
	OperatorID   int64               `json:"operator_id"`
	OperatorTime time.Time           `json:"operator_time"`
	Remark       string              `json:"remark"`
	RemarkStatus int8                `json:"remark_status"`
	CTime        time.Time           `json:"ctime"`
	MTime        time.Time           `json:"mtime"`
}

// IsPass is.
func (r *RealnameApply) IsPass() bool {
	switch r.Status {
	case RealnameApplyStatusPass:
		return true
	default:
		return false
	}
}

// RealnameApplyImage is.
type RealnameApplyImage struct {
	ID      int64
	IMGData string
	CTime   time.Time
	MTime   time.Time
}

// RealnameCapture is.
type RealnameCapture struct {
	Code      int
	CodeCTime time.Time
	Times     []time.Time
}

// RealnameAlipayApply is
type RealnameAlipayApply struct {
	ID       int64               `json:"id"`
	MID      int64               `json:"mid"`
	Realname string              `json:"realname"`
	Card     string              `json:"card"`
	IMG      string              `json:"img"`
	Status   RealnameApplyStatus `json:"status"`
	Reason   string              `json:"reason"`
	Bizno    string              `json:"bizno"`
	CTime    time.Time           `json:"ctime"`
	MTime    time.Time           `json:"mtime"`
}

// IsPass is.
func (r *RealnameAlipayApply) IsPass() bool {
	switch r.Status {
	case RealnameApplyStatusPass:
		return true
	default:
		return false
	}
}

// RealnameAlipayInfo is
type RealnameAlipayInfo struct {
	Bizno string
}

const (
	// RealnameCountryChina is.
	RealnameCountryChina = 0
	// RealnameCardTypeIdentity is.
	RealnameCardTypeIdentity = 0
)

// RealnameAdultType is.
type RealnameAdultType uint8

const (
	// RealnameAdultTypeFalse is.
	RealnameAdultTypeFalse RealnameAdultType = iota // 未成年
	// RealnameAdultTypeTrue is.
	RealnameAdultTypeTrue // 已成年
	//RealnameAdultTypeUnknown is.
	RealnameAdultTypeUnknown // 未知(未绑定身份证)
)

// http param

// ParamRealnameCheck is.
type ParamRealnameCheck struct {
	MID      int64  `form:"mid" validate:"required"`
	CardType int8   `form:"card_type" default:"-1"`
	CardCode string `form:"card_code" validate:"required"`
}

// ParamRealnameSyncImage is.
type ParamRealnameSyncImage struct {
	Data string `form:"data" validate:"required"`
}

// ParamRealnameTelCaptureCheck is.
type ParamRealnameTelCaptureCheck struct {
	MID     int64 `form:"mid" validate:"required"`
	Capture int   `form:"capture" validate:"required"`
}
