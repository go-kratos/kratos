package model

// CaseApplyModifyLog .
type CaseApplyModifyLog struct {
	CID     int64 `json:"case_id"`
	AType   int8  `json:"apply_type"`
	OReason int8  `json:"origin_reason"`
	AReason int8  `json:"apply_reason"`
	Num     int
}
