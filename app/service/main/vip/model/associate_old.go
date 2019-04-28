package model

import "go-common/library/time"

// VipOrderActivityRecord vip order activity record.
type VipOrderActivityRecord struct {
	ID             int64  `json:"id"`
	Mid            int64  `json:"mid"`
	OrderNO        string `json:"order_no"`
	ProductID      string `json:"product_id"`
	Months         int32  `json:"months"`
	PanelType      string `json:"panel_type"`
	AssociateState int8   `json:"associate_state"`
	Ctime          time.Time
	Mtime          time.Time
}
