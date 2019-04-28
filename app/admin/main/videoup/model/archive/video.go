package archive

import "go-common/library/time"

const (
	// VideoUploadInfo 视频上传完成
	VideoUploadInfo = int8(0)
	// VideoXcodeSDFail 视频转码失败
	VideoXcodeSDFail = int8(1)
	// VideoXcodeSDFinish 一转完成
	VideoXcodeSDFinish = int8(2)
	// VideoXcodeHDFail 二转失败
	VideoXcodeHDFail = int8(3)
	// VideoXcodeHDFinish 二转完成
	VideoXcodeHDFinish = int8(4)
	// VideoDispatchRunning 正在分发
	VideoDispatchRunning = int8(5)
	// VideoDispatchFinish 分发完成
	VideoDispatchFinish = int8(6)
	// VideoStatusOpen 视频开放浏览
	VideoStatusOpen = int16(0)
	// VideoStatusAccess 视频会员可见
	VideoStatusAccess = int16(10000)
	// VideoStatusWait 视频待审
	VideoStatusWait = int16(-1)
	// VideoStatusRecycle 视频被打回
	VideoStatusRecycle = int16(-2)
	// VideoStatusLock 视频被锁定
	VideoStatusLock = int16(-4)
	// VideoStatusXcodeFail 视频转码失败
	VideoStatusXcodeFail = int16(-16)
	// VideoStatusSubmit 视频创建已提交
	VideoStatusSubmit = int16(-30)
	// VideoStatusDelete 视频被删除
	VideoStatusDelete = int16(-100)
	// XcodeFailZero 转码失败
	XcodeFailZero = 0
)

//XcodeStateNames xcode name.
var (
	XcodeStateNames = map[int8]string{
		VideoUploadInfo:      "上传成功",
		VideoXcodeSDFail:     "一转失败",
		VideoXcodeSDFinish:   "一转成功",
		VideoXcodeHDFail:     "二转失败",
		VideoXcodeHDFinish:   "二转成功",
		VideoDispatchRunning: "分发中",
		VideoDispatchFinish:  "分发完成",
	}
)

// Video is archive_video model.
type Video struct {
	ID           int64     `json:"-"`
	Aid          int64     `json:"aid"`
	Title        string    `json:"title"`
	Desc         string    `json:"desc"`
	Filename     string    `json:"filename"`
	SrcType      string    `json:"-"`
	Cid          int64     `json:"cid"`
	Duration     int64     `json:"-"`
	Filesize     int64     `json:"-"`
	Resolutions  string    `json:"-"`
	Index        int       `json:"index"`
	Playurl      string    `json:"-"`
	Status       int16     `json:"status"`
	StatusDesc   string    `json:"status_desc"`
	FailCode     int8      `json:"fail_code"`
	FailDesc     string    `json:"fail_desc"`
	XcodeState   int8      `json:"xcode"`
	Attribute    int32     `json:"-"`
	RejectReason string    `json:"reject_reason"`
	WebLink      string    `json:"weblink"`
	CTime        time.Time `json:"ctime"`
	MTime        time.Time `json:"-"`
}

// AttrSet video Attr set
func (v *Video) AttrSet(attr int32, bit uint) {
	v.Attribute = v.Attribute&(^(1 << bit)) | (attr << bit)
}
