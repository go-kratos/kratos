package model

// ArgAllowanceList req.
type ArgAllowanceList struct {
	State int8 `form:"state"`
}

// ArgCouponPage req .
type ArgCouponPage struct {
	State int8 `form:"state"`
	Pn    int  `form:"pn"`
	Ps    int  `form:"ps"`
}

// ArgPrizeDraw struct .
type ArgPrizeDraw struct {
	CardType int8 `form:"card_type" validate:"min=0,gte=0,lte=2"`
}
