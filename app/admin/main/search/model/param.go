package model

// Pager .
type Pager struct {
	Pn int `form:"pn" validate:"min=1" default:"1"`
	Ps int `form:"ps" validate:"min=1" default:"10"`
}

// ParamMngBusiness .
type ParamMngBusiness struct {
	ID    int64  `form:"id"`
	Name  string `form:"name"`
	Desc  string `form:"desc"`
	Apps  string `form:"apps"`
	IsJob bool   `form:"is_job"`
	Pager
}

// ParamMngBusinessApp .
type ParamMngBusinessApp struct {
	Business string `form:"business"`
	App      string `form:"app"`
	IsJob    bool   `form:"is_job"`
	IncrWay  string `form:"incr_way"`
	IncrOpen bool   `form:"incr_open"`
}

// ParamMngAsset .
type ParamMngAsset struct {
	ID     int64  `form:"id"`
	Type   int    `form:"type"`
	Name   string `form:"name"`
	Config string `form:"config"`
	Desc   string `form:"desc"`
	Pager
}
