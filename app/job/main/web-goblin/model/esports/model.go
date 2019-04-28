package esports

// Contest contest.
type Contest struct {
	ID             int64  `json:"id"`
	Stime          int64  `json:"stime"`
	LiveRoom       int64  `json:"live_room"`
	HomeID         int64  `json:"home_id"`
	AwayID         int64  `json:"away_id"`
	SuccessTeam    int64  `json:"success_team"`
	SeasonTitle    string `json:"season_title"`
	SeasonSubTitle string `json:"season_sub_title"`
	Special        int    `json:"special"`
	SpecialName    string `json:"special_name"`
	SpecialTips    string `json:"special_tips"`
}

// Team team.
type Team struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	SubTitle string `json:"sub_title"`
}

// Arc  arc.
type Arc struct {
	ID        int64 `json:"id"`
	Aid       int64 `json:"aid"`
	Score     int64 `json:"score"`
	IsDeleted int   `json:"is_deleted"`
}
