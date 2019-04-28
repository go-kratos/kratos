package show

import (
	"fmt"
	"strings"

	"go-common/library/xstr"
)

const (
	_queryInsertSQL = "INSERT INTO search_web_query(sid,value) VALUES %s"
	_queryEditSQL   = "UPDATE search_web_query SET value = CASE %s END WHERE id IN (%s)"
)

//SearchWebQuery search web query
type SearchWebQuery struct {
	ID      int64  `json:"id" form:"id"`
	SID     int64  `json:"sid" form:"sid" gorm:"column:sid"`
	Value   string `json:"value" form:"value"`
	Deleted int    `json:"deleted" form:"deleted"`
}

// TableName .
func (a SearchWebQuery) TableName() string {
	return "search_web_query"
}

// BatchAddQuerySQL .
func BatchAddQuerySQL(sID int64, data []*SearchWebQuery) string {
	if len(data) == 0 {
		return ""
	}
	var rowStrings []string
	for _, v := range data {
		rowStrings = append(rowStrings, fmt.Sprintf("(%d,'%s')", sID, v.Value))
	}
	return fmt.Sprintf(_queryInsertSQL, strings.Join(rowStrings, ","))
}

// BatchEditQuerySQL .
func BatchEditQuerySQL(querys []*SearchWebQuery) string {
	if len(querys) == 0 {
		return ""
	}
	var (
		oidsStr string
		ids     []int64
	)
	for _, query := range querys {
		oidsStr = fmt.Sprintf("%s WHEN id = %d THEN '%s'", oidsStr, query.ID, query.Value)
		ids = append(ids, query.ID)
	}
	return fmt.Sprintf(_queryEditSQL, oidsStr, xstr.JoinInts(ids))
}
