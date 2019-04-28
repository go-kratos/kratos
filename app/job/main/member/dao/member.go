package dao

import (
	"context"
	"database/sql"
	"fmt"

	"go-common/app/job/main/member/model"
	smodel "go-common/app/service/main/member/model"
	"go-common/library/log"
)

const (
	_shard            = 100
	_SelBaseInfo      = "SELECT mid,name,sex,face,sign,rank,birthday,ctime,mtime FROM user_base_%02d WHERE mid=?"
	_SetOfficial      = "INSERT INTO user_official (mid,role,title,description) values (?,?,?,?) ON DUPLICATE KEY UPDATE role=VALUES(role),title=VALUES(title),description=VALUES(description)"
	_SetBaseInfo      = "INSERT INTO user_base_%02d(mid,name,sex,face,sign,rank) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE name=?,sex=?,face=?,sign=?,rank=?"
	_SetName          = "INSERT INTO user_base_%02d(mid,name) VALUES (?,?) ON DUPLICATE KEY UPDATE name=?"
	_SetSign          = "INSERT INTO user_base_%02d(mid,sign) VALUES (?,?) ON DUPLICATE KEY UPDATE sign=?"
	_InitBase         = "INSERT IGNORE INTO user_base_%02d(mid) VALUES (?)"
	_SetFace          = "UPDATE user_base_%02d SET face=? WHERE mid=?"
	_SelMoral         = "SELECT mid,moral,added,deducted,last_recover_date from user_moral where mid = ?"
	_recoverMoral     = `update user_moral set moral=moral+?,added=added+?, last_recover_date=? where mid = ? and last_recover_date != ? and moral < 7000 `
	_auditQueuingFace = `UPDATE user_property_review SET state=?, operator='system', remark=? WHERE property=1 AND state=10 AND is_monitor=0 AND new=?`
)

// hit get table suffix
func hit(id int64) int64 {
	return id % _shard
}

// BaseInfo info of user.
//用于检查数据，不需要拼url
func (d *Dao) BaseInfo(c context.Context, mid int64) (r *model.BaseInfo, err error) {
	r = &model.BaseInfo{}
	row := d.db.QueryRow(c, fmt.Sprintf(_SelBaseInfo, hit(mid)), mid)
	if err = row.Scan(&r.Mid, &r.Name, &r.Sex, &r.Face, &r.Sign, &r.Rank, &r.Birthday, &r.CTime, &r.MTime); err != nil {
		r = nil
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("row.Scan() error(%v) mid(%v)", err, mid)
	}
	return
}

// SelMoral select moral by mid from db
func (d *Dao) SelMoral(c context.Context, mid int64) (moral *smodel.Moral, err error) {
	moral = &smodel.Moral{}
	row := d.db.QueryRow(c, _SelMoral, mid)
	if err = row.Scan(&moral.Mid, &moral.Moral, &moral.Added, &moral.Deducted, &moral.LastRecoverDate); err != nil {
		if err == sql.ErrNoRows {
			moral = nil
			err = nil
			return
		}
		log.Error(" SelMoral row.Scan() error(%v) mid(%v)", err, mid)
	}
	return
}

// RecoverMoral set moral by mid from db
func (d *Dao) RecoverMoral(c context.Context, mid, delta, added int64, lastRecoverDate string) (rowsAffected int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _recoverMoral, delta, added, lastRecoverDate, mid, lastRecoverDate); err != nil {
		log.Error(" RecoverMoral d.db.Exec() error(%v) mid(%v)", err, mid)
		return
	}
	return res.RowsAffected()
}

// SetSign set sign.
func (d *Dao) SetSign(c context.Context, mid int64, sign string) (err error) {
	if _, err = d.db.Exec(c, fmt.Sprintf(_SetSign, hit(mid)), mid, sign, sign); err != nil {
		log.Error("setSign: tx.Exec(%d, %v) error(%v)", mid, sign, err)
		return
	}
	return
}

// SetOfficial set official.
func (d *Dao) SetOfficial(c context.Context, mid int64, role int8, title string, desc string) error {
	if _, err := d.db.Exec(c, _SetOfficial, mid, role, title, desc); err != nil {
		log.Error("_SetOfficial: tx.Exec(%d,%d,%s,%s) error(%v)", mid, role, title, desc, err)
		return err
	}
	return nil
}

// SetName set name.
func (d *Dao) SetName(c context.Context, mid int64, name string) (err error) {
	if len(name) <= 0 {
		log.Error("SetName: mid(%d) len(name)<=0 ", mid)
		return
	}
	if _, err = d.db.Exec(c, fmt.Sprintf(_SetName, hit(mid)), mid, name, name); err != nil {
		log.Error("setName: tx.Exec(%d, %v) error(%v)", mid, name, err)
		return
	}
	return
}

// SetBaseInfo set base info of user.
func (d *Dao) SetBaseInfo(c context.Context, r *model.BaseInfo) (err error) {
	if len(r.Name) <= 0 {
		log.Error("SetBaseInfo: mid(%d) len(r.Name)<=0 BaseInfo(%v)", r.Mid, r)
		return
	}
	if _, err = d.db.Exec(c, fmt.Sprintf(_SetBaseInfo, hit(r.Mid)), r.Mid, r.Name, r.Sex, r.Face, r.Sign, r.Rank, r.Name, r.Sex, r.Face, r.Sign, r.Rank); err != nil {
		log.Error("SetBaseInfo: db.Exec(%d, %v) error(%v)", r.Mid, r, err)
	}
	return
}

// SetFace set face.
func (d *Dao) SetFace(c context.Context, mid int64, face string) (err error) {
	if _, err = d.db.Exec(c, fmt.Sprintf(_SetFace, hit(mid)), face, mid); err != nil {
		log.Error("SetFace: tx.Exec(%v,%v) error(%v)", face, mid, err)
	}
	return
}

// InitBase init base info.
func (d *Dao) InitBase(c context.Context, mid int64) (err error) {
	if _, err = d.db.Exec(c, fmt.Sprintf(_InitBase, hit(mid)), mid); err != nil {
		log.Error("InitBase: tx.Exec(%d) error(%v)", mid, err)
	}
	return
}

// AuditQueuingFace auditQueuingFace.
func (d *Dao) AuditQueuingFace(c context.Context, face string, remark string, state int8) error {
	_, err := d.db.Exec(c, _auditQueuingFace, state, remark, face)
	if err != nil {
		log.Error("Failed to audit queuing face: %s, remark: %s, state: %d: %+v", face, remark, state, err)
		return err
	}
	return nil
}
