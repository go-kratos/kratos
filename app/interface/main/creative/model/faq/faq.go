package faq

import "go-common/library/time"

// PhoneFaqName const
const (
	PhoneFaqName             = "手机投稿FAQ"
	PhoneFaqQuesTypeID       = "7e8eb6dca628490b9f2c089c8c751329"
	PadFaqName               = "iPad投稿FAQ"
	PadFaqQuesTypeID         = "ffd32371f3b94e95ad41e5c387ea62a0"
	FaqUgcProtocolName       = "UGC付费最新协议内容"
	FaqUgcProtocolQuesTypeID = "b899fa132d9c4070b458d5898448f5e3"
)

// Faq str
type Faq struct {
	State            bool   `json:"state"`
	QuestionTypeID   string `json:"questionTypeId"`
	QuestionTypeName string `json:"questionTypeName"`
	URL              string `json:"url"`
}

// Detail help detail and search
type Detail struct {
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
