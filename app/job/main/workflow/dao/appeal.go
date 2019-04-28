package dao

import (
	"context"
	"time"

	"go-common/app/job/main/workflow/model"

	"github.com/jinzhu/gorm"
)

// consts for workflow business_state
const (
	BusStCreated    = int8(1) // 未处理
	BusStRead       = int8(2) // 已回复已读
	BusStAutoClosed = int8(5) // 过期自动关闭
	BusStNotRead    = int8(6) // 已回复未读
)

// Appeals .
func (d *Dao) Appeals(c context.Context, ids []int64) (appeals []*model.Appeal, err error) {
	err = d.ReadORM.Table("workflow_appeal").Where("id in (?)", ids).Find(&appeals).Error
	return
}

// SetAppealTransferState will close expired feedback 关闭超时的申诉 (用户未评价)
func (d *Dao) SetAppealTransferState(c context.Context, ids []int64, transferState int8) (err error) {
	err = d.WriteORM.Table("workflow_appeal").Where("id IN (?)", ids).Update("transfer_state", transferState).
		Update("ttime", time.Now().Format("2006-01-02 15:04:05")).Error
	return
}

// TxSetWeight db覆盖权重值
func (d *Dao) TxSetWeight(tx *gorm.DB, newWeight map[int64]int64) (err error) {
	for id, weight := range newWeight {
		if err = tx.Table("workflow_appeal").Where("id = ?", id).Update("weight", weight).Error; err != nil {
			return
		}
	}
	return
}

// SetAppealAssignState .
func (d *Dao) SetAppealAssignState(c context.Context, ids []int64, assignState int8) (err error) {
	return d.WriteORM.Table("workflow_appeal").Where("id in (?)", ids).Update("assign_state", assignState).Error
}

// LastEvent return last event of appeal_id
func (d *Dao) LastEvent(id int64) (e *model.Event, err error) {
	e = new(model.Event)
	err = d.ReadORM.Table("workflow_event").Where("appeal_id = ?", id).Last(e).Error
	return
}
