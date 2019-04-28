package model

// Bangumi bangumi struct.
type Bangumi struct {
	SeasonID      string `json:"season_id"`
	ShareURL      string `json:"share_url"`
	Title         string `json:"title"`
	IsFinish      string `json:"is_finish"`
	Favorites     string `json:"favorites"`
	NewestEpIndex string `json:"newest_ep_index"`
	LastEpIndex   string `json:"last_ep_index"`
	TotalCount    string `json:"total_count"`
	Cover         string `json:"cover"`
	Evaluate      string `json:"evaluate"`
	Brief         string `json:"brief"`
}
