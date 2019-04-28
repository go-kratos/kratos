package v1

// Base .
type Base struct {
	App      string `json:"mobi_app" form:"mobi_app"`
	Client   string `json:"platform" form:"platform"`
	Version  string `json:"version" form:"version"`
	Channel  string `json:"channel" form:"channel"`
	Location string `json:"location" form:"location"`
	QueryID  string `json:"query_id" form:"query_id"`
	Module   int    `json:"module_id" form:"module_id"`
	BUVID    string
}

// DataReport 埋点上报字段
type DataReport struct {
	Base
	SVID          int `json:"svid" form:"svid"`
	TotalDuration int `json:"total_duration" form:"total_duration"`
	PlayDuration  int `json:"duration" form:"duration"`
	DataType      int `json:"data_type" form:"data_type"`
	Page          int `json:"page_id" form:"page_id"`
}
