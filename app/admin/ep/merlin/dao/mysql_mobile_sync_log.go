package dao

import (
	"go-common/app/admin/ep/merlin/model"

	pkgerr "github.com/pkg/errors"
)

// InsertMobileSyncLog Insert Mobile Sync Log.
func (d *Dao) InsertMobileSyncLog(mobileSyncLog *model.MobileSyncLog) (err error) {
	return pkgerr.WithStack(d.db.Create(mobileSyncLog).Error)
}

// UpdateMobileSyncLog Update Mobile SyncLog.
func (d *Dao) UpdateMobileSyncLog(mobileSyncLog *model.MobileSyncLog) (err error) {
	return pkgerr.WithStack(d.db.Model(&model.MobileSyncLog{}).Where("uuid=?", mobileSyncLog.UUID).Update(mobileSyncLog).Error)
}

// FindMobileSyncLogStartStatus FindMobile SyncLog Start Status.
func (d *Dao) FindMobileSyncLogStartStatus() (startCount int, err error) {
	err = pkgerr.WithStack(d.db.Model(&model.MobileSyncLog{}).Where("status=0 and mtime>=DATE_SUB(NOW(),INTERVAL 10 MINUTE)").Count(&startCount).Error)
	return
}
