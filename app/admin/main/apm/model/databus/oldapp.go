package databus

import "go-common/library/time"

// TableName case tablename
func (*OldApp) TableName() string {
	return "app"
}

// OldApp app model
type OldApp struct {
	ID        int64     `gorm:"column:id" json:"id"`
	AppKey    string    `gorm:"column:app_key" json:"app_key"`
	AppSecret string    `gorm:"column:app_secret" json:"app_secret"`
	Cluster   string    `gorm:"column:cluster" json:"cluster"`
	Ctime     time.Time `gorm:"column:ctime" json:"ctime"`
	Mtime     time.Time `gorm:"column:mtime" json:"-"`
}
