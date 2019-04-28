package model

import (
	xtime "go-common/library/time"
)

// EmailState .
type EmailState int8

// const .
const (
	EmailStateSendNone EmailState = 1
	EmailStateSendSucc EmailState = 2
)

// MCNSignState .
type MCNSignState int8

// const .
const (
	// MCNSignStateNoApply 未申请
	MCNSignStateNoApply MCNSignState = 0
	// MCNSignStateOnReview 待审核
	MCNSignStateOnReview MCNSignState = 1
	// MCNSignStateOnReject 已驳回
	MCNSignStateOnReject MCNSignState = 2
	// MCNSignStateOnSign 已签约
	MCNSignStateOnSign MCNSignState = 10
	// MCNSignStateOnCooling 冷却中
	MCNSignStateOnCooling MCNSignState = 11
	// MCNSignStateOnExpire 已到期
	MCNSignStateOnExpire MCNSignState = 12
	// MCNSignStateOnBlock 已封禁
	MCNSignStateOnBlock MCNSignState = 13
	// MCNSignStateOnClear 已清退
	MCNSignStateOnClear MCNSignState = 14
	// MCNSignStateOnPreOpen 待开启
	MCNSignStateOnPreOpen MCNSignState = 15
	// MCNSignStateOnDelete 已移除
	MCNSignStateOnDelete MCNSignState = 100
)

// NotDealState .
func (mss MCNSignState) NotDealState() bool {
	if mss == MCNSignStateNoApply || mss == MCNSignStateOnReview || mss == MCNSignStateOnReject ||
		mss == MCNSignStateOnBlock || mss == MCNSignStateOnClear || mss == MCNSignStateOnDelete ||
		mss == MCNSignStateOnExpire {
		return true
	}
	return false
}

// MCNSignInfo .
type MCNSignInfo struct {
	SignID             int64        `json:"sign_id"`
	McnMid             int64        `json:"mcn_mid"`
	McnName            string       `json:"mcn_name"`
	CompanyName        string       `json:"company_name"`
	CompanyLicenseID   string       `json:"company_license_id"`
	CompanyLicenseLink string       `json:"company_license_link"`
	ContractLink       string       `json:"contract_link"`
	ContactName        string       `json:"contact_name"`
	ContactTitle       string       `json:"contact_title"`
	ContactPhone       string       `json:"contact_phone"`
	ContactIdcard      string       `json:"contact_idcard"`
	BeginDate          xtime.Time   `json:"begin_date"`
	EndDate            xtime.Time   `json:"end_date"`
	PayExpireState     int8         `json:"pay_expire_state"`
	State              MCNSignState `json:"state"`
	RejectTime         xtime.Time   `json:"reject_time"`
	RejectReason       string       `json:"reject_reason"`
	Ctime              xtime.Time   `json:"ctime"`
	Mtime              xtime.Time   `json:"mtime"`
}

// SignPayInfo  .
type SignPayInfo struct {
	SignPayID int64      `json:"sign_pay_id"`
	McnMid    int64      `json:"mcn_mid"`
	McnName   string     `json:"mcn_name"`
	SignID    int64      `json:"sign_id"`
	State     int8       `json:"state"`
	DueDate   xtime.Time `json:"due_date"`
	PayValue  int64      `json:"pay_value"` // thousand bit
}

// GetDueDate used for template
func (s *SignPayInfo) GetDueDate() string {
	return s.DueDate.Time().Format(TimeFormatDay)
}

// GetPayValue for template
func (s *SignPayInfo) GetPayValue() float64 {
	return float64(s.PayValue) / 1000.0
}

// MCNUPState .
type MCNUPState int8

// const .
const (
	// MCNUPStateNoAuthorize 未授权
	MCNUPStateNoAuthorize MCNUPState = 0
	// MCNUPStateOnRefuse 已拒绝
	MCNUPStateOnRefuse MCNUPState = 1
	// MCNUPStateOnReview 待审核
	MCNUPStateOnReview MCNUPState = 2
	// MCNSignStateOnReject 已驳回
	MCNUPStateOnReject MCNUPState = 3
	// MCNUPStateOnSign 已签约
	MCNUPStateOnSign MCNUPState = 10
	// MCNUPStateOnFreeze 已冻结
	MCNUPStateOnFreeze MCNUPState = 11
	// MCNUPStateOnExpire 已到期
	MCNUPStateOnExpire MCNUPState = 12
	// MCNUPStateOnBlock 已封禁
	MCNUPStateOnBlock MCNUPState = 13
	// MCNUPStateOnClear 已解约
	MCNUPStateOnClear MCNUPState = 14
	// MCNUPStateOnPreOpen 待开启
	MCNUPStateOnPreOpen MCNUPState = 15
	// MCNUPStateOnDelete 已删除
	MCNUPStateOnDelete MCNUPState = 100
)

// MCNUPInfo .
type MCNUPInfo struct {
	SignUpID        int64      `json:"sign_up_id"`
	SignID          int64      `json:"sign_id"`
	McnMid          int64      `json:"mcn_mid"`
	UpMid           int64      `json:"up_mid"`
	BeginDate       xtime.Time `json:"begin_date"`
	EndDate         xtime.Time `json:"end_date"`
	ContractLink    string     `json:"contract_link"`
	UpAuthLink      string     `json:"up_auth_link"`
	RejectTime      xtime.Time `json:"reject_time"`
	RejectReason    string     `json:"reject_reason"`
	State           MCNUPState `json:"state"`
	StateChangeTime xtime.Time `json:"state_change_time"`
	Ctime           xtime.Time `json:"ctime"`
	Mtime           xtime.Time `json:"mtime"`
	UpName          string     `json:"up_name"`
	FansCount       int64      `json:"fans_count"`
	ActiveTid       int64      `json:"active_tid"`
}

// NotDealState .
func (mus MCNUPState) NotDealState() bool {
	if mus == MCNUPStateNoAuthorize || mus == MCNUPStateOnRefuse || mus == MCNUPStateOnReview ||
		mus == MCNUPStateOnReject || mus == MCNUPStateOnFreeze || mus == MCNUPStateOnExpire ||
		mus == MCNUPStateOnBlock || mus == MCNUPStateOnClear || mus == MCNUPStateOnDelete {
		return true
	}
	return false
}
