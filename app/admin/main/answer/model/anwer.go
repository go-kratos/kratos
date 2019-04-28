package model

import "go-common/library/time"

// TypeInfo type info
type TypeInfo struct {
	ID        int64      `json:"id"`
	Parentid  int64      `json:"-" gorm:"column:parentid"`
	Name      string     `json:"name" gorm:"column:typename"`
	LabelName string     `json:"label_name" gorm:"column:lablename"`
	Subs      []*SubType `json:"subs"`
}

// TableName for gorm.
func (t *TypeInfo) TableName() string {
	return "ans_v3_question_type"
}

// SubType sub type info
type SubType struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	LabelName string `json:"-"`
}

// AnswerHistoryDB info.
type AnswerHistoryDB struct {
	ID                    int64     `json:"id"`
	Hid                   int64     `json:"hid"`
	Mid                   int64     `json:"mid"`
	StartTime             time.Time `json:"start_time"`
	StepOneErrTimes       int8      `json:"step_one_err_times"`
	StepOneCompleteTime   int64     `json:"step_one_complete_time"`
	StepExtraStartTime    time.Time `json:"step_extra_start_time"`
	StepExtraCompleteTime int64     `json:"step_extra_complete_time"`
	StepExtraScore        int64     `json:"step_extra_score"`
	StepTwoStartTime      time.Time `json:"step_two_start_time"`
	CompleteTime          time.Time `json:"complete_time"`
	CompleteResult        string    `json:"complete_result"`
	Score                 int8      `json:"score"`
	IsFirstPass           int8      `json:"is_first_pass"`
	IsPassCaptcha         int8      `json:"is_pass_captcha"`
	PassedLevel           int8      `json:"passed_level"`
	RankID                int       `json:"rank_id"`
	Ctime                 time.Time `json:"ctime"`
	Mtime                 time.Time `json:"mtime"`
}
