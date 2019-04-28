package live

import (
	"encoding/json"
)

type Room struct {
	UID              int64  `json:"uid,omitempty"`
	RoomID           int64  `json:"room_id,omitempty"`
	Title            string `json:"title,omitempty"`
	Cover            string `json:"cover,omitempty"`
	Uname            string `json:"uname,omitempty"`
	Face             string `json:"face,omitempty"`
	Online           int32  `json:"online,omitempty"`
	LiveStatus       int8   `json:"live_status,omitempty"`
	AreaV2ParentID   int64  `json:"area_v2_parent_id,omitempty"`
	AreaV2ParentName string `json:"area_v2_parent_name,omitempty"`
	AreaV2ID         int64  `json:"area_v2_id,omitempty"`
	AreaV2Name       string `json:"area_v2_name,omitempty"`
	BroadcastType    int    `json:"broadcast_type,omitempty"`
}

type Card struct {
	RoomID        int64  `json:"roomid,omitempty"`
	UID           int64  `json:"uid,omitempty"`
	Title         string `json:"title,omitempty"`
	Uname         string `json:"uname,omitempty"`
	ShowCover     string `json:"show_cover,omitempty"`
	Online        int32  `json:"online,omitempty"`
	LiveStatus    int8   `json:"live_status,omitempty"`
	BroadcastType int    `json:"broadcast_type,omitempty"`
}

type TopicHot struct {
	TID      int    `json:"topic_id"`
	TName    string `json:"topic_name"`
	Picture  string `json:"picture"`
	ImageURL string `json:"-"`
}

type TopicImage struct {
	ImageSrc    string `json:"image_src"`
	ImageWidth  int    `json:"image_width"`
	ImageHeight int    `json:"image_height"`
}

type DynamicHot struct {
	ID           int64    `json:"dynamic_id"`
	AuditStatus  int      `json:"audit_status"`
	DeleteStatus int      `json:"delete_status"`
	MID          int64    `json:"mid"`
	NickName     string   `json:"nick_name"`
	FaceImg      string   `json:"face_img"`
	RidType      int      `json:"rid_type"`
	RID          int64    `json:"rid"`
	ViewCount    int64    `json:"view_count"`
	CommentCount int64    `json:"comment_count"`
	RcmdReason   string   `json:"rcmd_reason"`
	DynamicText  string   `json:"dynamic_text"`
	ImgCount     int      `json:"img_count"`
	Imgs         []string `json:"imgs"`
}

func (t *TopicHot) TopicJSONChange() (err error) {
	var tmp TopicImage
	if err = json.Unmarshal([]byte(t.Picture), &tmp); err != nil {
		return
	}
	t.ImageURL = tmp.ImageSrc
	return
}
