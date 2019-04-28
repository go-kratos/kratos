package model

import "fmt"

const (
	_statFmt = "heartbeat_in_%s_%d"
)

// ErrReport def
type ErrReport struct {
	MobiApp string `form:"mobi_app" validate:"required"`
	Build   int64  `form:"build" validate:"required"`
	Ecode   int    `form:"error_code" validate:"required"`
}

// SuccReport def
type SuccReport struct {
	MobiApp string `json:"mobi_app"`
	Build   int64  `json:"build"`
}

// ToProm def.
func (v *SuccReport) ToProm() string {
	return fmt.Sprintf(_statFmt, v.MobiApp, v.Build)
}

// ToProm def.
func (v *ErrReport) ToProm() string {
	return fmt.Sprintf(_statFmt, v.MobiApp, v.Build)
}
