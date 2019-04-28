package need

import "go-common/library/time"

// type and states
const (
	TypeCancel     = 0
	TypeLike       = 1
	TypeDislike    = 2
	VerifyAccept   = 2
	VerifyReject   = 3
	VerifyObserved = 4
	NeedApply      = 5
	NeedVerify     = 6
	NeedReview     = 7
)

//VerifyType is
var (
	VerifyType = map[int]string{
		VerifyAccept:   "采纳",
		VerifyReject:   "驳回",
		VerifyObserved: "待观察",
		NeedApply:      "申请",
		NeedVerify:     "确认",
		NeedReview:     "审核",
	}
)

//TableName needs
func (*NInfo) TableName() string {
	return "needs"
}

//NInfo struct
type NInfo struct {
	ID            int64     `gorm:"column:id" json:"id"`
	Title         string    `gorm:"column:title" json:"title"`
	Content       string    `gorm:"column:content" json:"content"`
	Reporter      string    `gorm:"column:reporter" json:"reporter"`
	Status        int8      `gorm:"column:status" json:"status"`
	LikeCounts    int       `gorm:"column:like_counts" json:"like_counts"`
	DislikeCounts int       `gorm:"column:dislike_counts" json:"dislike_counts"`
	CTime         time.Time `gorm:"column:ctime" json:"ctime"`
	MTime         time.Time `gorm:"column:mtime" json:"mtime"`
	LikeState     int8      `gorm:"-" json:"like_state"`
}

//NAddReq add request struct
type NAddReq struct {
	Title   string `form:"title" validate:"required"`
	Content string `form:"content" validate:"required"`
}

// EmpResp is empty resp.
type EmpResp struct {
}

//NEditReq edit request struct
type NEditReq struct {
	ID      int64  `form:"id" validate:"required"`
	Title   string `form:"title"`
	Content string `form:"content"`
}

//NListReq is list request struct
type NListReq struct {
	Ps       int    `form:"ps" default:"20"`
	Pn       int    `form:"pn" default:"1"`
	Status   int    `form:"status"`
	Reporter string `form:"reporter"`
}

//NListResp is list resp struct
type NListResp struct {
	Data  []*NInfo `json:"data"`
	Total int64    `json:"total"`
}

//NVerifyReq is verify req struct
type NVerifyReq struct {
	ID     int64 `form:"id" validate:"required"`
	Status int   `form:"status" validate:"required"`
}

//TableName user_likes
func (*UserLikes) TableName() string {
	return "user_likes"
}

//UserLikes struct
type UserLikes struct {
	ID       int64     `gorm:"column:id" json:"id"`
	ReqID    int64     `gorm:"column:req_id" json:"req_id"`
	User     string    `gorm:"column:user" json:"user"`
	LikeType int8      `gorm:"column:like_type" json:"like_type"`
	CTime    time.Time `gorm:"column:ctime" json:"ctime"`
	MTime    time.Time `gorm:"column:mtime" json:"mtime"`
}

//Likereq is userlike req struct
type Likereq struct {
	ReqID    int64 `form:"req_id" validate:"required"`
	LikeType int8  `form:"like_type"`
}

//VoteListResp is vote resp struct
type VoteListResp struct {
	Data  []*UserLikes `json:"data"`
	Total int64        `json:"total"`
}
