package model

// RankType
const (
	MaxDurationRank = "max_duration"
	MinDurationRank = "min_duration"
	AvgDurationRank = "avg_duration"
	ErrorsRank      = "errors"
)

// VerifyRankType .
func VerifyRankType(rankType string) bool {
	switch rankType {
	case MaxDurationRank, MinDurationRank, AvgDurationRank, ErrorsRank:
		return true
	}
	return false
}
