package model

import (
	"fmt"
	"strings"
)

const _matchMapInsertSQL = "INSERT INTO es_matchs_map(mid,aid) VALUES %s"

// MatchMap .
type MatchMap struct {
	ID        int64 `json:"id"`
	Mid       int64 `json:"mid"`
	Aid       int64 `json:"aid"`
	IsDeleted int   `json:"is_deleted"`
}

// TableName es_year_map.
func (m MatchMap) TableName() string {
	return "es_matchs_map"
}

// BatchAddMachMapSQL .
func BatchAddMachMapSQL(data []*MatchMap) string {
	if len(data) == 0 {
		return ""
	}
	var rowStrings []string
	for _, v := range data {
		rowStrings = append(rowStrings, fmt.Sprintf("(%d,%d)", v.Mid, v.Aid))
	}
	return fmt.Sprintf(_matchMapInsertSQL, strings.Join(rowStrings, ","))
}
