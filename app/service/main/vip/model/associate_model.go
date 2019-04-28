package model

// associate bind state.
const (
	AssociateBindStateNone int32 = iota
	AssociateBindStateNotPurchase
	AssociateBindStatePurchased
)

// associate prize type.
const (
	AssociatePrizeTypeCode int8 = iota + 1
	AssociatePrizeTypeEleBag
)

// associate appid.
const (
	EleAppID = 32
)

// month type.
const (
	UnionOneMonth int32 = 1
	UnionOneYear  int32 = 12
)

// ele vip type
const (
	EleMonthVip int32 = 2
	EleYearVip  int32 = 4
)

// eleme grant remark.
const (
	ElemeGrantRemark = "宅e块联合会员"
)

// EleGrantVipDays eleme grant vip days.
var EleGrantVipDays = map[int32]int64{
	1:  31,
	12: 366,
}
