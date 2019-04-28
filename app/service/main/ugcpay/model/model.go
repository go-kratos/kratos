package model

// 各种状态枚举
const (
	StateValid   = "valid"
	StateInvalid = "invalid"

	BizAsset = "asset"
	BizElec  = "elec"

	AssetRelationPaid = "paid"
	AssetRelationNone = "none"
	AssetRelationMiss = "miss"

	OTypeArchive = "archive"

	OrderStateCreated          = "created"
	OrderStatePaying           = "paying"
	OrderStatePaid             = "paid"
	OrderStateFailed           = "failed"
	OrderStateClosed           = "closed"
	OrderStateExpired          = "expired"
	OrderStateSettled          = "settled"
	OrderStateRefunding        = "refunding"
	OrderStateRefunded         = "refunded"
	OrderStateSettledRefunding = "st_refunding"
	OrderStateSettledRefunded  = "st_refunded"
	OrderStateRefundFinished   = "ref_finished"
	OrderStateDupRefunded      = "dup_refunded"
	OrderStateBadDebt          = "bad_debt"

	PayStatePaying              = "PAYING"
	PayStateOverdue             = "OVERDUE"
	PayStateClosed              = "CLOSED"
	PayStateFail                = "FAIL"
	PayStateSuccess             = "SUCCESS"
	PayStateFinished            = "FINISHED"
	PayStatePaySuccessAndCancel = "PAY_SUCCESS_AND_CANCEL"
	PayStatePayCancel           = "PAY_CANCEL"
	PayStateRefund              = "REFUND_SUCCESS"

	RechargeShellSuccess = "success"
	RechargeShellFail    = "fail"
)
