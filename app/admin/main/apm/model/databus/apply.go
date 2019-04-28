package databus

import "go-common/library/time"

// TableName case tablename
func (*Apply) TableName() string {
	return "group_apply"
}

// Apply apply model
type Apply struct {
	ID          int       `gorm:"column:id" json:"id"`
	Group       string    `gorm:"column:group" json:"group"`
	Cluster     string    `gorm:"column:cluster" json:"cluster"`
	TopicRemark string    `gorm:"column:topic_remark" json:"topic_remark"`
	TopicID     int       `gorm:"column:topic_id" json:"topic_id"`
	TopicName   string    `gorm:"column:topic_name" json:"topic"`
	AppID       int       `gorm:"column:app_id" json:"app_id"`
	Project     string    `gorm:"column:project" json:"project"`
	Operation   int8      `gorm:"column:operation" json:"operation"`
	State       int8      `gorm:"column:state" json:"state"`
	Operator    string    `gorm:"column:operator" json:"operator"`
	Remark      string    `gorm:"column:remark" json:"remark"`
	Ctime       time.Time `gorm:"column:ctime" json:"ctime"`
	Mtime       time.Time `gorm:"column:mtime" json:"mtime"`
	// notify
	Gid        int       `gorm:"-" json:"notify_gid"`
	Nid        int64     `gorm:"-" json:"notify_id"`
	Offset     string    `gorm:"-" json:"notify_offset"`
	Nstate     int8      `gorm:"-" json:"notify_state"`
	Filter     int8      `gorm:"-" json:"notify_filter"`
	Concurrent int8      `gorm:"-" json:"notify_concurrent"`
	Callback   string    `gorm:"-" json:"notify_callback"`
	Filters    string    `gorm:"-" json:"-"`
	FilterList []*Filter `gorm:"-" json:"filters"`
	Zone       string    `gorm:"-" json:"notify_zone"`
}
