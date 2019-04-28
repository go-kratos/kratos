package blocked

import (
	xtime "go-common/library/time"
)

// const opinion
const (
	// vote_desc
	NoVote        = int8(0)
	BlockedVote   = int8(1)
	RuleVote      = int8(2)
	DeleteVote    = int8(3)
	BlockedDelete = int8(4)

	// vote_state
	VoteStateON  = int8(0)
	VoteStateOFF = int8(1)
	// attr
	AttrStateOFF = int8(0) // 匿名
	AttrStateOn  = int8(1) // 展示
)

// var opinion
var (
	VoteDesc = map[int8]string{
		NoVote:        "未投票",
		BlockedVote:   "违规",
		RuleVote:      "不违规",
		DeleteVote:    "弃权",
		BlockedDelete: "违规删除",
	}
	AttrDesc = map[int8]string{
		AttrStateOFF: "匿名",
		AttrStateOn:  "展示",
	}

	VoteStateDesc = map[int8]string{
		VoteStateON:  "正常",
		VoteStateOFF: "删除",
	}
)

// Opinion opinion struct.
type Opinion struct {
	ID            int64      `gorm:"column:id" json:"id"`
	VID           int64      `gorm:"column:vid" json:"vid"`
	CID           int64      `gorm:"column:cid" json:"cid"`
	MID           int64      `gorm:"column:mid" json:"mid"`
	OperID        int64      `gorm:"column:oper_id" json:"oper_id"`
	Vote          int8       `gorm:"column:vote" json:"vote"`
	State         int8       `gorm:"column:state" json:"state"`
	Attr          int8       `gorm:"column:attr" json:"attr"`
	Likes         int        `gorm:"column:likes" json:"likes"`
	Hates         int        `gorm:"column:hates" json:"hates"`
	Content       string     `gorm:"column:content" json:"content"`
	CTime         xtime.Time `gorm:"column:ctime" json:"ctime"`
	UName         string     `gorm:"-" json:"uname"`
	AttrDesc      string     `gorm:"-" json:"attr_desc"`
	VoteDesc      string     `gorm:"-" json:"vote_desc"`
	VoteStateDesc string     `gorm:"-" json:"vote_state_desc"`
	OPName        string     `gorm:"-" json:"oname"`
	Fans          int64      `gorm:"-" json:"fans"`
}

// TableName blocked_opinion tablename
func (*Opinion) TableName() string {
	return "blocked_opinion"
}

// OpinionList is Opinion list.
type OpinionList struct {
	Count int        `json:"count"`
	Order string     `json:"order"`
	Sort  string     `json:"sort"`
	PN    int        `json:"pn"`
	PS    int        `json:"ps"`
	IDs   []int64    `json:"-"`
	List  []*Opinion `json:"list"`
}

// OpinionCaseResult struct.
type OpinionCaseResult struct {
	CID     int64  `gorm:"column:cid"`
	MID     int64  `gorm:"column:mid"`
	VID     int64  `gorm:"column:mid"`
	Content string `gorm:"column:content"`
}
