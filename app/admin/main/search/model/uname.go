package model

// UNameParams search params.
type UNameParams struct {
	Bsp   *BasicSearchParams
	MIds  []int64 `form:"mids,split" params:"mids"`
	Sex   int64   `form:"sex" params:"sex" default:"-1"`
	Ranks []int64 `form:"ranks,split" params:"ranks"`
}
