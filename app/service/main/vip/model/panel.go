package model

import (
	"fmt"
	"math"
	"strconv"

	colapi "go-common/app/service/main/coupon/api"
	col "go-common/app/service/main/coupon/model"
	"go-common/library/time"
)

// vip_price_config suit_type
const (
	AllUser int8 = iota
	OldVIP
	NewVIP
	OldSubVIP
	NewSubVIP
	OldPackVIP
	NewPackVIP
)

// order type
const (
	NoRenew int8 = iota
	OtherRenew
	IOSRenew
)

// order type by month for vip_user_discount_history table
const (
	OneMonthSub int8 = iota + 1
	ThreeMonthSub
	OneYearSub
)

// const month
const (
	OneMonth   = int8(1)
	ThreeMonth = int8(3)
	OneYear    = int8(12)
)

// const vip_price_config beforeSuitType
const (
	All int8 = iota
	VIP
	Sub
	Pack
)

// const panel month sort
const (
	PanelMonthDESC int8 = iota
	PanelMonthASC
)

// const PanelType
const (
	PanelTypeNormal = "normal"
	PanelTypeFriend = "friend"
	PanelTypeCheck  = "check"
	PanelTypeEle    = "ele"
)

const (
	// PlatVipPriceConfigOther 其他平台
	PlatVipPriceConfigOther int64 = iota + 1
	// PlatVipPriceConfigIOS IOS平台
	PlatVipPriceConfigIOS
	// PlatVipPriceConfigIPADHD ipad hd平台
	PlatVipPriceConfigIPADHD
	// PlatVipPriceConfigFriendsGift 好友赠送
	PlatVipPriceConfigFriendsGift
	// PlatVipPriceConfigInternational 安卓国际版
	PlatVipPriceConfigInternational
	// PlatVipPriceConfigIphoneB iphone蓝版
	PlatVipPriceConfigIphoneB
	// PlatVipPriceConfigCheck 审核态价格
	PlatVipPriceConfigCheck = 20
)

// const select
const (
	PanelNotSelected int32 = iota
	PanelSelected
)

// VipPriceConfig price config.
type VipPriceConfig struct {
	ID          int64     `json:"id"`
	Plat        int64     `json:"platform"`
	PdName      string    `json:"product_name"`
	PdID        string    `json:"product_id"`
	SuitType    int8      `json:"suit_type"`
	TopSuitType int8      `json:"-"`
	Month       int16     `json:"month"`
	SubType     int8      `json:"sub_type"`
	OPrice      float64   `json:"original_price"`
	DPrice      float64   `json:"discount_price"`
	Selected    int32     `json:"selected"`
	Remark      string    `json:"remark"`
	Status      int8      `json:"status"`
	Forever     bool      `json:"-"`
	Operator    string    `json:"operator"`
	OpID        int64     `json:"oper_id"`
	Superscript string    `json:"superscript"`
	StartBuild  int64     `json:"start_build"`
	EndBuild    int64     `json:"end_build"`
	PanelType   string    `json:"panel_type"`
	CTime       time.Time `json:"ctime"`
	MTime       time.Time `json:"mtime"`
}

// VipPirceResp vip pirce resp.
type VipPirceResp struct {
	Vps         []*VipPanelInfo               `json:"price_list"`
	CouponInfo  *col.CouponAllowancePanelInfo `json:"coupon_info"`
	CouponSwith int8                          `json:"coupon_switch"`
	CodeSwitch  int8                          `json:"code_switch"`
	GiveSwitch  int8                          `json:"give_switch"`
	ExistCoupon int8                          `json:"exist_coupon"`
	Privileges  *PrivilegesResp               `json:"privileges"`
}

// VipPirceResp5 vip pirce resp.
type VipPirceResp5 struct {
	Vps         []*VipPanelInfo               `json:"price_list"`
	CouponInfo  *col.CouponAllowancePanelInfo `json:"coupon_info"`
	CouponSwith int8                          `json:"coupon_switch"`
	CodeSwitch  int8                          `json:"code_switch"`
	GiveSwitch  int8                          `json:"give_switch"`
	Privileges  map[int8]*PrivilegesResp      `json:"privileges"`
}

