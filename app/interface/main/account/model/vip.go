package model

import (
	col "go-common/app/service/main/coupon/model"
	vipv1 "go-common/app/service/main/vip/api"
	vipmol "go-common/app/service/main/vip/model"
	"go-common/library/time"
)

// vip tips.
const (
	PanelPosition int8 = iota + 1
	PgcPosition
)

const (
	// PlatAndroid is int8 for android.
	PlatAndroid = int8(0)
	// PlatIPhone is int8 for iphone.
	PlatIPhone = int8(1)
	// PlatIPad is int8 for ipad.
	PlatIPad = int8(2)
	// PlatWPhone is int8 for wphone.
	PlatWPhone = int8(3)
	// PlatAndroidG is int8 for Android Googleplay.
	PlatAndroidG = int8(4)
	// PlatIPhoneI is int8 for Iphone Global.
	PlatIPhoneI = int8(5)
	// PlatIPadI is int8 for IPAD Global.
	PlatIPadI = int8(6)
	// PlatAndroidTV is int8 for AndroidTV Global.
	PlatAndroidTV = int8(7)
	// PlatAndroidI is int8 for Android Global.
	PlatAndroidI = int8(8)
	// PlatIpadHD is int8 for IpadHD
	PlatIpadHD = int8(9)
	// PlatAndroidB is int8 for Android Blue.
	PlatAndroidB = int8(10)
	// PlatIPhoneB is int8 for Android Blue.
	PlatIPhoneB = int8(11)
)

// resource id .
const (
	ResourceBannerIPhone  = "2850"
	ResourceBannerAndroid = "2857"
	ResourceBannerIPad    = "2864"
	ResourceBuyIPhone     = "2898"
	ResourceBuyAndroid    = "2903"
	ResourceBuyIPad       = "2908"
)

// IsAndroid check plat is android or ipad.
func IsAndroid(plat int8) bool {
	return plat == PlatAndroid || plat == PlatAndroidG || plat == PlatAndroidI || plat == PlatAndroidB
}

// IsIOS check plat is iphone or ipad.
func IsIOS(plat int8) bool {
	return plat == PlatIPad || plat == PlatIPhone || plat == PlatIPadI || plat == PlatIPhoneI
}

// IsIPhone check plat is iphone.
func IsIPhone(plat int8) bool {
	return plat == PlatIPhone || plat == PlatIPhoneI
}

// IsIPad check plat is pad.
func IsIPad(plat int8) bool {
	return plat == PlatIPad || plat == PlatIPadI || plat == PlatIpadHD
}

// IsIPhoneB check plat is iphone_b.
func IsIPhoneB(plat int8) bool {
	return plat == PlatIPhoneB
}

// Plat .
func Plat(mobiApp, device string) int8 {
	switch mobiApp {
	case "iphone":
		if device == "pad" {
			return PlatIPad
		}
		return PlatIPhone
	case "white":
		return PlatIPhone
	case "ipad":
		return PlatIpadHD
	case "android":
		return PlatAndroid
	case "win", "winphone":
		return PlatWPhone
	case "android_G":
		return PlatAndroidG
	case "android_i":
		return PlatAndroidI
	case "iphone_i":
		if device == "pad" {
			return PlatIPadI
		}
		return PlatIPhoneI
	case "ipad_i":
		return PlatIPadI
	case "android_tv":
		return PlatAndroidTV
	case "android_b":
		return PlatAndroidB
	case "iphone_b":
		return PlatIPhoneB
	}
	return PlatIPhone
}

// VIPInfo vip info.
type VIPInfo struct {
	Mid       int64  `json:"mid"`
	Type      int8   `json:"vipType"`
	DueDate   int64  `json:"vipDueDate"`
	DueMsec   int64  `json:"vipSurplusMsec"`
	DueRemark string `json:"dueRemark"`
	Status    int8   `json:"accessStatus"`
	VipStatus int8   `json:"vipStatus"`
}

// TipsReq tips request.
type TipsReq struct {
	Version  int64  `form:"build"`
	Platform string `form:"platform" validate:"required"`
	Position int8   `form:"position" default:"1"`
}

//CodeInfoReq code info request
type CodeInfoReq struct {
	Appkey    string    `form:"appkey" validate:"required"`
	Sign      string    `form:"sign"`
	Ts        time.Time `form:"ts"`
	StartTime time.Time `form:"start_time" validate:"required"`
	EndTime   time.Time `form:"end_time" validate:"required"`
	Cursor    int64     `form:"cursor"`
}

// VipPanelRes .
type VipPanelRes struct {
	Device    string `form:"device"`
	MobiApp   string `form:"mobi_app"`
	Platform  string `form:"platform" default:"pc"`
	SortTP    int8   `form:"sort_type"`
	PanelType string `form:"panel_type" default:"normal"`
	Month     int32  `form:"month"`
	SubType   int32  `form:"order_type"`
	Build     int64  `form:"build"`
}

// ArgVipCoupon req.
type ArgVipCoupon struct {
	ID int64 `form:"id" validate:"required,min=1,gte=1"`
}

// ArgVipCancelPay req.
type ArgVipCancelPay struct {
	CouponToken string `form:"coupon_token" validate:"required"`
}

// coupon cancel explain
const (
	CouponCancelExplain = "解锁成功,请重新选择劵信息"
)

// const for vip
const (
	MobiAppIphone = iota + 1
	MobiAppIpad
	MobiAppPC
	MobiAppANDROID
)

