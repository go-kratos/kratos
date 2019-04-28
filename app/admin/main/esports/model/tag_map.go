package model

import (
	"fmt"
	"strings"
)

const _tagMapInsertSQL = "INSERT INTO es_tags_map(tid,aid) VALUES %s"

// TagMap .
type TagMap struct {
	ID        int64 `json:"id"`
	Tid       int64 `json:"tid"`
	Aid       int64 `json:"aid"`
	IsDeleted int   `json:"is_deleted"`
}

// TableName es_year_map.
func (t TagMap) TableName() string {
	return "es_tags_map"
}

// BatchAddTagMapSQL .
func BatchAddTagMapSQL(data []*TagMap) string {
	if len(data) == 0 {
		return ""
	}
	var rowStrings []string
	for _, v := range data {
		rowStrings = append(rowStrings, fmt.Sprintf("(%d,%d)", v.Tid, v.Aid))
	}
	return fmt.Sprintf(_tagMapInsertSQL, strings.Join(rowStrings, ","))
}
