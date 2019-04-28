package model

import (
	"fmt"
	"strings"
)

const _yearMapInsertSQL = "INSERT INTO es_year_map(year,aid) VALUES %s"

// YearMap .
type YearMap struct {
	ID        int64 `json:"id"`
	Year      int64 `json:"year"`
	Aid       int64 `json:"aid"`
	IsDeleted int   `json:"is_deleted"`
}

// TableName es_year_map.
func (y YearMap) TableName() string {
	return "es_year_map"
}

// BatchAddYearMapSQL .
func BatchAddYearMapSQL(data []*YearMap) string {
	if len(data) == 0 {
		return ""
	}
	var rowStrings []string
	for _, v := range data {
		rowStrings = append(rowStrings, fmt.Sprintf("(%d,%d)", v.Year, v.Aid))
	}
	return fmt.Sprintf(_yearMapInsertSQL, strings.Join(rowStrings, ","))
}
