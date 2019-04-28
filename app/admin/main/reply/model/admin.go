package model

const (
	// AdminIsNotReport 不是举报记录
	AdminIsNotReport int32 = 0
	// AdminIsReport 是举报记录
	AdminIsReport = 1
	// AdminIsNotNew 不是该评论最新的
	AdminIsNotNew int32 = 0
	// AdminIsNew 是该评论最新的
	AdminIsNew = 1
	// AdminOperDelete 管理员删除评论
	AdminOperDelete int32 = 0
	// AdminOperDeleteByReport 管理员通过举报删除 NOTE Deprecated
	AdminOperDeleteByReport = 1
	// AdminOperIgnoreReport 管理员忽略举报 NOTE Deprecated
	AdminOperIgnoreReport = 2
	// AdminOperRecover 管理员恢复评论
	AdminOperRecover = 3
	// AdminOperEdit 管理员编辑内容
	AdminOperEdit = 4
	// AdminOperPass 管理通过待审
	AdminOperPass = 5
	// AdminOperSubState 修改主题的状态
	AdminOperSubState = 6
	// AdminOperSubTop 置顶评论
	AdminOperSubTop = 7
	// AdminOperSubMid 修改主题的mid
	AdminOperSubMid = 8
	// AdminOperRptIgnore1 举报一审忽略
	AdminOperRptIgnore1 = 9
	// AdminOperRptIgnore2 举报二审忽略
	AdminOperRptIgnore2 = 10
	// AdminOperRptDel1 举报一审删除
	AdminOperRptDel1 = 11
	// AdminOperRptDel2 举报二审删除
	AdminOperRptDel2 = 12
	// AdminOperRptRecover1 举报一审恢复
	AdminOperRptRecover1 = 13
	// AdminOperRptRecover2 举报二审恢复
	AdminOperRptRecover2 = 14
	// AdminOperActionSet 对点赞点踩设置
	AdminOperActionSet = 15
	// AdminOperDeleteUp up主删除评论
	AdminOperDeleteUp = 16
	// AdminOperDeleteUser 用户删除评论
	AdminOperDeleteUser = 17
	// AdminOperDeleteAssist 协管删除评论
	AdminOperDeleteAssist = 18
	// AdminOperSubMonitor 设置监控状态
	AdminOperSubMonitor = 19
	// AdminOperRptTransfer1 转一审
	AdminOperRptTransfer1 = 20
	// AdminOperRptTransfer2 转二审
	AdminOperRptTransfer2 = 21
	// AdminOperRptTransferArbitration 移交仲裁
	AdminOperRptTransferArbitration = 22
	// AdminOperRptStateSet 设置举报状态
	AdminOperRptStateSet = 23
	// AdminOperSubMonitorOpen 先发后审打开
	AdminOperSubMonitorOpen = 24
	// AdminOperSubMonitorClose 先发后审关闭
	AdminOperSubMonitorClose = 25
	// AdminOperSubAuditOpen 先审后发打开
	AdminOperSubAuditOpen = 26
	// AdminOperSubAuditClose 先审后发关闭
	AdminOperSubAuditClose = 27
	// AdminOperMarkSpam 标记为垃圾
	AdminOperMarkSpam = 28
)

// AdminLog log.
type AdminLog struct {
	ID       int64  `json:"id"`
	Type     int8   `json:"type"`
	Oid      int64  `json:"oid"`
	ReplyID  int64  `json:"reply_id"`
	AdminID  int64  `json:"admin_id"`
	Result   string `json:"result"`
	Remark   string `json:"remark"`
	IsNew    int8   `json:"is_new"`
	IsReport int8   `json:"is_report"`
	State    int8   `json:"state"`
	CTime    string `json:"ctime"`
	MTime    string `json:"mtime"`
}

// SearchAdminLog log.
type SearchAdminLog struct {
	Type      int8   `json:"type"`
	Oid       int64  `json:"oid"`
	ReplyID   int64  `json:"rpid"`
	AdminID   int64  `json:"adminid"`
	AdminName string `json:"admin_name"`
	Result    string `json:"opresult"`
	Remark    string `json:"opremark"`
	State     int8   `json:"state"`
	CTime     string `json:"opctime"`
}
