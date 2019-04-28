package model

import (
	"fmt"
	"strings"

	"go-common/library/xstr"
)

const (
	_cDataInsertSQL = "INSERT INTO es_contests_data(cid,url,point_data) VALUES %s"
	_cDataEditSQL   = "UPDATE es_contests_data SET url = CASE %s END,point_data = CASE %s END WHERE id IN (%s)"
)

// ContestData .
type ContestData struct {
	ID        int64  `json:"id"`
	CID       int64  `json:"cid" gorm:"column:cid"`
	URL       string `json:"url"`
	PointData int    `json:"point_data"`
	IsDeleted int    `json:"is_deleted"`
}

// TableName es_contests_data.
func (t ContestData) TableName() string {
	return "es_contests_data"
}

// BatchAddCDataSQL .
func BatchAddCDataSQL(cID int64, data []*ContestData) string {
	if len(data) == 0 {
		return ""
	}
	var rowStrings []string
	for _, v := range data {
		rowStrings = append(rowStrings, fmt.Sprintf("(%d,'%s',%d)", cID, v.URL, v.PointData))
	}
	return fmt.Sprintf(_cDataInsertSQL, strings.Join(rowStrings, ","))
}

// BatchEditCDataSQL .
func BatchEditCDataSQL(cDatas []*ContestData) string {
	if len(cDatas) == 0 {
		return ""
	}
	var (
		urlStr, pDataStr string
		ids              []int64
	)
	for _, module := range cDatas {
		urlStr = fmt.Sprintf("%s WHEN id = %d THEN '%s'", urlStr, module.ID, module.URL)
		pDataStr = fmt.Sprintf("%s WHEN id = %d THEN '%d'", pDataStr, module.ID, module.PointData)
		ids = append(ids, module.ID)
	}
	return fmt.Sprintf(_cDataEditSQL, urlStr, pDataStr, xstr.JoinInts(ids))
}
