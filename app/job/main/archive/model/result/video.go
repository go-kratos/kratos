package result

const (
	// video xcode and dispatch state.
	VideoUploadInfo      = 0
	VideoXcodeSDFail     = 1
	VideoXcodeSDFinish   = 2
	VideoXcodeHDFail     = 3
	VideoXcodeHDFinish   = 4
	VideoDispatchRunning = 5
	VideoDispatchFinish  = 6

	XcodeFailZero = 0

	// video status.
	VideoStatusOpen      = int16(0)
	VideoStatusAccess    = int16(10000)
	VideoStatusWait      = int16(-1)
	VideoStatusRecicle   = int16(-2)
	VideoStatusLock      = int16(-4)
	VideoStatusXcodeFail = int16(-16)
	VideoStatusSubmit    = int16(-30)
	VideoStatusDelete    = int16(-100)
)

type VideoUpInfo struct {
	Table  string
	Action string
	Nw     *Video
	Old    *Video
}

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
	Status      int16  `json:"status"`
	CTime       string `json:"ctime"`
	MTime       string `json:"mtime"`
}
