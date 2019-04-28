package model

// SearchGroup .
type SearchGroup struct {
	Order  string      `json:"order"`
	Sort   string      `json:"sort"`
	Page   *Pager      `json:"page"`
	Result []*GroupRes `json:"result"`
}

// GroupRes .
type GroupRes struct {
	ID  int64 `json:"id"`
	OID int64 `json:"oid"`
}

// Pager .
type Pager struct {
	Num   int64 `json:"num"`
	Size  int64 `json:"size"`
	Total int64 `json:"total"`
}

// SearchChall .
type SearchChall struct {
	Order  string      `json:"order"`
	Sort   string      `json:"sort"`
	Page   *Pager      `json:"page"`
	Result []*ChallRes `json:"result"`
}

// ChallRes .
type ChallRes struct {
	ID  int64 `json:"id"`
	GID int64 `json:"gid"`
	MID int64 `json:"mid"`
	OID int64 `json:"oid"`
}
