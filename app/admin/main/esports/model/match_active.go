package model

import (
	"fmt"
	"strings"

	"go-common/library/xstr"
)

const (
	_moduleInsertSQL = "INSERT INTO es_matchs_module(ma_id,name,oids) VALUES %s"
	_moduleEditSQL   = "UPDATE es_matchs_module  SET name = CASE %s END,oids = CASE %s END WHERE id IN (%s)"
)

// ParamMA .
type ParamMA struct {
	MatchActive
	Modules string `json:"-" form:"modules"`
	Adid    int64  `json:"-" form:"adid" validate:"required"`
}

// MatchActive .
type MatchActive struct {
	ID           int64  `json:"id" form:"id"`
	Sid          int64  `json:"sid" form:"sid" validate:"required"`
	Mid          int64  `json:"mid" form:"mid" validate:"required"`
	Background   string `json:"background" form:"background"`
	BackColor    string `json:"back_color" form:"back_color"`
	ColorStep    string `json:"color_step" form:"color_step"`
	LiveID       int64  `json:"live_id" form:"live_id" validate:"required"`
	Intr         string `json:"intr" form:"intr"`
	Focus        string `json:"focus" form:"focus"`
	URL          string `json:"url" form:"url"`
	Status       int    `json:"status" form:"status"`
	H5Background string `json:"h5_background" form:"h5_background"`
	H5BackColor  string `json:"h5_back_color" form:"h5_back_color"`
	H5Focus      string `json:"h5_focus" form:"h5_focus"`
	H5URL        string `json:"h5_url" form:"h5_url"`
	IntrLogo     string `json:"intr_logo" form:"intr_logo"`
	IntrTitle    string `json:"intr_title" form:"intr_title"`
	IntrText     string `json:"intr_text" form:"intr_text"`
}

// Module .
type Module struct {
	ID     int64  `json:"id"`
	MaID   int64  `json:"ma_id"`
	Name   string `json:"name"`
	Oids   string `json:"oids"`
	Status int    `json:"-" form:"status"`
}

// MatchModule .
type MatchModule struct {
	*MatchActive
	Modules        []*Module `json:"modules"`
	MatchTitle     string    `json:"match_title"`
	MatchSubTitle  string    `json:"match_sub_title"`
	SeasonTitle    string    `json:"season_title"`
	SeasonSubTitle string    `json:"season_sub_title"`
}

// TableName es_matchs_module.
func (t Module) TableName() string {
	return "es_matchs_module"
}

// TableName es_matchs_active.
func (t MatchActive) TableName() string {
	return "es_matchs_active"
}

// BatchAddModuleSQL .
func BatchAddModuleSQL(maID int64, data []*Module) string {
	if len(data) == 0 {
		return ""
	}
	var rowStrings []string
	for _, v := range data {
		rowStrings = append(rowStrings, fmt.Sprintf("(%d,'%s','%s')", maID, v.Name, v.Oids))
	}
	return fmt.Sprintf(_moduleInsertSQL, strings.Join(rowStrings, ","))
}

// BatchEditModuleSQL .
func BatchEditModuleSQL(mapModuel []*Module) string {
	if len(mapModuel) == 0 {
		return ""
	}
	var (
		nameStr, oidsStr string
		ids              []int64
	)
	for _, module := range mapModuel {
		nameStr = fmt.Sprintf("%s WHEN id = %d THEN '%s'", nameStr, module.ID, module.Name)
		oidsStr = fmt.Sprintf("%s WHEN id = %d THEN '%s'", oidsStr, module.ID, module.Oids)
		ids = append(ids, module.ID)
	}
	return fmt.Sprintf(_moduleEditSQL, nameStr, oidsStr, xstr.JoinInts(ids))
}
