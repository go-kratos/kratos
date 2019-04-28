package databus

import "go-common/library/time"

// TableName case tablename
func (*App) TableName() string {
	return "app2"
}

// App app model
type App struct {
	ID        int       `gorm:"column:id" json:"id"`
	AppKey    string    `gorm:"column:app_key" json:"app_key"`
	AppSecret string    `gorm:"column:app_secret" json:"app_secret"`
	Project   string    `gorm:"column:project" json:"cluster"`
	Remark    string    `gorm:"column:remark" json:"remark"`
	Ctime     time.Time `gorm:"column:ctime" json:"ctime"`
	Mtime     time.Time `gorm:"column:mtime" json:"-"`
}
