package live

type Room struct {
	UID              int64  `json:"uid"`
	RoomID           int64  `json:"room_id"`
	Title            string `json:"title"`
	Cover            string `json:"cover"`
	Uname            string `json:"uname"`
	Face             string `json:"face"`
	Online           int    `json:"online"`
	Area             string `json:"area"`
	AreaID           int    `json:"area_id"`
	LiveStatus       int    `json:"live_status"`
	AreaV2ID         int64  `json:"area_v2_id"`
	AreaV2Name       string `json:"area_v2_name"`
	AreaV2ParentID   int64  `json:"area_v2_parent_id"`
	AreaV2ParentName string `json:"area_v2_parent_name"`
	BroadcastType    int    `json:"broadcast_type,omitempty"`
}

type Feed struct {
	RoomID int64  `json:"room_id"`
	Face   string `json:"face"`
}

type Card struct {
	RoomID        int64  `json:"roomid,omitempty"`
	UID           int64  `json:"uid,omitempty"`
	Title         string `json:"title,omitempty"`
	Uname         string `json:"uname,omitempty"`
	ShowCover     string `json:"show_cover,omitempty"`
	Online        int64  `json:"online,omitempty"`
	Attentions    int64  `json:"attentions,omitempty"`
	LiveStatus    int    `json:"live_status,omitempty"`
	BroadcastType int    `json:"broadcast_type,omitempty"`
}
