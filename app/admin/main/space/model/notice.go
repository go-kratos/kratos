package model

import "fmt"

// const .
const (
	NoticeLogID              = 241
	NoticeForbid             = 1
	NoticeNoForbid           = 0
	NoticeClear              = "clear"
	NoticeClearAndForbid     = "clear_forbid"
	NoticeUnForbid           = "unforbid"
	NoticeTypeClear          = 1
	NoticeTypeClearAndForbid = 2
	NoticeTypeUnForbid       = 3
)

// NoticeUpArg .
type NoticeUpArg struct {
	Mid   int64  `form:"mid" validate:"min=1"`
	Type  int    `form:"type" validate:"min=1,max=3"`
	UID   int64  `form:"-"`
	Uname string `form:"-"`
}

// Notice .
type Notice struct {
	ID       int64  `json:"id" form:"id"`
	Mid      int64  `json:"mid" form:"mid" validate:"required"`
	Notice   string `json:"notice"`
	IsForbid int    `json:"is_forbid"`
}

// TableName notice
func (c *Notice) TableName() string {
	return fmt.Sprintf("member_up_notice%d", c.Mid%10)
}
