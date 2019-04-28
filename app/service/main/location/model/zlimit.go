package model

// auth const
const (
	Forbidden     = int64(1)
	Allow         = int64(2)
	Formal        = int64(3)
	Pay           = int64(4)
	AllowDown     = int64(1)
	ForbiddenDown = int64(0)

	AuthOK    = 1
	AuthNotOK = 0
)

// for prom
var (
	PlayAuth = map[int64]string{
		Forbidden: "play_forbidden",
		Allow:     "play_allown",
		Formal:    "play_formal_member",
		Pay:       "play_pay_member",
	}

	DownAuth = map[int64]string{
		ForbiddenDown: "down_forbidden",
		AllowDown:     "down_allown",
	}
)

// Auth for auth result
type Auth struct {
	Play int64 `json:"play"`
	Down int64 `json:"down"`
}
