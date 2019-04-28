package model

const (
	// TypeSeason ...
	TypeSeason = "season"
	// TypeItems ...
	TypeItems = "items"
)

// SugList sug list
type SugList struct {
	SeasonId   int64   `json:"season_id"`
	SeasonName string  `json:"season_name"`
	ItemsID    int64   `json:"items_id"`
	ItemsName  string  `json:"items_name"`
	PicURL     string  `json:"head_url"`
	Score      float64 `json:"score"`
	SugURL     string  `json:"pic_url"`
}

// Season .
type Season struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

// Item .
type Item struct {
	ItemsID int64
	Name    string
	Brief   string
	Img     string
}

// Match .
type Match struct {
	SeasonID int64
	ItemsID  int64
	Type     int
}

// Items .
type Items struct {
	ItemsID int64    `json:"itemsId"`
	Name    string   `json:"name"`
	Img     []string `json:"img"`
}

// ItemsList .
type ItemsList struct {
	Total    int     `json:"total"`
	PageNum  int     `json:"pageNum"`
	PageSize int     `json:"pageSize"`
	List     []Items `json:"list"`
}

// HTTPResponse .
type HTTPResponse struct {
	Code int       `json:"code"`
	Data ItemsList `json:"data"`
}
