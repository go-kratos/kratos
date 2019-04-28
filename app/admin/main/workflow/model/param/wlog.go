package param

// TargetLogParam log/target/list api param
type TargetLogParam struct {
	Target int64 `form:"target" validate:"required,min=1"`
	Module int8  `form:"module" validate:"required,min=1"`
}

// LogListParam .
type LogListParam struct {
	Business  int64   `form:"business" validate:"required"`
	AdminIDs  []int64 `form:"adminid,split"`
	OIDs      []int64 `form:"oid,split"`
	MIDs      []int64 `form:"up_mid,split"`
	TypeIDs   []int64 `form:"typeids,split"`
	CTimeFrom string  `form:"optctime_from"`
	CTimeTo   string  `form:"optctime_to"`
	Order     string  `form:"order" default:"opt_ctime"`
	Sort      string  `form:"sort_order" default:"desc"`
	PN        int64   `form:"page" default:"1"`
	PS        int64   `form:"ps" default:"50"`
}
