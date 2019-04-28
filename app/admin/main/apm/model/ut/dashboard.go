package ut

import (
	"go-common/library/time"
)

// PCurveReq is
type PCurveReq struct {
	User      string `form:"user"`
	Path      string `form:"path"`
	StartTime int64  `form:"start_time"`
	EndTime   int64  `form:"end_time"`
}

// PCurveResp .
type PCurveResp struct {
	Pkg        string    `gorm:"column:pkg" json:"pkg,omitempty"`
	Coverage   float64   `gorm:"column:coverage" json:"coverage"`
	Assertions int64     `gorm:"column:assertions" json:"assertions"`
	Failures   int64     `gorm:"column:failures" json:"failures"`
	Panics     int64     `gorm:"column:panics" json:"panics"`
	Passed     int64     `gorm:"column:passed" json:"passed"`
	PassRate   float64   `gorm:"-" json:"pass_rate"`
	MTime      time.Time `gorm:"column:mtime" json:"mtime"`
}

// PCurveDetailResp .
type PCurveDetailResp struct {
	PCurveResp
	Username string `json:"username"`
}
