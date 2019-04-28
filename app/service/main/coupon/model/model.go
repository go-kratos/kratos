package model

// coupon use state.
const (
	UseFaild int8 = iota
	UseSuccess
)

// coupon state.
const (
	NotUsed = iota
	InUse
	Used
	Expire
	Block
)

// coupon state.
const (
	WaitPay = iota
	InPay
	PaySuccess
	PayFaild
)

// max salary count.
const (
	MaxSalaryCount = 100
)

// blance change type
const (
	VipSalary int8 = iota + 1
	SystemAdminSalary
	Consume
	ConsumeFaildBack
)

// coupon type
const (
	CouponVideo = iota + 1
	CouponCartoon
	CouponAllowance
)

//allowance origin
const (
	AllowanceNone = iota
	AllowanceSystemAdmin
	AllowanceBusinessReceive
	AllowanceBusinessNewYear
	AllowanceCodeOpen
)

// batch state
const (
	BatchStateNormal = iota
	BatchStateBlock
)

// coupon disables explains
const (
	CouponHadBlock             = "代金券已被冻结"
	CouponFullAmountDissatisfy = "未达到满额条件"
	CouponNotInUsableTime      = "当前不在有效期内"
	CouponInUse                = "已绑定在其他未支付订单,点击解锁"
	CouponPlatformExplain      = "当前平台不可使用"
	CouponProductExplain       = "当前商品不可使用"
)

// coupon scope explains
const (
	ScopeNoLimit    = "不限使用平台"
	ScopePlatFmt    = "仅限%s端，"
	ScopeProductFmt = "购买%s%s大会员时使用"
)

// coupon send message
const (
	ReceiveMessageTitle = "大会员代金券到账通知"
	ReceiveMessage      = "大会员代金券已到账，快到“我的代金券”看看吧！IOS端需要在网页使用。#{传送门}{\"https://account.bilibili.com/account/big/voucher\"}"
)

// device code
const (
	DeviceIOS int = iota + 1
	DeviceIPAD
	DevicePC
	DeviceANDROID
	DeviceIPADHD
	DeviceIOSBLUE
	DeviceANDROIDBLUE
	DevicePUBLIC
)

// PlatformByCode device name map.
var PlatformByCode = map[int]string{
	DeviceIOS:     "ios",
	DeviceIPAD:    "ipad",
	DevicePC:      "网页",
	DeviceANDROID: "Android",
}

// coupon format
const (
	CouponFullAmountLimit = "满%s元可用"
	CouponAllowanceName   = "大会员代金券"
)

// coupon seleted
const (
	Seleted = 1
)

// allowance change type
const (
	AllowanceSalary int8 = iota + 1
	AllowanceConsume
	AllowanceCancel
	AllowanceConsumeSuccess
	AllowanceConsumeFaild
	AllowanceReceive
)

// allowance notify pay status
const (
	AllowanceUseFaild int8 = iota
	AllowanceUseSuccess
)

// allowance able state
const (
	AllowanceDisables int8 = iota
	AllowanceUsable
)

//PlatformByName .
var PlatformByName = map[string]int{
	"ios":       DeviceIOS,
	"ios_b":     DeviceIOS,
	"ipad":      DeviceIPAD,
	"ipadhd":    DeviceIPAD,
	"pc":        DevicePC,
	"public":    DevicePC,
	"android":   DeviceANDROID,
	"android_b": DeviceANDROID,
}

//PlatfromMapping .
var PlatfromMapping = map[int]int{
	DeviceIPADHD:      DeviceIPAD,
	DeviceIOSBLUE:     DeviceIOS,
	DeviceANDROIDBLUE: DeviceANDROID,
	DevicePUBLIC:      DevicePC,
}

// coupon tip.
const (
	CouponTipNotUse      = "不使用代金券"
	CouponTipChooseOther = "选中其他商品有惊喜"
	CouponTipUse         = "抵扣%.2f元"
	CouponTipInUse       = "有代金券被锁定"
)
