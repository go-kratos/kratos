package model

import "go-common/library/time"

// VipOrderActivityRecord vip record.
type VipOrderActivityRecord struct {
	ID             int64
	Mid            int64
	OrderNo        string
	ProductID      string
	Months         int32
	PanelType      string
	AssociateState int32
	Ctime          time.Time
	Mtime          time.Time
}
