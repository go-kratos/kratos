package model

// Emoji Emoji
type Emoji struct {
	ID        int64  `json:"id"`
	PackageID int64  `json:"package_id"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	State     int32  `json:"state"`
	Sort      int32  `json:"sort"`
	Remark    string `json:"remark"`
}

// EmojiPackage  emoji package
type EmojiPackage struct {
	ID     int64    `json:"id"`
	Name   string   `json:"name"`
	URL    string   `json:"url"`
	State  int32    `json:"state"`
	Sort   int32    `json:"sort"`
	Remark string   `json:"remark"`
	Emojis []*Emoji `json:"emojis"`
}
