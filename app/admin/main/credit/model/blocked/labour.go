package blocked

import xtime "go-common/library/time"

// LabourAnswerLog is blocked_labour_answer_log model.
type LabourAnswerLog struct {
	ID      int64      `gorm:"column:id" json:"id"`
	UID     int64      `gorm:"column:mid" json:"uid"`
	Score   int16      `gorm:"column:score" json:"score"`
	Content string     `gorm:"column:content" json:"content"`
	Stime   xtime.Time `gorm:"column:start_time" json:"start_time"`
	CTime   xtime.Time `gorm:"column:ctime" json:"-"`
	MTime   xtime.Time `gorm:"column:mtime" json:"-"`
}

// TableName blocked_labour_answer_log tablename
func (*LabourAnswerLog) TableName() string {
	return "blocked_labour_answer_log"
}

// LabourQuestion is blocked_labour_question model.
type LabourQuestion struct {
	ID         int64      `gorm:"column:id" json:"id"`
	Question   string     `gorm:"column:question" json:"question"`
	Ans        int8       `gorm:"column:ans" json:"ans"`
	AVID       int64      `gorm:"column:av_id" json:"av_id"`
	Status     int8       `gorm:"column:status" json:"status"`
	Source     int8       `gorm:"column:source" json:"source"`
	IsDel      int8       `gorm:"column:isdel" json:"isdel"`
	Total      int64      `gorm:"column:total" json:"total"`
	RightTotal int64      `gorm:"column:right_total" json:"right_total"`
	OperID     int        `gorm:"column:oper_id" json:"oper_id"`
	CTime      xtime.Time `gorm:"column:ctime" json:"-"`
	MTime      xtime.Time `gorm:"column:mtime" json:"-"`
	OPName     string     `gorm:"-" json:"oname"`
}

// TableName blocked_labour_question tablename
func (*LabourQuestion) TableName() string {
	return "blocked_labour_question"
}
