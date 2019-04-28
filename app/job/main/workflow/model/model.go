package model

var (
	// ListBefore .
	ListBefore = 0 // 入队列前
	// ListAfter .
	ListAfter = 1 // 出队列后
	// ListIng .
	ListIng = 2 // 队列中
	// SysAssignType .
	SysAssignType = 1 // 系统指派类型
	// ADealType .
	ADealType = 0 // 审核处理类型
	// FDealType .
	FDealType = 1 // 反馈处理类型
	// PDealType .
	PDealType = 2 // 工作台处理类型
	// FListBeforeStates .
	FListBeforeStates = "2" // 反馈入队列前state状态
	// FListBeforeBusinessStates .
	FListBeforeBusinessStates = "1" // 反馈入队列前business_state状态
	// AListBeforeStates .
	AListBeforeStates = "0" // 审核入队列前state状态
	// FListAfterStates .
	FListAfterStates = "2" // 反馈入队列后state状态
	// FListAfterBusinessStates .
	FListAfterBusinessStates = "1" // 反馈入队列后business_state状态
	// AListAfterStates .
	AListAfterStates = "15" // 审核入队列后state状态
)

// SearchParams .
type SearchParams struct {
	Business            string
	States              string
	BusinessStates      string
	AssigneeAdminIDs    string
	AssigneeAdminIDsNot string
	MtimeTo             string
}
