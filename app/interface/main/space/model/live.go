package model

// Live .
type Live struct {
	RoomStatus    int    `json:"roomStatus"`
	LiveStatus    int    `json:"liveStatus"`
	URL           string `json:"url"`
	Title         string `json:"title"`
	Cover         string `json:"cover"`
	Online        int    `json:"online"`
	RoomID        int64  `json:"roomid"`
	RoundStatus   int    `json:"roundStatus"`
	BroadcastType int    `json:"broadcast_type"`
}
