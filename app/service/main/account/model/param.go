package model

// ParamMid is.
type ParamMid struct {
	Mid int64 `form:"mid" validate:"gt=0,required"`
}

// ParamMids is.
type ParamMids struct {
	Mids []int64 `form:"mids,split" validate:"gt=0,dive,gt=0"`
}

// ParamNames is.
type ParamNames struct {
	Names []string `form:"names,split" validate:"gt=0,dive,gt=0"`
}

// ParamModify is.
type ParamModify struct {
	Mid          int64  `form:"mid" validate:"gt=0,required"`
	ModifiedAttr string `form:"modifiedAttr" validate:"gt=0,required"`
}

// ParamMsg is.
type ParamMsg struct {
	// by notify
	Msg string `form:"msg"`
}
