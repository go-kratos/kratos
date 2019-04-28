package feedback

import "go-common/library/time"

// Feedback feedback session.
type Feedback struct {
	Session *Session `json:"session"`
	Tag     *Tag     `json:"tag"`
}

// Reply feedback reply.
type Reply struct {
	ReplyID string    `json:"reply_id"`
	Type    int8      `json:"type"`
	Content string    `json:"content"`
	ImgURL  string    `json:"img_url"`
	CTime   time.Time `json:"ctime"`
}

// Tag feedback tags.
type Tag struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Platform string `json:"-"`
}

// Session  for Feedback.
type Session struct {
	SessionID int64     `json:"id"`
	Content   string    `json:"content"`
	ImgURL    string    `json:"img_url"`
	State     int8      `json:"state"`
	CTime     time.Time `json:"ctime"`
}

// TagList list tag.
type TagList struct {
	Platforms []*Platform `json:"platforms"`
	Limit     int         `json:"limit"`
}

// Platform for tag info.
type Platform struct {
	EN   string `json:"en"`
	ZH   string `json:"zh"`
	Tags []*Tag `json:"tags"`
}
