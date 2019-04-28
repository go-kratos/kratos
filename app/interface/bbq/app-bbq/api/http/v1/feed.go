package v1

import (
	"go-common/app/interface/bbq/app-bbq/model"
	video "go-common/app/service/bbq/video/api/grpc/v1"
	bm "go-common/library/net/http/blademaster"
)

// HotReply 热评
type HotReply struct {
	Hots []*model.Reply `json:"hots,omitempty"`
}

// SvDetail one video detail
type SvDetail struct {
	VideoResponse
	CursorValue string   `json:"cursor_value"` // 透传给客户端，标记在列表中的位置
	ElapsedTime int64    `json:"elapsed_time"` // 从发布到现在时间
	HotReply    HotReply `json:"hot_reply"`    // 热评
}

// FeedListRequest feed/list request
type FeedListRequest struct {
	MID    int64
	Device *bm.Device
	BUVID  string
	Mark   string `json:"mark" form:"mark"`
	Page   int    `json:"page" form:"page" validate:"required"`
	Qn     int64  `json:"qn" form:"qn" validate:"required"`
}

// FeedListResponse feed/list request
type FeedListResponse struct {
	Mark    string      `json:"mark" form:"mark"`
	HasMore bool        `json:"has_more" form:"has_more"`
	List    []*SvDetail `json:"list,omitempty" form:"list"`
	RecList []*SvDetail `json:"rec_list,omitempty" form:"list"`
}

// FeedUpdateNumResponse feed/list request
type FeedUpdateNumResponse struct {
	Num int64 `json:"num"`
}

// SpaceSvListRequest feed/list request
// 所有在空间中的视频列表，都复用该请求，同理回包
type SpaceSvListRequest struct {
	MID        int64
	Size       int
	Device     *bm.Device
	DeviceID   string `json:"device_id" form:"device_id"`
	Qn         int64  `json:"qn" form:"qn" validate:"required"`
	UpMid      int64  `json:"up_mid" form:"up_mid" validate:"required"`
	CursorPrev string `json:"cursor_prev" form:"cursor_prev"` // CursorValue
	CursorNext string `json:"cursor_next" form:"cursor_next"`
}

// SpaceSvListResponse feed/list request
type SpaceSvListResponse struct {
	HasMore     bool                    `json:"has_more" form:"has_more"`
	List        []*SvDetail             `json:"list,omitempty" form:"list"`
	PrepareList []*video.UploadingVideo `json:"prepare_list,omitempty"`
}
