package model

// IdxSeason is the struct of season in zone index page
type IdxSeason struct {
	SeasonID int64  `json:"season_id"`
	Title    string `json:"title"`
	Cover    string `json:"cover"`
	Upinfo   string `json:"upinfo"`
}

// IdxData is the zone index data struct in http response
type IdxData struct {
	List  []*IdxSeason `json:"list"`
	Pager *IdxPager    `json:"pager"`
}

// IdxPager is the pager struct to return in zone index page
type IdxPager struct {
	CurrentPage int `json:"current_page"`
	TotalItems  int `json:"total_items"`
	TotalPages  int `json:"total_pages"`
	PageSize    int `json:"page_size"`
}
