package model

import "go-common/library/time"

// HelpList help list
type HelpList struct {
	Last               bool   `json:"last"`
	ParentTypeID       string `json:"parentTypeId"`
	QuestionTypeDesc   string `json:"questionTypeDesc"`
	QuestionTypeID     string `json:"questionTypeId"`
	QuestionTypeName   string `json:"questionTypeName"`
	QuestionTypeStatus int    `json:"questionTypeStatus"`
	SortNo             int    `json:"sortNo"`
	TypeLevel          int    `json:"typeLevel"`
}

// HelpDeatil help deatil and search
type HelpDeatil struct {
	AllTypeName      string    `json:"allTypeName"`
	AnswerDesc       string    `json:"answerDesc"`
	AnswerFlag       int       `json:"answerFlag"`
	AnswerID         string    `json:"answerId"`
	AnswerImg        string    `json:"answerImg"`
	AnswerTxt        string    `json:"answerTxt"`
	AuditStatus      int       `json:"auditStatus"`
	CompanyID        string    `json:"companyId"`
	CreateID         string    `json:"createId"`
	CreateTime       time.Time `json:"createTime"`
	DocID            string    `json:"docId"`
	LinkFlag         int       `json:"linkFlag"`
	MatchFlag        int       `json:"matchFlag"`
	QuestionID       string    `json:"questionId"`
	QuestionTitle    string    `json:"questionTitle"`
	QuestionTypeID   string    `json:"questionTypeId"`
	QuestionTypeName string    `json:"questionTypeName"`
	UpdateID         string    `json:"updateId"`
	UpdateTime       time.Time `json:"updateTime"`
	UsedFlag         int       `json:"usedFlag"`
}
