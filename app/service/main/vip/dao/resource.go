package dao

import (
	"context"
	xsql "database/sql"
	"fmt"
	"strings"

	"go-common/app/service/main/vip/model"
	"go-common/library/database/sql"
	"go-common/library/time"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_defaultsize = 10000
)

//SelNewResourcePool select new resource pool by id
func (d *Dao) SelNewResourcePool(c context.Context, id int64) (r *model.VipResourcePool, err error) {
	var row = d.db.QueryRow(c, _selResourcePoolByIDSQL, id)
	r = new(model.VipResourcePool)
	if err = row.Scan(&r.ID, &r.PoolName, &r.BusinessID, &r.Reason, &r.CodeExpireTime, &r.StartTime, &r.EndTime, &r.Contacts, &r.ContactsNumber, &r.Ctime, &r.Mtime); err != nil {
		if err == sql.ErrNoRows {
			r = nil
			err = nil
		} else {
			err = errors.WithStack(err)
			d.errProm.Incr("row_scan_db")
			return
		}
	}
	return
}

//SelNewBusiness select newdb businessInfo by id.
func (d *Dao) SelNewBusiness(c context.Context, id int64) (r *model.VipBusinessInfo, err error) {
	var row = d.db.QueryRow(c, _selBusinessByIDSQL, id)
	r = new(model.VipBusinessInfo)
	if err = row.Scan(&r.ID, &r.BusinessName, &r.BusinessType, &r.Status, &r.AppKey, &r.Secret, &r.Contacts, &r.ContactsNumber, &r.Ctime, &r.Mtime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
		} else {
			err = errors.WithStack(err)
			d.errProm.Incr("row_scan_db")
		}
	}
	return
}

//SelNewBusinessByAppkey select newdb businessInfo by id.
func (d *Dao) SelNewBusinessByAppkey(c context.Context, appkey string) (r *model.VipBusinessInfo, err error) {
	var row = d.db.QueryRow(c, _selBusinessByAppkeySQL, appkey)
	r = new(model.VipBusinessInfo)
	if err = row.Scan(&r.ID, &r.BusinessName, &r.BusinessType, &r.Status, &r.AppKey, &r.Secret, &r.Contacts, &r.ContactsNumber, &r.Ctime, &r.Mtime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
		} else {
			err = errors.WithStack(err)
			d.errProm.Incr("row_scan_db")
		}
	}
	return
}

//SelCode sel code.
func (d *Dao) SelCode(c context.Context, codeStr string) (code *model.VipResourceCode, err error) {
	row := d.db.QueryRow(c, _selCodeSQL, codeStr)
	code = new(model.VipResourceCode)
	if err = row.Scan(&code.ID, &code.BatchCodeID, &code.Status, &code.Code, &code.Mid, &code.UseTime, &code.RelationID); err != nil {
		if sql.ErrNoRows == err {
			c = nil
			err = nil
			return
		}
		err = errors.WithStack(err)
		d.errProm.Incr("db_scan")
		return
	}
	return
}

//SelCodes sel codes
func (d *Dao) SelCodes(c context.Context, codes []string) (cs []*model.VipResourceCode, err error) {
	var (
		rows *sql.Rows
	)
	if len(codes) <= 0 {
		return
	}
	if rows, err = d.db.Query(c, fmt.Sprintf(_selCodesSQL, strings.Join(codes, "','"))); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		code := new(model.VipResourceCode)
		if err = rows.Scan(&code.ID, &code.BatchCodeID, &code.Status, &code.Code, &code.Mid, &code.UseTime, &code.RelationID); err != nil {
			if sql.ErrNoRows == err {
				c = nil
				err = nil
				return
			}
			err = errors.WithStack(err)
			d.errProm.Incr("db_scan")
			return
		}
		cs = append(cs, code)
	}
	err = rows.Err()
	return
}

// SelBatchCode set batch code.
func (d *Dao) SelBatchCode(c context.Context, id int64) (bc *model.VipResourceBatchCode, err error) {
	row := d.db.QueryRow(c, _selBatchCodeSQL, id)
	bc = new(model.VipResourceBatchCode)
	if err = row.Scan(&bc.ID, &bc.BusinessID, &bc.PoolID, &bc.Status, &bc.BatchName, &bc.Reason, &bc.Unit, &bc.Count, &bc.SurplusCount, &bc.Price, &bc.StartTime, &bc.EndTime, &bc.Type, &bc.LimitDay, &bc.MaxCount); err != nil {
		if sql.ErrNoRows == err {
			bc = nil
			err = nil
			return
		}
		err = errors.WithStack(err)
		d.errProm.Incr("db_scan")
		return
	}
	return
}

