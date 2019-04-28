package model

// DmDateParams .
type DmDateParams struct {
	Bsp       *BasicSearchParams
	Oid       int64  `form:"oid" params:"oid" default:"-1"`
	Month     string `form:"month" params:"month" default:""`
	MonthFrom string `form:"month_from" params:"month_from" default:""`
	MonthTo   string `form:"month_to" params:"month_to" default:""`
}
