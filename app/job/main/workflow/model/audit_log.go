package model

// AuditLog .
type AuditLog struct {
	Order  string    `json:"order"`
	Sort   string    `json:"sort"`
	Page   *Pager    `json:"page"`
	Result []*LogRes `json:"result"`
}

// LogRes .
type LogRes struct {
	Int1 int64 `json:"int_1"`
	Oid  int64 `json:"oid"`
}
