package model

// HistoryData ...
type HistoryData struct {
	RoomId int64 `json:"oid"`
}

// user info platform for wallet info
const (
	PlatformPc      = "pc"
	PlatformIos     = "ios"
	PlatformAndroid = "android"
)