//MobiAppByName .
var MobiAppByName = map[string]int{
	"iphone":  MobiAppIphone,
	"ipad":    MobiAppIpad,
	"pc":      MobiAppPC,
	"android": MobiAppANDROID,
}

// MobiAppPlat .
func MobiAppPlat(mobiApp string) (p int) {
	p = MobiAppByName[mobiApp]
	if p == 0 {
		// def pc.
		p = MobiAppPC
	}
	return
}

// ArgVipPanel arg panel.
type ArgVipPanel struct {
	Device    string `form:"device"`
	Build     int64  `form:"build"`
	MobiApp   string `form:"mobi_app"`
	Platform  string `form:"platform" default:"pc"`
	SortTP    int8   `form:"sort_type"`
	PanelType string `form:"panel_type" default:"normal"`
	Mid       int64
	IP        string
}

// VipPanelResp vip panel resp.
type VipPanelResp struct {
	Vps        []*vipmol.VipPanelInfo          `json:"price_list"`
	CodeSwitch int8                            `json:"code_switch"`
	GiveSwitch int8                            `json:"give_switch"`
	Privileges map[int8]*vipmol.PrivilegesResp `json:"privileges,omitempty"`
	TipInfo    *vipmol.TipsResp                `json:"tip_info,omitempty"`
	UserInfo   *vipmol.VipPanelExplain         `json:"user_info,omitempty"`
}

// VipPanelV8Resp vip panel v8 resp.
type VipPanelV8Resp struct {
	Vps           []*vipmol.VipPanelInfo          `json:"price_list"`
	CouponInfo    *col.CouponAllowancePanelInfo   `json:"coupon_info,omitempty"`
	CouponSwith   int8                            `json:"coupon_switch,omitempty"`
	CodeSwitch    int8                            `json:"code_switch"`
	GiveSwitch    int8                            `json:"give_switch"`
	Privileges    map[int8]*vipmol.PrivilegesResp `json:"privileges,omitempty"`
	TipInfo       *vipmol.TipsResp                `json:"tip_info,omitempty"`
	UserInfo      *vipmol.VipPanelExplain         `json:"user_info,omitempty"`
	AssociateVips []*vipmol.AssociateVipResp      `json:"associate_vips,omitempty"`
}

// VipPanelRespV9 vip panel resp v9.
type VipPanelRespV9 struct {
	Vps           []*vipv1.ModelVipPanelInfo          `json:"price_list,omitempty"`
	Coupon        *vipv1.CouponBySuitIDReply          `json:"coupon,omitempty"`
	CouponSwith   int32                               `json:"coupon_switch"`
	CodeSwitch    int32                               `json:"code_switch"`
	GiveSwitch    int32                               `json:"give_switch"`
	Privileges    map[int32]*vipv1.ModelPrivilegeResp `json:"privileges,omitempty"`
	TipInfo       *vipmol.TipsResp                    `json:"tip_info,omitempty"`
	UserInfo      *vipmol.VipPanelExplain             `json:"user_info,omitempty"`
	AssociateVips []*vipmol.AssociateVipResp          `json:"associate_vips,omitempty"`
}

// ManagerResp manager resp.
type ManagerResp struct {
	JointlyInfo []*vipmol.JointlyResp `json:"jointly_info"`
}

//ArgCreateOrder2 .
type ArgCreateOrder2 struct {
	Month       int32  `form:"months" validate:"required,min=1,gte=1"`
	Platform    string `form:"platform"`
	MobiApp     string `form:"mobi_app"`
	Device      string `form:"device"`
	AppID       int64  `form:"appId"`
	AppSubID    string `form:"appSubId"`
	OrderType   int8   `form:"orderType"`
	Dtype       int8   `form:"dtype"`
	ReturnURL   string `form:"returnUrl"`
	CouponToken string `form:"coupon_token"`
	Bmid        int64  `form:"bmid"`
	PanelType   string `form:"panel_type" default:"normal"`
	Build       int64  `form:"build"`
	IP          string
	Mid         int64
}

//ArgCreateAssociateOrder create asoociate order .
type ArgCreateAssociateOrder struct {
	Month       int32  `form:"months" validate:"required,min=1,gte=1"`
	Platform    string `form:"platform" default:"pc"`
	MobiApp     string `form:"mobi_app"`
	Device      string `form:"device"`
	AppID       int64  `form:"appId"`
	AppSubID    string `form:"appSubId"`
	OrderType   int8   `form:"orderType"`
	Dtype       int8   `form:"dtype"`
	ReturnURL   string `form:"returnUrl"`
	CouponToken string `form:"coupon_token"`
	Bmid        int64  `form:"bmid"`
	PanelType   string `form:"panel_type" default:"normal"`
	Build       int64  `form:"build"`
	IP          string
	Mid         int64
}

// ArgResource struct .
type ArgResource struct {
	MID     int64
	ResIDs  string
	Plat    int8   `form:"plat"`
	Build   int    `form:"build" validate:"required"`
	MobiApp string `form:"mobi_app" validate:"required"`
	Device  string `form:"device"`
	Buvid   string
	IP      string
	Network string `form:"network"`
	Channel string
}

// ArgCouponBySuitID coupon by suit id.
type ArgCouponBySuitID struct {
	Mid       int64
	Sid       int64  `form:"id" validate:"required,min=1,gte=1"`
	Platform  string `form:"platform" default:"pc"`
	MobiApp   string `form:"mobi_app"`
	Device    string `form:"device"`
	PanelType string `form:"panel_type" default:"normal"`
	Build     int64  `form:"build"`
}
