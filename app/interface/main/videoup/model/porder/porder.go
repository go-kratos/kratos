package porder

// OfficialIndustryMaps map
var OfficialIndustryMaps = map[int64]int64{
	1: 1,
}

// Config str
type Config struct {
	ID   int64  `json:"id"`
	Tp   int8   `json:"type"`
	Name string `json:"name"`
}

// Game str
type Game struct {
	GameBaseID int64  `json:"game_base_id"`
	GameName   string `json:"game_name"`
	Source     int8   `json:"source"`
}
