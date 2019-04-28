package archive

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

	// video relation state
	VideoRelationBind = int16(0)
)

type VideoUpInfo struct {
	Nw  *Video
	Old *Video
}

type Video struct {
	ID          int64  `json:"id"`
	Aid         int64  `json:"aid"`
	Title       string `json:"eptitle"`
	Desc        string `json:"description"`
	Filename    string `json:"filename"`
	SrcType     string `json:"src_type"`
	Cid         int64  `json:"cid"`
	Duration    int64  `json:"duration"`
	Filesize    int64  `json:"filesize"`
	Resolutions string `json:"resolutions"`
	Index       int    `json:"index_order"`
	CTime       string `json:"ctime"`
	MTime       string `json:"mtime"`
	Status      int16  `json:"status"`
	State       int16  `json:"state"`
	Playurl     string `json:"playurl"`
	Attribute   int32  `json:"attribute"`
	FailCode    int8   `json:"failinfo"`
	XcodeState  int8   `json:"xcode_state"`
	WebLink     string `json:"weblink"`
	Dimensions  string `json:"dimensions"`
}
