package model

import (
	"encoding/json"
	"fmt"
	xtime "time"

	arccli "go-common/app/service/main/archive/api"
	"go-common/library/time"

	"github.com/siddontang/go-mysql/mysql"
)

// label related params
const (
	ParamTypeid  = "typeid"
	ParamUgctime = "pubtime"
	UgcLabel     = 2
	PgcLabel     = 1
)

// TpLabel def.
type TpLabel struct {
	Category  int    `json:"-"`
	Param     string `json:"param"`
	ParamName string `json:"param_name"`
}

// ReqLabel def.
type ReqLabel struct {
	Category int    `form:"category" validate:"required"`
	Param    string `form:"param" validate:"required"` // pubtime for time labels, typeid for type labels
	Title    string `form:"title"`
	ID       int    `form:"id"`
}

// LabelDB is the index label in DB
type LabelDB struct {
	LabelCore
	Mtime time.Time `json:"Mtime"`
}

// SameType tells whether the given label has the exact same type with the V
func (v *LabelDB) SameType(given *LabelDB) bool {
	return v.Category == given.Category && v.Param == given.Param && v.CatType == given.CatType
}

// LabelCore is core of Label
type LabelCore struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Param     string `json:"param"`
	ParamName string `json:"param_name"`
	Value     string `json:"value"`
	Category  int32  `json:"category"`
	CatType   int    `json:"cat_type"`
	Valid     int    `json:"valid"`
	Position  int    `json:"position"`
}

// LabelList is used to list in TV CMS
type LabelList struct {
	LabelCore
	Mtime string `json:"mtime"`
	Stime string `json:"stime,omitempty"`
	Etime string `json:"etime,omitempty"`
}

// PgcCondResp is pgc condition response structure
type PgcCondResp struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Result  *PgcCond `json:"result"`
}

// PgcCond def.
type PgcCond struct {
	Filter []*Cond `json:"filter"`
}

// Cond def.
type Cond struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Value []*CondV `json:"value"`
}

// CondV def.
type CondV struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// UgcTime is used to add time labels for ugc
type UgcTime struct {
	UTime
	Category int32  `form:"category" validate:"required"`
	Name     string `form:"name" validate:"required"`
}

// EditUgcTime def.
type EditUgcTime struct {
	ID   int64  `form:"id" validate:"required"`
	Name string `form:"name" validate:"required"`
	UTime
}

// UTime is used for storage in DB by json
type UTime struct {
	Stime int64 `form:"stime" validate:"required" json:"stime"`
	Etime int64 `form:"etime" validate:"required" json:"etime"`
}

// TimeV picks time value in Json
func (tm *UTime) TimeV() string {
	timeV, _ := json.Marshal(tm)
	return string(timeV)
}

// ToList transforms LabelDB to LabelList
func (v *LabelDB) ToList() *LabelList {
	res := &LabelList{
		LabelCore: v.LabelCore,
		Mtime:     v.Mtime.Time().Format(mysql.TimeFormat),
	}
	if v.CatType == UgcLabel && v.Param == ParamUgctime && v.Value != "" {
		utime := UTime{}
		if err := json.Unmarshal([]byte(v.Value), &utime); err != nil {
			return res
		}
		res.Stime = xtime.Unix(utime.Stime, 0).Format(mysql.TimeFormat)
		res.Etime = xtime.Unix(utime.Etime, 0).Format(mysql.TimeFormat)
	}
	return res
}

// TableName tv_rank
func (v LabelDB) TableName() string {
	return "tv_label"
}

// FromArcTp def.
func (v *LabelDB) FromArcTp(tp *arccli.Tp, paramName string) {
	v.LabelCore = LabelCore{
		Name:      tp.Name,
		Value:     fmt.Sprintf("%d", tp.ID),
		Category:  tp.Pid,
		Param:     ParamTypeid,
		ParamName: paramName,
		CatType:   UgcLabel,
		Valid:     1,
	}
}

// FromPgcCond def.
func (v *LabelDB) FromPgcCond(value *CondV, cond *Cond, category int32) {
	v.LabelCore = LabelCore{
		Name:      value.Name,
		Value:     value.ID,
		Category:  category,
		Param:     cond.ID,
		ParamName: cond.Name,
		CatType:   PgcLabel,
		Valid:     1,
	}
}

// FromUgcTime def.
func (v *LabelDB) FromUgcTime(tm *UgcTime, paramName string) {
	v.LabelCore = LabelCore{
		Name:      tm.Name,
		Value:     tm.TimeV(),
		Category:  tm.Category,
		Param:     ParamUgctime,
		ParamName: paramName,
		CatType:   UgcLabel,
		Valid:     1,
	}
}
