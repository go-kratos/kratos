package model

// Contest .
type Contest struct {
	ID            int64  `json:"id" form:"id"`
	GameStage     string `json:"game_stage" form:"game_stage" validate:"required"`
	Stime         int64  `json:"stime" form:"stime"`
	Etime         int64  `json:"etime" form:"etime"`
	HomeID        int64  `json:"home_id" form:"home_id"`
	AwayID        int64  `json:"away_id" form:"away_id"`
	HomeScore     int64  `json:"home_score" form:"home_score"`
	AwayScore     int64  `json:"away_score" form:"away_score"`
	LiveRoom      int64  `json:"live_room" form:"live_room"`
	Aid           int64  `json:"aid" form:"aid"`
	Collection    int64  `json:"collection" form:"collection"`
	GameState     int    `json:"game_state" form:"game_state"`
	Dic           string `json:"dic" form:"dic"`
	Status        int    `json:"status" form:"status"`
	Sid           int64  `json:"sid" form:"sid" validate:"required"`
	Mid           int64  `json:"mid" form:"mid" validate:"required"`
	Special       int    `json:"special" form:"special"`
	SuccessTeam   int64  `json:"success_team" form:"success_team"`
	SpecialName   string `json:"special_name" form:"special_name"`
	SpecialTips   string `json:"special_tips" form:"special_tips"`
	SpecialImage  string `json:"special_image" form:"special_image"`
	Playback      string `json:"playback" form:"playback"`
	CollectionURL string `json:"collection_url" form:"collection_url"`
	LiveURL       string `json:"live_url" form:"live_url"`
	DataType      int64  `json:"data_type" form:"data_type"`
	Data          string `json:"-" form:"data" gorm:"-"`
	Adid          int64  `json:"-" form:"adid"  gorm:"-" validate:"required"`
	MatchID       int64  `json:"match_id" form:"match_id"`
}

// ContestInfo .
type ContestInfo struct {
	*Contest
	Games       []*Game        `json:"games"`
	HomeName    string         `json:"home_name"`
	AwayName    string         `json:"away_name"`
	SuccessName string         `json:"success_name" form:"success_name"`
	Data        []*ContestData `json:"data"`
}

// TableName es_contests
func (c Contest) TableName() string {
	return "es_contests"
}
