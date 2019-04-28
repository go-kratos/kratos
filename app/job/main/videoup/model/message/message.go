package message

import (
	"encoding/json"
	"go-common/app/job/main/videoup/model/archive"
)

//RouteVideocovers routes.
const (
	// bvc
	RouteVideocovers     = "videocovers"
	RouteBFSVideocovers  = "bfs_videocovers"
	RouteUploadInfo      = "upload_info"
	RouteXcodeSdFinish   = "xcode_sd_finish"
	RouteXcodeSDFail     = "xcode_sd_fail"
	RouteXcodeHDFinish   = "xcode_hd_finish"
	RouteXcodeHDFail     = "xcode_hd_fail"
	RouteDispatchRunning = "dispatch_running"
	RouteDispatchFinish  = "dispatch_finish"
	RouteVideoshotpv     = "bfs_videoshotpv"
	// videoup
	RouteSyncCid          = "sync_cid"
	RouteFirstRound       = "first_round"
	RouteUGCFirstRound    = "ugc_first_round"
	RouteSecondRound      = "second_round"
	RouteAddArchive       = "add_archive"
	RouteModifyArchive    = "modify_archive"
	RouteModifyVideo      = "modify_video"
	RouteDeleteArchive    = "delete_archive"
	RouteDeleteVideo      = "delete_video"
	RouteDelayOpen        = "delay_open"
	RouteAutoOpen         = "auto_open"
	RouteForceSync        = "force_sync"
	RouteFirstRoundForbid = "first_round_forbid"
	RoutePostFirstRound   = "post_first_round"
	// bvc video_capable
	CanPlay    = 0
	CanNotPlay = 1
)

// BvcVideo from bvc video info.
type BvcVideo struct {
	Route     string `json:"route"`
	Filename  string `json:"filename"`
	Timestamp int64  `json:"timestamp"`
	// covers
	Count     int    `json:"count,omitempty"`
	URLFormat string `json:"url_format,omitempty"`
	Deadline  int64  `json:"deadline,omitempty"`
	// video
	Filesize    int64    `json:"filesize,omitempty"`
	Duration    int64    `json:"duration,omitempty"`
	Width       int64    `json:"width,omitempty"`
	Height      int64    `json:"height,omitempty"`
	Rotate      int8     `json:"rotate,omitempty"`
	PlayURL     string   `json:"playurl,omitempty"`
	FailInfo    string   `json:"failinfo,omitempty"`
	Resolutions string   `json:"resolutions,omitempty"`
	BinURL      string   `json:"bin_url"`
	ImgURLs     []string `json:"img_urls"`
}

// ArcResult archive result databus message
type ArcResult struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// Videoup from videoup api.
type Videoup struct {
	Route     string `json:"route"`
	Filename  string `json:"filename"`
	Timestamp int64  `json:"timestamp"`
	// cid
	Cid int64 `json:"cid,omitempty"`
	// encode
	Xcode          int8   `json:"xcode,omitempty"`
	EncodePurpose  string `json:"encode_purpose,omitempty"`
	EncodeRegionID int16  `json:"encode_region_id,omitempty"`
	Status         int16  `json:"status,omitempty"`
	// modify archive
	Aid         int64 `json:"aid,omitempty"`
	EditArchive bool  `json:"edit_archive,omitempty"`
	EditVideo   bool  `json:"edit_video,omitempty"`
	// second_round
	Reply        int  `json:"reply,omitempty"`
	IsSendNotify bool `json:"send_notify,omitempty"`
	// ChangeTypeID
	ChangeTypeID bool `json:"change_typeid,omitempty"`
	// ChangeCopyright
	ChangeCopyright bool `json:"change_copyright,omitempty"`
	// ChangeCover
	ChangeCover bool `json:"change_cover,omitempty"`
	// ChangeTitle
	ChangeTitle bool `json:"change_title,omitempty"`

	MissionID   int64 `json:"mission_id,omitempty"`
	AdminChange bool  `json:"admin_change,omitempty"`
}

// BlogCardMsg 粉丝动态databus消息
type BlogCardMsg struct {
	Card *archive.BlogCard `json:"card"`
}

// StatMsg from archive stat.
type StatMsg struct {
	Type      string
	ID        int64
	Count     int
	Timestamp int64
}
