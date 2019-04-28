package search

import (
	"go-common/app/interface/main/tv/model"
)

// ResultResponse def .
type ResultResponse struct {
	Page       int    `json:"page"`
	Pagesize   int    `json:"pagesize"`
	NumResults int    `json:"numResults"`
	NumPages   int    `json:"numPages"`
	Seid       string `json:"seid"`
}

type pageinfo struct {
	Tvpgc *Page `json:"tvpgc"`
	Tvugc *Page `json:"tvugc"`
}

// Page struct .
type Page struct {
	NumResult int `json:"numResults"`
	Total     int `json:"total"`
	Pages     int `json:"pages"`
}

// RespAll def .
type RespAll struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	*ResultResponse
	PageInfo *pageinfo  `json:"pageinfo"`
	Result   *AllResult `json:"result"`
}

// AllResult def .
type AllResult struct {
	Pgc []*PgcResult `json:"tvpgc"`
	Ugc []*UgcResult `json:"tvugc"`
}

// RespPgc def .
type RespPgc struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	*ResultResponse
	Result []*PgcResult `json:"result"`
}

// RespUgc def .
type RespUgc struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	*ResultResponse
	Result []*UgcResult `json:"result"`
}

// UgcResult def .
type UgcResult struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Cover       string `json:"cover"`
	Description string `json:"description"`
	Pubtime     int    `json:"pubtime"`
	Category    int    `json:"category"`
}

// PgcResult def .
type PgcResult struct {
	*UgcResult
	CV         string             `json:"cv"`
	Staff      string             `json:"staff"`
	CornerMark *model.SnVipCorner `json:"cornermark"`
}

// ToCommon transform pgc to common .
func (p *PgcResult) ToCommon() (res *CommonResult) {
	return &CommonResult{
		PgcResult: p,
		Type:      "pgc",
	}
}

// ToCommon transform pgc to common .
func (p *UgcResult) ToCommon() (res *CommonResult) {
	res = &CommonResult{}
	res.PgcResult = &PgcResult{
		UgcResult: p,
	}
	res.Type = "ugc"
	return
}

// CommonResult is the common result for both pgc & ugc .
type CommonResult struct {
	*PgcResult
	Type string `json:"type"`
}

// ReqSearch def .
type ReqSearch struct {
	SearchType string `form:"search_type" validate:"required"`
	Order      string `form:"order"`
	Category   int    `form:"category"`
	Platform   string `form:"platform"  validate:"required"`
	Build      string `form:"build"  validate:"required"`
	MobiAPP    string `form:"mobi_app"`
	Device     string `form:"device"`
	Keyword    string `form:"keyword"  validate:"required"`
	Page       int    `form:"page"  validate:"required,min=1"`
	Pagesize   int    `form:"pagesize"`
}
