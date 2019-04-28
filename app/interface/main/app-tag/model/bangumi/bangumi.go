package bangumi

type Bangumi struct {
	SeasonId     string `json:"season_id"`
	Spid         string `json:"spid"`
	Title        string `json:"title"`
	Brief        string `json:"brief"`
	Cover        string `json:"cover"`
	Evaluate     string `json:"evaluate"`
	TotalCount   string `json:"total_count"`
	PlayCount    string `json:"play_count"`
	DanmakuCount string `json:"danmaku_count"`
	Finish       string `json:"is_finish"`
	Badge        string `json:"badge"`
	SeasonStatus int    `json:"season_status"`
	Favorites    string `json:"favorites"`
	NewEp        struct {
		Aid    string `json:"av_id"`
		Cover  string `json:"cover"`
		Index  string `json:"index"`
		UpTime string `json:"update_time"`
	} `json:"new_ep"`
}

type SeasonInfo struct {
	SeasonID   int64 `json:"season_id"`
	SeasonType int   `json:"season_type"`
	EpisodeID  int   `json:"episode_id"`
}
