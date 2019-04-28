package archive

import (
	xtime "go-common/library/time"
)

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
	// VideoStatusRecicle 视频被打回
	VideoStatusRecicle = int16(-2)
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
	// SrcTypeForVupload 合作方嵌套
	SrcTypeForVupload = "vupload"
	// SrcTypeForQQ 腾讯视频
	SrcTypeForQQ = "qq"
	// SrcTypeForHunan 湖南
	SrcTypeForHunan = "hunan"
	// SrcTypeForSohu 搜狐
	SrcTypeForSohu = "sohu"
)

var (
	// XcodeFailCodes is bvc message mapping int value.
	//http://git.bilibili.co/bili_xcode/bili_xcode_docs/blob/master/%E7%B3%BB%E7%BB%9F%E6%95%B0%E6%8D%AE/%E8%BD%AC%E7%A0%81%E9%94%99%E8%AF%AF%E5%8E%9F%E5%9B%A0.md
	XcodeFailCodes = map[string]int8{
		"FileDataUnrecognized":        1,  // 上传文件不是视频
		"VideoTrackAbsent":            2,  // 没有视频轨
		"AudioTrackAbsent":            3,  // 没有音频轨
		"VideoTrackEmpty":             4,  // 视频轨无有效内容
		"AudioTrackEmpty":             5,  // 音频轨无有效内容
		"DurationOverflow":            6,  // 视频过长
		"VideoTooNarrow":              7,  // 画面太窄
		"VideoTooFlat":                8,  // 画面太扁
		"DataCorrupted":               9,  // 文件损坏
		"WatermarkDownloadFail":       10, // 水印图片损坏
		"DurationUnderflow":           11, // 可检测到的时长不足一秒
		"StreamDataCorrupted":         12, // 文件编码数据错误
		"IncorrectDataPackaging":      13, // 文件的封包数据错误
		"UntolerableTimestampJump":    14, // 文件中时间戳有跳变
		"UntolerableTimestampStretch": 15, // 文件中时间戳异常
		"AACDataCorrupted":            16, // AAC音频数据错误

	}
	// XcodeFailMsgs is int value mapping comment.
	XcodeFailMsgs = map[int8]string{
		1:  "文件格式错误，请检查是否上传了错误文件并尝试重新上传",
		2:  "无视频轨，请补充视频轨并重新压制上传",
		3:  "无音频轨，请补充音频轨并重新压制上传",
		4:  "视频轨无有效内容，请补充缺失的视频数据重新压制上传",
		5:  "音频轨无有效内容，请补充缺失的音频数据重新压制上传",
		6:  "单个视频时长超过10小时，请剪辑后通过分P上传",
		7:  "视频画面过窄，请纵向裁剪视频后重新上传",
		8:  "视频画面过扁，请横向裁剪视频后重新上传",
		9:  "视频数据有误，请重新编码后重新上传",
		10: "水印图片损坏",
		11: "单个视频时长不足1秒，请检查视频时长并尝试重新上传",
		12: "文件编码数据错误",
		13: "文件封包数据错误，请重新压制后上传",
		14: "视频时间戳有异常，请修正后重新压制上传",
		15: "视频时间戳有异常，请检查音视频数据并重新压制上传",
		16: "AAC音频数据错误，请重新使用AAC编码后上传",
	}
)

// Video is archive_video model.
type Video struct {
	ID          int64
	Filename    string
	Cid         int64
	Aid         int64
	Title       string
	Desc        string
	SrcType     string
	Duration    int64
	Filesize    int64
	Resolutions string
	Playurl     string
	FailCode    int8
	Index       int
	Attribute   int32
	XcodeState  int8
	Status      int16
	WebLink     string
	Dimensions  string
	CTime       xtime.Time
	MTime       xtime.Time
}

//AuditParam is from video audit
type AuditParam struct {
	IsAudit bool
}

// AttrVal get attribute value.
func (v *Video) AttrVal(bit uint) int32 {
	return (v.Attribute >> bit) & int32(1)
}

// AttrSet set attribute value.
func (v *Video) AttrSet(vl int32, bit uint) {
	v.Attribute = v.Attribute&(^(1 << bit)) | (vl << bit)
}
