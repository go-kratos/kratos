package model

// Achieve is
type Achieve struct {
	Award string `json:"award"`
	Mid   int64  `json:"mid"`
}

// AchieveGetReply is
type AchieveGetReply struct {
	AwardToken string `json:"award_token"`
}
