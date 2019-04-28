package model

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
	// MCNSignStateOnBlock 封禁
	MCNSignStateOnBlock MCNSignState = 13
	// MCNSignStateOnClear 清退
	MCNSignStateOnClear MCNSignState = 14
	// MCNSignStateOnPreOpen 待开启
	MCNSignStateOnPreOpen MCNSignState = 15
	// MCNSignStateOnDelete 移除
	MCNSignStateOnDelete MCNSignState = 100
)

// MCNUPState .
type MCNUPState int8

// const .
const (
	// MCNUPStateNoAuthorize 未授权
	MCNUPStateNoAuthorize MCNUPState = 0
	// MCNUPStateOnRefuse 已拒绝
	MCNUPStateOnRefuse MCNUPState = 1
	// MCNUPStateOnReview 审核中
	MCNUPStateOnReview MCNUPState = 2
	// MCNUPStateOnReject 已驳回
	MCNUPStateOnReject MCNUPState = 3
	// MCNUPStateOnSign 已签约
	MCNUPStateOnSign MCNUPState = 10
	// MCNUPStateOnCooling 已冻结
	MCNUPStateOnCooling MCNUPState = 11
	// MCNUPStateOnExpire 已到期
	MCNUPStateOnExpire MCNUPState = 12
	// MCNUPStateOnBlock 封禁
	MCNUPStateOnBlock MCNUPState = 13
	// MCNUPStateOnClear 已解约
	MCNUPStateOnClear MCNUPState = 14
	// MCNUPStateOnPreOpen 待开启
	MCNUPStateOnPreOpen MCNUPState = 15
	// MCNUPStateOnDelete 删除
	MCNUPStateOnDelete MCNUPState = 100
)
