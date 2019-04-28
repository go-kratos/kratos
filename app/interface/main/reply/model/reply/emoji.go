package reply

// Emoji Emoji
type Emoji struct {
	ID        int64  `json:"id"`
	PackageID int64  `json:"-"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	State     int8   `json:"state"`
	Remark    string `json:"remark"`
}

// EmojiPackage EmojiPackage
type EmojiPackage struct {
	ID     int64    `json:"pid"`
	Name   string   `json:"pname"`
	State  int8     `json:"pstate"`
	URL    string   `json:"purl"`
	Emojis []*Emoji `json:"emojis"`
}
