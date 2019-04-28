package audit

import (
	"context"
	xtime "time"

	"go-common/app/interface/main/tv/model"
	xsql "go-common/library/database/sql"
	"go-common/library/time"
)

const (
	_updateCont  = "UPDATE `tv_content` SET `state` = ?, `valid` = ?, `reason` = ?, `inject_time` = ? WHERE `epid` = ? AND `is_deleted` = 0"
	_updateSea   = "UPDATE `tv_ep_season` SET `check` = ?, `valid` = ?, `reason` = ?, `inject_time` = ? WHERE `id` = ? AND `is_deleted` = 0"
	_updateVideo = "UPDATE `ugc_video` SET `result` = ?, `valid` = ?, `reason` = ? , `inject_time` = ? WHERE `cid` = ? AND `deleted` = 0"
	_updateArc   = "UPDATE `ugc_archive` SET `result` = ?, `valid` = ?, `reason` = ?, `inject_time` = ? WHERE `aid` = ? AND `deleted` = 0"
)

// BeginTran def.
func (d *Dao) BeginTran(c context.Context) (tx *xsql.Tx, err error) {
	return d.db.Begin(c)
}

// UpdateVideo .
func (d *Dao) UpdateVideo(c context.Context, v *model.AuditOp, tx *xsql.Tx) (err error) {
	now := time.Time(xtime.Now().Unix())
	_, err = tx.Exec(_updateVideo, v.Result, v.Valid, v.AuditMsg, now, v.KID)
	return
}

// UpdateArc .
func (d *Dao) UpdateArc(c context.Context, v *model.AuditOp, tx *xsql.Tx) (err error) {
	now := time.Time(xtime.Now().Unix())
	_, err = tx.Exec(_updateArc, v.Result, v.Valid, v.AuditMsg, now, v.KID)
	return
}

// UpdateCont .
func (d *Dao) UpdateCont(c context.Context, val *model.AuditOp, tx *xsql.Tx) (err error) {
	now := time.Time(xtime.Now().Unix())
	_, err = tx.Exec(_updateCont, val.Result, val.Valid, val.AuditMsg, now, val.KID)
	return
}

// UpdateSea .
func (d *Dao) UpdateSea(c context.Context, val *model.AuditOp, tx *xsql.Tx) (err error) {
	now := time.Time(xtime.Now().Unix())
	_, err = tx.Exec(_updateSea, val.Result, val.Valid, val.AuditMsg, now, val.KID)
	return
}