//SelBatchCodes sel batchcodes
func (d *Dao) SelBatchCodes(c context.Context, ids []int64) (bcs []*model.VipResourceBatchCode, err error) {
	var (
		rows *sql.Rows
	)
	if len(ids) <= 0 {
		return
	}
	if rows, err = d.db.Query(c, fmt.Sprintf(_selBatchCodesSQL, xstr.JoinInts(ids))); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		bc := new(model.VipResourceBatchCode)
		if err = rows.Scan(&bc.ID, &bc.BusinessID, &bc.PoolID, &bc.Status, &bc.BatchName, &bc.Reason, &bc.Unit, &bc.Count, &bc.SurplusCount, &bc.Price, &bc.StartTime, &bc.EndTime, &bc.Type, &bc.LimitDay, &bc.MaxCount); err != nil {
			if sql.ErrNoRows == err {
				bc = nil
				err = nil
				return
			}
			err = errors.WithStack(err)
			d.errProm.Incr("db_scan")
			return
		}
		bcs = append(bcs, bc)
	}
	err = rows.Err()
	return
}

// SelBatchCodesByBisID set batch codes by business id.
func (d *Dao) SelBatchCodesByBisID(c context.Context, bisID int64) (bcs []*model.VipResourceBatchCode, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _selBatchCodeByBisSQL, bisID); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		bc := new(model.VipResourceBatchCode)
		if err = rows.Scan(&bc.ID, &bc.BusinessID, &bc.PoolID, &bc.Status, &bc.BatchName, &bc.Reason, &bc.Unit, &bc.Count, &bc.SurplusCount, &bc.Price, &bc.StartTime, &bc.EndTime, &bc.Type, &bc.LimitDay, &bc.MaxCount); err != nil {
			if sql.ErrNoRows == err {
				bc = nil
				err = nil
				return
			}
			err = errors.WithStack(err)
			d.errProm.Incr("db_scan")
			return
		}
		bcs = append(bcs, bc)
	}
	err = rows.Err()
	return
}

// SelCodeOpened set code open.
func (d *Dao) SelCodeOpened(c context.Context, bisIDs []int64, arg *model.ArgCodeOpened) (cs []*model.CodeInfoResp, err error) {
	var rows *sql.Rows
	size := d.c.Property.CodeOpenedSearchSize
	if size == 0 {
		size = _defaultsize
	}
	if rows, err = d.db.Query(c, fmt.Sprintf(_selCodeOpenedSQL, xstr.JoinInts(bisIDs), arg.Cursor, arg.StartTime.Time().Format("2006-01-02 15:04:05"), arg.EndTime.Time().Format("2006-01-02 15:04:05"), size)); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		code := new(model.CodeInfoResp)
		if err = rows.Scan(&code.ID, &code.Code, &code.UserTime); err != nil {
			if sql.ErrNoRows == err {
				c = nil
				err = nil
				return
			}
			err = errors.WithStack(err)
			d.errProm.Incr("db_scan")
			return
		}
		cs = append(cs, code)
	}
	err = rows.Err()
	return
}

//TxUpdateCode tx update code.
func (d *Dao) TxUpdateCode(tx *sql.Tx, id int64, mid int64, useTime time.Time) (eff int64, err error) {
	var result xsql.Result
	if result, err = tx.Exec(_updateCodeSQL, mid, useTime, id); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("db_update")
		return
	}
	if eff, err = result.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("db_row_affected")
		return
	}
	return
}

//TxUpdateCodeStatus tx update code status.
func (d *Dao) TxUpdateCodeStatus(tx *sql.Tx, id int64, status int8) (eff int64, err error) {
	var result xsql.Result
	if result, err = tx.Exec(_updateCodeStatusSQL, status, id); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("db_update")
		return
	}
	if eff, err = result.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("db_row_affected")
		return
	}
	return
}

//TxUpdateBatchCode tx update batch code.
func (d *Dao) TxUpdateBatchCode(tx *sql.Tx, id int64, sc int32) (eff int64, err error) {
	var result xsql.Result
	if result, err = tx.Exec(_updateBatchCodeSQL, sc, id); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("db_batch_code")
		return
	}
	if eff, err = result.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("db_row_affected")
		return
	}
	return
}

//SelCodesByBMid sel codes by bmid
func (d *Dao) SelCodesByBMid(c context.Context, mid int64) (cs []string, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _selCodesByBmidSQL, mid); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("db_query")
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r string
		if err = rows.Scan(&r); err != nil {
			if err == sql.ErrNoRows {
				err = nil
				return
			}
			err = errors.WithStack(err)
			d.errProm.Incr("db_rows_scan")
			return
		}
		cs = append(cs, r)
	}
	err = rows.Err()
	return
}

//SelActives sel active data.
func (d *Dao) SelActives(c context.Context, relations []string) (rep []*model.VipActiveShow, err error) {
	if len(relations) <= 0 {
		return
	}
	var rows *sql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_selActiveSQL, strings.Join(relations, "','"))); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.VipActiveShow)
		if err = rows.Scan(&r.ID, &r.Type, &r.ProductName, &r.ProductPic, &r.RelationID, &r.BusID, &r.ProductDetail, &r.UseType); err != nil {
			err = errors.WithStack(err)
			return
		}
		rep = append(rep, r)
	}
	err = rows.Err()
	return
}

//SelBatchCount sel batch count
func (d *Dao) SelBatchCount(c context.Context, batchCodeID, mid int64) (count int64, err error) {
	row := d.db.QueryRow(c, _selBatchCountSQL, mid, batchCodeID)
	if err = row.Scan(&count); err != nil {
		err = errors.WithStack(err)
	}
	return
}
