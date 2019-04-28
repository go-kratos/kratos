package model

// AuthToken auth token.
type AuthToken struct {
	DeviceID  string  `json:"device_id"`  // 客户端唯一key
	RoomID    string  `json:"room_id"`    // 业务房间号
	AccessKey string  `json:"access_key"` // access key用来获取mid
	Platform  string  `json:"platform"`   // 平台, android/ios/h5/web
	MobiApp   string  `json:"mobi_app"`   // mobi_app
	Build     int32   `json:"build"`      // build
	Accepts   []int32 `json:"accepts"`    // accept operations
	// 兼容goim-chat
	Aid int64 `json:"aid"`
	Cid int64 `json:"roomid"`
}
