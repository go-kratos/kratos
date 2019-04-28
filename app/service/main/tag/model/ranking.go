package model

// RankingRegion ranking region or region tag
type RankingRegion struct {
	Tid   int64  `json:"target_id"`
	Tname string `json:"target_name"`
	IsTag int64  `json:"is_tag"`
}

// RankingBangumi ranking bangumi
type RankingBangumi struct {
	ID          int64  `json:"id"`
	SeasonID    int64  `json:"season_id"`
	URL         string `json:"url"`
	Cover       string `json:"cover"`
	Color       string `json:"color"`
	Name        string `json:"title"`
	PlayCount   int64  `json:"play_count"`
	FollowCount int64  `json:"follow_count"`
	Tags        []*Tag `json:"tags"`
}

// BangumiInfo bangumi info
type BangumiInfo struct {
	URL       string `json:"url"`
	Favorites int64  `json:"favorites"`
	PlayCount int64  `json:"play_count"`
	SeasonID  int64  `json:"season_id"`
	Title     string `json:"title"`
	Tags      []*Tag `json:"tags"`
}
