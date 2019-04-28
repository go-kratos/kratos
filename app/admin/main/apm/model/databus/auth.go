package databus

import (
	"go-common/library/time"
)

// TableName case tablename
func (*Group) TableName() string {
	return "auth2"
}

// Group group model
type Group struct {
	ID         int       `gorm:"column:id" json:"id"`
	Group      string    `gorm:"column:group" json:"group"`
	AppID      int       `gorm:"column:app_id" json:"app_id"`
	AppKey     string    `gorm:"-" json:"app_key"`
	Project    string    `gorm:"-" json:"project"`
	TopicID    int       `gorm:"column:topic_id" json:"topic_id"`
	Topic      string    `gorm:"-" json:"topic"`
	Cluster    string    `gorm:"-" json:"cluster"`
	Operation  int8      `gorm:"column:operation" json:"operation"`
	IsDelete   int8      `gorm:"column:is_delete" json:"is_delete"`
	Remark     string    `gorm:"column:remark" json:"remark"`
	Alarm      int8      `gorm:"column:alarm;default:1" json:"alarm"`
	Percentage string    `gorm:"column:percentage" json:"percentage"`
	Number     int       `gorm:"column:number" json:"number"`
	Ctime      time.Time `gorm:"column:ctime" json:"ctime"`
	Mtime      time.Time `gorm:"column:mtime" json:"mtime"`
	Callback   string    `gorm:"-" json:"notify_callback"`
	Concurrent string    `gorm:"-" json:"notify_concurrent"`
	Filter     int8      `gorm:"-" json:"notify_filter"`
	Filters    string    `gorm:"-" json:"-"`
	FilterList []*Filter `gorm:"-" json:"filters"`
	State      int8      `gorm:"-" json:"notify_state"`
	Gid        int64     `gorm:"-" json:"notify_gid"`
	Nid        int64     `gorm:"-" json:"notify_id"`
	Zone       string    `gorm:"-" json:"notify_zone"`
}

//Alarm alarm
type Alarm struct {
	Group      string `json:"group"`
	Project    string `json:"project"`
	Alarm      int8   `json:"alarm"`
	Percentage string `json:"percentage"`
}

//Alarms alarms
type Alarms struct {
	Cluster    string    `json:"cluster"`
	Topic      string    `json:"topic"`
	Group      string    `json:"group"`
	Project    string    `json:"project"`
	Alarm      int8      `json:"alarm"`
	Percentage string    `json:"percentage"`
	Diff       []*Record `json:"diff"`
}

// Record diff
type Record struct {
	Partition int32 `json:"partition"`
	Diff      int64 `json:"diff"`
	New       int64 `json:"new"`
}
