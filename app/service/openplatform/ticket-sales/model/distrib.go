package model

//DistOrderArg 分销订单同步入参
type DistOrderArg struct {
	Oid       uint64 `form:"oid" validate:"required"`
	CmAmount  uint64 `form:"cm_amount" validate:"min=1,required"`
	CmMethod  int64  `form:"cm_method" validate:"min=1,required"`
	CmPrice   uint64 `form:"cm_price" validate:"min=1,required"`
	Duid      uint64 `form:"dist_user" validate:"min=0"`
	Stat      int64  `form:"status" validate:"min=0"`
	Pid       uint64 `form:"pid" validate:"min=0"`
	Count     uint64 `form:"count" validate:"min=1,required"`
	Sid       uint64 `form:"sid" validate:"min=0"`
	Type      int64  `form:"type" validate:"min=0"`
	RefStat   int64  `form:"refund_status" validate:"min=0"`
	PayAmount uint64 `form:"payment_amount" validate:"min=0"`
	Serial    string `form:"serial" validate:"required"`
}

//DistOrderGetArg 分销订单查询入参
type DistOrderGetArg struct {
	Oid uint64 `form:"oid" validate:"required"`
}
