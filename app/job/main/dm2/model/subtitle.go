package model

// SubtitleStatus .
type SubtitleStatus uint8

// SubtitleStatus
const (
	SubtitleStatusUnknown SubtitleStatus = iota
	SubtitleStatusDraft
	SubtitleStatusToAudit
	SubtitleStatusAuditBack
	SubtitleStatusRemove
	SubtitleStatusPublish
	SubtitleStatusCheckToAudit
	SubtitleStatusCheckPublish
)

// Subtitle .
type Subtitle struct {
	ID            int64          `json:"id"`
	Oid           int64          `json:"oid"`
	Type          int32          `json:"type"`
	Lan           uint8          `json:"lan"`
	Aid           int64          `json:"aid"`
	Mid           int64          `json:"mid"`
	UpMid         int64          `json:"up_mid"`
	Status        SubtitleStatus `json:"status"`
	SubtitleURL   string         `json:"subtitle_url"`
	PubTime       int64          `json:"pub_time"`
	RejectComment string         `json:"reject_comment"`
}

// SubtitlePub .
type SubtitlePub struct {
	Oid        int64
	Type       int32
	Lan        uint8
	SubtitleID int64
	IsDelete   bool
}

// SubtitleItem .
type SubtitleItem struct {
	From     float64 `json:"from"`
	To       float64 `json:"to"`
	Location int8    `json:"location"`
	Content  string  `json:"content"`
}

// SubtitleBody .
type SubtitleBody struct {
	FontSize        float64         `json:"font_size,omitempty"`
	FontColor       string          `json:"font_color,omitempty"`
	BackgroundAlpha float64         `json:"background_alpha,omitempty"`
	BackgroundColor string          `json:"background_color,omitempty"`
	Stroke          string          `json:"Stroke,omitempty"`
	Bodys           []*SubtitleItem `json:"body"`
}

// SubtitleAuditMsg .
type SubtitleAuditMsg struct {
	SubtitleID int64 `json:"subtitle_id"`
	Oid        int64 `json:"oid"`
}
