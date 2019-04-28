package search

import (
	"encoding/json"

	"go-common/app/interface/main/tv/model"
	"go-common/app/interface/main/tv/model/thirdp"
	"go-common/library/time"

	"github.com/siddontang/go-mysql/mysql"
)

const (
	_pgcType     = 1
	_ugcType     = 2
	_hotestOrder = 2
	// AllLabel is the string that means all data
	AllLabel = "-1" // it means all for index
)

// EsPage def.
type EsPage struct {
	Num     int `json:"num"`
	Size    int `json:"size"`
	Total   int `json:"total"`
	PageNum int `json:"page_num"`
}

// GetPageNb def.
func (v *EsPage) GetPageNb() {
	if v.Total%v.Size == 0 {
		v.PageNum = v.Total / v.Size
	} else {
		v.PageNum = v.Total/v.Size + 1
	}
}

// EsPgcResult def.
type EsPgcResult struct {
	Page   *EsPage     `json:"page"`
	Result []*PgcEsMdl `json:"result"`
}

// EsUgcResult def.
type EsUgcResult struct {
	Page   *EsPage     `json:"page"`
	Result []*UgcEsMdl `json:"result"`
}

// UgcEsMdl def.
type UgcEsMdl struct {
	AID    int64 `json:"aid"`
	TypeID int32 `json:"typeid"`
}

// PgcEsMdl def.
type PgcEsMdl struct {
	SeasonID int64 `json:"season_id"`
}

// EsCard def.
type EsCard struct {
	model.Card
	DataType int `json:"data_type"` // 1=pgc, 2=ugc
}

// EsPager def.
type EsPager struct {
	Title  string    `json:"title"`
	Page   *EsPage   `json:"page"`
	Result []*EsCard `json:"result"`
}

// FromPgc def.
func (v *EsCard) FromPgc(card *model.Card) {
	v.Card = *card
	v.DataType = _pgcType
}

// FromUgc def.
func (v *EsCard) FromUgc(card *model.Card) {
	v.Card = *card
	v.DataType = _ugcType
}

// ReqEsPn is part of the es request for page related cfg
type ReqEsPn struct {
	Ps    int `form:"ps" default:"50"`
	Pn    int `form:"pn" default:"1"`
	Order int `form:"order" validate:"required"`
	Sort  int `form:"sort" default:"0"`
}

// ReqPgcIdx is bm request for pgc index
type ReqPgcIdx struct {
	SeasonType    int    `form:"category" validate:"required,max=5"`
	ProducerID    int    `form:"producer_id"`
	Year          string `form:"year"`
	StyleID       int    `form:"style_id"`
	PubDate       string `form:"pub_date"`
	SeasonMonth   int    `form:"season_month" validate:"max=10"`
	SeasonStatus  []int  `form:"season_status,split"`
	Copyright     []int  `form:"copyright,split"`
	IsFinish      string `form:"is_finish"`
	Area          []int  `form:"area,split"`
	SeasonVersion int    `form:"season_version"`
	ReqEsPn
}

// IsDefault def.
func (v *ReqPgcIdx) IsDefault() bool {
	return v.Order == _hotestOrder && v.ProducerID <= 0 &&
		v.IsAllStr(v.Year) && v.StyleID <= 0 && v.IsAllStr(v.PubDate) && v.SeasonMonth <= 0 &&
		v.IsAll(v.SeasonStatus) && v.IsAll(v.Copyright) && v.IsAllStr(v.IsFinish) && v.IsAll(v.Area) &&
		v.SeasonVersion <= 0
}

// IsDefault def.
func (v *ReqUgcIdx) IsDefault() bool {
	return v.Order == _hotestOrder && v.SecondTID <= 0 && (v.PubTime == "" || v.PubTime == AllLabel)
}

// Title def.
func (v *ReqPgcIdx) Title() string {
	return thirdp.PgcCat(v.SeasonType)
}

// PgcOrder treats the order to get the responding field
func (v *ReqPgcIdx) PgcOrder() (field string) {
	switch v.Order {
	case 0: // update time
		return "latest_time"
	case 1:
		return "dm_count"
	case 2:
		return "play_count"
	case 3: // follow
		return "fav_count"
	case 4:
		return "score"
	case 5: // for others
		return "pub_time"
	case 6: // for movie
		return "release_date"
	default:
		return "play_count"
	}
}

// IdxSort treats the sort
func IdxSort(sortV int) (sort string) {
	switch sortV {
	case 0:
		return "desc"
	case 1:
		return "asc"
	default:
		return "desc"
	}
}

// ReqUgcIdx is bm request for ugc index
type ReqUgcIdx struct {
	ParentTID int32  `form:"category" validate:"required"`
	SecondTID int32  `form:"typeid"`
	PubTime   string `form:"pubtime"`
	ReqEsPn
}

// TimeStr picks json str from request and returns the time range struct in String format
func (v *ReqUgcIdx) TimeStr() (res *UgcTime, err error) {
	var pubTimeV = &UgcTimeV{}
	if err = json.Unmarshal([]byte(v.PubTime), &pubTimeV); err != nil {
		return
	}
	res = &UgcTime{
		STime: pubTimeV.STime.Time().Format(mysql.TimeFormat),
		ETime: pubTimeV.ETime.Time().Format(mysql.TimeFormat),
	}
	return
}

// UgcTimeV def. INT64 for http
type UgcTimeV struct {
	STime time.Time `json:"stime"`
	ETime time.Time `json:"etime"`
}

// UgcTime def. STRING for ES
type UgcTime struct {
	STime string `json:"stime"`
	ETime string `json:"etime"`
}

// SrvUgcIdx is the request treated by service for DAO, like type_id field and pub_time
type SrvUgcIdx struct {
	TIDs    []int32
	PubTime *UgcTime
	ReqEsPn
}

// UgcOrder treats the order to get the responding field
func (v *SrvUgcIdx) UgcOrder() (field string) {
	switch v.Order {
	case 1: // update time
		return "pubtime"
	case 2:
		return "click"
	default:
		return "click"
	}
}

// ReqIdxInterv is used for index intervention treatment
type ReqIdxInterv struct {
	EsIDs    []int64
	Category int
	IsPGC    bool
	Pn       int
}

// FromPGC def.
func (v *ReqIdxInterv) FromPGC(sids []int64, req *ReqPgcIdx) {
	v.EsIDs = sids
	v.Category = req.SeasonType
	v.IsPGC = true
	v.Pn = req.Pn
}

// FromUGC def.
func (v *ReqIdxInterv) FromUGC(aids []int64, req *ReqUgcIdx) {
	v.EsIDs = aids
	v.Category = int(req.ParentTID)
	v.IsPGC = false
	v.Pn = req.Pn
}

// IdxIntervSave def.
type IdxIntervSave struct {
	Pgc map[int][]int64
	Ugc map[int][]int64
}

// IsAll checks whether a slice of int means all data
func (v *ReqPgcIdx) IsAll(params []int) (res bool) {
	if len(params) == 0 {
		return true
	}
	for _, v := range params {
		if v < 0 {
			return true
		}
	}
	return false
}

// IsAllStr checks whether a string means all data
func (v *ReqPgcIdx) IsAllStr(duration string) (res bool) {
	if duration == "" {
		return true
	}
	if duration == AllLabel {
		return true
	}
	return false
}
