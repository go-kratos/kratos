package model

// LotteryRes ..
type LotteryRes struct {
	RIDs    []int64 `json:"rids"`
	Offset  int64   `json:"offset"`
	HasMore int     `json:"has_more"`
}

// VoteBIZArchive .
type VoteBIZArchive struct {
	Aid int64 `json:"aid"`
}
