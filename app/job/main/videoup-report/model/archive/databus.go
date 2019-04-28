package archive

import "encoding/json"

const (
	//RouteFirstRound 一转
	RouteFirstRound = "first_round"
	//RouteUGCFirstRound 一转
	RouteUGCFirstRound = "ugc_first_round"
	//RouteSecondRound 二转
	RouteSecondRound = "second_round"
	//RouteAddArchive 新增稿件
	RouteAddArchive = "add_archive"
	//RouteModifyArchive 稿件编辑
	RouteModifyArchive = "modify_archive"
	//RouteAutoOpen 自动开放
	RouteAutoOpen = "auto_open"
	//RouteDelayOpen 定时开放
	RouteDelayOpen = "delay_open"
	//RoutePostFirstRound first_round后续处理
	RoutePostFirstRound = "post_first_round"
)

// Message databus message
type Message struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

//VideoupMsg msg
type VideoupMsg struct {
	Route     string `json:"route"`
	Filename  string `json:"filename"`
	Timestamp int64  `json:"timestamp"`
	// cid
	Cid       int64  `json:"cid,omitempty"`
	DMIndex   string `json:"dm_index,omitempty"`
	SendEmail bool   `json:"send_email"`
	// encode
	Xcode          int8         `json:"xcode"`
	EncodePurpose  string       `json:"encode_purpose,omitempty"`
	EncodeRegionID int16        `json:"encode_region_id,omitempty"`
	EncodeTypeID   int16        `json:"encode_type_id,omitempty"`
	VideoDesign    *VideoDesign `json:"video_design,omitempty"`
	Status         int16        `json:"status,omitempty"`
	// add or modify archive
	Aid         int64 `json:"aid,omitempty"`
	EditArchive bool  `json:"edit_archive,omitempty"`
	EditVideo   bool  `json:"edit_video,omitempty"`
	// ChangeTypeID
	ChangeTypeID bool `json:"change_typeid"`
	// ChangeCopyright
	ChangeCopyright bool `json:"change_copyright"`
	// ChangeCover
	ChangeCover bool `json:"change_cover"`
	// ChangeTitle
	ChangeTitle bool `json:"change_title"`
	// Notify
	Notify bool `json:"send_notify"`
	// MissionID
	MissionID int64 `json:"mission_id,omitempty"`
	// AdminChange
	AdminChange bool   `json:"admin_change,omitempty"`
	FromList    string `json:"from_list"`
	TagChange   bool   `json:"tag_change,omitempty"`
	AddVideos   bool   `json:"add_videos,omitempty"`
}

//VideoDesign 自定义马赛克和水印
type VideoDesign struct {
	Mosaic    []*Mosaic    `json:"mosaic,omitempty"`
	WaterMark []*WaterMark `json:"watermark,omitempty"`
}

//Mosaic 马赛克
type Mosaic struct {
	X     int64 `json:"x"`
	Y     int64 `json:"y"`
	W     int64 `json:"w"`
	H     int64 `json:"h"`
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

//WaterMark 水印
type WaterMark struct {
	LOC   int8   `json:"loc,omitempty"`
	URL   string `json:"url,omitempty"`
	MD5   string `json:"md5,omitempty"`
	Start int64  `json:"start,omitempty"`
	End   int64  `json:"end,omitempty"`
	X     int64  `json:"x,omitempty"`
	Y     int64  `json:"y,omitempty"`
}
