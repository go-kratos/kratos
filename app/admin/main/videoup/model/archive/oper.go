package archive

import (
	"fmt"

	"go-common/app/admin/main/videoup/model/utils"
)

// .
const (
	// OperTypeMission 活动id被修改
	OperTypeMission = int8(1)
	// OperTypeTag tag被修改
	OperTypeTag = int8(2)
	// OperTypeCopyright 版权类型被修改
	OperTypeCopyright = int8(3)
	// OperTypeTypeID 分区ID被修改
	OperTypeTypeID = int8(4)
	// OperTypeRejectReason 打回理由被修改
	OperTypeRejectReason = int8(5)
	// OperTypeForwardID 转车跳转被修改
	OperTypeForwardID = int8(6)
	// OperTypeFlowID 私单类型被修改
	OperTypeFlowID = int8(7)
	// OperTypeDelay 定时发布被修改
	OperTypeDelay = int8(8)
	// OperTypePtime 发布时间被修改
	OperTypePtime = int8(10)
	// OperTypeAccess 可见属性被修改
	OperTypeAccess = int8(11)
	// OperTypeAduitReason 审核理由被修改
	OperTypeAduitReason = int8(12)
	// OperTypeRecicleTag 打回理由被修改
	OperTypeRecicleTag = int8(13)
	// OperTypeTaskID 任务ID被修改
	OperTypeTaskID = int8(14)
	// OperTypeOpenTag 通过Tag被修改
	OperTypeOpenTag = int8(15)
	// OperTypeDynamic 动态描述被修改
	OperTypeDynamic = int8(16)
	OperNotify      = int8(17)
	//私单
	OperPorderIndustryID = int8(18)
	OperPorderOfficial   = int8(19)
	OperPorderBrandID    = int8(20)
	OperPorderBrandName  = int8(21)
	OperPorderShowType   = int8(22)
	OperPorderAdvertiser = int8(23)
	OperPorderAgent      = int8(24)
	OperPorderShowFront  = int8(25)
	//频道回查属性
	OperFlowAttrNoChannel = int8(26)
	OperFlowAttrNoHot     = int8(27)

	// OperStyleOne 操作展示类型1：[%s]从[%v]设为[%v]
	OperStyleOne = int8(1)
	// OperStyleTwo 操作展示类型2：[%s]%v:%v
	OperStyleTwo = int8(2)
)

var (
	//FlowOperType flow oper id
	FlowOperType = map[int64]int8{
		FlowGroupNoChannel: OperFlowAttrNoChannel,
		FlowGroupNoHot:     OperFlowAttrNoHot,
	}
	_operType = map[int8]string{
		OperTypeMission:       "活动ID",
		OperTypeTag:           "TAG内容",
		OperTypeCopyright:     "投稿类型",
		OperTypeTypeID:        "分区类型",
		OperTypeRejectReason:  "回查理由",
		OperTypeForwardID:     "撞车跳转",
		OperTypeFlowID:        "流量TAG",
		OperTypeDelay:         "定时发布",
		OperTypePtime:         "发布时间",
		OperTypeAccess:        "可见属性",
		OperTypeAduitReason:   "审核理由",
		OperTypeRecicleTag:    "打回Tag",
		OperTypeTaskID:        "任务ID",
		OperTypeOpenTag:       "通过Tag",
		OperTypeDynamic:       "动态描述",
		OperNotify:            "系统通知",
		OperPorderIndustryID:  "推广行业",
		OperPorderOfficial:    "是否官方",
		OperPorderBrandID:     "推广品牌ID",
		OperPorderBrandName:   "推广品牌",
		OperPorderShowType:    "推广形式",
		OperPorderAdvertiser:  "广告主",
		OperPorderAgent:       "代理商",
		OperPorderShowFront:   "是否前端展示",
		OperFlowAttrNoChannel: "频道禁止",
		OperFlowAttrNoHot:     "热门禁止",
	}
)

// ArcOper archive oper.
type ArcOper struct {
	ID        int64
	Aid       int64
	UID       int64
	TypeID    int16
	State     int16
	Content   string
	Round     int8
	Attribute int32
	LastID    int64
	Remark    string
}

// VideoOper video oper.
type VideoOper struct {
	ID        int64            `json:"id"`
	AID       int64            `json:"aid"`
	UID       int64            `json:"uid"`
	VID       int64            `json:"vid"`
	Status    int16            `json:"status"`
	Content   string           `json:"content"`
	Attribute int32            `json:"attribute"`
	LastID    int64            `json:"last_id"`
	Remark    string           `json:"remark"`
	CTime     utils.FormatTime `json:"ctime"`
}

// Operformat oper format.
func Operformat(tagID int8, old, new interface{}, style int8) (cont string) {
	var template string
	switch style {
	case OperStyleOne:
		template = "[%s]从[%v]设为[%v]"
	case OperStyleTwo:
		template = "[%s]%v:%v"
	}
	cont = fmt.Sprintf(template, _operType[tagID], old, new)
	return
}

// AccessState  get orange state
func AccessState(state int8, access int16) (newState int16) {
	if NormalState(state) && access == AccessMember {
		newState = access
		return
	}
	newState = int16(state)
	return
}
