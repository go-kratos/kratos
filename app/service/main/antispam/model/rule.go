package model

const (
	// LimitTypeDefault .
	LimitTypeDefault = "limit"
	// LimitTypeRestrict .
	LimitTypeRestrict = "restrict"
	// LimitTypeBlack .
	LimitTypeBlack = "black"
	// LimitTypeWhite .
	LimitTypeWhite = "white"

	// LimitScopeLocal .
	LimitScopeLocal = "local"
	// LimitScopeGlobal .
	LimitScopeGlobal = "global"
)

// AggregateRule .
type AggregateRule struct {
	Area      string `json:"area"`
	LimitType string `json:"limit_type"`

	GlobalAllowedCounts int64 `json:"global_allowed_counts"`
	LocalAllowedCounts  int64 `json:"local_allowed_counts"`
	GlobalDurationSec   int64 `json:"global_dur"`
	LocalDurationSec    int64 `json:"local_dur"`
}

// Rule .
type Rule struct {
	ID            int64  `json:"id"`
	Area          string `json:"area"`
	AllowedCounts int64  `json:"allowed_counts"`
	LimitType     string `json:"limit_type"`
	LimitScope    string `json:"limit_scope"`
	DurationSec   int64  `json:"dur"`
}
