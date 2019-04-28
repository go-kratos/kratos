package model

// BusinessAttr .
type BusinessAttr struct {
	ID         int64  `json:"id"`
	Bid        int    `json:"bid"`
	Name       string `json:"name"`
	DealType   int    `json:"deal_type"`
	ExpireTime int    `json:"expire_time"`
	AssignType int    `json:"assign_type"`
	AssignMax  int    `json:"assign_max"`
	GroupType  int    `json:"group_type"`
}
