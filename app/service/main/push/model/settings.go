package model

const (
	// UserSettingArchive up主新投稿提醒
	UserSettingArchive = 1
	// UserSettingLive 主播开播提醒
	UserSettingLive = 2
)

// Settings .
var Settings = map[int]int{
	UserSettingArchive: SwitchOn,
	UserSettingLive:    SwitchOn,
}
