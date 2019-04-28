package dao

import (
	"go-common/app/admin/ep/merlin/model"
	"go-common/library/ecode"

	pkgerr "github.com/pkg/errors"
)

// InsertSnapshotRecord Insert Snapshot Record.
func (d *Dao) InsertSnapshotRecord(snapshotRecord *model.SnapshotRecord) (err error) {
	return pkgerr.WithStack(d.db.Create(snapshotRecord).Error)
}

// FindSnapshotRecord Find Snapshot Record.
func (d *Dao) FindSnapshotRecord(machineID int64) (snapshotRecord *model.SnapshotRecord, err error) {
	snapshotRecord = &model.SnapshotRecord{}
	if err = d.db.Where("machine_id = ?", machineID).First(snapshotRecord).Error; err == ecode.NothingFound {
		err = nil
	}
	return
}

// UpdateSnapshotRecordStatus Update Snapshot Record Status.
func (d *Dao) UpdateSnapshotRecordStatus(machineID int64, status string) (err error) {
	return pkgerr.WithStack(d.db.Model(&model.SnapshotRecord{}).Where("machine_id = ?", machineID).Update("status", status).Error)
}

// UpdateSnapshotRecord Update Snapshot Record.
func (d *Dao) UpdateSnapshotRecord(snapshotRecord *model.SnapshotRecord) (err error) {
	return pkgerr.WithStack(d.db.Model(&model.SnapshotRecord{}).Where("machine_id = ?", snapshotRecord.MachineID).Update(snapshotRecord).Error)
}

// FindSnapshotStatusInDoingOver2Hours Find Snapshot Status In Doing Over 2 Hours.
func (d *Dao) FindSnapshotStatusInDoingOver2Hours() (snapshotRecords []*model.SnapshotRecord, err error) {
	if err = d.db.Where("status = ? and unix_timestamp(now())> unix_timestamp(date_add(mtime, interval 2 hour))", model.SnapshotDoing).Find(&snapshotRecords).Error; err == ecode.NothingFound {
		err = nil
	}
	return
}
