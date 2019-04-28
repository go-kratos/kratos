package model

import "time"

// Report etc.
type Report struct {
	ID          int64
	Name        string
	DateVersion string
	Val         int64
	Ctime       time.Time
}

// ReportDto etc.
type ReportDto struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	DateVersion string `json:"date_version"`
	Val         int64  `json:"val"`
	Ctime       int64  `json:"ctime"`
}

// ReportPage def.
type ReportPage struct {
	TotalCount int          `json:"total_count"`
	Pn         int          `json:"pn"`
	Ps         int          `json:"ps"`
	Items      []*ReportDto `json:"items"`
}

const (
	//BlockCount block count
	BlockCount = "封禁总数"
	// SecurityLoginCount security login count
	SecurityLoginCount = "二次验证总数"
)
