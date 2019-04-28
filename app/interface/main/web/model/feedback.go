package model

// Feedback feedback param struct.
type Feedback struct {
	Aid     int64
	Mid     int64
	TagID   int64
	Buvid   string
	Content *Content
	Browser string
	Version string
	Email   string
	QQ      string
	Other   string
}

// Content Content struct.
type Content struct {
	Reason string `json:"reason"`
	URL    string `json:"url"`
}
