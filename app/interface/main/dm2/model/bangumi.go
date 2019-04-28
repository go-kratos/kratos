package model

import "encoding/json"

// Season 番剧季度对象
type Season struct {
	AllowDownload string          `json:"allow_download,omitempty"`
	SeasonID      string          `json:"season_id"`
	IsJump        int             `json:"is_jump"`
	EpisodeStatus int             `json:"episode_status"`
	Title         string          `json:"title"`
	Cover         string          `json:"cover"`
	IsFinish      string          `json:"is_finish"`
	IsStarted     int             `json:"is_started"`
	NewestEpID    string          `json:"newest_ep_id"`
	NewestEpIndex string          `json:"newest_ep_index"`
	TotalCount    string          `json:"total_count"`
	Weekday       string          `json:"weekday"`
	Evaluate      string          `json:"evaluate"`
	Bp            json.RawMessage `json:"rank,omitempty"`
	UserSeason    *struct {
		Attention string `json:"attention"`
	} `json:"user_season,omitempty"`
}
