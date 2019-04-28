package model

import (
	"time"

	xtime "go-common/library/time"
)

// UserInfo represents user info.
type UserInfo struct {
	ID            int32      `json:"id"`              // 用户信息表
	Mid           int64      `json:"mid"`             // 用户mid
	Ver           int32      `json:"ver"`             // 版本控制
	VipType       int8       `json:"vip_type"`        // tv-vip类型:1.vip 2.年费vip
	PayType       int8       `json:"pay_type"`        // tv-vip购买类型:0.正常购买 1.连续包月
	PayChannelId  string     `json:"pay_channel_id"`  // 自动续费渠道:wechat,alipay
	Status        int8       `json:"status"`          // tv-vip状态:0:过期 1:未过期
	OverdueTime   xtime.Time `json:"overdue_time"`    // tv-vip过期时间
	RecentPayTime xtime.Time `json:"recent_pay_time"` // tv-vip最近开通时间
	Ctime         xtime.Time `json:"ctime"`           // 创建时间
	Mtime         xtime.Time `json:"mtime"`           // 修改时间
}

// IsEmpty returns true if user id equals -1.
func (ui *UserInfo) IsEmpty() bool {
	return ui.ID == -1
}

// IsExpired returns true if user is expired vip.
func (ui *UserInfo) IsExpired() bool {
	return ui.OverdueTime < xtime.Time(time.Now().Unix())
}

// MarkExpired sets user status to expired status.
func (ui *UserInfo) MarkExpired() {
	ui.Status = 0
}

// IsVip returns true if user is vip.
func (ui *UserInfo) IsVip() bool {
	if ui.IsEmpty() {
		return false
	}
	if ui.IsExpired() {
		return false
	}
	return ui.Status == 1
}

// IsContracted returns true if user buys contracted package.
func (ui *UserInfo) IsContracted() bool {
	return ui.PayType == 1
}

// CopyFromPayOrder copies fileds from pay order.
func (ui *UserInfo) CopyFromPayOrder(po *PayOrder) {
	ui.VipType = VipTypeVip
	ui.PayChannelId = po.PaymentType
	ui.RecentPayTime = xtime.Time(time.Now().Unix())
}

// CopyFromPanel copies field from panel.
func (ui *UserInfo) CopyFromPanel(p *PanelPriceConfig) {
	if p.SubType == SubTypeContract {
		ui.PayType = VipPayTypeSub
		return
	}
	ui.PayType = VipPayTypeNormal
}
