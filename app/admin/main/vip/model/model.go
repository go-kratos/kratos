package model

// tips status.
const (
	WaitShowTips = iota + 1
	EffectiveTips
	ExpireTips
)

// PlatformByCode .
var PlatformByCode = map[int]string{
	DeviceIOS:     "ios",
	DeviceIPAD:    "ipad",
	DevicePC:      "pc",
	DeviceANDROID: "android",
}

// const for vip
const (
	DeviceIOS = iota + 1
	DeviceIPAD
	DevicePC
	DeviceANDROID
)

// push progress status
const (
	NotStart = iota + 1
	Starting
	Started
)

// const .
const (
	UnDisable = iota
	Disable
)

// const .
const (
	Normal = iota + 1
	Fail
)

// tips judge type .
const (
	VersionTypeNone int8 = iota
	VersionMoreThan
	VersionEqual
	VersionLessThan
)

// Delete state
const (
	Delete = 1
)

// PageInfo common page info.
type PageInfo struct {
	Count       int         `json:"count"`
	CurrentPage int         `json:"currentPage,omitempty"`
	Item        interface{} `json:"item"`
}

// UserChangeHistoryReq user change history request.
type UserChangeHistoryReq struct {
	Mid             int64  `form:"mid"`
	ChangeType      int8   `form:"change_type"`
	StartChangeTime int64  `form:"startchangetime"`
	EndChangeTime   int64  `form:"endchangetime"`
	BatchID         int64  `form:"batch_id"`
	RelationID      string `form:"relation_id"`
	Pn              int    `form:"pn"`
	Ps              int    `form:"ps"`
}

// PushPlatformMap .
var PushPlatformMap = map[string]string{
	"1": "Android",
	"2": "iPhone",
	"3": "iPad",
}

// PushPlatformNameMap .
var PushPlatformNameMap = map[string]string{
	"Android": "1",
	"iPhone":  "2",
	"iPad":    "3",
}

// ConditionMap .
var ConditionMap = map[string]string{
	"gte": ">=",
	"lte": "<=",
	"eq":  "=",
	"neq": "!=",
}

// ConditionNameMap .
var ConditionNameMap = map[string]string{
	">=": "gte",
	"<=": "lte",
	"=":  "eq",
	"!=": "neq",
}

// privilege resources state.
const (
	DisablePrivilege = iota
	NormalPrivilege
)

// privilege resources type.
const (
	WebResources = iota
	AppResources
)

// jointly state.
const (
	WillEffect int8 = iota + 1
	Effect
	LoseEffect
)

// order type
const (
	NormalOrder = iota
	AutoOrder
	IAPAutoOrder
)

// pay order status.
const (
	PAYING = iota + 1
	SUCCESS
	FAILED
	SIGN
	UNSIGN
	REFUNDING
	REFUNDED
)

// order type
const (
	General int8 = iota
	AutoRenew
)
