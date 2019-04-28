package v1

import (
	user "go-common/app/service/bbq/user/api"
)

//LoginRequest 登陆
type LoginRequest struct {
	NewTag int8 `json:"new_tag" form:"new_tag"`
}

//PhoneCheckResponse ...
type PhoneCheckResponse struct {
	TELStatus int32 `json:"tel_status"`
}

// SpaceUserProfileRequest ...
type SpaceUserProfileRequest struct {
	Upmid int64 `json:"up_mid" form:"up_mid" validate:"required"`
}

// NumResponse 空返回值
type NumResponse struct {
	Num int64 `json:"num"`
}

//UserRelationRequest .
type UserRelationRequest struct {
	UPMID int64 `json:"up_mid" form:"up_mid" validate:"required"`
	// 见上述RelationAction
	Action int32 `json:"action" form:"action"`
}

// UserRelationListResponse 关注、粉丝、拉黑列表结构
type UserRelationListResponse struct {
	HasMore bool        `json:"has_more"`
	List    []*UserInfo `json:"list,omitempty"`
}

//UserLikeAddRequest .
type UserLikeAddRequest struct {
	SVID int64 `json:"svid" form:"svid" validate:"required"`
}

//UserLikeCancelRequest .
type UserLikeCancelRequest struct {
	SVID int64 `json:"svid" form:"svid" validate:"required"`
}

//InviteCodeRequest .
type InviteCodeRequest struct {
	Num    int64  `json:"num" form:"num" validate:"required"`
	Type   string `json:"type"  form:"type" validate:"required"`
	Digit  int64  `json:"digit" form:"digit" validate:"required"`
	Author int64  `json:"author" form:"author" validate:"required"`
}

//CheckInviteCodeRequest .
type CheckInviteCodeRequest struct {
	Code     int64  `json:"code" form:"code" validate:"required"`
	DeviceID string `json:"device_id"  form:"device_id" validate:"required"`
	Uname    string `json:"uname" form:"uname"`
}

// UserInfo 用户相关信息，统一提供对外结构
type UserInfo struct {
	user.UserBase
	user.UserStat
	FollowState int8   `json:"follow_state"`
	CursorValue string `json:"cursor_value"`
}

// UnLikeReq 不感兴趣
type UnLikeReq struct {
	MID  int64 `json:"mid" form:"mid"`
	SVID int64 `json:"svid" form:"svid"`
}
