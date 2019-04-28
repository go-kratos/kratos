package blocked

import xtime "go-common/library/time"

// Reason is blocked_reason model.
type Reason struct {
	ID      int        `gorm:"column:id" json:"id"`
	Content string     `gorm:"column:content" json:"content"`
	Reason  string     `gorm:"column:reason" json:"reason"`
	Status  int8       `gorm:"column:status" json:"status"`
	OperID  int        `gorm:"column:oper_id" json:"oper_id"`
	CTime   xtime.Time `gorm:"column:ctime" json:"-"`
	MTime   xtime.Time `gorm:"column:mtime" json:"-"`
	OPName  string     `gorm:"-" json:"oname"`
}

// TableName publish tablename
func (*Reason) TableName() string {
	return "blocked_reason"
}
