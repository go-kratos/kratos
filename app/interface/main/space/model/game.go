package model

// Game game struct.
type Game struct {
	Website string `json:"website"`
	Image   string `json:"image"`
	Name    string `json:"name"`
}

// AppGame app game struct.
type AppGame struct {
	GameBaseID int64   `json:"game_base_id"`
	GameName   string  `json:"game_name"`
	GameIcon   string  `json:"game_icon"`
	Grade      float64 `json:"grade"`
	DetailURL  string  `json:"detail_url"`
}
