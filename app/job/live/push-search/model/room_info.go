package model

import (
	"encoding/json"
)

// NotifyInfo notify info.
type ApRoomNotifyInfo struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

type LiveDatabusAttention struct {
	Topic	string `json:"topic"`
	MsgId	string `json:"msg_id"`
	MsgContent *AttentionNotifyInfo `json:"msg_content"`
}

// NotifyInfo notify info.
type AttentionNotifyInfo struct {
	Uid      int64  `json:"uid"`
	UpUid    int64  `json:"up_uid"`
	ExtInfo  *ExInfo `json:"ext_info"`
}

type ExInfo struct {
	UpUidFans int `json:"up_uid_fans"`
}

type LiveDatabus struct {
	Topic	string `json:"topic"`
	MsgId	string `json:"msg_id"`
	MsgContent string `json:"msg_content"`
}

type UnameNotifyInfo struct{
	Uid int64 `json:"uid"`
	Uname string `json:"uname"`
	Identification int `json:"identification"`
}

type TableField struct {
	RoomId         int    `json:"roomid"`
	ShortId        int    `json:"short_id"`
	Uid            int64  `json:"uid"`
	UName          string `json:"uname"`
	Area           int    `json:"area"`
	Title          string `json:"title"`
	Tag            string `json:"tags"`
	MTime          string `json:"mtime"`
	CTime          string `json:"ctime"`
	TryTime        string `json:"try_time"`
	Cover          string `json:"cover"`
	UserCover      string `json:"user_cover"`
	LockStatus     string    `json:"lock_status"`
	HiddenStatus   string    `json:"hidden_status"`
	Attentions     int    `json:"attentions"`
	Online         int    `json:"online"`
	LiveTime       string `json:"live_time"`
	AreaV2Id       int    `json:"area_v2_id"`
	AreaV2Name	   string `json:"area_v2_name"`
	AreaV2ParentId int    `json:"area_v2_parent_id"`
	Virtual        int    `json:"virtual"`
	RoundStatus    int    `json:"round_status"`
	OnFlag         int    `json:"on_flag"`
}

type DataMap struct {
	Action string
	Table  string
	New    *TableField
	Old    *TableField
}