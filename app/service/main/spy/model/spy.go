package model

// UserScore rpc return for getting user score.
type UserScore struct {
	Mid   int64 `json:"mid"`
	Score int8  `json:"score"`
}
