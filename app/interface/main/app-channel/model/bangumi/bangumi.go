package bangumi

type Season struct {
	Aid           int64  `json:"aid,omitempty"`
	SeasonID      int64  `json:"season_id,omitempty"`
	SeasonStatus  int8   `json:"season_status,omitempty"`
	Title         string `json:"title,omitempty"`
	Index         string `json:"index,omitempty"`
	IndexTitle    string `json:"index_title,omitempty"`
	Cover         string `json:"cover,omitempty"`
	Badge         string `json:"badge,omitempty"`
	PlayCount     int    `json:"play_count,omitempty"`
	Favorites     int    `json:"favorites,omitempty"`
	SeasonType    int8   `json:"season_type,omitempty"`
	TypeBadge     string `json:"type_badge,omitempty"`
	IsFinish      int8   `json:"is_finish,omitempty"`
	TotalCount    int    `json:"total_count,omitempty"`
	NewestEpIndex string `json:"newest_ep_index,omitempty"`
	SeasonCover   string `json:"season_cover,omitempty"`
	UpdateDesc    string `json:"update_desc,omitempty"`
}

type Update struct {
	SquareCover string `json:"square_cover"`
	Title       string `json:"title"`
	Updates     int    `json:"updates"`
}
