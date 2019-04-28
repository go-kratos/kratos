package model

import (
	"strconv"
	"time"

	"go-common/library/log"
	xtime "go-common/library/time"
)

// UserChangeHistory 会员变动流水.
type UserChangeHistory struct {
	ID         int32      `json:"id"`          // vip开通历史
	Mid        int64      `json:"mid"`         // 用户mid
	ChangeType int8       `json:"change_type"` // 变更类型(1:充值开通 2:系统发放 3:活动赠送 4:重复领取扣除)
	ChangeTime xtime.Time `json:"change_time"` // 变更时间
	OrderNo    string     `json:"order_no"`    // 关联订单号
	Days       int32      `json:"days"`        // 开通天数
	OperatorId string     `json:"operator_id"` // 操作人id
	Remark     string     `json:"remark"`      // 备注
	Ctime      xtime.Time `json:"ctime"`       // 创建时间
	Mtime      xtime.Time `json:"mtime"`       // 修改时间
}

func (uc *UserChangeHistory) orderType2ChangeType(orderType int8) int8 {
	var ct int8
	switch orderType {
	case PayOrderTypeNormal:
		ct = UserChangeTypeRecharge
	case PayOrderTypeSub:
		ct = UserChangeTypeSystem
	default:
		log.Error("uc.CopyFromPayOrder() err(UnknownOrderType) orderType(%d)", orderType)
		ct = UserChangeTypeRecharge
	}
	return ct
}

// CopyFromPayOrder copies fields from pay order.
func (uc *UserChangeHistory) CopyFromPayOrder(po *PayOrder) {
	uc.Mid = po.Mid
	uc.OrderNo = po.OrderNo
	uc.Days = int32(po.BuyMonths) * 31
	uc.OperatorId = strconv.Itoa(int(po.Mid))
	uc.ChangeTime = xtime.Time(time.Now().Unix())
	uc.ChangeType = uc.orderType2ChangeType(po.OrderType)
}
