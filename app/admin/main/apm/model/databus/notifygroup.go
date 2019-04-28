package databus

import "go-common/library/time"

// TableName case tablename
func (*NotifyGroup) TableName() string {
	return "group_apply"
}

// NotifyGroup apply model
type NotifyGroup struct {
	ID           int       `gorm:"column:id" json:"notify_id"`
	Gid          int       `gorm:"column:gid" json:"notify_gid"`
	Offset       string    `gorm:"column:offset" json:"notify_offset"`
	State        int8      `gorm:"column:state" json:"notify_state"`
	Filter       int8      `gorm:"column:filter" json:"notify_filter"`
	Concurrent   int8      `gorm:"column:concurrent" json:"notify_concurrent"`
	Callback     string    `gorm:"column:callback" json:"notify_callback"`
	Ctime        time.Time `gorm:"column:ctime" json:"notify_ctime"`
	Mtime        time.Time `gorm:"column:mtime" json:"notify_mtime"`
	GGroup       string    `gorm:"column:group" json:"group"`
	GCluster     string    `gorm:"column:cluster" json:"cluster"`
	GTopicRemark string    `gorm:"column:topic_remark" json:"topic_remark"`
	GTopicID     int       `gorm:"column:topic_id" json:"topic_id"`
	GTopicName   string    `gorm:"column:topic_name" json:"topic"`
	GAppID       int       `gorm:"column:app_id" json:"app_id"`
	GState       int8      `gorm:"column:gstate" json:"state"`
	GOperation   int8      `gorm:"column:operation" json:"operation"`
	GOperator    string    `gorm:"column:operator" json:"operator"`
	GRemark      string    `gorm:"column:remark" json:"remark"`
	Filters      string    `gorm:"column:filters" json:"filters"`
	GCtime       time.Time `gorm:"column:gctime" json:"ctime"`
	GMtime       time.Time `gorm:"column:gmtime" json:"mtime"`
}
