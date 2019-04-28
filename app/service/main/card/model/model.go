package model

// card type
const (
	CardTypeNone int32 = iota
	CardTypeVip
	CardTypeFree
)

// card is hot
const (
	CardNotHot int32 = iota
	CardIsHot
)

// CardTypeNameMap card name map.
var CardTypeNameMap = map[int32]string{
	CardTypeNone: "",
	CardTypeVip:  "大会员",
	CardTypeFree: "免费",
}
