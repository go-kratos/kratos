package api

// int的action
const (
	FollowAdd    int32 = 1
	FollowCancel int32 = 2
	BlackAdd     int32 = 3
	BlackCancel  int32 = 4
)

// relation list type
const (
	Follow int32 = 1
	Fan    int32 = 2
	Black  int32 = 4
)

//ForbidRequest ..
type ForbidRequest struct {
	MID        uint64 `json:"mid" form:"mid" validate:"required,gt=0"`
	ExpireTime uint64 `json:"expire_time" form:"expire_time" validate:"required,gt=0"`
}

//ReleaseRequest ..
type ReleaseRequest struct {
	MID uint64 `json:"mid" form:"mid" validate:"required,gt=0"`
}

// CmsTagRequest 修改cms_tag的请求
type CmsTagRequest struct {
	Mid    int64 `json:"mid" form:"mid" validate:"required"`
	CmsTag int64 `json:"cms_tag" form:"cms_tag"`
}
