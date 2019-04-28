package blocked

import xtime "go-common/library/time"

const (
	// NoticeStateOpen state open.
	NoticeStateOpen = int8(0)
	// NoticeStateClose state close.
	NoticeStateClose = int8(1)
)

var (
	// NoticeStateDesc state open or close.
	NoticeStateDesc = map[int8]string{
		NoticeStateOpen:  "启用",
		NoticeStateClose: "已删除",
	}
)

// Notice notice struct.
type Notice struct {
	ID         int64      `gorm:"column:id" json:"id"`
	Content    string     `gorm:"column:content" json:"content"`
	URL        string     `gorm:"column:url" json:"url"`
	Status     int8       `gorm:"column:status" json:"status"`
	OperID     int64      `gorm:"column:oper_id" json:"oper_id"`
	Ctime      xtime.Time `gorm:"column:ctime" json:"-"`
	Mtime      xtime.Time `gorm:"column:mtime" json:"-"`
	StatusDesc string     `gorm:"-" json:"status_desc"`
	OPName     string     `gorm:"-" json:"oname"`
}

// NoticeList is notice list.
type NoticeList struct {
	Count int64     `json:"total_count"`
	Pn    int       `json:"pn"`
	Ps    int       `json:"ps"`
	List  []*Notice `json:"list"`
}

// TableName notice tablename
func (*Notice) TableName() string {
	return "blocked_notice"
}
