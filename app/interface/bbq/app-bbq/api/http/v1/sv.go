package v1

import (
	"go-common/app/interface/bbq/app-bbq/model"
	user "go-common/app/service/bbq/user/api"
	bm "go-common/library/net/http/blademaster"
)

// SvListReq 列表参数
type SvListReq struct {
	MID      int64
	Base     *Base
	RemoteIP string
	PageSize int64 `form:"pagesize" validate:"gt=0,max=20,required"`
}

// SvDetailReq 视频详情参数
type SvDetailReq struct {
	SVID int64 `form:"svid" validate:"gt=0,required"`
}

// SvDetailResponse 视频详情返回
type SvDetailResponse struct {
	IsLike bool `json:"is_like"`
}

//PlayListReq 批量获取playurl接口
type PlayListReq struct {
	CIDs     string `form:"cids" validate:"required"`
	MID      int64
	Device   *bm.Device
	RemoteIP string
}

// VideoResponse 返回视频结构
type VideoResponse struct {
	model.SvInfo
	model.SvStInfo
	SVID      int64         `json:"svid"`
	Tag       string        `json:"tag"`                 // 首选标签
	Tags      []VideoTag    `json:"tags,omitempty"`      // 所有标签
	UserInfo  user.UserBase `json:"user_info,omitempty"` // 用户信息
	Play      VideoPlay     `json:"play,omitempty"`      // 播放信息
	IsLike    bool          `json:"is_like"`             // 是否点赞
	QueryID   string        `json:"query_id"`            // 推荐播放埋点字段，勿删！
	Extension string        `json:"extension"`           // 扩展字段，如title_extra等
}

//VideoTag 视屏标签
type VideoTag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type int16  `json:"type"`
}

// VideoPlay playinfo
type VideoPlay struct {
	SVID           int64       `json:"svid"`
	ExpireTime     int64       `json:"expire_time"`     //过期时间
	FileInfo       []*FileInfo `json:"file_info"`       //分片信息
	Quality        int64       `json:"quality"`         //清晰度
	SupportQuality []int64     `json:"support_quality"` //支持清晰度
	URL            string      `json:"url"`             //基础url
	CurrentTime    int64       `json:"current_time"`    //当前时间戳
}

// CVideo bvc playurl结构
type CVideo struct {
	CID                int64                  `json:"cid"`
	ExpireTime         int64                  `json:"expire_time"`
	FileInfo           map[string][]*FileInfo `json:"file_info,omitempty"`
	Fnval              int64                  `json:"fnval"`
	Fnver              int64                  `json:"fnver"`
	Quality            int64                  `json:"quality"`
	SupportDescription []string               `json:"support_description,omitempty"`
	SupportFormats     []string               `json:"support_formats,omitempty"`
	SupportQuality     []int64                `json:"support_quality,omitempty"`
	URL                string                 `json:"url"`
	VideoCodeCID       int64                  `json:"video_codecid"`
	VideoProject       bool                   `json:"video_project"`
}

// FileInfo bvc fileinfo
type FileInfo struct {
	Ahead      string `json:"ahead"`
	FileSize   int64  `json:"filesize"`
	TimeLength int64  `json:"timelength"`
	Vhead      string `json:"vhead"`
	Path       string `json:"path"`
	URL        string `json:"url"`
	URLBc      string `json:"url_bc"`
}

// SvStatRes 视频统计数据返回
type SvStatRes struct {
	model.SvStInfo
	IsLike      bool `json:"is_like"`
	FollowState int8 `json:"follow_state"`
}

//SvRelReq 相关推荐请求
type SvRelReq struct {
	SVID       int64  `form:"svid"`
	QueryID    string `form:"query_id"`
	MID        int64
	APP        string
	APPVersion string
	Limit      int32
	Offset     int32
	BUVID      string
}
