package archive

import (
	"sync"
)

const (
	//VideoUploadInfo 转码 创建上传
	VideoUploadInfo = 0
	//VideoXcodeSDFail 一转失败
	VideoXcodeSDFail = 1
	//VideoXcodeSDFinish 一转成功
	VideoXcodeSDFinish = 2
	//VideoXcodeHDFail 二转失败
	VideoXcodeHDFail = 3
	//VideoXcodeHDFinish 二转成功
	VideoXcodeHDFinish = 4
	//VideoDispatchRunning 分发中
	VideoDispatchRunning = 5
	//VideoDispatchFinish 分发成功
	VideoDispatchFinish = 6

	//XcodeFailZero fail zero
	XcodeFailZero = 0

	//VideoStatusOpen 开放浏览
	VideoStatusOpen = int16(0)
	//VideoStatusAccess 会员可见
	VideoStatusAccess = int16(10000)
	//VideoStatusWait 待审
	VideoStatusWait = int16(-1)
	//VideoStatusRecicle 打回
	VideoStatusRecicle = int16(-2)
	//VideoStatusLock 锁定
	VideoStatusLock = int16(-4)
	//VideoStatusXcodeFail 转码失败
	VideoStatusXcodeFail = int16(-16)
	//VideoStatusSubmit 创建提交
	VideoStatusSubmit = int16(-30)
	//VideoStatusDelete 删除
	VideoStatusDelete = int16(-100)

	// VideoStatusRecycle video status which be recycled
	VideoStatusRecycle = int16(-2)

	//VideoRelationBind video relation state
	VideoRelationBind = int16(0)
)

//VideoUpInfo  info
type VideoUpInfo struct {
	Nw  *Video
	Old *Video
}

// Video struct
type Video struct {
	ID          int64  `json:"id"`
	Filename    string `json:"filename"`
	Cid         int64  `json:"cid"`
	Aid         int64  `json:"aid"`
	Title       string `json:"eptitle"`
	Desc        string `json:"description"`
	SrcType     string `json:"src_type"`
	Duration    int64  `json:"duration"`
	Filesize    int64  `json:"filesize"`
	Resolutions string `json:"resolutions"`
	Playurl     string `json:"playurl"`
	FailCode    int8   `json:"failinfo"`
	Index       int    `json:"index_order"`
	Attribute   int32  `json:"attribute"`
	XcodeState  int8   `json:"xcode_state"`
	State       int8   `json:"state"`
	Status      int16  `json:"status"`
	CTime       string `json:"ctime"`
	MTime       string `json:"mtime"`
}

// VideoAuditCache video audit count
type VideoAuditCache struct {
	Data map[int16]map[string]int
	sync.Mutex
}

// XcodeTimeCache store video xcode time list
type XcodeTimeCache struct {
	Data map[int8][]int
	sync.Mutex
}

// AttrVal get attribute value.
func (v *Video) AttrVal(bit uint) int32 {
	return (v.Attribute >> bit) & int32(1)
}
