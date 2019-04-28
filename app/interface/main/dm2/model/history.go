package model

// SearchPage page struct from search
type SearchPage struct {
	Num   int64 `json:"num"`
	Size  int64 `json:"size"`
	Total int64 `json:"total"`
}

// SearchHistoryResult history dm from search
type SearchHistoryResult struct {
	Page   *SearchPage `json:"page"`
	Result []*struct {
		ID int64 `json:"id"`
	} `json:"result"`
}

// SearchHistoryIdxResult history date index
type SearchHistoryIdxResult struct {
	Page   *SearchPage `json:"page"`
	Result []*struct {
		Date string `json:"date"`
	} `json:"result"`
}
