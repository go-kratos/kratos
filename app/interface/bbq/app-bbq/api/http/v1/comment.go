package v1

import (
	"go-common/app/interface/bbq/app-bbq/model"
)

//CommentList 评论列表
type CommentList struct {
	OReply *model.CursorRes `json:"oreply,omitempty"`
}

//CommentCursorReq 游标获取评论请求参数
type CommentCursorReq struct {
	SvID   int64  `form:"oid" validate:"gt=0,required"`
	RpID   int64  `form:"rpid"`
	Sort   int16  `form:"sort"`
	MaxID  int64  `form:"max_id"`
	MinID  int64  `form:"min_id"`
	Size   int64  `form:"size" validate:"min=0"`
	Plat   int64  `form:"plat"`
	Build  int64  `form:"build"`
	Access string `form:"access_key"`
	Type   int64
	MID    int64
}

//CommentAddReq 发表评论请求参数
type CommentAddReq struct {
	SvID      int64  `form:"oid" validate:"gt=0,required"`
	Root      int64  `form:"root"`
	Parent    int64  `form:"parent"`
	At        string `form:"at"`
	Message   string `form:"message" validate:"required"`
	Plat      int16  `form:"plat"`
	Device    string `form:"device"`
	Code      string `form:"code"`
	Type      int64
	AccessKey string
}

//CommentLikeReq 评论点怎/取消请求参数
type CommentLikeReq struct {
	SvID      int64 `form:"oid" validate:"gt=0,required"`
	RpID      int64 `form:"rpid" validate:"gt=0,required"`
	Action    int16 `form:"action" validate:"min=0"`
	Type      int64
	AccessKey string
}

//CommentReportReq 评论举报参数
type CommentReportReq struct {
	SvID      int64  `form:"oid" validate:"required"`
	RpID      int64  `form:"rpid" validate:"gt=0,required"`
	Reason    int16  `form:"reason" validate:"min=0"`
	Content   string `form:"content" validate:"min=2,max=200"`
	Type      int64
	AccessKey string
}

//CommentListReq 评论列表请求参数
type CommentListReq struct {
	SvID   int64  `form:"oid" validate:"gt=0,required"`
	Sort   int16  `form:"sort"`
	NoHot  int16  `form:"nohot"`
	Pn     int64  `form:"pn"`
	Ps     int64  `form:"ps"`
	Plat   int64  `form:"plat"`
	Build  int64  `form:"build"`
	Access string `form:"access_key"`
	Type   int64
	MID    int64
}

//CommentSubCursorReq 游标获取子回复及子评论定位参数
type CommentSubCursorReq struct {
	SvID   int64 `form:"oid" validate:"gt=0,required"`
	Sort   int16 `form:"sort" validate:"min=0"`
	Root   int64 `form:"root" validate:"gt=0,required"`
	RpID   int64 `form:"rpid" validate:"min=0"`
	Size   int64 `form:"size" validate:"min=0"`
	MaxID  int64 `form:"max_id" validate:"min=0"`
	MinID  int64 `form:"min_id" validate:"min=0"`
	Type   int64
	Plat   int64  `form:"plat"`
	Build  int64  `form:"build"`
	Access string `form:"access_key"`
}
