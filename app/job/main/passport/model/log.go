package model

// LoginLog login log.
type LoginLog struct {
	Mid       int64  `json:"mid"`
	Timestamp int64  `json:"timestamp"`
	LoginIP   int64  `json:"loginip"`
	Type      int64  `json:"type"`
	Server    string `json:"server"`
}

// PwdLog pwd log.
type PwdLog struct {
	ID        int64  `json:"id"`
	Mid       int64  `json:"mid"`
	Timestamp int64  `json:"timestamp"`
	IP        int64  `json:"ip"`
	OldPwd    string `json:"old_pwd"`
	OldSalt   string `json:"old_salt"`
	NewPwd    string `json:"new_pwd"`
	NewSalt   string `json:"new_salt"`
}
