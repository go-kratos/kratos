package dao

import (
	"context"
	"go-common/app/service/bbq/user/internal/model"
	"go-common/library/log"
)

const (
	_insertforbidUser = "insert into forbid_user (`mid`, `expire_time`, `forbid_status`) values (?, ?, ?) on duplicate key update `expire_time` = VALUES(`expire_time`), `forbid_status` = VALUES(`forbid_status`)"
)

//ForbidUser .
func (d *Dao) ForbidUser(c context.Context, mid uint64, exTime uint64) (err error) {
	if _, err = d.db.Exec(c, _insertforbidUser, mid, exTime, model.ForbiddenStatus); err != nil {
		log.Errorw(c, "event", "ForbidUser", "err", err)
		return
	}
	return
}

//ReleaseUser ..
func (d *Dao) ReleaseUser(c context.Context, mid uint64) (err error) {
	if _, err = d.db.Exec(c, _insertforbidUser, mid, 0, model.NormalStatus); err != nil {
		log.Errorw(c, "event", "ReleaseUser", "err", err)
		return
	}
	return
}
