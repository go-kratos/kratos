package model

// review state const.
const (
	ReviewStateWait = iota
	ReviewStatePass
	ReviewStateNoPass
	ReviewStateArchived
	ReviewStateQueuing = 10
)

// review property const.
const (
	ReviewProperty = iota
	ReviewPropertyFace
	ReviewPropertySign
	ReviewPropertyName
)

// UserPropertyReview is.
type UserPropertyReview struct {
	Mid       int64
	Old       string
	New       string
	State     int8
	Property  int8
	IsMonitor bool
	Extra     string
}
