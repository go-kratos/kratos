package blocked

import (
	"fmt"
)

// const msg
const (
	MsgTypeDelOpinion = int8(1)
	MsgTypeGetJuryer  = int8(2)
	MsgTypeAppealSucc = int8(3)
	MsgTypeAppealFail = int8(4)
)

var _msg = map[int8]map[string]string{
	MsgTypeDelOpinion: {
		"title": "观点违规警告",
		"content": `你在%d号案件中发布的观点“%s”，违反#{《“众议观点”使用守则》}{"https://www.bilibili.com/blackroom/notice/44"}，已被管理员删除。
		请遵守相关守则，合理发布观点，多次违反将被取消风纪委员资格。`,
	},
	MsgTypeGetJuryer: {
		"title": "获得风纪委员资格",
		"content": `恭喜您获得%d天风纪委员资格！风纪委员应遵守以下原则：
		"1. 在了解举报案件背景后，公正客观投票。对不了解或难以判断的案件，可以选择弃权。
		"2. 以身作则，不在举报案件相关视频、评论下讨论或发布不相关内容。相关违规举报被落实处罚后，将会失去风纪委员资格。`,
	},
	MsgTypeAppealSucc: {
		"title":   "申诉处理通知",
		"content": `经复核，您对案件%d的申诉成功，相应惩罚将被撤销。很抱歉给您带来了不便。感谢您对社区工作的理解与支持。请继续遵守社区规范，共同维护良好的社区氛围！`,
	},
	MsgTypeAppealFail: {
		"title":   "申诉处理通知",
		"content": `经复核，您对案件%d的申诉未能通过。请遵守社区规范，共同维护良好的社区氛围！`,
	},
}

// SysMsg msg struct
type SysMsg struct {
	Type        int8
	MID         int64
	Day         int
	CID         int64
	CaseContent string
	RemoteIP    string
}

// MsgInfo get msg info
func MsgInfo(msg *SysMsg) (title, content string) {
	switch msg.Type {
	case MsgTypeDelOpinion:
		title = _msg[MsgTypeDelOpinion]["title"]
		content = fmt.Sprintf(_msg[MsgTypeDelOpinion]["content"], msg.CID, msg.CaseContent)
	case MsgTypeGetJuryer:
		title = _msg[MsgTypeGetJuryer]["title"]
		content = fmt.Sprintf(_msg[MsgTypeGetJuryer]["content"], msg.Day)
	case MsgTypeAppealSucc:
		title = _msg[MsgTypeAppealSucc]["title"]
		content = fmt.Sprintf(_msg[MsgTypeAppealSucc]["content"], msg.CID)
	case MsgTypeAppealFail:
		title = _msg[MsgTypeAppealFail]["title"]
		content = fmt.Sprintf(_msg[MsgTypeAppealFail]["content"], msg.CID)
	}
	return
}
