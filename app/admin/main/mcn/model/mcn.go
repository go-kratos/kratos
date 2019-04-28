package model

// MCNSignState .
type MCNSignState int8

// const .
const (
	// MCNSignStateUnKnown 未知状态
	MCNSignStateUnKnown MCNSignState = -1
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

func (state MCNSignState) String() string {
	switch state {
	case MCNSignStateNoApply:
		return "未申请"
	case MCNSignStateOnReview:
		return "待审核"
	case MCNSignStateOnReject:
		return "已驳回"
	case MCNSignStateOnSign:
		return "已签约"
	case MCNSignStateOnCooling:
		return "冷却中"
	case MCNSignStateOnExpire:
		return "已到期"
	case MCNSignStateOnBlock:
		return "已封禁"
	case MCNSignStateOnClear:
		return "已清退"
	case MCNSignStateOnPreOpen:
		return "待开启"
	case MCNSignStateOnDelete:
		return "已移除"
	default:
		return ""
	}
}

// MCNSignAction .
type MCNSignAction int8

const (
	// MCNSignActionEntry 录入
	MCNSignActionEntry MCNSignAction = 0
	// MCNSignActionApply 申请
	MCNSignActionApply MCNSignAction = 1
	// MCNSignActionReject 驳回
	MCNSignActionReject MCNSignAction = 2
	// MCNSignActionPass 通过
	MCNSignActionPass MCNSignAction = 10
	// MCNSignActionBlock 封禁
	MCNSignActionBlock MCNSignAction = 13
	// MCNSignActionClear 清退
	MCNSignActionClear MCNSignAction = 14
	// MCNSignActionRenew 续约
	MCNSignActionRenew MCNSignAction = 16
	// MCNSignActionRestore 恢复
	MCNSignActionRestore MCNSignAction = 17
	// MCNSignActionPermit 签约用户权限变更（只用于记日志）
	MCNSignActionPermit MCNSignAction = 99
	// MCNSignActionDelete 移除
	MCNSignActionDelete MCNSignAction = 100
)

// GetState .
func (action MCNSignAction) GetState(oldState MCNSignState) MCNSignState {
	switch action {
	// MCNSignActionEntry 录入
	case MCNSignActionEntry:
		return MCNSignStateNoApply
	// MCNSignActionApply 申请
	case MCNSignActionApply:
		return MCNSignStateOnReview
	// MCNSignActionReject 驳回
	case MCNSignActionReject:
		return MCNSignStateOnReject
	// MCNSignActionPass 通过
	case MCNSignActionPass:
		return MCNSignStateOnSign
	// MCNSignActionBlock 封禁
	case MCNSignActionBlock:
		return MCNSignStateOnBlock
	// MCNSignActionClear 清退
	case MCNSignActionClear:
		return MCNSignStateOnClear
	// MCNSignActionRenew 续约
	case MCNSignActionRenew:
		return MCNSignStateOnSign
	// MCNSignActionRestore 恢复
	case MCNSignActionRestore:
		switch oldState {
		case MCNSignStateOnBlock:
			return MCNSignStateOnSign
		case MCNSignStateOnClear:
			return MCNSignStateNoApply
		}
	// MCNSignActionDelete 移除
	case MCNSignActionDelete:
		return MCNSignStateOnDelete
	}
	return MCNSignState(MCNSignStateUnKnown)
}

// NotRejectState .
func (state MCNSignState) NotRejectState() bool {
	return state != MCNSignStateOnReject
}

// NotRightAction .
func (action MCNSignAction) NotRightAction() bool {
	return action == MCNSignActionReject || action == MCNSignActionPass || action == MCNSignActionDelete
}

// IsOnReviewState .
func (state MCNSignState) IsOnReviewState(action MCNSignAction) bool {
	return state != MCNSignStateOnReview && action != MCNSignActionDelete
}

// IsRenewalState .
func (state MCNSignState) IsRenewalState() bool {
	return state != MCNSignStateOnSign
}

// GetmsgType .
func (action MCNSignAction) GetmsgType(oldState MCNSignState) MSGType {
	switch action {
	// MCNSignActionEntry 录入
	case MCNSignActionEntry:
		return MSGType(0)
	// MCNSignActionApply 申请
	case MCNSignActionApply:
		return MSGType(0)
	// MCNSignActionReject 驳回
	case MCNSignActionReject:
		return McnSignNoApplyPass
	// MCNSignActionPass 通过
	case MCNSignActionPass:
		return McnSignApplyPass
		// MCNSignActionBlock 封禁
	case MCNSignActionBlock:
		return McnBackstageBlock
		// MCNSignActionClear 清退
	case MCNSignActionClear:
		return McnBackstageClose
		// MCNSignActionRenew 续约
	case MCNSignActionRenew:
		return McnRenewcontract
		// MCNSignActionRestore 恢复
	case MCNSignActionRestore:
		switch oldState {
		case MCNSignStateOnBlock:
			return McnAccountRestore
		case MCNSignStateOnClear:
			return MSGType(0)
		}
		// MCNSignActionDelete 移除
	case MCNSignActionDelete:
		return MSGType(0)
	}
	return MSGType(0)
}

func (action MCNSignAction) String() string {
	switch action {
	case MCNSignActionEntry:
		return "mcn录入"
	case MCNSignActionApply:
		return "前台mcn申请"
	case MCNSignActionReject:
		return "驳回申请"
	case MCNSignActionPass:
		return "申请通过"
	case MCNSignActionBlock:
		return "封禁mcn"
	case MCNSignActionClear:
		return "清退mcn"
	case MCNSignActionRenew:
		return "续约mcn"
	case MCNSignActionRestore:
		return "恢复mcn"
	case MCNSignActionPermit:
		return "mcn权限变更"
	case MCNSignActionDelete:
		return "移除mcn"
	default:
		return ""
	}
}

// MCNUPState .
type MCNUPState int8

// const .
const (
	// MCNUPStateUnKnown 未知状态
	MCNUPStateUnKnown MCNUPState = -1
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

func (state MCNUPState) String() string {
	switch state {
	case MCNUPStateNoAuthorize:
		return "未授权"
	case MCNUPStateOnRefuse:
		return "已拒绝"
	case MCNUPStateOnReview:
		return "待审核"
	case MCNUPStateOnReject:
		return "已驳回"
	case MCNUPStateOnSign:
		return "已签约"
	case MCNUPStateOnFreeze:
		return "已冻结"
	case MCNUPStateOnExpire:
		return "已到期"
	case MCNUPStateOnBlock:
		return "已封禁"
	case MCNUPStateOnClear:
		return "已解约"
	case MCNUPStateOnDelete:
		return "已删除"
	default:
		return ""
	}
}

// MCNUPAction .
type MCNUPAction int8

const (
	// MCNUPActionBind 发起绑定
	MCNUPActionBind MCNUPAction = 0
	// MCNUPActionReject 运营驳回
	MCNUPActionReject MCNUPAction = 3
	// MCNUPActionAgree up主同意
	MCNUPActionAgree MCNUPAction = 4
	// MCNUPActionRefuse up主拒绝
	MCNUPActionRefuse MCNUPAction = 5
	// MCNUPActionPass 通过
	MCNUPActionPass MCNUPAction = 10
	// MCNUPActionFreeze 冻结
	MCNUPActionFreeze MCNUPAction = 11
	// MCNUPActionRelease 解约
	MCNUPActionRelease MCNUPAction = 14
	// MCNUPActionRestore 恢复
	MCNUPActionRestore MCNUPAction = 16
)

// GetState .
func (action MCNUPAction) GetState() MCNUPState {
	switch action {
	// MCNUPActionBind 发起绑定
	case MCNUPActionBind:
		return MCNUPStateNoAuthorize
		// MCNUPActionReject 运营驳回
	case MCNUPActionReject:
		return MCNUPStateOnReject
		// MCNUPActionAgree up主同意
	case MCNUPActionAgree:
		return MCNUPStateOnReview
		// MCNUPActionRefuse up主拒绝
	case MCNUPActionRefuse:
		return MCNUPStateOnRefuse
		// MCNUPActionPass 通过
	case MCNUPActionPass:
		return MCNUPStateOnSign
		// MCNUPActionFreeze 冻结
	case MCNUPActionFreeze:
		return MCNUPStateOnFreeze
		// MCNUPActionRelease 解约
	case MCNUPActionRelease:
		return MCNUPStateOnClear
		// MCNUPActionRestore 恢复
	case MCNUPActionRestore:
		return MCNUPStateOnSign
	}
	return MCNUPState(MCNUPStateUnKnown)
}

// GetmsgType .
func (action MCNUPAction) GetmsgType(isMcn bool) MSGType {
	switch {
	// MCNUPActionBind 发起绑定
	case action == MCNUPActionBind:
		return McnUpBindAuthApply
		// MCNUPActionRefuse up主拒绝
	case action == MCNUPActionRefuse:
		return McnUpBindAuthApplyRefuse
		// MCNUPActionAgree up主同意
	case action == MCNUPActionAgree:
		return McnUpBindAuthReview
		// MCNUPActionReject 运营驳回
	case action == MCNUPActionReject && isMcn:
		return McnUpBindAuthApplyNoPass
	case action == MCNUPActionReject && !isMcn:
		return UpMcnBindAuthApplyNoPass
		// MCNUPActionPass 通过
	case action == MCNUPActionPass && isMcn:
		return McnUpBindAuthApplyPass
	case action == MCNUPActionPass && !isMcn:
		return UpMcnBindAuthApplyPass
	// MCNUPActionFreeze 冻结
	case action == MCNUPActionFreeze && isMcn:
		return McnUpRelationFreeze
	case action == MCNUPActionFreeze && !isMcn:
		return UpMcnRelationFreeze
		// MCNUPActionRelease 解约
	case action == MCNUPActionRelease && isMcn:
		return McnUpRelationRelease
	case action == MCNUPActionRelease && !isMcn:
		return UpMcnRelationRelease
		// MCNUPActionRestore 恢复
	case action == MCNUPActionRestore && isMcn:
		return MSGType(0)
	case action == MCNUPActionRestore && !isMcn:
		return MSGType(0)
	}

	return MSGType(0)
}

func (action MCNUPAction) String() string {
	switch action {
	case MCNUPActionBind:
		return "mcn发起绑定"
	case MCNUPActionReject:
		return "运营驳回"
	case MCNUPActionAgree:
		return "up主同意"
	case MCNUPActionRefuse:
		return "up主拒绝"
	case MCNUPActionPass:
		return "审核通过"
	case MCNUPActionFreeze:
		return "up主申请冻结"
	case MCNUPActionRelease:
		return "up主和mcn相互解约"
	case MCNUPActionRestore:
		return "恢复up主和mcn的合同"
	default:
		return ""
	}
}

// NotRightAction .
func (action MCNUPAction) NotRightAction() bool {
	return action == MCNUPActionReject || action == MCNUPActionPass
}

// NoRejectState .
func (action MCNUPAction) NoRejectState() bool {
	return action != MCNUPActionReject
}

// NotRightState .
func (state MCNUPState) NotRightState() bool {
	return state == MCNUPStateOnReject || state == MCNUPStateOnSign
}

// IsOnReviewState .
func (state MCNUPState) IsOnReviewState() bool {
	return state != MCNUPStateOnReview
}

// MCNSignCycleAction .
type MCNSignCycleAction int8

// const .
const (
	// MCNSignCycleActionUp 变更
	MCNSignCycleActionUp MCNSignCycleAction = iota
	// MCNSignCycleActionAdd 新增
	MCNSignCycleActionAdd
	// MCNSignCycleActionDel 删除
	MCNSignCycleActionDel
)

func (act MCNSignCycleAction) String() string {
	switch act {
	case MCNSignCycleActionUp:
		return "变更"
	case MCNSignCycleActionAdd:
		return "新增"
	case MCNSignCycleActionDel:
		return "删除"
	default:
		return ""
	}
}

// MCNPayState .
type MCNPayState int8

// const .
const (
	// MCNPayNo 未支付
	MCNPayNo MCNPayState = 0
	// MCNPayed 已支付
	MCNPayed MCNPayState = 1
	// MCNPayDel 已删除
	MCNPayDel MCNPayState = 100
)

func (mps MCNPayState) String() string {
	switch mps {
	case MCNPayNo:
		return "未支付"
	case MCNPayed:
		return "已支付"
	case MCNPayDel:
		return "已删除"
	default:
		return ""
	}
}

// MCNUPRecommendState .
type MCNUPRecommendState int8

// const .
const (
	// MCNUPRecommendStateUnKnown 未知状态
	MCNUPRecommendStateUnKnown MCNUPRecommendState = 0
	// MCNUPRecommendStateOff 未推荐
	MCNUPRecommendStateOff MCNUPRecommendState = 1
	// MCNUPRecommendStateOn  推荐中
	MCNUPRecommendStateOn MCNUPRecommendState = 2
	// MCNUPRecommendStateBan 禁止推荐
	MCNUPRecommendStateBan MCNUPRecommendState = 3
	// MCNUPRecommendStateDel 移除中
	MCNUPRecommendStateDel MCNUPRecommendState = 100
)

func (state MCNUPRecommendState) String() string {
	switch state {
	case MCNUPRecommendStateOff:
		return "未推荐"
	case MCNUPRecommendStateOn:
		return "推荐中"
	case MCNUPRecommendStateBan:
		return "禁止推荐"
	case MCNUPRecommendStateDel:
		return "移除中"
	default:
		return "未知状态"
	}
}

// MCNUPRecommendSource .
type MCNUPRecommendSource int8

// const .
const (
	// MCNUPRecommendSourceUnKnown 未知来源
	MCNUPRecommendSourceUnKnown MCNUPRecommendSource = iota
	// MCNUPRecommendSourceAuto 自动添加(大数据)
	MCNUPRecommendSourceAuto
	// MCNUPRecommendStateManual  手动添加
	MCNUPRecommendStateManual
)

func (source MCNUPRecommendSource) String() string {
	switch source {
	case MCNUPRecommendSourceAuto:
		return "自动添加(大数据)"
	case MCNUPRecommendStateManual:
		return "手动添加"
	default:
		return "未知来源"
	}
}

// MCNUPRecommendAction .
type MCNUPRecommendAction int8

// const .
const (
	// MCNUPRecommendActionOn 推荐
	MCNUPRecommendActionOn MCNUPRecommendAction = iota + 1
	// MCNUPRecommendActionBan 禁止推荐
	MCNUPRecommendActionBan
	// MCNUPRecommendActionRestore 恢复
	MCNUPRecommendActionRestore
	// MCNUPRecommendActionAdd 手动添加
	MCNUPRecommendActionAdd
	// MCNUPRecommendActionDel  移除
	MCNUPRecommendActionDel
)

// GetState .
func (action MCNUPRecommendAction) GetState() MCNUPRecommendState {
	switch action {
	// MCNUPRecommendActionOn 推荐
	case MCNUPRecommendActionOn:
		return MCNUPRecommendStateOn
		// MCNUPRecommendActionBan 禁止推荐
	case MCNUPRecommendActionBan:
		return MCNUPRecommendStateBan
		// MCNUPRecommendActionRestore 恢复
	case MCNUPRecommendActionRestore:
		return MCNUPRecommendStateOff
		// MCNUPRecommendActionAdd 手动添加
	case MCNUPRecommendActionAdd:
		return MCNUPRecommendStateOff
		// MCNUPRecommendActionDel 移除
	case MCNUPRecommendActionDel:
		return MCNUPRecommendStateDel
	}
	return MCNUPRecommendStateUnKnown
}

func (action MCNUPRecommendAction) String() string {
	switch action {
	case MCNUPRecommendActionOn:
		return "推荐"
	case MCNUPRecommendActionBan:
		return "禁止推荐"
	case MCNUPRecommendActionRestore:
		return "恢复"
	case MCNUPRecommendActionAdd:
		return "手动添加"
	case MCNUPRecommendActionDel:
		return "移除"
	default:
		return ""
	}
}

// MCNUPPermissionState .
type MCNUPPermissionState int8

// const .
const (
	// MCNUPPermissionStateUnKnown 未知状态
	MCNUPPermissionStateUnKnown MCNUPPermissionState = -1
	// MCNUPPermissionStateNoAuthorize 待Up主同意
	MCNUPPermissionStateNoAuthorize MCNUPPermissionState = 0
	// MCNUPStateOnRefuse Up主拒绝
	MCNUPPermissionStateOnRefuse MCNUPPermissionState = 1
	// MCNUPPermissionStateReview 待审中
	MCNUPPermissionStateReview MCNUPPermissionState = 2
	// MCNUPPermissionStatePass  已通过
	MCNUPPermissionStatePass MCNUPPermissionState = 3
	// MCNUPPermissionStateFail 已驳回
	MCNUPPermissionStateFail MCNUPPermissionState = 4
	// MCNUPPermissionStateDel 已删除
	MCNUPPermissionStateDel MCNUPPermissionState = 100
)

func (state MCNUPPermissionState) String() string {
	switch state {
	case MCNUPPermissionStateNoAuthorize:
		return "待Up主同意"
	case MCNUPPermissionStateOnRefuse:
		return "Up主拒绝"
	case MCNUPPermissionStateReview:
		return "待审中"
	case MCNUPPermissionStatePass:
		return "已通过"
	case MCNUPPermissionStateFail:
		return "已驳回"
	case MCNUPPermissionStateDel:
		return "已删除"
	default:
		return "未知状态"
	}
}

// MCNUPPermissionAction .
type MCNUPPermissionAction int8

// const .
const (
	// MCNUPPermissionActionOn 通过
	MCNUPPermissionActionOn MCNUPPermissionAction = iota + 1
	// MCNUPPermissionActionFail 驳回
	MCNUPPermissionActionFail
	// MCNUPPermissionActionDel  移除
	MCNUPPermissionActionDel
)

// NotRightAction .
func (action MCNUPPermissionAction) NotRightAction() bool {
	return action == MCNUPPermissionActionOn || action == MCNUPPermissionActionFail
}

// GetState .
func (action MCNUPPermissionAction) GetState() MCNUPPermissionState {
	switch action {
	// MCNUPPermissionActionOn 通过
	case MCNUPPermissionActionOn:
		return MCNUPPermissionStatePass
		// MCNUPPermissionActionFail 驳回
	case MCNUPPermissionActionFail:
		return MCNUPPermissionStateFail
		// MCNUPPermissionActionDel 移除
	case MCNUPPermissionActionDel:
		return MCNUPPermissionStateDel
	}
	return MCNUPPermissionStateUnKnown
}

func (action MCNUPPermissionAction) String() string {
	switch action {
	case MCNUPPermissionActionOn:
		return "通过"
	case MCNUPPermissionActionFail:
		return "驳回"
	case MCNUPPermissionActionDel:
		return "移除"
	default:
		return ""
	}
}

// AttrBasePermit .
type AttrBasePermit uint

// const .
const (
	AttrBasePermitBit   AttrBasePermit = 0 // 基础权限
	AttrDataPermitBit   AttrBasePermit = 1 // 数据权限
	AttrRecPermitBit    AttrBasePermit = 2 // 推荐权限
	AttrDepartPermitBit AttrBasePermit = 3 // 起飞权限
)

// PermitMap .
var PermitMap = map[AttrBasePermit]struct{}{
	AttrBasePermitBit:   {},
	AttrDataPermitBit:   {},
	AttrRecPermitBit:    {},
	AttrDepartPermitBit: {},
}

func (p AttrBasePermit) String() string {
	switch p {
	case AttrBasePermitBit:
		return "基础权限"
	case AttrDataPermitBit:
		return "数据权限"
	case AttrRecPermitBit:
		return "推荐权限"
	case AttrDepartPermitBit:
		return "起飞权限"
	default:
		return ""
	}
}
