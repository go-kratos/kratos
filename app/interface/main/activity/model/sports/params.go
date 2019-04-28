package sports

// ParamQq qq param
type ParamQq struct {
	Tp    int64  `form:"tp"`
	Route string `form:"route"`
}

// ParamNews qq news param.
type ParamNews struct {
	Route string `form:"route" validate:"required"`
}
