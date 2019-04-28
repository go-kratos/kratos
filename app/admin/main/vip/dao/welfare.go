package dao

import (
	"fmt"
	"strings"
	"time"

	"go-common/app/admin/main/vip/model"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

const (
	_welfareTypeTable       = "vip_welfare_type"
	_welfareTable           = "vip_welfare"
	_codeBatchTable         = "vip_welfare_code_batch"
	_welfareCodeTable       = "vip_welfare_code"
	_notDelete              = 0
	_nobody                 = 0
	_noTid                  = 0
	_batchInsertWelfareCode = "INSERT INTO vip_welfare_code (bid, wid, code) VALUES %s"
)

// WelfareTypeAdd add welfare type
func (d *Dao) WelfareTypeAdd(wt *model.WelfareType) (err error) {
	if err = d.vip.Save(wt).Error; err != nil {
		err = errors.Wrapf(err, "WelfareTypeAdd(%+v)", wt)
	}
	return
}

// WelfareTypeUpd update welfare type
func (d *Dao) WelfareTypeUpd(wt *model.WelfareType) (err error) {
	if err = d.vip.Table(_welfareTypeTable).Where("id = ? and state = ?", wt.ID, _notDelete).
		Update(map[string]interface{}{
			"oper_id":   wt.OperID,
			"oper_name": wt.OperName,
			"state":     wt.State,
			"name":      wt.Name,
		}).Error; err != nil {
		err = errors.Wrapf(err, "WelfareTypeUpd(%+v)", wt)
	}
	return
}

// WelfareTypeState delete welfare type
func (d *Dao) WelfareTypeState(tx *gorm.DB, id, state, operId int, operName string) (err error) {
	if err = tx.Table(_welfareTypeTable).Where("id = ?", id).
		Update(map[string]interface{}{
			"oper_id":   operId,
			"oper_name": operName,
			"state":     state,
		}).Error; err != nil {
		err = errors.Wrapf(err, "WelfareTypeState id(%v) state(%v)", id, state)
	}
	return
}

// WelfareTypeList get welfare type list
func (d *Dao) WelfareTypeList() (wts []*model.WelfareTypeRes, err error) {
	if err = d.vip.Table(_welfareTypeTable).Where("state = ?", _notDelete).Find(&wts).Error; err != nil {
		err = errors.Wrapf(err, "WelfareTypeList")
	}
	return
}

// WelfareAdd add welfare
func (d *Dao) WelfareAdd(wt *model.Welfare) (err error) {
	if err = d.vip.Save(wt).Error; err != nil {
		err = errors.Wrapf(err, "WelfareAdd(%+v)", wt)
	}
	return
}

// WelfareUpd update welfare
func (d *Dao) WelfareUpd(wt *model.WelfareReq) (err error) {
	if err = d.vip.Table(_welfareTable).Where("id = ? and state = ?", wt.ID, _notDelete).
		Update(map[string]interface{}{
			"welfare_name": wt.WelfareName,
			"welfare_desc": wt.WelfareDesc,
			"homepage_uri": wt.HomepageUri,
			"backdrop_uri": wt.BackdropUri,
			"recommend":    wt.Recommend,
			"rank":         wt.Rank,
			"tid":          wt.Tid,
			"stime":        wt.Stime,
			"etime":        wt.Etime,
			"usage_form":   wt.UsageForm,
			"receive_rate": wt.ReceiveRate,
			"receive_uri":  wt.ReceiveUri,
			"vip_type":     wt.VipType,
			"oper_id":      wt.OperID,
			"oper_name":    wt.OperName,
		}).Error; err != nil {
		err = errors.Wrapf(err, "WelfareUpd(%+v)", wt)
	}
	return
}

// WelfareState delete welfare
func (d *Dao) WelfareState(id, state, operId int, operName string) (err error) {
	if err = d.vip.Table(_welfareTable).Where("id = ?", id).
		Update(map[string]interface{}{
			"oper_id":   operId,
			"oper_name": operName,
			"state":     state,
		}).Error; err != nil {
		err = errors.Wrapf(err, "WelfareState id(%v) state(%v)", id, state)
	}
	return
}

// ResetWelfareTid reset welfare tid to 0
func (d *Dao) ResetWelfareTid(tx *gorm.DB, tid int) (err error) {
	if err = tx.Table(_welfareTable).Where("tid = ?", tid).Update("tid", _noTid).Error; err != nil {
		err = errors.Wrapf(err, "ResetWelfareTid(%v)", tid)
	}
	return
}

// WelfareList get welfare list
func (d *Dao) WelfareList(tid int) (ws []*model.WelfareRes, err error) {
	db := d.vip.Table(_welfareTable)
	if tid != 0 {
		db = db.Where("tid = ?", tid)
	}
	if err = db.Where("state = ?", _notDelete).Order("recommend desc, rank").Find(&ws).Error; err != nil {
		err = errors.Wrapf(err, "WelfareList(%+v)", tid)
	}
	return
}

// WelfareBatchSave add welfare batch
func (d *Dao) WelfareBatchSave(wcb *model.WelfareCodeBatch) (err error) {
	if err = d.vip.Table(_codeBatchTable).Save(wcb).Error; err != nil {
		err = errors.Wrapf(err, "WelfareBatchSave(%+v)", wcb)
	}
	return
}

// WelfareBatchList get welfare list
func (d *Dao) WelfareBatchList(wid int) (wbs []*model.WelfareBatchRes, err error) {
	if err = d.vip.Table(_codeBatchTable).Where("wid = ? and state = ?", wid, _notDelete).Find(&wbs).Error; err != nil {
		err = errors.Wrapf(err, "WelfareBatchList(%+v)", wid)
	}
	return
}

// WelfareBatchState delete welfare batch
func (d *Dao) WelfareBatchState(tx *gorm.DB, id, state, operId int, operName string) (err error) {
	if err = tx.Table(_codeBatchTable).Where("id = ?", id).
		Update(map[string]interface{}{
			"oper_id":   operId,
			"oper_name": operName,
			"state":     state,
		}).Error; err != nil {
		err = errors.Wrapf(err, "WelfareBatchState(%+v)", id)
	}

	return
}

// WelfareCodeBatchInsert insert welfare batch code
func (d *Dao) WelfareCodeBatchInsert(wcs []*model.WelfareCode) (err error) {
	log.Info("WelfareCodeBatchInsert start time (%s)", time.Now())
	var (
		buff    = make([]*model.WelfareCode, 2000)
		buffEnd = 0
	)
	for _, wc := range wcs {
		buff[buffEnd] = wc
		buffEnd++
		if buffEnd >= 2000 {
			buffEnd = 0
			stmt, valueArgs := getBatchInsertSQL(buff)
			if err = d.vip.Exec(stmt, valueArgs...).Error; err != nil {
				return
			}
		}
	}
	if buffEnd > 0 {
		stmt, valueArgs := getBatchInsertSQL(buff[:buffEnd])
		buffEnd = 0
		if err = d.vip.Exec(stmt, valueArgs...).Error; err != nil {
			return
		}
	}
	log.Info("WelfareCodeBatchInsert end time (%s)", time.Now())
	return
}

// WelfareCodeStatus delete welfare batch code
func (d *Dao) WelfareCodeStatus(tx *gorm.DB, bid, state int) (err error) {
	if err = tx.Table(_welfareCodeTable).Where("bid = ? and mid = ?", bid, _nobody).
		Update("state", state).Error; err != nil {
		err = errors.Wrapf(err, "WelfareCodeStatus(%+v) (%+v)", bid, state)
	}

	return
}

func getBatchInsertSQL(buff []*model.WelfareCode) (stmt string, valueArgs []interface{}) {
	values := []string{}
	for _, b := range buff {
		values = append(values, "(?,?,?)")
		valueArgs = append(valueArgs, b.Bid, b.Wid, b.Code)
	}
	stmt = fmt.Sprintf(_batchInsertWelfareCode, strings.Join(values, ","))
	return
}
