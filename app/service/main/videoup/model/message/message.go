package message

const (
	//RouteSyncCid cid同步
	RouteSyncCid = "sync_cid"
	//RouteFirstRound 一审
	RouteFirstRound = "first_round"
	//RoutePGCSubmit pgc提交
	RoutePGCSubmit = "pgc_submit"
	//RouteDRMSubmit drm提交
	RouteDRMSubmit = "drm_submit"
	//RouteUGCSubmit ugc提交
	RouteUGCSubmit = "ugc_submit"
	//RouteSecondRound 二审
	RouteSecondRound = "second_round"
	//RouteAddArchive 投稿
	RouteAddArchive = "add_archive"
	//RouteModifyArchive 编辑稿件
	RouteModifyArchive = "modify_archive"
	//RouteModifyVideo 编辑视频
	RouteModifyVideo = "modify_video"
	//RouteUserDelete 用户删除  NOTE: after change this route by delete_video
	RouteUserDelete = "user_delete"
	//RouteDeleteVideo 删除视频
	RouteDeleteVideo = "delete_video"
	//RouteDeleteArchive 删除稿件
	RouteDeleteArchive = "delete_archive"
	//RouteForceSync  同步稿件库
	RouteForceSync = "force_sync"
)

//Videoup  messgae
type Videoup struct {
	Route     string `json:"route"`
	Filename  string `json:"filename"`
	Timestamp int64  `json:"timestamp"`
	// cid
	Cid     int64  `json:"cid,omitempty"`
	DMIndex string `json:"dm_index,omitempty"`
	UpFrom  int8   `json:"up_from"`
	// encode
	Xcode          int8   `json:"xcode"`
	EncodePurpose  string `json:"encode_purpose,omitempty"`
	EncodeRegionID int16  `json:"encode_region_id,omitempty"`
	VideoDesign    struct {
		Mosaic    []*Mosaic  `json:"mosaic,omitempty"`
		WaterMark *WaterMark `json:"watermark,omitempty"`
	} `json:"video_design"`
	Status int16 `json:"status,omitempty"`
	// add or modify archive
	Aid         int64 `json:"aid,omitempty"`
	EditArchive bool  `json:"edit_archive,omitempty"`
	EditVideo   bool  `json:"edit_video,omitempty"`
	// MissionID
	MissionID int64 `json:"mission_id,omitempty"`
	// pgc submit
	Submit       int       `json:"submit"`
	TagChange    bool      `json:"tag_change,omitempty"`
	AddVideos    bool      `json:"add_videos,omitempty"`
	ChangeTypeID bool      `json:"change_typeid,omitempty"`
	StaffBox     *StaffBox `json:"staff_box,omitempty"`
}

//Mosaic message
type Mosaic struct {
	X     int64 `json:"x"`
	Y     int64 `json:"y"`
	W     int64 `json:"w"`
	H     int64 `json:"h"`
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

//WaterMark message
type WaterMark struct {
	URL   string `json:"url"`
	MD5   string `json:"md5"`
	Start int64  `json:"start"`
	End   int64  `json:"end"`
	X     int64  `json:"x"`
	Y     int64  `json:"y"`
}
