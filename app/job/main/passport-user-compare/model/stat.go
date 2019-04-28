package model

// Stat stat
type Stat struct {
	Total     int64 `json:"total"`
	ErrorType int64 `json:"error_type"`
}

// ErrorFix error fix
type ErrorFix struct {
	Action    string `json:"action"`
	Mid       int64  `json:"mid"`
	ErrorType int64  `json:"error_type"`
}
