package live

// Live for live
type Live struct {
	Mid    int64  `json:"mid,omitempty"`
	RoomID int64  `json:"roomid,omitempty"`
	Title  string `json:"title,omitempty"`
}

// Room for live
type Room struct {
	UID              int64  `json:"uid,omitempty"`
	RoomID           int64  `json:"room_id,omitempty"`
	Title            string `json:"title,omitempty"`
	Cover            string `json:"cover,omitempty"`
	Uname            string `json:"uname,omitempty"`
	Face             string `json:"face,omitempty"`
	Online           int    `json:"online,omitempty"`
	Area             string `json:"area,omitempty"`
	AreaID           int    `json:"area_id,omitempty"`
	LiveStatus       int    `json:"live_status,omitempty"`
	AreaV2ID         int64  `json:"area_v2_id,omitempty"`
	AreaV2Name       string `json:"area_v2_name,omitempty"`
	AreaV2ParentID   int64  `json:"area_v2_parent_id,omitempty"`
	AreaV2ParentName string `json:"area_v2_parent_name,omitempty"`
	BroadcastType    int    `json:"broadcast_type,omitempty"`
}

// Status for live
type Status struct {
	LiveStatus int `json:"live_status,omitempty"`
}

// Glory for live
type Glory struct {
	ID        string `json:"id,omitempty"`
	UID       string `json:"uid,omitempty"`
	GID       string `json:"gid,omitempty"`
	On        string `json:"on,omitempty"`
	GloryInfo *struct {
		ID       string `json:"id,omitempty"`
		Name     string `json:"name,omitempty"`
		Cover    string `json:"pic_url,omitempty"`
		Level    string `json:"level,omitempty"`
		Activity string `json:"activity,omitempty"`
		URI      string `json:"jump_url,omitempty"`
	} `json:"glory_info,omitempty"`
	UserInfo *struct {
		UID        int64  `json:"uid,omitempty"`
		Name       string `json:"uname,omitempty"`
		Face       string `json:"face,omitempty"`
		Level      int    `json:"rlevel,omitempty"`
		LevelColor int64  `json:"rlevel_color,omitempty"`
	} `json:"user_info,omitempty"`
	Version int `json:"version,omitempty"`
}

// Exp for live
type Exp struct {
	Level  int `json:"user_level,omitempty"`
	Master *struct {
		Level int   `json:"level,omitempty"`
		Color int64 `json:"color,omitempty"`
	} `json:"master_level,omitempty"`
	Color int64 `json:"color,omitempty"`
}

// RoomInfo for live
type RoomInfo struct {
	RoomID        int64  `json:"roomid"`
	ShortID       int64  `json:"short_id"`
	Title         string `json:"title,omitempty"`
	Cover         string `json:"cover,omitempty"`
	UserCover     string `json:"user_cover,omitempty"`
	URI           string `json:"uri,omitempty"`
	Mid           int64  `json:"uid,omitempty"`
	Name          string `json:"uname,omitempty"`
	TagName       string `json:"area_v2_name,omitempty"`
	Status        int    `json:"live_status"`
	BroadcastType int    `json:"broadcast_type,omitempty"`
}
