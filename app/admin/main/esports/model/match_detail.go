package model

// MatchDetail .
type MatchDetail struct {
	ID           int64  `json:"id" form:"id"`
	MaID         int64  `json:"ma_id" form:"ma_id" validate:"required"`
	GameType     int64  `json:"game_type" form:"game_type" validate:"required"`
	Stime        int64  `json:"stime" form:"stime" validate:"required"`
	Etime        int64  `json:"etime" form:"etime" validate:"required"`
	GameStage    string `json:"game_stage" form:"game_stage" validate:"required"`
	KnockoutType int64  `json:"knockout_type" form:"knockout_type"`
	WinnerType   int64  `json:"winner_type" form:"winner_type"`
	ScoreID      int64  `json:"score_id" form:"score_id"`
	Status       int    `json:"status" form:"status"`
	Online       int    `json:"online" form:"online"`
}

// TableName .
func (t MatchDetail) TableName() string {
	return "es_matchs_detail"
}
