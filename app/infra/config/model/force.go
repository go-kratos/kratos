package model

import "go-common/library/time"

// Force ...
type Force struct {
	ID       int64     `json:"id"`
	AppID    int64     `json:"app_id"`
	HostName string    `json:"hostname"`
	IP       string    `json:"ip"`
	Version  int64     `json:"version"`
	Operator string    `json:"operator"`
	Ctime    time.Time `json:"ctime"`
	Mtime    time.Time `json:"mtime"`
}

// TableName build.
func (Force) TableName() string {
	return "force"
}
