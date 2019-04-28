package up

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-common/app/admin/main/mcn/model"
	xsql "go-common/library/database/sql"
	xtime "go-common/library/time"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_inMcnSignEntrySQL  = "INSERT mcn_sign(mcn_mid,begin_date,end_date,permission) VALUES (?,?,?,?)"
	_inMcnSignPaySQL    = "INSERT mcn_sign_pay(mcn_mid,sign_id,due_date,pay_value,note) VALUES (?,?,?,?,?)"
	_upMcnSignOPSQL     = "UPDATE mcn_sign SET state = ?, reject_time = ?, reject_reason = ? WHERE id = ?"
	_upMcnUpOPSQL       = "UPDATE mcn_up SET state = ?, reject_time = ?, reject_reason = ?, state_change_time = ? WHERE id = ?"
	_mcnUpPermitOPSQL   = "UPDATE mcn_up SET permission = ?, up_auth_link = ? WHERE sign_id = ? AND mcn_mid = ? AND up_mid = ?"
	_upPermitApplyOPSQL = "UPDATE mcn_up_permission_apply SET state = ?, reject_reason = ?, reject_time = ?, admin_id = ?, admin_name = ? WHERE id = ?"
	_selMcnSignsSQL     = `SELECT id,mcn_mid,company_name,company_license_id,company_license_link,contract_link,contact_name,
	contact_title,contact_idcard,contact_phone,begin_date,end_date,reject_reason,reject_time,state,ctime,mtime,permission FROM mcn_sign WHERE state = ? ORDER BY mtime DESC limit ?,?`
	_selMcnSignSQL = `SELECT id,mcn_mid,company_name,company_license_id,company_license_link,contract_link,contact_name,
	contact_title,contact_idcard,contact_phone,begin_date,end_date,reject_reason,reject_time,state,ctime,mtime,permission FROM mcn_sign WHERE id = ?`
	_selMcnSignPayMapSQL      = "SELECT id,sign_id,due_date,pay_value,state FROM mcn_sign_pay WHERE sign_id IN (%s) AND state IN (0,1)"
	_selMcnUpsSQL             = "SELECT id,sign_id,mcn_mid,up_mid,begin_date,end_date,contract_link,up_auth_link,reject_reason,reject_time,state,ctime,mtime,up_type,site_link,confirm_time,publication_price,permission FROM mcn_up WHERE state = ? ORDER BY mtime DESC limit ?,?"
	_selMcnUpSQL              = "SELECT id,sign_id,mcn_mid,up_mid,begin_date,end_date,contract_link,up_auth_link,reject_reason,reject_time,state,ctime,mtime,up_type,site_link,confirm_time,publication_price,permission FROM mcn_up WHERE id = ?"
	_selMcnSignCountUQTimeSQL = "SELECT COUNT(1) FROM mcn_sign WHERE mcn_mid = ? AND state IN (0,1,2,10,13,15) AND (end_date >= ? OR end_date >= ?) AND begin_date <= ?"
	_mcnSignTotalSQL          = "SELECT COUNT(1) FROM mcn_sign WHERE state = ?"
	_mcnUpTotalSQL            = "SELECT COUNT(1) FROM mcn_up WHERE state = ?"
	_mcnSignNoStateCountSQL   = "SELECT COUNT(1) FROM mcn_sign WHERE mcn_mid = ? AND state IN (0,1,2,10,13,15)"
	_mcnUpPermitReviewsSQL    = "SELECT id,mcn_mid,up_mid,sign_id,new_permission,old_permission,reject_reason,reject_time,state,ctime,mtime,admin_id,admin_name,up_auth_link FROM mcn_up_permission_apply WHERE state = ? ORDER BY mtime DESC limit ?,?"
	_mcnUpPermitReviewSQL     = "SELECT id,mcn_mid,up_mid,sign_id,new_permission,old_permission,reject_reason,reject_time,state,ctime,mtime,admin_id,admin_name,up_auth_link FROM mcn_up_permission_apply WHERE id = ?"
	_mcnUpPermitTotalSQL      = "SELECT COUNT(1) FROM mcn_up_permission_apply WHERE state = ?"
)

// TxAddMcnSignEntry .
func (d *Dao) TxAddMcnSignEntry(tx *xsql.Tx, mcnMid int64, beginDate, endDate string, permission uint32) (lastID int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_inMcnSignEntrySQL, mcnMid, beginDate, endDate, permission); err != nil {
		return lastID, err
	}
	if lastID, err = res.LastInsertId(); err != nil {
		return lastID, errors.Errorf("res.LastInsertId(%d,%s,%s,%d) error(%+v)", mcnMid, beginDate, endDate, permission, err)
	}
	return lastID, nil
}

