package model

// game play and type
const (
	PlatPc     = 1
	PlatMobile = 2
	TypeMOBA   = 1
	TypeACT    = 2
	TypeFPS    = 3
	TypeFTG    = 4
	TypeRTS    = 5
	TypeRPG    = 6
)

// game plat map and type map
var (
	PlatMap = map[int]int{
		PlatPc:     PlatPc,
		PlatMobile: PlatMobile,
	}
	TypeMap = map[int]int{
		TypeMOBA: TypeMOBA,
		TypeACT:  TypeACT,
		TypeFPS:  TypeFPS,
		TypeFTG:  TypeFTG,
		TypeRTS:  TypeRTS,
		TypeRPG:  TypeRPG,
	}
)

// Game .
type Game struct {
	ID         int64  `json:"id" form:"id"`
	Title      string `json:"title" form:"title" validate:"required"`
	SubTitle   string `json:"sub_title" form:"sub_title"`
	ETitle     string `json:"e_title" form:"e_title"`
	Plat       int    `json:"plat" form:"plat"`
	Type       int    `json:"type" form:"type"`
	Logo       string `json:"logo" form:"logo" validate:"required"`
	Publisher  string `json:"publisher" form:"publisher"`
	Operations string `json:"operations" form:"operations"`
	PbTime     int64  `json:"pb_time" form:"pb_time"`
	Dic        string `json:"dic" form:"dic"`
	Status     int    `json:"status" form:"status"`
}

// TableName es_game
func (g Game) TableName() string {
	return "es_games"
}
