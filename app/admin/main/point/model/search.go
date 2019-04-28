package model

// SearchPage struct.
type SearchPage struct {
	PN    int `json:"num"`
	PS    int `json:"size"`
	Total int `json:"total"`
}

// SearchData search result detail.
type SearchData struct {
	Order  string          `json:"order"`
	Sort   string          `json:"sort"`
	Page   *SearchPage     `json:"page"`
	Result []*PointHistory `json:"result"`
}
