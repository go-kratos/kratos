package model

import accmdl "go-common/app/service/main/account/model"

// TelInfo def.
type TelInfo struct {
	Mid      int64  `json:"mid"`
	Tel      string `json:"tel"`
	JoinIP   string `json:"join_ip"`
	JoinTime int64  `json:"join_time"`
}

// AuditInfo is.
type AuditInfo struct {
	Mid      int64 `json:"mid"`
	BindTel  bool  `json:"bind_tel"`
	BindMail bool  `json:"bind_mail"`
	Rank     int64 `json:"rank"`
	SpaceSta int64 `json:"spacesta"`
}

// ProfileInfo profile info.
type ProfileInfo = accmdl.Profile

// TelRiskInfo tel risk info.
type TelRiskInfo struct {
	TelRiskLevel    int8
	RestoreHistory  *UserEventHistory
	UnicomGiftState int
}

// UnicomGiftState unicom gift state.
type UnicomGiftState struct {
	State int `json:"state"`
}
