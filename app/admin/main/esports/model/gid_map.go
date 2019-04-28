package model

import (
	"fmt"
	"strings"
)

// gid map
const (
	TypeMatch        = 1
	TypeSeason       = 2
	TypeContest      = 3
	TypeTeam         = 4
	TypeArc          = 5
	_gidMapInsertSQL = "INSERT INTO es_gid_map(`type`,`oid`,`gid`) VALUES %s"
)

// GIDMap .
type GIDMap struct {
	ID        int64 `json:"id"`
	Type      int   `json:"type"`
	Oid       int64 `json:"oid"`
	Gid       int64 `json:"gid"`
	IsDeleted int   `json:"is_deleted"`
}

// TableName .
func (g GIDMap) TableName() string {
	return "es_gid_map"
}

// GidBatchAddSQL .
func GidBatchAddSQL(gidMap []*GIDMap) string {
	if len(gidMap) == 0 {
		return ""
	}
	var rowStrings []string
	for _, v := range gidMap {
		rowStrings = append(rowStrings, fmt.Sprintf("(%d,%d,%d)", v.Type, v.Oid, v.Gid))
	}
	return fmt.Sprintf(_gidMapInsertSQL, strings.Join(rowStrings, ","))
}
