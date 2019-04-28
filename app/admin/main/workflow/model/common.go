package model

// Page common page struct for list api
type Page struct {
	Num   int `json:"num"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

// CommonResponse .
type CommonResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
}

// CommonExtraDataResponse .
type CommonExtraDataResponse struct {
	*CommonResponse
	Data map[string]interface{} `json:"data"` //map[gid]interface{}
}

// SourceQueryResponse
type SourceQueryResponse struct {
	*CommonResponse
	Data map[string]interface{} `json:"data"`
}
