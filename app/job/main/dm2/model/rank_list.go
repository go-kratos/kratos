package model

// RankRecentRegion  分区排行榜
type RankRecentRegion struct {
	Aid    int64               `json:"aid"`
	Mid    int64               `json:"mid"`
	Others []*RankRecentRegion `json:"others"`
}

// RankRecentResp .
type RankRecentResp struct {
	Code int32               `json:"code"`
	List []*RankRecentRegion `json:"list"`
	Num  int32               `json:"num"`
}
