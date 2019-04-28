package dao

import (
	"go-common/app/admin/ep/merlin/model"

	pkgerr "github.com/pkg/errors"
)

// InsertMobileMachineErrorLog Insert Mobile Machine Error Log.
func (d *Dao) InsertMobileMachineErrorLog(mobileMachineErrorLog *model.MobileMachineErrorLog) (err error) {
	return pkgerr.WithStack(d.db.Create(mobileMachineErrorLog).Error)
}

// FindMobileMachineErrorLog Find Mobile Machine Error Log.
func (d *Dao) FindMobileMachineErrorLog(queryRequest *model.QueryMobileMachineErrorLogRequest) (total int64, mobileMachineErrorLogs []*model.MobileMachineErrorLog, err error) {
	cdb := d.db.Model(&model.MobileMachineErrorLog{}).Where("machine_id=?", queryRequest.MachineID).Order("ID desc").Offset((queryRequest.PageNum - 1) * queryRequest.PageSize).Limit(queryRequest.PageSize)
	if err = pkgerr.WithStack(cdb.Find(&mobileMachineErrorLogs).Error); err != nil {
		return
	}
	if err = pkgerr.WithStack(d.db.Model(&model.MobileMachineErrorLog{}).Where("machine_id=?", queryRequest.MachineID).Count(&total).Error); err != nil {
		return
	}
	return
}
