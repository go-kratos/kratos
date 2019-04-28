package business

import (
	"encoding/json"

	"go-common/app/admin/main/aegis/model/middleware"
	"go-common/library/log"
	xtime "go-common/library/time"
)

// BizCFG .
type BizCFG struct {
	ID         int64      `json:"id" gorm:"primary_key" form:"id"`
	BusinessID int64      `json:"business_id" gorm:"column:business_id" form:"business_id"`
	TP         int8       `json:"type" gorm:"column:type" form:"type"`
	Config     string     `json:"config" gorm:"column:config" form:"config"`
	State      int64      `json:"state" gorm:"column:state" form:"state"`
	Ctime      xtime.Time `json:"ctime" gorm:"column:ctime"`
	Mtime      xtime.Time `json:"mtime" gorm:"column:mtime"`
}

// TableName .
func (t *BizCFG) TableName() string {
	return "business_config"
}

// FormatMngBID .
func (t *BizCFG) FormatMngBID() (biz int64, roles map[int64][]int64, err error) {
	var (
		cfgs []map[string]int64
	)

	if t.TP != TypeManagerBID || t.Config == "" {
		return
	}

	biz = t.BusinessID
	roles = map[int64][]int64{}
	if err = json.Unmarshal([]byte(t.Config), &cfgs); err != nil {
		log.Error("FormatMngBID json.Unmarshal(%+v) error(%v)", t, err)
		return
	}

	for _, item := range cfgs {
		if item[MngBIDMID] <= 0 || item[MngBIDFlow] <= 0 {
			continue
		}
		if _, exist := roles[item[MngBIDMID]]; !exist {
			roles[item[MngBIDMID]] = []int64{}
		}
		roles[item[MngBIDMID]] = append(roles[item[MngBIDMID]], item[MngBIDFlow])
	}

	return
}

// FormatBizBID .
func (t *BizCFG) FormatBizBID() (biz int64, roles map[string]int64, err error) {
	if t.TP != TypeBizBID || t.Config == "" {
		return
	}

	biz = t.BusinessID
	if err = json.Unmarshal([]byte(t.Config), &roles); err != nil {
		log.Error("FormatBizBID json.Unmarshal(%+v) error(%v)", t, err)
	}
	return
}

//FormatAggregate 聚合配置
func (t *BizCFG) FormatAggregate() (cfgs []*middleware.Aggregate, err error) {
	if t.TP != TypeMiddleware || t.Config == "" {
		return
	}

	cfgs = []*middleware.Aggregate{}
	if err = json.Unmarshal([]byte(t.Config), &cfgs); err != nil {
		log.Error("FormatAggregate json.Unmarshal error(%+v)", err)
	}
	return
}

// ManagerRole .
type ManagerRole struct {
	FlowID int64 `json:"flow_id"`
	MngBid int64 `json:"manager_bid"`
}
