package dao

import (
	"time"

	"go-common/app/admin/ep/merlin/model"

	pkgerr "github.com/pkg/errors"
)

// InsertApplicationRecord insert application record.
func (d *Dao) InsertApplicationRecord(ar *model.ApplicationRecord) (err error) {
	return pkgerr.WithStack(d.db.Create(&ar).Error)
}

// FindApplicationRecordsByApplicant find application record by applicant.
func (d *Dao) FindApplicationRecordsByApplicant(applicant string, pn, ps int) (total int64, ars []*model.ApplicationRecord, err error) {
	cdb := d.db.Model(&model.ApplicationRecord{}).Where("applicant=?", applicant).Order("ID desc").Offset((pn - 1) * ps).Limit(ps)
	if err = pkgerr.WithStack(cdb.Find(&ars).Error); err != nil {
		return
	}
	if err = pkgerr.WithStack(d.db.Model(&model.ApplicationRecord{}).Where("applicant=?", applicant).Count(&total).Error); err != nil {
		return
	}
	return
}

// FindApplicationRecordsByAuditor find application record by auditor.
func (d *Dao) FindApplicationRecordsByAuditor(auditor string, pn, ps int) (total int64, ars []*model.ApplicationRecord, err error) {
	cdb := d.db.Model(&model.ApplicationRecord{}).Where("auditor=?", auditor).Order("ID desc").Offset((pn - 1) * ps).Limit(ps)
	if err = pkgerr.WithStack(cdb.Find(&ars).Error); err != nil {
		return
	}
	if err = pkgerr.WithStack(d.db.Model(&model.ApplicationRecord{}).Where("auditor=?", auditor).Count(&total).Error); err != nil {
		return
	}
	return
}

// FindApplicationRecordsByMachineID find application records by machineId.
func (d *Dao) FindApplicationRecordsByMachineID(machineID int64, pn, ps int) (total int64, ars []*model.ApplicationRecord, err error) {
	cdb := d.db.Model(&model.ApplicationRecord{}).Where("machine_id=?", machineID).Order("ID desc").Offset((pn - 1) * ps).Limit(ps)
	if err = pkgerr.WithStack(cdb.Find(&ars).Error); err != nil {
		return
	}
	if err = pkgerr.WithStack(d.db.Model(&model.ApplicationRecord{}).Where("machine_id=?", machineID).Count(&total).Error); err != nil {
		return
	}
	return
}

// FindApplicationRecordsByID find application records by id.
func (d *Dao) FindApplicationRecordsByID(auditID int64) (ar *model.ApplicationRecord, err error) {
	ar = &model.ApplicationRecord{}
	err = pkgerr.WithStack(d.db.Model(&model.ApplicationRecord{}).Where("ID=?", auditID).First(ar).Error)
	return
}

// UpdateAuditStatus update audit status.
func (d *Dao) UpdateAuditStatus(auditID int64, status string) (err error) {
	return pkgerr.WithStack(d.db.Model(&model.ApplicationRecord{}).Where("ID=?", auditID).Update("STATUS", status).Error)
}

// UpdateAuditStatusAndComment update audit status and comment.
func (d *Dao) UpdateAuditStatusAndComment(auditID int64, status, comment string) (err error) {
	return pkgerr.WithStack(d.db.Model(&model.ApplicationRecord{}).Where("ID=?", auditID).Update("STATUS", status).Update("AUDITOR_COMMENT", comment).Error)
}

// InsertApplicationRecordAndUpdateMachineDelayStatus insert application record and update machine delay status.
func (d *Dao) InsertApplicationRecordAndUpdateMachineDelayStatus(ar *model.ApplicationRecord, machineID int64, delayStatus int) (err error) {
	tx := d.db.Begin()
	if err = tx.Error; err != nil {
		return
	}

	if err = tx.Create(&ar).Error; err != nil {
		tx.Rollback()
		return
	}
	if err = tx.Model(&model.Machine{}).Where("ID=?", machineID).Update("DELAY_STATUS", delayStatus).Error; err != nil {
		tx.Rollback()
		return
	}

	err = tx.Commit().Error
	return
}

// UpdateAuditStatusAndUpdateMachineDelayStatus update audit status and update machine delay status.
func (d *Dao) UpdateAuditStatusAndUpdateMachineDelayStatus(machineID, auditID int64, delayStatus int, applyStatus string) (err error) {
	tx := d.db.Begin()
	if err = tx.Error; err != nil {
		return
	}

	if err = tx.Model(&model.ApplicationRecord{}).Where("ID=?", auditID).Update("STATUS", applyStatus).Error; err != nil {
		tx.Rollback()
		return
	}

	if err = tx.Model(&model.Machine{}).Where("ID=?", machineID).Update("DELAY_STATUS", delayStatus).Error; err != nil {
		tx.Rollback()
		return
	}

	err = tx.Commit().Error
	return
}

// UpdateAuditStatusAndUpdateMachineDelayStatusComment update audit status and update machine delay status comment.
func (d *Dao) UpdateAuditStatusAndUpdateMachineDelayStatusComment(machineID, auditID int64, delayStatus int, applyStatus, comment string) (err error) {
	tx := d.db.Begin()
	if err = tx.Error; err != nil {
		return
	}

	if err = tx.Model(&model.ApplicationRecord{}).Where("ID=?", auditID).Update("STATUS", applyStatus).Update("AUDITOR_COMMENT", comment).Error; err != nil {
		tx.Rollback()
		return
	}

	if err = tx.Model(&model.Machine{}).Where("ID=?", machineID).Update("DELAY_STATUS", delayStatus).Error; err != nil {
		tx.Rollback()
		return
	}

	err = tx.Commit().Error
	return
}

// UpdateAuditStatusAndUpdateMachineEndTime update audit status and update machine end time.
func (d *Dao) UpdateAuditStatusAndUpdateMachineEndTime(machineID, auditID int64, delayStatus int, applyStatus string, applyEndTime time.Time, comment string) (err error) {
	tx := d.db.Begin()
	if err = tx.Error; err != nil {
		return
	}

	if err = tx.Model(&model.ApplicationRecord{}).Where("ID=?", auditID).Update("STATUS", applyStatus).Update("AUDITOR_COMMENT", comment).Error; err != nil {
		tx.Rollback()
		return
	}

	if err = tx.Model(&model.Machine{}).Where("ID=?", machineID).Update("DELAY_STATUS", delayStatus).Update("END_TIME", applyEndTime).Error; err != nil {
		tx.Rollback()
		return
	}

	err = tx.Commit().Error
	return
}
