package model

// VideoSearch 视频搜索json
type VideoSearch struct {
	ID      int32  `json:"id"`
	Title   string `json:"title"`
	Play    int32  `json:"play"`
	Review  int32  `json:"review"`
	PubTime int32  `json:"pubtime"`
}

// UserSearch 用户搜索json
type UserSearch struct {
	ID      int64  `json:"id"`
	Uname   string `json:"uname"`
	Fans    int64  `json:"fans"`
	Article int64  `json:"article"`
}

// Sug Sug文本结构
type Sug struct {
	ID    int64  `json:"id"`
	Term  string `json:"term"`
	Type  string `json:"type"`
	Score int64  `json:"score"`
}
