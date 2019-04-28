package model

//MessageValue  php 格式
type MessageValue struct {
	Topic      string `json:"topic"`
	MsgID      string `json:"msg_id"`
	MsgContent string `json:"msg_content"`
}

//MessageWithoutMsgId  golang投递
type MessageWithoutMsgId struct {
	Topic string `json:"topic"`
	Value string `json:"value"`
}

//DMSendMessageContent
type DMSendMsgContent struct {
	RoomId int64 `json:"room_id"`
	//Uid       int64     `json:"uid"`
	//Uname     string    `json:"uname"`
	//UserLevel int64     `json:"user_level"`
	//Color     string    `json:"color"`
	//Content       string  `json:"content"`
}

type GiftSendMsgContent struct {
	Body BodyMsg `json:"body"`
}

//BodyMsg
type BodyMsg struct {
	MsgID    string `json:"msg_id"`
	Uid      int64  `json:"uid"`
	Ruid     int64  `json:"ruid"`
	RoomId   int64  `json:"roomid"`
	GiftId   int64  `json:"giftid"`
	PayCoin  int64  `json:"pay_coin"`
	Num      int64  `json:"num"`
	CoinType string `json:"coinType"`
}

// GuardBuyMessageContent
type GuardBuyMessageContent struct {
	Uid       int64  `json:"uid"`
	Ruid      int64  `json:"ruid"`
	RoomId    int64  `json:"roomid"`
	Privilege int64  `json:"privilege"`
	Coin      int64  `json:"coin"`
	Num       int64  `json:"num"`
	Type      string `json:"type"`
	Platform  string `json:"platform"`
	IsNew     bool   `json:"is_new"`
}

type TopicCommonMsg struct {
	MsgId  string `json:"msg_id"`
	RoomId int64  `json:"room_id"`
	Value  int64  `json:"value"`
	Cycle  int64  `json:"cycle"`
	Type   int64  `json:"type"`
}

type LiveRoomTagMsg struct {
	MsgId      string `json:"msg_id"`
	RoomId     int64  `json:"room_id"`
	TagId      int64  `json:"tag_id"`
	TagSubId   int64  `json:"tag_sub_id"`
	TagValue   int64  `json:"tag_value"`
	TagExt     string `json:"tag_ext"`
	ExpireTime int64  `json:"expire_time"`
}

type LiveRankListMsg struct {
	RankId     int64           `json:"rank_id"`
	RankType   string          `json:"rank_type"`
	RankList   map[int64]int64 `json:"rank_list"`
	ExpireTime int64           `json:"expire_time"`
}
