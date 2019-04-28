package game

type Info struct {
	GameBaseID  int64   `json:"game_base_id,omitempty"`
	IsOnline    bool    `json:"is_online,omitempty"`
	GameName    string  `json:"game_name,omitempty"`
	GameIcon    string  `json:"game_icon,omitempty"`
	GameStatus  int     `json:"game_status,omitempty"`
	GameLink    string  `json:"game_link,omitempty"`
	GradeStatus int     `json:"grade_status,omitempty"`
	Grade       float64 `json:"grade,omitempty"`
	BookNum     int64   `json:"book_num,omitempty"`
}