// TxAddMcnSignPay .
func (d *Dao) TxAddMcnSignPay(tx *xsql.Tx, mcnMid, signID, payValue int64, dueDate, note string) (rows int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_inMcnSignPaySQL, mcnMid, signID, dueDate, payValue, note); err != nil {
		return rows, err
	}
	return res.RowsAffected()
}

// UpMcnSignOP .
func (d *Dao) UpMcnSignOP(c context.Context, signID int64, state int8, rejectTime xtime.Time, rejectReason string) (rows int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _upMcnSignOPSQL, state, rejectTime, rejectReason, signID); err != nil {
		return rows, err
	}
	return res.RowsAffected()
}

// UpMcnUpOP .
func (d *Dao) UpMcnUpOP(c context.Context, signUpID int64, state int8, rejectTime xtime.Time, rejectReason string) (rows int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _upMcnUpOPSQL, state, rejectTime, rejectReason, time.Now(), signUpID); err != nil {
		return rows, err
	}
	return res.RowsAffected()
}

// TxMcnUpPermitOP .
func (d *Dao) TxMcnUpPermitOP(tx *xsql.Tx, signID, mcnMid, upMid int64, permission uint32, upAuthLink string) (rows int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_mcnUpPermitOPSQL, permission, upAuthLink, signID, mcnMid, upMid); err != nil {
		return rows, err
	}
	return res.RowsAffected()
}

// TxUpPermitApplyOP .
func (d *Dao) TxUpPermitApplyOP(tx *xsql.Tx, arg *model.McnUpPermissionApply) (rows int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_upPermitApplyOPSQL, arg.State, arg.RejectReason, arg.RejectTime, arg.AdminID, arg.AdminName, arg.ID); err != nil {
		return rows, err
	}
	return res.RowsAffected()
}

// McnSigns .
func (d *Dao) McnSigns(c context.Context, arg *model.MCNSignStateReq) (ms []*model.MCNSignInfoReply, err error) {
	var rows *xsql.Rows
	limit, offset := arg.PageArg.CheckPageValidation()
	if rows, err = d.db.Query(c, _selMcnSignsSQL, arg.State, offset, limit); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := new(model.MCNSignInfoReply)
		if err = rows.Scan(&m.SignID, &m.McnMid, &m.CompanyName, &m.CompanyLicenseID, &m.CompanyLicenseLink, &m.ContractLink, &m.ContactName,
			&m.ContactTitle, &m.ContactIdcard, &m.ContactPhone, &m.BeginDate, &m.EndDate, &m.RejectReason, &m.RejectTime, &m.State, &m.Ctime, &m.Mtime, &m.Permission); err != nil {
			return
		}
		ms = append(ms, m)
	}
	err = rows.Err()
	return
}

// McnSign .
func (d *Dao) McnSign(c context.Context, signID int64) (m *model.MCNSignInfoReply, err error) {
	row := d.db.QueryRow(c, _selMcnSignSQL, signID)
	m = new(model.MCNSignInfoReply)
	if err = row.Scan(&m.SignID, &m.McnMid, &m.CompanyName, &m.CompanyLicenseID, &m.CompanyLicenseLink, &m.ContractLink, &m.ContactName,
		&m.ContactTitle, &m.ContactIdcard, &m.ContactPhone, &m.BeginDate, &m.EndDate, &m.RejectReason, &m.RejectTime, &m.State, &m.Ctime, &m.Mtime, &m.Permission); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			m = nil
			return
		}
	}
	return
}

// McnSignPayMap .
func (d *Dao) McnSignPayMap(c context.Context, signIDs []int64) (sm map[int64][]*model.SignPayInfoReply, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_selMcnSignPayMapSQL, xstr.JoinInts(signIDs))); err != nil {
		return
	}
	defer rows.Close()
	sm = make(map[int64][]*model.SignPayInfoReply)
	for rows.Next() {
		s := new(model.SignPayInfoReply)
		if err = rows.Scan(&s.SignPayID, &s.SignID, &s.DueDate, &s.PayValue, &s.State); err != nil {
			return
		}
		sm[s.SignID] = append(sm[s.SignID], s)
	}
	err = rows.Err()
	return
}

