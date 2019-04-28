package model

// Spam .
const (
	SpamBlack    = 52001
	SpamOverflow = 52002
	SpamRestrict = 52005

	LiveDanmuMsgTypeNormal = 0
)

// LiveDanmu .
type LiveDanmu struct {
	RoomID    int64  `json:"room_id"`
	UID       int64  `json:"uid"`
	UName     string `json:"uname"`
	UserLevel int32  `json:"user_level"`
	Color     int32  `json:"color"`
	Content   string `json:"content"`
	Time      int64  `json:"time"`
	MsgType   int32  `json:"msg_type"`
}

// BnjLiveConfig .
type BnjLiveConfig struct {
	DanmuDtarTime string `json:"danmu_start_time"`
	CommentID     int64  `json:"comment_id"`
	RoomID        int64  `json:"room_id"`
}
