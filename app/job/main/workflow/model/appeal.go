package model

import "go-common/library/time"

// appeal state
const (
	AuditStatePending   = 0
	AuditStateEffective = 1
	AuditStateInvalid   = 2

	TransferStatePendingSystemNotReply = 0
	TransferStatePendingSystemReply    = 1
	TransferStateAdminReplyReaded      = 2
	TransferStateAdminClosed           = 3
	TransferStateUserResolved          = 4
	TransferStateAutoClosedExpire      = 5
	TransferStateAdminReplyNotReaded   = 6
	TransferStateUserClosed            = 7
	TransferStatePassClosed            = 8

	AssignStateNotDispatch = 0
	AssignStatePushed      = 1
	AssignStatePoped       = 2
	AssignStateReAudit     = 3
)

// Appeal orm struct of table workflow_appeal
type Appeal struct {
	ApID          int64     `json:"id" gorm:"column:id"`
	Rid           int8      `json:"rid" gorm:"column:rid"`
	Tid           int32     `json:"tid" gorm:"column:tid"`
	Bid           int8      `json:"bid" gorm:"column:bid"`
	Mid           int64     `json:"mid" gorm:"column:mid"`
	Oid           int64     `json:"oid" gorm:"column:oid"`
	AuditState    int8      `json:"audit_state" gorm:"column:audit_state"`
	TransferState int8      `json:"transfer_state" gorm:"column:transfer_state"`
	AssignState   int8      `json:"assign_state" gorm:"column:assign_state"`
	Weight        int64     `json:"weight" gorm:"column:weight"`
	AuditAdmin    int64     `json:"audit_admin" gorm:"column:audit_admin"`
	TransferAdmin int64     `json:"transfer_admin" gorm:"column:transfer_admin"`
	DTime         time.Time `json:"dtime" gorm:"column:dtime"`
	TTime         time.Time `json:"ttime" gorm:"column:ttime"`
	CTime         time.Time `json:"ctime" gorm:"column:ctime"`
	MTime         time.Time `json:"mtime" gorm:"column:mtime"`
}
