package databus

import (
	"go-common/library/time"
)

// TableName case tablename
func (*Topic) TableName() string {
	return "topic"
}

// Topic topic
type Topic struct {
	ID      int       `gorm:"column:id" json:"id"`
	Topic   string    `gorm:"column:topic" json:"topic"`
	Cluster string    `gorm:"column:cluster" json:"cluster"`
	Remark  string    `gorm:"column:remark" json:"remark"`
	Ctime   time.Time `gorm:"column:ctime" json:"ctime"`
	Mtime   time.Time `gorm:"column:mtime" json:"-"`
}
