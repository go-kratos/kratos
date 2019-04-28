package model

import "encoding/json"

// Message databus
type Message struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// databus const
const (
	RouteFirstRoundForbid = "first_round_forbid"
	RouteSecondRound      = "second_round"
	RouteAutoOpen         = "auto_open"
	RouteDelayOpen        = "delay_open"
	RouteDeleteArchive    = "delete_archive"
	RouteForceSync        = "force_sync"
)

// Videoup message for videoup2BVC
type Videoup struct {
	Route     string `json:"route"`
	Timestamp int64  `json:"timestamp"`
	Aid       int64  `json:"aid"`
}
