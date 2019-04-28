package model

import (
	"fmt"
	"strings"

	"go-common/library/time"
)

const (
	//_blacklistBatchAddSQL add blacklist batch sql
	_blacklistBatchAddSQL    = "INSERT INTO blacklist (`mid`) VALUES %s"
	_blacklistBatchUpdateSQL = "UPDATE blacklist SET `status` = 0 WHERE id in (%s)"
	//StatusBlack blacklist
	StatusBlack = 0
)

// Blacklist .
type Blacklist struct {
	ID     int64     `json:"id" form:"id"`
	Mid    int64     `json:"mid" form:"mid"`
	Status int       `json:"status" form:"status"`
	Mtime  time.Time `json:"mtime"`
}

//BlacklistPager blacklist pager
type BlacklistPager struct {
	Item []*Blacklist
	Page Page
}

// BlacklistBatchAddSQL .
func BlacklistBatchAddSQL(mids []int64) string {
	if len(mids) == 0 {
		return ""
	}
	var rowStrings []string
	for _, mid := range mids {
		rowStrings = append(rowStrings, fmt.Sprintf("(%d)", mid))
	}
	return fmt.Sprintf(_blacklistBatchAddSQL, strings.Join(rowStrings, ","))
}

// BlacklistBatchUpdateSQL .
func BlacklistBatchUpdateSQL(mids []int64) string {
	if len(mids) == 0 {
		return ""
	}
	var rowStrings []string
	for _, mid := range mids {
		rowStrings = append(rowStrings, fmt.Sprintf("%d", mid))
	}
	return fmt.Sprintf(_blacklistBatchUpdateSQL, strings.Join(rowStrings, ","))
}

// TableName .
func (a Blacklist) TableName() string {
	return "blacklist"
}
