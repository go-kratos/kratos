package archive

import "go-common/library/time"

// VideoStatus
const (
	VideoUploadInfo      = int8(0)
	VideoXcodeSDFail     = int8(1)
	VideoXcodeSDFinish   = int8(2)
	VideoXcodeHDFail     = int8(3)
	VideoXcodeHDFinish   = int8(4)
	VideoDispatchRunning = int8(5)
	VideoDispatchFinish  = int8(6)
	VideoStatusOpen      = int16(0)
	VideoStatusAccess    = int16(10000)
	VideoStatusWait      = int16(-1)
	VideoStatusRecicle   = int16(-2)
	VideoStatusLock      = int16(-4)
	VideoStatusXcodeFail = int16(-16)
	VideoStatusSubmit    = int16(-30)
	VideoStatusDelete    = int16(-100)
	XcodeFailZero        = 0
)

// Video is archive_video model.
type Video struct {
	// ID          int64     `json:"-"`
	Aid      int64  `json:"aid"`
	Title    string `json:"title"`
	Desc     string `json:"desc"`
	Filename string `json:"filename"`
	// SrcType     string    `json:"-"`
	// Cid         int64     `json:"-"`
	// Duration    int64     `json:"-"`
	// Filesize    int64     `json:"-"`
	// Resolutions string    `json:"-"`
	Index int `json:"index"`
	// Playurl     string    `json:"-"`
	Status     int16  `json:"status"`
	StatusDesc string `json:"status_desc"`
	FailCode   int8   `json:"fail_code"`
	FailDesc   string `json:"fail_desc"`
	// XcodeState  int8      `json:"-"`
	// Attribute   int32     `json:"-"`
	CTime time.Time `json:"ctime"`
	// MTime       time.Time `json:"-"`
}
