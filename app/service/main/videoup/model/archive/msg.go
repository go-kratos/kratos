package archive

import (
	"fmt"
)

// const .
const (
	MSG_1  = int(1)
	MSG_2  = int(2)
	MSG_3  = int(3)
	MSG_4  = int(4)
	MSG_5  = int(5)
	MSG_6  = int(6)
	MSG_7  = int(7)
	MSG_8  = int(8)
	MSG_9  = int(9)
	MSG_10 = int(10)
	MSG_11 = int(11)
	MSG_12 = int(12)
	MSG_13 = int(13)
)

// MSG .
type MSG struct {
	To      string
	MSGID   int
	Code    string
	Title   string
	Content string
}

// ArgMsg .
type ArgMsg struct {
	MSGID int
	Apply *ApplyParam
}

// MsgInfo .
func (arg *ArgMsg) MsgInfo(msg *MSG) (mids []int64, title, content, code string) {
	switch arg.MSGID {
	/*
		1up主xxxx邀请您作为xxxx参与稿件《xxxx》（avxxxxx）的多人合作，在网页端创作中心查看吧点击查看
		2参与者xxxx已接受您的邀请作为xxxx参与稿件《xxxx》（avxxxxx）的多人合作，在网页端创作中心查看吧点击查看
		3参与者xxxx已拒绝您的邀请作为xxxx参与稿件《xxxx》（avxxxxx）的多人合作，在网页端创作中心查看吧点击查看

		4up主xxxx申请您参与稿件《xxxx》（avxxxx）的参与类型由xxxx变更为xxxx，在网页端创作中心查看吧点击查看
		5参与者已同意您申请其参与稿件《xxxx》（avxxxx）的参与类型由xxxx变更为xxxx，在网页端创作中心查看吧点击查看
		6参与者已拒绝您申请其参与稿件《xxxx》（avxxxx）的参与类型由xxxx变更为xxxx，在网页端创作中心查看吧点击查看

		7合作者xxxx申请终止其作为稿件《xxxx》（avxxxxx）的xxxx的多人合作，在网页端创作中心查看吧点击查看
		8up主xxxx已终止参与者xxxx作为稿件《xxxx》（avxxxx）的xxxx的多人合作，在网页端创作中心查看吧点击查看

		9up主xxxx申请终止您作为稿件xxxxxx（aid：xxxxx）的xxx的多人合作，在网页端创作中心查看吧点击查看
		10 up主xxxx已终止参与者xxxx作为稿件《xxxx》（avxxxx）的xxxx的多人合作，在网页端创作中心查看吧点击查看
		11 up主xxxx未能终止参与者xxxx作为稿件《xxxx》（avxxxx）的xxxx的多人合作，在网页端创作中心查看吧点击查看

		12 管理员解除了您与参与者xxx在稿件《xxxx》《avxxxx》中合作关系，在网页端创作中心查看吧点击查看
		13 管理员解除了您与up主xxxx在稿件《xxxx》《avxxxx》中合作关系，在网页端创作中心查看吧点击查看
	*/
	case MSG_1:
		//content = 'up主%s邀请您作为%s参与稿件《%s》（av%d）的多人合作，在网页端创作中心查看吧#{点击查看}{"https://member.bilibili.com/v2#/upload-manager/article/cooperation"}”'
		return []int64{arg.Apply.ApplyStaffMID}, msg.Title, fmt.Sprintf(msg.Content, arg.Apply.UpName, arg.Apply.ApplyTitle, arg.Apply.Archive.Title, arg.Apply.Archive.Aid), msg.Code
	case MSG_2:
		//content = '参与者%s已接受您的邀请作为%s参与稿件《%s》（av%d）的多人合作，在网页端创作中心查看吧#{点击查看}{"https://member.bilibili.com/v2#/upload-manager/article"}”'
		return []int64{arg.Apply.ApplyUpMID}, msg.Title, fmt.Sprintf(msg.Content, arg.Apply.StaffName, arg.Apply.ApplyTitle, arg.Apply.Archive.Title, arg.Apply.Archive.Aid), msg.Code
	case MSG_3:
		//content = '参与者%s已拒绝您的邀请作为%s参与稿件《%s》（av%d）的多人合作，在网页端创作中心查看吧#{点击查看}{"https://member.bilibili.com/v2#/upload-manager/article"}”'
		return []int64{arg.Apply.ApplyUpMID}, msg.Title, fmt.Sprintf(msg.Content, arg.Apply.StaffName, arg.Apply.ApplyTitle, arg.Apply.Archive.Title, arg.Apply.Archive.Aid), msg.Code

	case MSG_4:
		//content = 'up主%s申请您参与稿件《%s》（av%d）的参与类型由%s变更为%s，在网页端创作中心查看吧#{点击查看}{"https://member.bilibili.com/v2#/upload-manager/article/cooperation"}”'
		return []int64{arg.Apply.ApplyStaffMID}, msg.Title, fmt.Sprintf(msg.Content, arg.Apply.UpName, arg.Apply.Archive.Title, arg.Apply.Archive.Aid, arg.Apply.StaffTitle, arg.Apply.ApplyTitle), msg.Code
	case MSG_5:
		//content = '参与者%s已同意您申请其参与稿件《%s》（av%d）的参与类型由%s变更为%s，在网页端创作中心查看吧#{点击查看}{"https://member.bilibili.com/v2#/upload-manager/article"}”'
		return []int64{arg.Apply.ApplyUpMID}, msg.Title, fmt.Sprintf(msg.Content, arg.Apply.StaffName, arg.Apply.Archive.Title, arg.Apply.Archive.Aid, arg.Apply.StaffTitle, arg.Apply.ApplyTitle), msg.Code
	case MSG_6:
		//content = '参与者%s已拒绝您申请其参与稿件《%s》（av%d）的参与类型由%s变更为%s，在网页端创作中心查看吧#{点击查看}{"https://member.bilibili.com/v2#/upload-manager/article"}”'
		return []int64{arg.Apply.ApplyUpMID}, msg.Title, fmt.Sprintf(msg.Content, arg.Apply.StaffName, arg.Apply.Archive.Title, arg.Apply.Archive.Aid, arg.Apply.StaffTitle, arg.Apply.OldTitle), msg.Code

	case MSG_7:
		// content = '合作者%s申请终止其作为稿件《%s》（av%d）的%s的多人合作，在网页端创作中心查看吧#{点击查看}{"https://member.bilibili.com/v2#/upload-manager/article"}”'
		return []int64{arg.Apply.ApplyUpMID}, msg.Title, fmt.Sprintf(msg.Content, arg.Apply.StaffName, arg.Apply.Archive.Title, arg.Apply.Archive.Aid, arg.Apply.ApplyTitle), msg.Code
	case MSG_8:
		//content = 'up主%s已终止参与者%s作为稿件《%s》（av%d）的%s的多人合作，在网页端创作中心查看吧#{点击查看}{"https://member.bilibili.com/v2#/upload-manager/article/cooperation"}”'
		return []int64{arg.Apply.ApplyStaffMID}, msg.Title, fmt.Sprintf(msg.Content, arg.Apply.UpName, arg.Apply.StaffName, arg.Apply.Archive.Title, arg.Apply.Archive.Aid, arg.Apply.StaffTitle), msg.Code

	case MSG_9:
		//content = 'up主%s申请终止您作为稿件%s（aid：%d）的%s的多人合作，在网页端创作中心查看吧#{点击查看}{"https://member.bilibili.com/v2#/upload-manager/article/cooperation"}”'
		return []int64{arg.Apply.ApplyStaffMID}, msg.Title, fmt.Sprintf(msg.Content, arg.Apply.UpName, arg.Apply.Archive.Title, arg.Apply.Archive.Aid, arg.Apply.StaffTitle), msg.Code
	case MSG_10:
		//content = 'up主%s已终止参与者%s作为稿件《%s》（av%d）的%s的多人合作，在网页端创作中心查看吧#{点击查看}{"https://member.bilibili.com/v2#/upload-manager/article"}”'
		return []int64{arg.Apply.ApplyUpMID}, msg.Title, fmt.Sprintf(msg.Content, arg.Apply.UpName, arg.Apply.StaffName, arg.Apply.Archive.Title, arg.Apply.Archive.Aid, arg.Apply.StaffTitle), msg.Code
	case MSG_11:
		// content = 'up主%s未能终止参与者%s作为稿件《%s》（av%d）的%s的多人合作，在网页端创作中心查看吧#{点击查看}{"https://member.bilibili.com/v2#/upload-manager/article"}”'
		return []int64{arg.Apply.ApplyUpMID}, msg.Title, fmt.Sprintf(msg.Content, arg.Apply.UpName, arg.Apply.StaffName, arg.Apply.Archive.Title, arg.Apply.Archive.Aid, arg.Apply.StaffTitle), msg.Code
	case MSG_12:
		//content = '管理员解除了您与参与者%s在稿件《%s》《av%d》中合作关系，在网页端创作中心查看吧#{点击查看}{"https://member.bilibili.com/v2#/upload-manager/article"}”'
		return []int64{arg.Apply.ApplyUpMID}, msg.Title, fmt.Sprintf(msg.Content, arg.Apply.StaffsName, arg.Apply.Archive.Title, arg.Apply.Archive.Aid), msg.Code
	case MSG_13:
		//content = '管理员解除了您与UP主%s在稿件《%s》《av%d》中合作关系，在网页端创作中心查看吧#{点击查看}{"https://member.bilibili.com/v2#/upload-manager/article"}”'
		return []int64{arg.Apply.ApplyStaffMID}, msg.Title, fmt.Sprintf(msg.Content, arg.Apply.UpName, arg.Apply.Archive.Title, arg.Apply.Archive.Aid), msg.Code
	}
	return nil, "", "", ""
}
