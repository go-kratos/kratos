package model

// SubtitleReportAddParam .
type SubtitleReportAddParam struct {
	Oid        int64   `form:"oid" validate:"required"`
	SubtitleID int64   `form:"subtitle_id" validate:"required"`
	Tid        int64   `form:"tid" validate:"required"`
	MetaData   string  `form:"metadata" validate:"lt=300"`
	From       float64 `form:"from"`
	To         float64 `form:"to" validate:"required"`
	Content    string  `form:"content" validate:"required"`
}
