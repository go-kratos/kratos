package blocked

import (
	xtime "go-common/library/time"
)

// const publish
const (
	// ptype.
	PublishOfficial      = int8(1) // 官方公告
	PublishWeekCommunity = int8(2) // 社区周报
	PublishFeatureBuild  = int8(3) // 功能建设
	PublishHotCommunity  = int8(4) // 社区热点
	// stick_status.
	PublishStickON  = int8(1) // 置顶
	PublishStickOFF = int8(0) // 不置顶
)

// var publish
var (
	PTypeDesc = map[int8]string{
		PublishOfficial:      "官方公告",
		PublishWeekCommunity: "社区周报",
		PublishFeatureBuild:  "功能建设",
		PublishHotCommunity:  "社区热点",
	}
	SStatusDesc = map[int8]string{
		PublishStickON:  "置顶",
		PublishStickOFF: "不置顶",
	}
)

// Publish is blocked_publish model.
type Publish struct {
	ID                int64      `gorm:"column:id"  json:"id"`
	Title             string     `gorm:"column:title" json:"title"`
	SubTitle          string     `gorm:"column:sub_title" json:"sub_title"`
	Type              int8       `gorm:"column:ptype" json:"type"`
	PublishStatus     int8       `gorm:"column:publish_status" json:"publish_status"`
	StickStatus       int8       `gorm:"column:stick_status" json:"stick_status"`
	Status            int8       `gorm:"column:status" json:"status"`
	Content           string     `gorm:"column:content" json:"content"`
	URL               string     `gorm:"column:url" json:"url"`
	OPID              int64      `gorm:"column:oper_id" json:"oper_id"`
	ShowTime          xtime.Time `gorm:"column:show_time" json:"show_time"`
	PublishTypeDesc   string     `gorm:"-" json:"publish_type_desc"`
	PublishStatusDesc string     `gorm:"-" json:"publish_status_desc"`
	StickStatusDesc   string     `gorm:"-" json:"stick_status_desc"`
	OPName            string     `gorm:"-" json:"oname"`
	CTime             xtime.Time `gorm:"column:ctime" json:"-"`
	MTime             xtime.Time `gorm:"column:mtime" json:"-"`
}

// PublishList is publish list.
type PublishList struct {
	Count int        `json:"count"`
	Order string     `json:"order"`
	Sort  string     `json:"sort"`
	PN    int        `json:"pn"`
	PS    int        `json:"ps"`
	IDs   []int64    `json:"-"`
	List  []*Publish `json:"list"`
}

// TableName publish tablename
func (*Publish) TableName() string {
	return "blocked_publish"
}
