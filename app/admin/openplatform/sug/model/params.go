package model

// SourceSearch params
type SourceSearch struct {
	Type     string `form:"type"`
	Keyword  string `form:"keyword"`
	PageSize int    `form:"pageSize"`
	PageNum  int    `form:"pageNum"`
}

// Search params
type Search struct {
	SeasonID int64 `form:"season_id"`
	ItemsID  int64 `form:"items_id"`
}

// MatchOperate params
type MatchOperate struct {
	OpType   int8  `form:"op_type"`
	SeasonID int64 `form:"season_id"`
	ItemsID  int64 `form:"items_id"`
}
