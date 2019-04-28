package databus

import "go-common/library/time"

// TableName case tablename
func (*Notify) TableName() string {
	return "notify"
}

// Notify apply model
type Notify struct {
	ID         int       `gorm:"column:id" json:"id"`
	Gid        int       `gorm:"column:gid" json:"gid"`
	Offset     string    `gorm:"column:offset" json:"offset"`
	State      int8      `gorm:"column:state" json:"state"`
	Filter     int8      `gorm:"column:filter" json:"filter"`
	Concurrent int8      `gorm:"column:concurrent" json:"concurrent"`
	Callback   string    `gorm:"column:callback" json:"callback"`
	Zone       string    `gorm:"column:zone" json:"zone"`
	Ctime      time.Time `gorm:"column:ctime" json:"ctime"`
	Mtime      time.Time `gorm:"column:mtime" json:"mtime"`
}
