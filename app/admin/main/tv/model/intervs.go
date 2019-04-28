package model

import (
	"encoding/json"
	"go-common/library/time"

	"github.com/jinzhu/gorm"
	"github.com/siddontang/go-mysql/mysql"
)

// Tv_rank table related const params
const (
	_RankCategory = 5 // _RankCategory 模块干预
	RankIdxBase   = 5 // index page intervention base, pgc=5+1, ugc=5+2
)

// SimpleRank represents the table TV_RANK, but with only necessary fields for the front-end
type SimpleRank struct {
	Title      string `json:"title"`
	Source     int    `json:"source"`
	SourceName string `json:"source_name"`
	Mtime      string `json:"mtime"`
	Pubdate    string `json:"pubdate"`
	RankCore
}

// RankCore def
type RankCore struct {
	Rank     int64 `json:"rank"`
	ID       int64 `json:"id"`
	ContID   int64 `json:"cid" gorm:"column:cont_id"`
	ContType int   `json:"cont_type"`
	Position int   `json:"position"`
}

// RankError represents the invalid season info
type RankError struct {
	ID       int `json:"id"`
	SeasonID int `json:"season_id"`
}

// RankList is the output format for intervention list
type RankList struct {
	List []*SimpleRank `json:"list"`
}

// Rank represents the table TV_RANK
type Rank struct {
	Title     string
	ModuleID  int64 `gorm:"column:module_id"`
	Category  int8
	IsDeleted int8
	Mtime     time.Time
	RankCore
}

// TableName tv_rank
func (c SimpleRank) TableName() string {
	return "tv_rank"
}

// TableName tv_rank
func (v Rank) TableName() string {
	return "tv_rank"
}

// BeComplete transforms a simpleRank to Complete rank in order to create it in DB
func (c SimpleRank) BeComplete(req *IntervPubReq, title string, position int) (res *Rank) {
	res = &Rank{
		Title:    title,
		RankCore: c.RankCore,
	}
	res.Position = position
	if req.ModuleID > 0 {
		res.Category = _RankCategory
		res.ModuleID = req.ModuleID
		return
	}
	if req.Rank > 0 {
		res.Rank = req.Rank
	}
	if req.Category > 0 {
		res.Category = int8(req.Category)
	}
	return
}

type catName func(int) string                     // translate pgc category to CN name
type tpParName func(int32) (string, int32, error) // translate ugc type to its parent tid and parent's name

// BeSimpleSn def.
func (v *Rank) BeSimpleSn(sn *TVEpSeason, translate catName) *SimpleRank {
	return &SimpleRank{
		RankCore:   v.RankCore,
		Title:      sn.Title,
		Source:     sn.Category,
		SourceName: translate(sn.Category),
		Pubdate:    sn.PlayTime.Time().Format(mysql.TimeFormat),
		Mtime:      v.Mtime.Time().Format(mysql.TimeFormat),
	}
}

// BeSimpleArc def.
func (v *Rank) BeSimpleArc(arc *SimpleArc, translate tpParName) (res *SimpleRank) {
	res = &SimpleRank{
		RankCore: v.RankCore,
		Title:    arc.Title,
		Mtime:    v.Mtime.Time().Format(mysql.TimeFormat),
		Pubdate:  arc.Pubtime.Time().Format(mysql.TimeFormat),
	}
	if pname, pid, err := translate(arc.TypeID); err == nil {
		res.Source = int(pid)
		res.SourceName = pname
	}
	return
}

//BeError transforms a rank to rankError
func (v Rank) BeError() *RankError {
	return &RankError{
		ID:       int(v.ID),
		SeasonID: int(v.ContID),
	}
}

// RankListReq is rank list request
type RankListReq struct {
	Rank     int64 `form:"rank" validate:"min=0"`
	Category int64 `form:"category" validate:"required,min=1"`
}

// RankPubReq is rank publish request
type RankPubReq struct {
	RankListReq
	Intervs string `form:"intervs" validate:"required"`
}

// ModListReq is mod list request
type ModListReq struct {
	ModuleID int64 `form:"module_id" validate:"required,min=1"`
}

// ModPubReq is mod publish request
type ModPubReq struct {
	ModListReq
	Intervs string `form:"intervs" validate:"required"`
}

// IdxListReq is index list request
type IdxListReq struct {
	TypeID   int64 `form:"type_id" validate:"required,min=1"`
	RankType int64 `form:"rank_type" validate:"required,min=1,max=2"` // 1=pgc, 2=ugc
}

// IdxPubReq is index publish request.
type IdxPubReq struct {
	IdxListReq
	Intervs string `form:"intervs" validate:"required"`
}

// IntervListReq is common request for interv list.
type IntervListReq struct {
	Rank     int64
	Category int64
	ModuleID int64
}

// IntervPubReq is common request for interv publish.
type IntervPubReq struct {
	IntervListReq
	Items []*SimpleRank
}

// FromRank builds the request with rank & category params
func (v *IntervListReq) FromRank(rank *RankListReq) {
	v.Rank = rank.Rank
	v.Category = rank.Category
	v.ModuleID = 0
}

// FromRank def.
func (v *IntervPubReq) FromRank(rank *RankPubReq) (err error) {
	v.IntervListReq.FromRank(&rank.RankListReq)
	return json.Unmarshal([]byte(rank.Intervs), &v.Items)
}

// FromMod builds the request with module params
func (v *IntervListReq) FromMod(mod *ModListReq) {
	v.Rank = 0
	v.Category = _RankCategory
	v.ModuleID = mod.ModuleID
}

// FromMod builds the request with module params
func (v *IntervPubReq) FromMod(mod *ModPubReq) (err error) {
	v.IntervListReq.FromMod(&mod.ModListReq)
	return json.Unmarshal([]byte(mod.Intervs), &v.Items)
}

// FromIndex builds the request with index params
func (v *IntervListReq) FromIndex(idx *IdxListReq) {
	v.Rank = idx.TypeID                     // category id, pgc or ugc type id
	v.Category = idx.RankType + RankIdxBase // 6 or 7
	v.ModuleID = 0
}

// IsIdx tells whether this request is from index
func (v *IntervListReq) IsIdx() bool {
	return v.Category > RankIdxBase
}

// FromIndex def.
func (v *IntervPubReq) FromIndex(idx *IdxPubReq) (err error) {
	v.IntervListReq.FromIndex(&idx.IdxListReq)
	return json.Unmarshal([]byte(idx.Intervs), &v.Items)
}

// BuildDB builds the db from the intervention request
func (v *IntervListReq) BuildDB(db *gorm.DB) (newDB *gorm.DB) {
	newDB = db.Model(Rank{}).Where("is_deleted = 0")
	if v.ModuleID == 0 { // index or rank
		newDB = newDB.Where("rank = ?", v.Rank).Where("category = ?", v.Category)
	} else {
		newDB = newDB.Where("module_id = ?", v.ModuleID)
	}
	return
}
