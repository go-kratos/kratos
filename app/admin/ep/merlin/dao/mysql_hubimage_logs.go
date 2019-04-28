package dao

import (
	"strconv"

	"go-common/app/admin/ep/merlin/model"
	"go-common/library/ecode"

	pkgerr "github.com/pkg/errors"
)

// InsertHubImageLog Insert Hub Image Log.
func (d *Dao) InsertHubImageLog(hubImageLog *model.HubImageLog) (err error) {
	return pkgerr.WithStack(d.db.Create(hubImageLog).Error)
}

// UpdateHubImageLog Update Hub Image Log.
func (d *Dao) UpdateHubImageLog(hubImageLog *model.HubImageLog) (err error) {
	return pkgerr.WithStack(d.db.Model(&model.HubImageLog{}).Where("id=?", hubImageLog.ID).Updates(hubImageLog).Error)
}

// UpdateHubImageLogStatus Update Hub Image Log Status.
func (d *Dao) UpdateHubImageLogStatus(logID int64, status int) (err error) {
	return pkgerr.WithStack(d.db.Model(&model.HubImageLog{}).Where("id=?", logID).Updates(map[string]interface{}{"status": status}).Error)
}

// FindHubImageLogByImageTag Find Hub Image Log By Image Tag.
func (d *Dao) FindHubImageLogByImageTag(imageTag string) (hubImageLog *model.HubImageLog, err error) {
	hubImageLog = &model.HubImageLog{}
	if err = d.db.Where("imagetag = ?", imageTag).First(hubImageLog).Error; err == ecode.NothingFound {
		err = nil
	}
	return
}

// FindHubImageLogByMachineID Find Hub Image Log By MachineID .
func (d *Dao) FindHubImageLogByMachineID(machineID int64) (hubImageLogs []*model.HubImageLog, err error) {
	if err = d.db.Where("machine_id = ? ", machineID).Order("id DESC").Find(&hubImageLogs).Error; err == ecode.NothingFound {
		err = nil
	}
	return
}

// FindHubMachine2ImageLog Find Hub Machine to Image Log .
func (d *Dao) FindHubMachine2ImageLog(queryRequest *model.QueryMachine2ImageLogRequest) (total int64, hubImageLogs []*model.HubImageLog, err error) {
	cdb := d.db.Model(&model.HubImageLog{}).Where("machine_id=? and operate_type = ?", queryRequest.MachineID, strconv.Itoa(model.ImageMachine2Image)).Order("ID desc").Offset((queryRequest.PageNum - 1) * queryRequest.PageSize).Limit(queryRequest.PageSize)
	if err = pkgerr.WithStack(cdb.Find(&hubImageLogs).Error); err != nil {
		return
	}
	if err = pkgerr.WithStack(d.db.Model(&model.HubImageLog{}).Where("machine_id=? and operate_type = ?", queryRequest.MachineID, strconv.Itoa(model.ImageMachine2Image)).Count(&total).Error); err != nil {
		return
	}
	return
}

// UpdateHubImageLogStatusInDoingStatus Update Hub Image Log Status in doing status.
func (d *Dao) UpdateHubImageLogStatusInDoingStatus(machineID int64, status int) (err error) {
	return pkgerr.WithStack(d.db.Model(&model.HubImageLog{}).Where("machine_id=? and status=? and operate_type = ?", machineID, model.ImageInit, model.ImageMachine2Image).Updates(map[string]interface{}{"status": status}).Error)
}
