package model

import (
	"fmt"
	"strconv"

	"go-common/library/database/elastic"
	"go-common/library/time"
)

const (
	_TimeFormat  = ("2006-01-02 15:04:05")
	_mtimeDesc   = "2"
	_pubtimeDesc = "1"
	_pagesize    = 20
)

// EPRes def.
type EPRes struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	STitle   string `json:"season_title" gorm:"column:stitle"`
	Category int    `json:"category"`
	SeasonID int    `json:"season_id"`
	EPID     int    `json:"epid" gorm:"column:epid"`
	State    int    `json:"state"`
	Valid    int    `json:"valid"`
	Reason   string `json:"reason"`
}

// EPResDB defines the result structure of ep audit
type EPResDB struct {
	EPRes
	InjectTime time.Time `json:"inject_time" gorm:"column:inject_time"`
	CTime      time.Time `json:"ctime" gorm:"column:ctime"`
}

// EPResItem def.
type EPResItem struct {
	EPRes
	InjectTime string `json:"inject_time"`
	CTime      string `json:"ctime"`
}

// ToItem def.
func (v *EPResDB) ToItem() *EPResItem {
	res := &EPResItem{
		EPRes: v.EPRes,
		CTime: v.CTime.Time().Format(_TimeFormat),
	}
	switch v.State {
	case 3: // passed
		res.State = 1
	case 4: // rejected
		res.State = 2
	default:
		res.State = 0
	}
	if v.InjectTime > 0 {
		res.InjectTime = v.InjectTime.Time().Format(_TimeFormat)
	}
	return res
}

