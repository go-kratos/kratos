package model

//BNDatabus 拜年祭投递消息
type BNDatabus struct {
	Roomid    int64  `json:"room_id"`
	UID       int64  `json:"uid"`
	Uname     string `json:"uname"`
	UserLever int64  `json:"user_level"`
	Color     int64  `json:"color"`
	Msg       string `json:"content"`
	Time      int64  `json:"time"`
	MsgID     string `json:"msg_id"`
	MsgType   int64  `json:"msg_type"`
}
