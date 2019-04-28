package bangumi

import "encoding/json"

// Season struct
type Season struct {
	AllowDownload string `json:"allow_download,omitempty"`
	SeasonID      string `json:"season_id,omitempty"`
	IsJump        int    `json:"is_jump,omitempty"`
	Title         string `json:"title,omitempty"`
	Cover         string `json:"cover,omitempty"`
	IsFinish      string `json:"is_finish,omitempty"`
	NewestEpID    string `json:"newest_ep_id,omitempty"`
	NewestEpIndex string `json:"newest_ep_index,omitempty"`
	TotalCount    string `json:"total_count,omitempty"`
	Weekday       string `json:"weekday,omitempty"`
	UserSeason    *struct {
		Attention string `json:"attention,omitempty"`
	} `json:"user_season,omitempty"`
	Player *struct {
		Aid  int64  `json:"aid,omitempty"`
		Vid  string `json:"vid,omitempty"`
		Cid  int64  `json:"cid,omitempty"`
		From string `json:"from,omitempty"`
	} `json:"player,omitempty"`
}

// Movie struct
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

// Card for bangumi.
type Card struct {
	SeasonID    int64      `json:"season_id"`
	SeasonType  int        `json:"season_type"`
	IsFollow    int        `json:"is_follow"`
	IsSelection int        `json:"is_selection"`
	Badge       string     `json:"badge"`
	BadgeType   int        `json:"badge_type"`
	Episodes    []*Episode `json:"episodes"`
}

// Episode for bangumi card.
type Episode struct {
	ID         int64  `json:"id"`
	Badge      string `json:"badge"`
	BadgeType  int    `json:"badge_type"`
	Status     int    `json:"status"`
	Cover      string `json:"cover"`
	Index      string `json:"index"`
	IndexTitle string `json:"index_title"`
}
