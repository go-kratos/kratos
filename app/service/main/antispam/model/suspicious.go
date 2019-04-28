package model

// Suspicious .
type Suspicious struct {
	Id       int64  `json:"id"`
	SenderId int64  `json:"sender_id"`
	Content  string `json:"content"`
	Area     string `json:"area"`
	OId      int64  `json:"oid"`
}

// GetArea .
func (susp *Suspicious) GetArea() string {
	return susp.Area
}

// GetSenderID .
func (susp *Suspicious) GetSenderID() int64 {
	return susp.SenderId
}

// GetID .
func (susp *Suspicious) GetID() int64 {
	return susp.Id
}

// GetOID .
func (susp *Suspicious) GetOID() int64 {
	return susp.OId
}

// GetContent .
func (susp *Suspicious) GetContent() string {
	return susp.Content
}

// SuspiciousResp .
type SuspiciousResp struct {
	Area      string `json:"-"`
	Content   string `json:"content"`
	LimitType string `json:"susp_type"`
}
