package bangumi

import (
	"encoding/json"

	"go-common/app/interface/main/app-interface/model"
)

// Season for bangumi.
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

// Movie for bangumi.
type Movie struct {
	AllowDownload int8 `json:"allow_download,omitempty"`
	MovieStatus   int  `json:"movie_status"`
	PayUser       *struct {
		Status int `json:"status"`
	} `json:"pay_user"`
	Payment json.RawMessage `json:"payment,omitempty"`
	Season  *struct {
		Actor         json.RawMessage `json:"actor,omitempty"`
		Area          string          `json:"area"`
		Areas         json.RawMessage `json:"areas,omitempty"`
		Cover         string          `json:"cover"`
		Evaluate      string          `json:"evaluate"`
		PubTime       string          `json:"pub_time"`
		SeasonID      int             `json:"season_id"`
		Tags          json.RawMessage `json:"tags,omitempty"`
		Title         string          `json:"title"`
		TotalDuration int             `json:"total_duration"`
	} `json:"season"`
	TrailerAid  int             `json:"trailer_aid"`
	VideoLength int             `json:"video_length"`
	VipQuality  int             `json:"vip_quality"`
	Activity    json.RawMessage `json:"activity,omitempty"`
	List        []struct {
		Cid      int64  `json:"cid"`
		HasAlias bool   `json:"has_alias"`
		Page     int    `json:"page"`
		Part     string `json:"part"`
		Type     string `json:"type"`
		Vid      string `json:"vid"`
	} `json:"list,omitempty"`
}

// Recommend for bangumi.
type Recommend struct {
	Aid    int64  `json:"aid"`
	Cover  string `json:"cover"`
	Status int    `json:"status"`
	Title  string `json:"title"`
}

// Card for bangumi.
type Card struct {
	SeasonID       int64      `json:"season_id"`
	SeasonType     int        `json:"season_type"`
	IsFollow       int        `json:"is_follow"`
	IsSelection    int        `json:"is_selection"`
	Episodes       []*Episode `json:"episodes"`
	SeasonTypeName string     `json:"season_type_name"`
}

// Episode for bangumi card.
type Episode struct {
	ID         int64                `json:"id"`
	Status     int                  `json:"status"`
	Cover      string               `json:"cover"`
	Index      string               `json:"index"`
	IndexTitle string               `json:"index_title"`
	Badges     []*model.ReasonStyle `json:"badges"`
}
