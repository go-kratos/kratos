package monitor

import (
	"go-common/library/time"
)

// TableName define table name
func (*Monitor) TableName() string {
	return "monitor"
}

// Monitor .
type Monitor struct {
	ID        int64     `gorm:"column:id" json:"id"`
	AppID     string    `gorm:"column:app_id" json:"app_id"`
	Interface string    `gorm:"column:interface" json:"interface"`
	Count     int64     `gorm:"column:count" json:"count"`
	Cost      int64     `gorm:"column:cost" json:"cost"`
	CTime     time.Time `gorm:"column:ctime" json:"ctime"`
	MTime     time.Time `gorm:"column:mtime" json:"mtime"`
	TempName  string    `gorm:"-" json:"temp_name"`
}

// Data .
type Data struct {
	Interface string   `json:"interface"`
	Counts    []int64  `json:"counts"`
	Costs     []int64  `json:"costs"`
	Times     []string `json:"times"`
}

// MoniRet .
type MoniRet struct {
	XAxis []string `json:"xAxis"`
	Items []*Items `json:"items"`
}

// Items .
type Items struct {
	Interface string  `json:"interface"`
	YAxis     []int64 `json:"yAxis"`
}
