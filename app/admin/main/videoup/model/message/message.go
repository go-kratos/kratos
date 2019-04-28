package message

// videoup route
const (
	// videoup
	RouteSyncCid       = "sync_cid"
	RouteFirstRound    = "first_round"
	RouteUGCFirstRound = "ugc_first_round"
	RouteSecondRound   = "second_round"
	RouteAddArchive    = "add_archive"
	RouteModifyArchive = "modify_archive"
	RouteDeleteVideo   = "delete_video"
	RouteDeleteArchive = "delete_archive"
	RouteForceSync     = "force_sync"
)

// Videoup msg
type Videoup struct {
	Route     string `json:"route"`
	Fans      int64  `json:"fans,omitempty"`
	Filename  string `json:"filename"`
	Timestamp int64  `json:"timestamp"`
	// cid
	Cid     int64  `json:"cid,omitempty"`
	DMIndex string `json:"dm_index,omitempty"`
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
	SendEmail   bool  `json:"send_email"`
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
}

//VideoDesign video design
type VideoDesign struct {
	Mosaic    []*Mosaic    `json:"mosaic,omitempty"`
	WaterMark []*WaterMark `json:"watermark,omitempty"`
}

//Mosaic mosaic
type Mosaic struct {
	X     int64 `json:"x"`
	Y     int64 `json:"y"`
	W     int64 `json:"w"`
	H     int64 `json:"h"`
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

//WaterMark watermark
type WaterMark struct {
	LOC   int8   `json:"loc,omitempty"`
	URL   string `json:"url,omitempty"`
	MD5   string `json:"md5,omitempty"`
	Start int64  `json:"start,omitempty"`
	End   int64  `json:"end,omitempty"`
	X     int64  `json:"x,omitempty"`
	Y     int64  `json:"y,omitempty"`
}
