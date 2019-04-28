package model

import (
	"fmt"

	"go-common/library/xstr"
)

const _treeEditSQL = "UPDATE es_matchs_tree  SET mid = CASE %s END  WHERE id IN (%s)"

// TreeListParam .
type TreeListParam struct {
	MadID int64 `form:"mad_id" validate:"required"`
}

// TreeEditParam .
type TreeEditParam struct {
	MadID int64  `form:"mad_id" validate:"required"`
	Nodes string `form:"nodes" validate:"required"`
}

// TreeDelParam .
type TreeDelParam struct {
	MadID int64  `form:"mad_id" validate:"required"`
	IDs   string `form:"ids" validate:"required"`
}

// Tree .
type Tree struct {
	ID        int64 `json:"id" form:"id"`
	MaID      int64 `json:"ma_id,omitempty" form:"ma_id" validate:"required"`
	MadID     int64 `json:"mad_id,omitempty" form:"mad_id" validate:"required"`
	Pid       int64 `json:"pid" form:"pid"`
	RootID    int64 `json:"root_id" form:"root_id"`
	GameRank  int64 `json:"game_rank,omitempty" form:"game_rank" validate:"required"`
	Mid       int64 `json:"mid" form:"mid"`
	IsDeleted int   `json:"is_deleted,omitempty" form:"is_deleted"`
}

// TreeList .
type TreeList struct {
	*Tree
	*ContestInfo
}

// TreeDetailList .
type TreeDetailList struct {
	Detail *MatchDetail  `json:"detail"`
	Tree   [][]*TreeList `json:"tree"`
}

// TableName .
func (t Tree) TableName() string {
	return "es_matchs_tree"
}

// BatchEditTreeSQL .
func BatchEditTreeSQL(nodes map[int64]int64) string {
	if len(nodes) == 0 {
		return ""
	}
	var (
		caseStr string
		ids     []int64
	)
	for id, mid := range nodes {
		caseStr = fmt.Sprintf("%s WHEN id = %d THEN %d", caseStr, id, mid)
		ids = append(ids, id)
	}
	return fmt.Sprintf(_treeEditSQL, caseStr, xstr.JoinInts(ids))
}