// VipPirceRespV9 vip pirce resp v9.
type VipPirceRespV9 struct {
	Vps         []*VipPanelInfo                      `json:"price_list"`
	Coupon      *colapi.UsableAllowanceCouponV2Reply `json:"coupon"`
	CouponSwith int8                                 `json:"coupon_switch"`
	CodeSwitch  int8                                 `json:"code_switch"`
	GiveSwitch  int8                                 `json:"give_switch"`
	Privileges  map[int8]*PrivilegesResp             `json:"privileges"`
}

// VipDPriceConfig price discount config.
type VipDPriceConfig struct {
	ID         int64     `json:"id"`
	PdID       string    `json:"product_id"`
	DPrice     float64   `json:"discount_price"`
	STime      time.Time `json:"stime"`
	ETime      time.Time `json:"etime"`
	Remark     string    `json:"remark"`
	Operator   string    `json:"operator"`
	OpID       int64     `json:"oper_id"`
	CTime      time.Time `json:"ctime"`
	MTime      time.Time `json:"mtime"`
	FirstPrice float64   `json:"first_price"`
}

// DoTopSuitType .
func (vpc *VipPriceConfig) DoTopSuitType() {
	switch vpc.SuitType {
	case OldPackVIP, NewPackVIP:
		vpc.TopSuitType = Pack
	case OldSubVIP, NewSubVIP:
		vpc.TopSuitType = Sub
	case OldVIP, NewVIP:
		vpc.TopSuitType = VIP
	case AllUser:
		vpc.TopSuitType = All
	}
}

// DoCheckRealPrice ,
func (vpc *VipPriceConfig) DoCheckRealPrice(mvp map[int64]*VipDPriceConfig) {
	if vp, ok := mvp[vpc.ID]; ok {
		vpc.PdID = vp.PdID
		vpc.DPrice = vp.DPrice
		vpc.Remark = vp.Remark
		if vp.FirstPrice > 0 && vpc.SubType == AutoRenew {
			vpc.DPrice = vp.FirstPrice
		}
	}
	if vpc.DPrice == 0 {
		vpc.DPrice = vpc.OPrice
	}
}

// DoSubMonthKey .
func (vpc *VipPriceConfig) DoSubMonthKey() string {
	return fmt.Sprintf("%d%d", vpc.Month, vpc.SubType)
}

// FormatRate .
func (vpc *VipPriceConfig) FormatRate() string {
	if vpc.DPrice == 0 {
		return ""
	}
	if vpc.DPrice/vpc.OPrice == 1 {
		return ""
	}
	return strconv.FormatFloat(math.Floor((vpc.DPrice/vpc.OPrice)*100)/10, 'f', -1, 64) + "折"
}

// DoPayOrderTypeKey .
func (po *PayOrder) DoPayOrderTypeKey() string {
	if po.OrderType == IOSRenew {
		po.OrderType = OtherRenew
	}
	return fmt.Sprintf("%d%d", po.BuyMonths, po.OrderType)
}

// IsSub .
func (po *PayOrder) IsSub() bool {
	return po.OrderType == OtherRenew || po.OrderType == IOSRenew
}

// VipPirce vip pirce.
type VipPirce struct {
	Panel  *VipPanelInfo            `json:"pirce_info"`
	Coupon *col.CouponAllowanceInfo `json:"coupon_info"`
}

// VipPanelExplain vip panel explain.
type VipPanelExplain struct {
	BackgroundURL string `json:"background_url"`
	Explain       string `json:"user_explain"`
}

// FilterBuild filter price build .
func (vpc *VipPriceConfig) FilterBuild(build int64) bool {
	if (vpc.StartBuild != 0 && vpc.StartBuild > build) || (vpc.EndBuild != 0 && vpc.EndBuild < build) {
		return false
	}
	return true
}

// ArgProductLimit args product limit.
type ArgProductLimit struct {
	Mid       int64
	Months    int32
	PanelType string
}
