package blocked

import xtime "go-common/library/time"

// AutoCase is blocked_auto_case model.
type AutoCase struct {
	ID          int64      `gorm:"column:id" json:"id"`
	Platform    int8       `gorm:"column:platform" json:"platform"`
	OPID        int64      `gorm:"column:oper_id" json:"oper_id"`
	ReasonStr   string     `gorm:"column:reasons" json:"-"`
	Reasons     []int64    `gorm:"-" json:"reasons"`
	ReportScore int        `gorm:"column:report_score" json:"report_score"`
	Likes       int        `gorm:"column:likes" json:"likes"`
	CTime       xtime.Time `gorm:"column:ctime" json:"ctime"`
	MTime       xtime.Time `gorm:"column:mtime" json:"mtime"`
}

// TableName AutoCase tablename
func (*AutoCase) TableName() string {
	return "blocked_auto_case"
}