// McnUps .
func (d *Dao) McnUps(c context.Context, arg *model.MCNUPStateReq) (ups []*model.MCNUPInfoReply, err error) {
	var rows *xsql.Rows
	limit, offset := arg.PageArg.CheckPageValidation()
	if rows, err = d.db.Query(c, _selMcnUpsSQL, arg.State, offset, limit); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		up := new(model.MCNUPInfoReply)
		if err = rows.Scan(&up.SignUpID, &up.SignID, &up.McnMid, &up.UpMid, &up.BeginDate, &up.EndDate, &up.ContractLink,
			&up.UpAuthLink, &up.RejectReason, &up.RejectTime, &up.State, &up.Ctime, &up.Mtime, &up.UPType, &up.SiteLink,
			&up.ConfirmTime, &up.PubPrice, &up.Permission); err != nil {
			return
		}
		ups = append(ups, up)
	}
	err = rows.Err()
	return
}

// McnUp .
func (d *Dao) McnUp(c context.Context, signUpID int64) (up *model.MCNUPInfoReply, err error) {
	row := d.db.QueryRow(c, _selMcnUpSQL, signUpID)
	up = new(model.MCNUPInfoReply)
	if err = row.Scan(&up.SignUpID, &up.SignID, &up.McnMid, &up.UpMid, &up.BeginDate, &up.EndDate, &up.ContractLink,
		&up.UpAuthLink, &up.RejectReason, &up.RejectTime, &up.State, &up.Ctime, &up.Mtime, &up.UPType, &up.SiteLink,
		&up.ConfirmTime, &up.PubPrice, &up.Permission); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			up = nil
			return
		}
	}
	return
}

// McnSignCountUQTime .
func (d *Dao) McnSignCountUQTime(c context.Context, mcnMid int64, stime, etime xtime.Time) (count int64, err error) {
	row := d.db.QueryRow(c, _selMcnSignCountUQTimeSQL, mcnMid, stime, etime, etime)
	if err = row.Scan(&count); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			return
		}
	}
	return
}

// McnSignTotal .
func (d *Dao) McnSignTotal(c context.Context, arg *model.MCNSignStateReq) (count int64, err error) {
	row := d.db.QueryRow(c, _mcnSignTotalSQL, arg.State)
	if err = row.Scan(&count); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			return
		}
	}
	return
}

// McnUpTotal .
func (d *Dao) McnUpTotal(c context.Context, arg *model.MCNUPStateReq) (count int64, err error) {
	row := d.db.QueryRow(c, _mcnUpTotalSQL, arg.State)
	if err = row.Scan(&count); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			return
		}
	}
	return
}

// McnSignNoOKState .
func (d *Dao) McnSignNoOKState(c context.Context, mcnMid int64) (count int64, err error) {
	row := d.db.QueryRow(c, _mcnSignNoStateCountSQL, mcnMid)
	if err = row.Scan(&count); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			return
		}
	}
	return
}

// McnUpPermits .
func (d *Dao) McnUpPermits(c context.Context, arg *model.MCNUPPermitStateReq) (ms []*model.McnUpPermissionApply, err error) {
	var rows *xsql.Rows
	limit, offset := arg.PageArg.CheckPageValidation()
	if rows, err = d.db.Query(c, _mcnUpPermitReviewsSQL, arg.State, offset, limit); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := new(model.McnUpPermissionApply)
		if err = rows.Scan(&m.ID, &m.McnMid, &m.UpMid, &m.SignID, &m.NewPermission, &m.OldPermission, &m.RejectReason,
			&m.RejectTime, &m.State, &m.Ctime, &m.Mtime, &m.AdminID, &m.AdminName, &m.UpAuthLink); err != nil {
			return
		}
		ms = append(ms, m)
	}
	err = rows.Err()
	return
}

// McnUpPermit .
func (d *Dao) McnUpPermit(c context.Context, id int64) (m *model.McnUpPermissionApply, err error) {
	row := d.db.QueryRow(c, _mcnUpPermitReviewSQL, id)
	m = new(model.McnUpPermissionApply)
	if err = row.Scan(&m.ID, &m.McnMid, &m.UpMid, &m.SignID, &m.NewPermission, &m.OldPermission, &m.RejectReason,
		&m.RejectTime, &m.State, &m.Ctime, &m.Mtime, &m.AdminID, &m.AdminName, &m.UpAuthLink); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			return
		}
	}
	return
}

// McnUpPermitTotal .
func (d *Dao) McnUpPermitTotal(c context.Context, arg *model.MCNUPPermitStateReq) (count int64, err error) {
	row := d.db.QueryRow(c, _mcnUpPermitTotalSQL, arg.State)
	if err = row.Scan(&count); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			return
		}
	}
	return
}
