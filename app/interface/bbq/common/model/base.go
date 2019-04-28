package model

// Base .
type Base struct {
	App           string `json:"mobi_app" form:"mobi_app"`
	Client        string `json:"platform" form:"platform"`
	Version       string `json:"version" form:"version"`
	Channel       string `json:"channel" form:"channel"`
	Location      string `json:"location" form:"location"`
	QueryID       string `json:"query_id" form:"query_id"`
	BUVID         string `json:"buvid"`
	SVID          int    `json:"svid" form:"svid"`
	TotalDuration int    `json:"total_duration" form:"total_duration"`
	PlayDuration  int    `json:"duration" form:"duration"`
	DataType      int    `json:"data_type" form:"data_type"`
	From          string `json:"from" form:"from"`
	PFrom         string `json:"pfrom" form:"pfrom"`
	FromID        string `json:"from_id" form:"from_id"`
	PFromID       string `json:"pfrom_id" form:"pfrom_id"`
}
