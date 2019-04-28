package bws

// GameRes .
const (
	GameResWin  = 1
	GameResFail = 2
)

// ParamPoints points param
type ParamPoints struct {
	Bid int64 `form:"bid" validate:"required"`
	Tp  int64 `form:"tp"`
}

// ParamID point or achievements id param
type ParamID struct {
	Bid int64  `form:"bid" validate:"required"`
	ID  int64  `form:"id"`
	Day string `form:"day"`
}

// ParamAward point or achievements id param
type ParamAward struct {
	Bid int64  `form:"bid" validate:"required"`
	Aid int64  `form:"aid" validate:"required"`
	Key string `form:"key"`
	Mid int64  `form:"mid"`
}

// ParamBinding binding param
type ParamBinding struct {
	Bid int64  `form:"bid" validate:"required"`
	Key string `form:"key" validate:"required"`
}

// ParamUnlock .
type ParamUnlock struct {
	Bid        int64  `form:"bid" validate:"required"`
	Pid        int64  `form:"pid" validate:"required"`
	Key        string `form:"key"`
	Mid        int64  `form:"mid"`
	GameResult int    `form:"game_result"`
}
