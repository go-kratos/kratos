package model

// TelBindLog bind log
type TelBindLog struct {
	ID        int64  `json:"id"`
	Mid       int64  `json:"mid"`
	Tel       string `json:"tel"`
	Timestamp int64  `json:"timestamp"`
}

// EmailBindLog bind log
type EmailBindLog struct {
	ID        int64  `json:"id"`
	Mid       int64  `json:"mid"`
	Email     string `json:"email"`
	Timestamp int64  `json:"timestamp"`
}
