package param

// AddBusAttrParam describe params of add or update business attrs
type AddBusAttrParam struct {
	ID           int64  `form:"id" json:"id"`
	Bid          int64  `form:"bid" json:"bid"`
	Name         string `form:"name" json:"name"`
	DealType     int8   `form:"deal_type" json:"deal_type"`
	ExpireTime   int64  `form:"expire_time" json:"expire_time"`
	AssignType   int8   `form:"assign_type" json:"assign_type"`
	AssignMax    int8   `form:"assign_max" json:"assign_max"`
	GroupType    int8   `form:"group_type" json:"group_type"`
	BusinessName string `form:"business_name" json:"business_name"`
}

// BusAttrButtonSwitch .
type BusAttrButtonSwitch struct {
	Bid    int8  `form:"bid" json:"bid" validate:"required"`
	Index  uint8 `form:"index" json:"index" validate:"min=0,max=7"`
	Switch uint8 `form:"switch" json:"switch" validate:"min=0,max=1"`
}

// BusAttrButtonShortCut .
type BusAttrButtonShortCut struct {
	Bid      int8   `form:"bid" json:"bid" validate:"required"`
	Index    int8   `form:"index" json:"index" validate:"min=0,max=7"`
	ShortCut string `form:"short_cut" json:"short_cut" validate:"required"`
}

// BusAttrExtAPI .
type BusAttrExtAPI struct {
	Bid         int8   `form:"bid" json:"bid" validate:"required"`
	ExternalAPI string `form:"external_api" json:"external_api"`
}

// BlockInfo .
type BlockInfo struct {
	Mid int64 `form:"mid" json:"mid" validate:"required"`
}

// Source .
type Source struct {
	Bid int8 `form:"bid" json:"bid" validate:"required"`
}