// Page represents the standard page structure
type Page struct {
	Num   int `json:"num"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

// SnCount def.
type SnCount struct {
	Count int
}

// TableName tv_content
func (*EPResDB) TableName() string {
	return "tv_content"
}

// EPResultPager def.
type EPResultPager struct {
	Items []*EPResItem `json:"items"`
	Page  *Page        `json:"page"`
}

// SeasonRes def.
type SeasonRes struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Check    int    `json:"check"`
	Category int    `json:"category"`
	Valid    int    `json:"valid"`
	Reason   string `json:"reason"`
}

// SeasonResDB defines the result structure of ep audit
type SeasonResDB struct {
	SeasonRes
	InjectTime time.Time `json:"inject_time" gorm:"column:inject_time"`
	CTime      time.Time `json:"ctime" gorm:"column:ctime"`
}

// SeasonResItem def.
type SeasonResItem struct {
	SeasonRes
	InjectTime string `json:"inject_time"`
	CTime      string `json:"ctime"`
}

// ToItem Transforms to item
func (v *SeasonResDB) ToItem() *SeasonResItem {
	res := &SeasonResItem{
		CTime:     v.CTime.Time().Format(_TimeFormat),
		SeasonRes: v.SeasonRes,
	}
	switch v.Check {
	case 0: // reject
		res.Check = 2
	case 1: // passed
		res.Check = 1
	default:
		res.Check = 0
	}
	if v.InjectTime > 0 {
		res.InjectTime = v.InjectTime.Time().Format(_TimeFormat)
	}
	return res
}

// TableName tv_content
func (*SeasonResDB) TableName() string {
	return "tv_ep_season"
}

// SeasonResultPager def.
type SeasonResultPager struct {
	Items []*SeasonResItem `json:"items"`
	Page  *Page            `json:"page"`
}

// ReqArcCons def.
type ReqArcCons struct {
	Order     int    `form:"order" json:"order;Min(1)" default:"1"` // 1 default, desc; 2 asc
	FirstCat  int32  `form:"first_cat"`
	SecondCat int32  `form:"second_cat"`
	Status    string `form:"status"`
	Title     string `form:"title"`
	AVID      int64  `form:"avid"`
	Pn        int    `form:"pn" default:"1"`
}

// ReqVideoCons def.
type ReqVideoCons struct {
	AVID   int64  `form:"avid" validate:"required"`
	Order  int    `form:"order" json:"order;Min(1)" default:"1"` // 1 default, desc; 2 asc
	Title  string `form:"title"`
	CID    int64  `form:"cid"`
	Status string `form:"status"`
	Pn     int    `form:"pn"`
}

// PageCfg def.
type PageCfg struct {
	Pn int `form:"pn" default:"1" json:"pn"`
	Ps int `form:"ps" default:"20" json:"ps"`
}

// ReqArcES def.
type ReqArcES struct {
	AID          string
	Mids         []int64
	Title        string
	Typeids      []int32
	Valid        string
	Result       string
	MtimeOrder   string
	PubtimeOrder string
	PageCfg
}

// MtimeSort def.
func (v *ReqArcES) MtimeSort() string {
	if v.MtimeOrder == _mtimeDesc {
		return elastic.OrderDesc
	}
	return elastic.OrderAsc
}

// PubtimeSort def.
func (v *ReqArcES) PubtimeSort() string {
	if v.PubtimeOrder == _pubtimeDesc {
		return elastic.OrderDesc
	}
	return elastic.OrderAsc
}

// FromArcListParam build
func (v *ReqArcES) FromArcListParam(param *ArcListParam, tids []int32) {
	v.Title = param.Title
	v.Valid = param.Valid
	v.MtimeOrder = strconv.Itoa(param.Order)
	v.PageCfg = param.PageCfg
	v.AID = param.CID
	v.Typeids = tids
	v.Result = "1"
}

// FromAuditConsult def.
func (v *ReqArcES) FromAuditConsult(param *ReqArcCons, tids []int32) {
	v.PubtimeOrder = strconv.Itoa(param.Order)
	v.Typeids = tids
	v.Result = param.Status
	v.Title = param.Title
	if param.AVID != 0 {
		v.AID = fmt.Sprintf("%d", param.AVID)
	}
	v.Pn = param.Pn
	v.Ps = _pagesize
}

// ArcRes def.
type ArcRes struct {
	AVID       int64  `json:"avid"`
	Title      string `json:"title"`
	FirstCat   string `json:"first_cat"`
	SecondCat  string `json:"second_cat"`
	Status     int    `json:"status"`
	InjectTime string `json:"inject_time"`
	PubTime    string `json:"pubtime"`
	Reason     string `json:"reason"`
}

// VideoRes def.
type VideoRes struct {
	CID        int64  `json:"cid"`
	Title      string `json:"title"`
	Page       int    `json:"page"`
	Status     int    `json:"status"`
	Ctime      string `json:"ctime"`
	InjectTime string `json:"inject_time"`
	Reason     string `json:"reason"`
}

// ArcResPager def.
type ArcResPager struct {
	Items []*ArcRes `json:"items"`
	Page  *Page     `json:"page"`
}

// VideoResPager def.
type VideoResPager struct {
	Items []*VideoRes `json:"items"`
	Page  *Page       `json:"page"`
}

// EsUgcConsult def.
type EsUgcConsult struct {
	Fields    []string            `json:"fields"`
	From      string              `json:"from"`
	Highlight bool                `json:"highlight"`
	Pn        int                 `json:"pn"`
	Ps        int                 `json:"ps"`
	Where     *Where              `json:"where,omitempty"`
	Order     []map[string]string `json:"order"`
}

// Where def.
type Where struct {
	Eq    map[string]string        `json:"eq,omitempty"`
	Or    map[string][]interface{} `json:"or,omitempty"`
	In    map[string][]int32       `json:"in,omitempty"`
	Range map[string]string        `json:"range,omitempty"`
}

// EsArc def.
type EsArc struct {
	AID     int64  `json:"aid"`
	MID     int64  `json:"mid"`
	Deleted int    `json:"deleted"`
	Mtime   string `json:"mtime"`
	Pubtime string `json:"pubtime"`
	Result  int    `json:"result"`
	Typeid  int32  `json:"typeid"`
	Valid   int    `json:"valid"`
}

// EsUgcResult def.
type EsUgcResult struct {
	Result []*EsArc
	Page   *Page
}
