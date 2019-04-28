package dao

import (
	"go-common/app/admin/ep/merlin/model"

	pkgerr "github.com/pkg/errors"
)

// InsertMailLog insert mail log.
func (d *Dao) InsertMailLog(ml *model.MailLog) (err error) {
	return pkgerr.WithStack(d.db.Create(&ml).Error)
}

// DelMailLog delete mail log.
func (d *Dao) DelMailLog(receiverName string) (err error) {
	err = pkgerr.WithStack(d.db.Where("receiver_name=?", receiverName).Delete(model.MailLog{}).Error)
	return
}

// FindMailLog find mail log.
func (d *Dao) FindMailLog(receiverName string) (mailLogs []*model.MailLog, err error) {
	err = pkgerr.WithStack(d.db.Where("receiver_name=?", receiverName).Find(&mailLogs).Error)
	return
}
