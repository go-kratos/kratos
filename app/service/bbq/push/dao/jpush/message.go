package jpush

// Message .
type Message struct {
	Title       string      `json:"title"`
	ContentType string      `json:"content_type"`
	MsgContent  string      `json:"msg_content"`
	Extras      interface{} `json:"extras"`
}
