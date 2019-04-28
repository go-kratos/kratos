package model

// open change
const (
	OpenChangeNONE int8 = iota
	PointChange
	Recharge
	System
	Active
	ReacquireDeduction
	ActiveCode
	SystemDeduction
)

//OpenChangeMap .
var OpenChangeMap = map[int8]string{
	OpenChangeNONE:     "none",
	PointChange:        "积分兑换",
	Recharge:           "充值开通",
	System:             "系统发放",
	Active:             "活动赠送",
	ReacquireDeduction: "重复领取扣除",
	ActiveCode:         "激活码",
	SystemDeduction:    "系统扣减",
}

// const for vip
const (
	PlatfromIOS = iota + 1
	PlatfromIPAD
	PlatfromPC
	PlatfromANDROID
	PlatfromIPADHD
	PlatfromIOSBLUE
	PlatfromANDROIDBLUE
	PlatfromPUBLIC
	PlatfromAutoRenewServer
	PlatfromANDROIDI //安卓国际版
)

// const for vip
const (
	DeviceIOS = iota + 1
	DeviceIPAD
	DevicePC
	DeviceANDROID
	DeviceIPADHD
	DEVICEIOSBLUE
	DEVICEANDROIDBLUE
	DEVICEPUBLIC
)

// const for vip
const (
	MobiAppIphone = iota + 1
	MobiAppIpad
	MobiAppPC
	MobiAppANDROID
)

//PlatformByName .
var PlatformByName = map[string]int{
	"ios":       PlatfromIOS,
	"ipad":      PlatfromIPAD,
	"pc":        PlatfromPC,
	"android":   PlatfromANDROID,
	"ipadhd":    PlatfromIPADHD,
	"ios_b":     PlatfromIOSBLUE,
	"android_b": PlatfromANDROIDBLUE,
	"public":    PlatfromPUBLIC,
}

//PlatformByCode .
var PlatformByCode = map[int]string{
	PlatfromIOS:     "ios",
	PlatfromIPAD:    "ipad",
	PlatfromPC:      "pc",
	PlatfromANDROID: "android",
}

//MobiAppByName .
var MobiAppByName = map[string]int{
	"iphone":  MobiAppIphone,
	"ipad":    MobiAppIpad,
	"pc":      MobiAppPC,
	"android": MobiAppANDROID,
}

//PayWayName payWay name
var PayWayName = map[int8]string{
	ALIPAY: "支付宝",
	WECHAT: "微信",
	BCION:  "B币",
	BANK:   "银行卡",
	PAYPAL: "paypal",
	IOSPAY: "iospay",
	QPAY:   "qpay",
}

// user discount history enum
const (
	FirstDiscountBuyVip int64 = iota + 1
)

//PayPlatform vip mapping platform
var PayPlatform = map[int]int8{
	DeviceIOS:     2,
	DeviceIPAD:    2,
	DevicePC:      1,
	DeviceANDROID: 1,
}

// vip pay remark
const (
	RemarkBuy  = "充值开通"
	RemarkGift = "好友赠送"
)

// business status
const (
	StatusOpen  = iota //有效
	StatusClose = 1    //无效
)

// business status
const (
	BizTypeIn  = iota //内部
	BizTypeOut = 1    //外部
)

//code status
const (
	CodeUnUse int8 = iota + 1
	CodeUse
	CodeFrozen
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

// user discount
const (
	UnUse int8 = iota
	Used
)

// IsAutoRenewed is auto renewed.
const (
	IsAutoRenewed int32 = 1
)

// bcoin salary status.
const (
	BcoinUnissued int8 = iota
	Grant
)

// vip status.
const (
	Expire int32 = iota
	NotExpired
	Freeze
	Block
)

//batch code status
const (
	Nomal = iota
	OnlyNotVip
)

// batch status
const (
	BatchNormal int8 = iota + 1
	BatchFrozen
)

// tips judge type .
const (
	VersionTypeNone int8 = iota
	VersionMoreThan
	VersionEqual
	VersionLessThan
)

// vip pay type.
const (
	NormalPay int32 = iota
	AutoRenewPay
)

// vip tips.
const (
	PanelPosition int8 = iota + 1
	PgcPosition
)

// switch.
const (
	SwitchClose int8 = iota
	SwitchOpen
)

const (
	// VipUserFirstDiscount 促销类型
	VipUserFirstDiscount = 1
)

// Discount status.
const (
	DiscountNotUse = iota
	DiscountUsed
)

// privilege type.
const (
	AllPrivilege int8 = iota
	OnlyAnnualPrivilege
)

// privilege resources type.
const (
	WebResources = iota
	AppResources
)

// privilege title.
const (
	PrivilegeTitle       = "大会员权益"
	AnnualPrivilegeTitle = "年度大会员权益"
)

// plat arg
const (
	DeviceIapdName  = "pad"
	MobiAppIpadName = "ipad"
)

// pay  service type
const (
	ServiceTypeNormal        = 0
	ServiceTypeInternational = 2
	ServiceTypePublic        = 1
	ServiceTypeAuto          = 7
	ServiceTypeIap           = 100
)

// pay sub type
const (
	PaySubTypeAuto = 1
)

// vip pay type.
const (
	NormalPayType int8 = iota
	AutoRenewPayType
	IapAutoRenewPayType
)

// pay showTitle.
const (
	NormalShowTitle    = "购买大会员"
	AutoRenewShowTitle = "购买大会员连续包月"
)

// vip panel user explain.
const (
	UserNotLoginExplain   = "点击头像登录或注册优惠价开通大会员"
	NotVipExplain         = "你还不是大会员,开通福利多多"
	ExpireVipExplain      = "大会员离你而去了,快来续期吧"
	YYYYDDVipExplain      = "%s到期,购买后有效期将顺延"
	WillExplainVipExplain = "只剩%d天大会员就要离开你而去啦,快来续期吧"
)

// pay param show content.
const (
	ShowContent = "购买%d个月大会员"
)

// pay channel id
const (
	IapPayChannelID = 100
)
