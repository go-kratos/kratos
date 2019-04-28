package model

import (
	xtime "go-common/library/time"
)

// PriceConfig represents price config of tv vip.
type PriceConfig struct {
	ID          int32      `json:"id"`           // 主键id
	Pid         int32      `json:"pid"`          // 父id，为空表示为原价信息
	Platform    int8       `json:"platform"`     // 类型: 1:tv安卓 2:公众号
	ProductName string     `json:"product_name"` // 产品展示名
	ProductId   string     `json:"product_id"`   // 产品id
	SuitType    int8       `json:"suit_type"`    // 适用人群: 0.所有用户 1.旧客 2.新客 3.续期旧客 4.续期新客 5.套餐旧客 6.套餐新客 10.主站vip专项
	Month       int32      `json:"month"`        // 月份单位
	SubType     int8       `json:"sub_type"`     // 订阅类型：0.其他，1.连续包月
	Price       int32      `json:"price"`        // 价格，pid为0表示原价,单位:分
	Selected    int8       `json:"selected"`     // 选中状态: 0.未选中，1.选中
	Remark      string     `json:"remark"`       // 促销tip
	Status      int8       `json:"status"`       // 状态，0:有效,1:失效
	Superscript string     `json:"superscript"`  // 角标
	Operator    string     `json:"operator"`     // 操作者
	OperId      int64      `json:"oper_id"`      // 操作者id
	Stime       xtime.Time `json:"stime"`        // 折扣开始时间
	Etime       xtime.Time `json:"etime"`        // 折扣结束时间
	Ctime       xtime.Time `json:"ctime"`        // 创建时间
	Mtime       xtime.Time `json:"mtime"`        // 最后修改时间

}

// PanelPriceConfig represents panel config of tv vip.
type PanelPriceConfig struct {
	PriceConfig
	MaxNum      int32 // 允许最大购买数量，-1 表示不限制
	OriginPrice int32 // 原价
}

// CopyFromPriceConfig copies fields from price config.
func (pi *PanelPriceConfig) CopyFromPriceConfig(pc *PriceConfig) {
	pi.ID = pc.ID
	pi.Pid = pc.Pid
	pi.Platform = pc.Platform
	pi.ProductName = pc.ProductName
	pi.ProductId = pc.ProductId
	pi.SuitType = pc.SuitType
	pi.Month = pc.Month
	pi.SubType = pc.SubType
	pi.Price = pc.Price
	pi.Selected = pc.Selected
	pi.Remark = pc.Remark
	pi.Status = pc.Status
	pi.Superscript = pc.Superscript
	pi.Operator = pc.Operator
	pi.OperId = pc.OperId
	pi.Stime = pc.Stime
	pi.Etime = pc.Etime
	pi.Ctime = pc.Ctime
	pi.Mtime = pc.Mtime
	pi.MaxNum = 1
}

// IsContracted returns true if panel is contracted package.
func (pi *PanelPriceConfig) IsContracted() bool {
	return pi.SubType == SubTypeContract
}

// PidOrId returns panel parent id or panel id.
func (pi *PanelPriceConfig) PidOrId() int32 {
	if pi.Pid != 0 {
		return pi.Pid
	}
	return pi.ID
}
