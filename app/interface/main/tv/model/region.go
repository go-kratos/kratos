package model

// Region .
type Region struct {
	PageID    int    `json:"id"`
	Title     string `json:"name"`
	IndexTid  int    `json:"index_tid"`
	IndexType int    `json:"index_type"`
}
