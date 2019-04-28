package model

// rank
var (
	RankMonth           = int64(1)
	RankWeek            = int64(2)
	RankYesterday       = int64(3)
	RankBeforeYesterday = int64(4)
)

// RankCategory .
type RankCategory struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Rank .
type Rank struct {
	Aid   int64 `json:"aid"`
	Score int64 `json:"score"`
}

// RankMeta .
type RankMeta struct {
	*Meta
	Attention bool  `json:"attention"`
	Score     int64 `json:"score"`
}

// RankResp .
type RankResp struct {
	Note string  `json:"note"`
	List []*Rank `json:"list"`
}
