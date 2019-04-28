package model

// Page is.
type Page struct {
	Num   int `json:"num"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

// CommonPagination is.
type CommonPagination struct {
	Page Page `json:"page"`
}

// MemberPagination is.
type MemberPagination struct {
	Members interface{} `json:"members"`
	*CommonPagination
}

// FaceRecordPagination is.
type FaceRecordPagination struct {
	Records interface{} `json:"records"`
	*CommonPagination
}
