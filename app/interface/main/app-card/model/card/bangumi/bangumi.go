package bangumi

type Season struct {
	Aid         int64  `json:"aid,omitempty"`
	SeasonID    int64  `json:"season_id,omitempty"`
	EpisodeID   string `json:"episode_id,omitempty"`
	Title       string `json:"title,omitempty"`
	Cover       string `json:"cover,omitempty"`
	PlayCount   int32  `json:"play_count,omitempty"`
	Favorites   int32  `json:"favorites,omitempty"`
	SeasonType  int8   `json:"season_type,omitempty"`
	TypeBadge   string `json:"type_badge,omitempty"`
	SeasonCover string `json:"season_cover,omitempty"`
	UpdateDesc  string `json:"update_desc,omitempty"`
}

type Update struct {
	SquareCover string `json:"square_cover"`
	Title       string `json:"title"`
	Updates     int    `json:"updates"`
}

type Moe struct {
	ID     int64  `json:"id,omitempty"`
	Title  string `json:"title,omitempty"`
	Cover  string `json:"cover,omitempty"`
	Link   string `json:"link,omitempty"`
	Desc   string `json:"desc,omitempty"`
	Badge  string `json:"badge,omitempty"`
	Square string `json:"square,omitempty"`
}
