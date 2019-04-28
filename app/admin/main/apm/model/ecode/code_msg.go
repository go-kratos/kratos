package codes

import (
	xtime "go-common/library/time"
)

// TableName case tablename
func (*CodeMsg) TableName() string {
	return "code_msg"
}

// CodeMsg ...
type CodeMsg struct {
	ID       int64      `gorm:"column:id" json:"cid"`
	CodeID   int64      `gorm:"column:code_id" json:"code_id"`
	Locale   string     `gorm:"column:locale" json:"locale"`
	Msg      string     `gorm:"column:msg" json:"msg"`
	Operator string     `gorm:"column:operator" json:"c_operator"`
	CTime    xtime.Time `gorm:"column:ctime" json:"c_ctime"`
	MTime    xtime.Time `gorm:"column:mtime" json:"c_mtime"`
}
