package model

import "time"

//DanmuSendMessage
type DanmuSendMessage struct {
	Topic      string                  `json:"topic"`
	MsgID      string                  `json:"msg_id"`
	MsgKey     string                  `json:"msg_key"`
	MsgContent DanmuSendMessageContent `json:"msg_content"`
}

//DanmuSendMessageContent
type DanmuSendMessageContent struct {
	RoomId    int64     `json:"room_id"`
	Uid       int64     `json:"uid"`
	Uname     string    `json:"uname"`
	UserLevel int64     `json:"user_level"`
	Color     string    `json:"color"`
	Msg       string    `json:"msg"`
	Time      time.Time `json:"time"`
}

//GiftSendMessage
type GiftSendMessage struct {
	Topic      string                 `json:"topic"`
	MsgID      string                 `json:"msg_id"`
	MsgKey     string                 `json:"msg_key"`
	MsgContent GiftSendMessageContent `json:"msg_content"`
}

//GiftSendMessageContent
type GiftSendMessageContent struct {
	Uid      int64  `json:"uid"`
	Ruid     int64  `json:"ruid"`
	RoomId   int64  `json:"roomid"`
	GiftId   int64  `json:"giftid"`
	GiftName string `json:"giftName"`
	PayCoin  int64  `json:"pay_coin"`
	Num      int64  `json:"num"`
	CoinType string `json:"coinType"`
}
