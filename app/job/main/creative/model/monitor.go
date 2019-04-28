package model

// BAP for wechat msg.
type BAP struct {
	UserName  string `json:"username"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	URL       string `json:"url"`
	Ty        string `json:"type"`
	Token     string `json:"token"`
	Signature string `json:"signature"`
	TimeStamp int64  `json:"timestamp"`
}
