package model

import (
	"fmt"
)

// MSGType .
type MSGType uint8

// const .
const (
	// McnSignApplyPass MCN申请MCN管理入口申请成功
	McnSignApplyPass = iota + 1
	// McnSignNoApplyPass MCN申请MCN管理入口申请未通过
	McnSignNoApplyPass
	// McnUpBindAuthApply MCN申请和up主绑定申请授权
	McnUpBindAuthApply
	// McnUpBindAuthReview MCN申请和up主绑定up主同意等待运营审核中
	McnUpBindAuthReview
	// McnUpBindAuthApplyPass MCN申请和up主绑定up主同意且运营通过
	McnUpBindAuthApplyPass
	// UpMcnBindAuthApplyPass up主申请和MCN绑定up主同意且运营通过
	UpMcnBindAuthApplyPass
	// McnUpBindAuthApplyNoPass  MCN申请和up主绑定up主同意但运营未通过
	McnUpBindAuthApplyNoPass
	// UpMcnBindAuthApplyNoPass up主申请和MCN绑定up主同意但运营未通过
	UpMcnBindAuthApplyNoPass
	// McnUpBindAuthApplyRefuse MCN申请和up主绑定被up主拒绝
	McnUpBindAuthApplyRefuse
	// UpMcnRelationFreeze  MCN和up主纠纷处理 - Up主和MCN关系冻结
	UpMcnRelationFreeze
	// McnUpRelationFreeze  MCN和up主纠纷处理 - MCN和Up主关系冻结
	McnUpRelationFreeze
	// UpMcnRelationRelease MCN和up主纠纷处理 - Up主和MCN提前解约
	UpMcnRelationRelease
	// McnUpRelationRelease MCN和up主纠纷处理 - MCN和Up主提前解约
	McnUpRelationRelease
	// McnBackstageBlock MCN违规账号封禁
	McnBackstageBlock
	// McnBackstageClose MCN违规账号清退
	McnBackstageClose
	// McnRenewcontract 续约合同
	McnRenewcontract
	// McnAccountRestore MCN账号恢复
	McnAccountRestore
	// McnPermissionOpen MCN新开权限
	McnPermissionOpen
	// McnPermissionClosed MCN权限关闭
	McnPermissionClosed
	// McnUpNotAgreeChangePermit UP主不同意授权变更
	McnUpNotAgreeChangePermit
	// McnOperNotAgreeChangePermit 运营不同意授权变更
	McnOperNotAgreeChangePermit
	// McnOperAgreeChangePermit 运营同意授权变更
	McnOperAgreeChangePermit
	// McnApplyUpChangePermit MCN申请和up主的权限修改
	McnApplyUpChangePermit
)

// MSG .
type MSG struct {
	MSGType MSGType
	Code    string
	Title   string
	Content string
}

// ArgMsg .
type ArgMsg struct {
	MSGType     MSGType
	MIDs        []int64
	McnName     string
	UpName      string
	McnMid      int64
	UpMid       int64
	CompanyName string
	Reason      string
	SignUpID    int64
	Permission  string
}

// MsgInfo .
func (arg *ArgMsg) MsgInfo(msg *MSG) (mids []int64, title, content, code string) {
	switch arg.MSGType {
	case McnSignApplyPass:
		return arg.MIDs, msg.Title, msg.Content, msg.Code
	case McnSignNoApplyPass:
		return arg.MIDs, msg.Title, fmt.Sprintf(msg.Content, arg.Reason), msg.Code
	case McnUpBindAuthApply:
		return arg.MIDs, msg.Title, fmt.Sprintf(msg.Content, arg.McnName, arg.McnMid, arg.CompanyName, arg.SignUpID), msg.Code
	case McnUpBindAuthReview:
		return arg.MIDs, msg.Title, fmt.Sprintf(msg.Content, arg.CompanyName, arg.McnName, arg.McnMid), msg.Code
	case McnUpBindAuthApplyPass:
		return arg.MIDs, msg.Title, fmt.Sprintf(msg.Content, arg.UpName, arg.UpMid), msg.Code
	case UpMcnBindAuthApplyPass:
		return arg.MIDs, msg.Title, fmt.Sprintf(msg.Content, arg.CompanyName, arg.McnName, arg.McnMid), msg.Code
	case McnUpBindAuthApplyNoPass:
		return arg.MIDs, msg.Title, fmt.Sprintf(msg.Content, arg.UpName, arg.UpMid, arg.Reason), msg.Code
	case UpMcnBindAuthApplyNoPass:
		return arg.MIDs, msg.Title, fmt.Sprintf(msg.Content, arg.CompanyName, arg.McnName, arg.McnMid, arg.Reason), msg.Code
	case McnUpBindAuthApplyRefuse:
		return arg.MIDs, msg.Title, fmt.Sprintf(msg.Content, arg.UpName, arg.UpMid), msg.Code
	case UpMcnRelationFreeze:
		return arg.MIDs, msg.Title, fmt.Sprintf(msg.Content, arg.CompanyName, arg.McnName, arg.McnMid), msg.Code
	case McnUpRelationFreeze:
		return arg.MIDs, msg.Title, fmt.Sprintf(msg.Content, arg.UpName, arg.UpMid), msg.Code
	case UpMcnRelationRelease:
		return arg.MIDs, msg.Title, fmt.Sprintf(msg.Content, arg.CompanyName, arg.McnName, arg.McnMid), msg.Code
	case McnUpRelationRelease:
		return arg.MIDs, msg.Title, fmt.Sprintf(msg.Content, arg.UpName, arg.UpMid), msg.Code
	case McnBackstageBlock:
		return arg.MIDs, msg.Title, msg.Content, msg.Code
	case McnBackstageClose:
		return arg.MIDs, msg.Title, msg.Content, msg.Code
	case McnRenewcontract:
		return arg.MIDs, msg.Title, msg.Content, msg.Code
	case McnAccountRestore:
		return arg.MIDs, msg.Title, msg.Content, msg.Code
	case McnPermissionOpen:
		return arg.MIDs, fmt.Sprintf(msg.Title, arg.Permission), fmt.Sprintf(msg.Content, arg.Permission), msg.Code
	case McnPermissionClosed:
		return arg.MIDs, fmt.Sprintf(msg.Title, arg.Permission), fmt.Sprintf(msg.Content, arg.Permission), msg.Code
	case McnUpNotAgreeChangePermit:
		return arg.MIDs, fmt.Sprintf(msg.Title, arg.UpName), fmt.Sprintf(msg.Content, arg.UpName), msg.Code
	case McnOperNotAgreeChangePermit:
		return arg.MIDs, fmt.Sprintf(msg.Title, arg.UpName), fmt.Sprintf(msg.Content, arg.UpName, arg.Reason), msg.Code
	case McnOperAgreeChangePermit:
		return arg.MIDs, fmt.Sprintf(msg.Title, arg.UpName), fmt.Sprintf(msg.Content, arg.UpName, arg.Permission), msg.Code
	case McnApplyUpChangePermit:
		return arg.MIDs, fmt.Sprintf(msg.Title, arg.McnName), fmt.Sprintf(msg.Content, arg.McnName, arg.Permission, arg.SignUpID), msg.Code
	}
	return nil, "", "", ""
}
