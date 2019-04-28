package model

// Pager pager
type Pager struct {
	Pn  int   `json:"num"`
	Ps  int   `json:"size"`
	Sum int64 `json:"total"`
}

// ListParser list parser
type ListParser struct {
	Unames string `form:"uname"`
	Bt     string `form:"bt"`
	Et     string `form:"et"`
	Sort   string `form:"sort" default:"desc"`
	Ps     int64  `form:"ps" default:"20"`
	Pn     int64  `form:"pn" default:"1"`
}
