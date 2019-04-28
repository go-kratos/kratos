package model

// Business business struct
type Business struct {
	Type   int32  `json:"type"`
	Name   string `json:"name"`
	Appkey string `json:"app_key"`
	Remark string `json:"remark"`
	Alias  string `json:"alias"`
}
