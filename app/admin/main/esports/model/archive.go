package model

import (
	"fmt"
	"strings"
)

const _arcBatchAddSQL = "INSERT INTO es_archives(`aid`) VALUES %s"

// Arc .
type Arc struct {
	ID        int64 `json:"id"`
	Aid       int64 `json:"aid"`
	IsDeleted int   `json:"is_deleted"`
}

// ArcAddParam .
type ArcAddParam struct {
	Aids     []int64 `form:"aids,split" validate:"dive,gt=1"`
	Gids     []int64 `form:"gids,split"`
	MatchIDs []int64 `form:"match_ids,split"`
	TeamIDs  []int64 `form:"team_ids,split"`
	TagIDs   []int64 `form:"tag_ids,split"`
	Years    []int64 `form:"years,split"`
}

// ArcImportParam .
type ArcImportParam struct {
	Aid      int64   `form:"aid" validate:"min=1"`
	Gids     []int64 `form:"gids,split"`
	MatchIDs []int64 `form:"match_ids,split"`
	TeamIDs  []int64 `form:"team_ids,split"`
	TagIDs   []int64 `form:"tag_ids,split"`
	Years    []int64 `form:"years,split"`
}

// ArcListParam .
type ArcListParam struct {
	Title     string `form:"title"`
	Aid       int64  `form:"aid"`
	TypeID    int64  `form:"type_id"`
	Copyright int    `form:"copyright"`
	State     string `form:"state"`
	Pn        int    `form:"pn"`
	Ps        int    `form:"ps"`
}

// SearchArc .
type SearchArc struct {
	Aid    int64   `json:"aid"`
	TypeID int64   `json:"typeid"`
	Title  string  `json:"title"`
	State  int64   `json:"state"`
	Mid    int64   `json:"mid"`
	Gid    []int64 `json:"gid"`
	Tags   []int64 `json:"tags"`
	Matchs []int64 `json:"matchs"`
	Teams  []int64 `json:"teams"`
	Year   []int64 `json:"year"`
}

// ArcResult .
type ArcResult struct {
	Aid    int64    `json:"aid"`
	TypeID int64    `json:"type_id"`
	Title  string   `json:"title"`
	State  int64    `json:"state"`
	Mid    int64    `json:"mid"`
	Uname  string   `json:"uname"`
	Games  []*Game  `json:"games"`
	Tags   []*Tag   `json:"tags"`
	Matchs []*Match `json:"matchs"`
	Teams  []*Team  `json:"teams"`
	Years  []int64  `json:"years"`
}

// ArcRelation .
type ArcRelation struct {
	AddGids     []*GIDMap
	UpAddGids   []int64
	UpDelGids   []int64
	AddMatchs   []*MatchMap
	UpAddMatchs []int64
	UpDelMatchs []int64
	AddTags     []*TagMap
	UpAddTags   []int64
	UpDelTags   []int64
	AddTeams    []*TeamMap
	UpAddTeams  []int64
	UpDelTeams  []int64
	AddYears    []*YearMap
	UpAddYears  []int64
	UpDelYears  []int64
}

// TableName .
func (a Arc) TableName() string {
	return "es_archives"
}

// ArcBatchAddSQL .
func ArcBatchAddSQL(aids []int64) string {
	if len(aids) == 0 {
		return ""
	}
	var rowStrings []string
	for _, aid := range aids {
		rowStrings = append(rowStrings, fmt.Sprintf("(%d)", aid))
	}
	return fmt.Sprintf(_arcBatchAddSQL, strings.Join(rowStrings, ","))
}
