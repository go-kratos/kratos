package model

// ReqTagTop request tag top arg.
type ReqTagTop struct {
	Tid    int64
	TName  string
	Mid    int64
	RealIP string
}

// ReqChannelSquare request channel square.
type ReqChannelSquare struct {
	Mid        int64  `form:"mid"`
	TagNumber  int32  `form:"tag_number" validate:"required,gt=0"`
	OidNumber  int32  `form:"oid_number"`
	Type       int32  `form:"type" validate:"required,gt=0"`
	Buvid      string `form:"buvid"`
	Build      int32  `form:"build"`
	LoginEvent int32  `form:"login_event"`
	DisplayID  int32  `form:"display_id"`
	Plat       int32  `form:"plat" validate:"required,gte=0,lte=8"`
	From       int32  `form:"from" validate:"required,gte=0"`
	RealIP     string
}

// ReqChannelResourceInfos request channel resource infos.
type ReqChannelResourceInfos struct {
	Oids []int64 `form:"oids,split" validate:"required,min=1,max=50,dive,gt=0"`
	Tids []int64 `form:"eids,split" validate:"required,min=1,max=50,dive,gt=0"`
	IDs  []int64 `form:"ids,split" validate:"required,min=1,max=50,dive,gt=0"`
	Bid  int32   `form:"bid" validate:"required,gt=0"`
}

// ReqChannelDetail req channel detail.
type ReqChannelDetail struct {
	Mid    int64  `form:"mid"`
	Tid    int64  `form:"tid"`
	TName  string `form:"tname"`
	From   int32  `form:"from" validate:"required,gte=0"`
	RealIP string
}
