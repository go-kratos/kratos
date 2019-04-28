package watermark

import (
	"time"
)

const (
	// TypeName 带用户名的水印.
	TypeName = 1
	// TypeUID 带uid的水印.
	TypeUID = 2
	// TypeNewName 用户名和logo位置为上下的水印.
	TypeNewName = 3
	// StatClose 未开启水印.
	StatClose = 0
	// StatOpen 开启水印.
	StatOpen = 1
	// StatPreview 预览水印(不写入数据库).
	StatPreview = 2
	// PosLeftTop 水印位置左上角.
	PosLeftTop = 1
	// PosRightTop 水印位置右上角.
	PosRightTop = 2
	// PosLeftBottom 水印位置左下角.
	PosLeftBottom = 3
	// PosRightBottom 水印位置右下角.
	PosRightBottom = 4
)

// Watermark watermark info.
type Watermark struct {
	ID    int64     `json:"id"`
	MID   int64     `json:"mid"`
	Uname string    `json:"uname"`
	State int8      `json:"state"`
	Ty    int8      `json:"type"`
	Pos   int8      `json:"position"`
	URL   string    `json:"url"`
	MD5   string    `json:"md5"`
	Info  string    `json:"info"`
	Tip   string    `json:"tip"`
	CTime time.Time `json:"ctime"`
	MTime time.Time `json:"mtime"`
}

//WatermarkParam set watermark param
type WatermarkParam struct {
	MID   int64
	State int8
	Ty    int8
	Pos   int8
	Sync  int8
	IP    string
}

// Image image width & height.
type Image struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// IsState check state.
func IsState(st int8) bool {
	return st == StatClose || st == StatOpen || st == StatPreview
}

// IsType check type.
func IsType(ty int8) bool {
	return ty == TypeName || ty == TypeUID || ty == TypeNewName
}

// IsPos check position.
func IsPos(pos int8) bool {
	return pos == PosLeftTop || pos == PosRightTop || pos == PosLeftBottom || pos == PosRightBottom
}

// Msg from passport.
type Msg struct {
	Action string    `json:"action"`
	Old    *UserInfo `json:"old"`
	New    *UserInfo `json:"new"`
}

// UserInfo user modify detail.
type UserInfo struct {
	MID    int64  `json:"mid"`
	Uname  string `json:"uname"`
	UserID string `json:"userid"`
}

//GenWatermark for wm api.
type GenWatermark struct {
	Location string `json:"location"`
	MD5      string `json:"md5"` // 文件的hash值
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}
