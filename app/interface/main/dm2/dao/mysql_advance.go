package dao

import (
	"context"
	"time"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	// advance comment
	_addAdvSQL     = "INSERT INTO dm_advancecomment (dm_inid,owner,mid,type,timestamp,mode,refund) VALUES (?,?,?,?,?,?,?)"
	_selAdvSQL     = "SELECT buy_id,owner,dm_inid,type,mode,mid,timestamp,refund FROM dm_advancecomment WHERE owner=? AND buy_id=?"
	_selAdvsSQL    = "SELECT buy_id,owner,dm_inid,type,mode,mid,timestamp,refund FROM dm_advancecomment WHERE owner=? ORDER BY buy_id DESC limit 100"
	_selAdvModeSQL = "SELECT type FROM dm_advancecomment WHERE dm_inid=? AND mid=? AND mode=?"
	_upAdvTypeSQL  = "UPDATE dm_advancecomment SET type=? WHERE buy_id=?"
	_delAdvSQL     = "DELETE FROM dm_advancecomment WHERE buy_id=?"
	_selAdvanceCmt = "SELECT buy_id,owner,dm_inid,type,mode,mid,timestamp,refund FROM dm_advancecomment WHERE dm_inid=? AND mid=? AND mode=?"
)

// AdvanceType get advance type by cid,mid and mode.
func (d *Dao) AdvanceType(c context.Context, cid int64, mid int64, mode string) (typ string, err error) {
	row := d.dbDM.QueryRow(c, _selAdvModeSQL, cid, mid, mode)
	if err = row.Scan(&typ); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// Advance 获取购买高级弹幕功能状态
func (d *Dao) Advance(c context.Context, mid, id int64) (adv *model.Advance, err error) {
	row := d.dbDM.QueryRow(c, _selAdvSQL, mid, id)
	adv = &model.Advance{}
	if err = row.Scan(&adv.ID, &adv.Owner, &adv.Cid, &adv.Type, &adv.Mode, &adv.Mid, &adv.Timestamp, &adv.Refund); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			adv = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// Advances 获取高级弹幕申请列表
func (d *Dao) Advances(c context.Context, owner int64) (res []*model.Advance, err error) {
	rows, err := d.dbDM.Query(c, _selAdvsSQL, owner)
	if err != nil {
		log.Error("d.dbDM.Query(%s,%d) error(%v)", _selAdvsSQL, owner, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		adv := &model.Advance{}
		if err = rows.Scan(&adv.ID, &adv.Owner, &adv.Cid, &adv.Type, &adv.Mode, &adv.Mid, &adv.Timestamp, &adv.Refund); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		res = append(res, adv)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// BuyAdvance 购买高级弹幕功能
func (d *Dao) BuyAdvance(c context.Context, mid, cid, owner, refund int64, typ, mode string) (id int64, err error) {
	now := time.Now().Unix()
	res, err := d.dbDM.Exec(c, _addAdvSQL, cid, owner, mid, typ, now, mode, refund)
	if err != nil {
		log.Error("d.dbDM.Exec(cid:%d,mid:%d,type:%s) error(%v)", cid, mid, typ, err)
		return
	}
	return res.LastInsertId()
}

// UpdateAdvType 更新购买高级弹幕类型
func (d *Dao) UpdateAdvType(c context.Context, id int64, typ string) (affect int64, err error) {
	res, err := d.dbDM.Exec(c, _upAdvTypeSQL, typ, id)
	if err != nil {
		log.Error("d.dbDM.Exec(%s,%s,%d) error(%v)", _upAdvTypeSQL, typ, id, err)
		return
	}
	return res.RowsAffected()
}

// DelAdvance 删除购买高级弹幕记录
func (d *Dao) DelAdvance(c context.Context, id int64) (affect int64, err error) {
	res, err := d.dbDM.Exec(c, _delAdvSQL, id)
	if err != nil {
		log.Error("d.dbDM.Exec(%s,%d) error(%v)", _delAdvSQL, id, err)
		return
	}
	return res.RowsAffected()
}

// AdvanceCmt get advance comment.
func (d *Dao) AdvanceCmt(c context.Context, oid, mid int64, mode string) (adv *model.AdvanceCmt, err error) {
	adv = &model.AdvanceCmt{}
	row := d.dbDM.QueryRow(c, _selAdvanceCmt, oid, mid, mode)
	if err = row.Scan(&adv.ID, &adv.Owner, &adv.Oid, &adv.Type, &adv.Mode, &adv.Mid, &adv.Timestamp, &adv.Refund); err != nil {
		if err == sql.ErrNoRows {
			adv = nil
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
	}
	return
}
