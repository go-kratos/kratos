package model

// SearchResult is
type SearchResult struct {
	Code int `json:"code"`
	Data struct {
		Debug string `json:"debug"`
		Order string `json:"order"`
		Page  struct {
			Num   int64 `json:"num"`
			Size  int64 `json:"size"`
			Total int64 `json:"total"`
		} `json:"page"`
		Result []struct {
			Action    string `json:"action"`
			Build     int64  `json:"build"`
			Business  int64  `json:"business"`
			Buvid     string `json:"buvid"`
			Ctime     string `json:"ctime"`
			ExtraData string `json:"extra_data"`
			IP        string `json:"ip"`
			Mid       int64  `json:"mid"`
			Oid       int64  `json:"oid"`
			Platform  string `json:"platform"`
			Type      int64  `json:"type"`
		} `json:"result"`
		Sort string `json:"sort"`
	} `json:"data"`
	Message string `json:"message"`
	TTL     int64  `json:"ttl"`
}
