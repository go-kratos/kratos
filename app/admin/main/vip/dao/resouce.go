package dao

import (
	"context"
	xsql "database/sql"
	"fmt"
	"strconv"
	"strings"

	"go-common/app/admin/main/vip/model"
	"go-common/library/database/sql"

	"github.com/pkg/errors"
)

// SelBatchCodeCount .
func (d *Dao) SelBatchCodeCount(c context.Context, arg *model.ArgBatchCode) (n int64, err error) {
	autoSQLStr := d.batchCodeAutoArg(arg)
	row := d.db.QueryRow(c, _selBatchCodeCountSQL+autoSQLStr)
	if err = row.Scan(&n); err != nil {
		if err == sql.ErrNoRows {
			n = 0
			err = nil
			return
		}
		err = errors.WithStack(err)
		return
	}
	return
}

func (d *Dao) selBatchCodeIDs(c context.Context, poolID, businessID, batchID int64) (ids []int64, err error) {
	var rows *sql.Rows
	autoSQL := d.batchCodeAutoArg(&model.ArgBatchCode{
		BusinessID: businessID,
		PoolID:     poolID,
		ID:         batchID,
	})
	if rows, err = d.db.Query(c, _selBatchCodeIDSQL+autoSQL); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("db_query")
		return
	}
	for rows.Next() {
		var r int64
		if err = rows.Scan(&r); err != nil {
			err = errors.WithStack(err)
			ids = nil
			d.errProm.Incr("db_scan")
			return
		}
		ids = append(ids, r)
	}
	return
}

func (d *Dao) codeAutoArgSQL(arg *model.ArgCode) string {
	autoSQLStr := ""
	if len(arg.BatchCodeIDs) > 0 {
		idStr := ""
		for _, v := range arg.BatchCodeIDs {
			idStr += "," + strconv.Itoa(int(v))
		}
		autoSQLStr += fmt.Sprintf(" AND batch_code_id in(%s) ", idStr[1:])
	}
	if arg.Status != 0 {
		autoSQLStr += fmt.Sprintf(" AND status = %v ", arg.Status)
	}
	if arg.UseStartTime > 0 && arg.UseEndTime > 0 {
		autoSQLStr += fmt.Sprintf(" AND use_time >='%v' AND use_time <= '%v'", arg.UseStartTime.Time().Format("2006-01-02 15:04:05"), arg.UseEndTime.Time().Format("2006-01-02 15:04:05"))
	}
	if len(arg.Code) > 0 {
		autoSQLStr += fmt.Sprintf(" AND code = '%v' ", arg.Code)
	}
	if arg.Mid != 0 {
		autoSQLStr += fmt.Sprintf(" AND mid = %v ", arg.Mid)
	}

	return autoSQLStr
}

// SelCode .
func (d *Dao) SelCode(c context.Context, arg *model.ArgCode, cursor int64, ps int) (res []*model.ResourceCode, err error) {
	var batchIDs []int64
	if batchIDs, err = d.selBatchCodeIDs(c, arg.PoolID, arg.BusinessID, arg.BatchCodeID); err != nil {
		err = errors.WithStack(err)
		return
	}
	arg.BatchCodeIDs = batchIDs
	argSQL := d.codeAutoArgSQL(arg)
	if ps <= 0 || ps > 2000 {
		ps = _defps
	}
	var rows *sql.Rows
	argSQL += fmt.Sprintf(" AND id>%v LIMIT %v", cursor, ps)
	if rows, err = d.db.Query(c, _selCodeSQL+argSQL); err != nil {
		err = errors.WithStack(err)
		return
	}
	for rows.Next() {
		r := new(model.ResourceCode)
		if err = rows.Scan(&r.ID, &r.BatchCodeID, &r.Status, &r.Code, &r.Mid, &r.UseTime, &r.Ctime); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return

}

// batchCodeAutoArg .
func (d *Dao) batchCodeAutoArg(arg *model.ArgBatchCode) string {
	autoSQLStr := ""
	if arg.BusinessID > 0 {
		autoSQLStr += fmt.Sprintf(" AND business_id=%v ", arg.BusinessID)
	}
	if arg.PoolID > 0 {
		autoSQLStr += fmt.Sprintf(" AND pool_id = %v ", arg.PoolID)
	}
	if arg.ID > 0 {
		autoSQLStr += fmt.Sprintf(" AND id = %v ", arg.ID)
	}
	if len(arg.Name) > 0 {
		autoSQLStr += " AND batch_name like '%" + arg.Name + "%'"
	}
	if arg.Status > 0 {
		autoSQLStr += fmt.Sprintf(" AND status = %v", arg.Status)
	}

	if arg.StartTime > 0 && arg.EndTime > 0 {
		autoSQLStr += fmt.Sprintf(" AND ctime >='%v' AND ctime <= '%v'", arg.StartTime.Time().Format("2006-01-02 15:04:05"), arg.EndTime.Time().Format("2006-01-02 15:04:05"))
	}
	return autoSQLStr
}

// SelBatchCodes .
func (d *Dao) SelBatchCodes(c context.Context, batchIds []int64) (res []*model.BatchCode, err error) {
	var (
		rows *sql.Rows
		ids  []string
	)
	if len(batchIds) <= 0 {
		return
	}
	for _, v := range batchIds {
		ids = append(ids, fmt.Sprintf("%v", v))
	}

	if rows, err = d.db.Query(c, fmt.Sprintf(_selBatchCodesSQL, strings.Join(ids, ","))); err != nil {
		err = errors.WithStack(err)
		return
	}
	for rows.Next() {
		r := new(model.BatchCode)
		if err = rows.Scan(&r.ID, &r.BusinessID, &r.PoolID, &r.Status, &r.Type, &r.LimitDay, &r.MaxCount, &r.BatchName, &r.Reason, &r.Unit, &r.Count, &r.SurplusCount, &r.Price, &r.StartTime, &r.EndTime, &r.Contacts, &r.ContactsNumber, &r.Ctime); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}

// SelBatchCode .
func (d *Dao) SelBatchCode(c context.Context, arg *model.ArgBatchCode, pn, ps int) (res []*model.BatchCode, err error) {
	var rows *sql.Rows
	autoSQLStr := d.batchCodeAutoArg(arg)
	if pn <= 0 {
		pn = _defpn
	}
	if ps <= 0 || ps > _maxps {
		ps = _defps
	}

	autoSQLStr += fmt.Sprintf(" LIMIT %v,%v", (pn-1)*ps, ps)

	if rows, err = d.db.Query(c, _selBatchCodeSQL+autoSQLStr); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.BatchCode)
		if err = rows.Scan(&r.ID, &r.BusinessID, &r.PoolID, &r.Status, &r.Type, &r.LimitDay, &r.MaxCount, &r.BatchName, &r.Reason, &r.Unit, &r.Count, &r.SurplusCount, &r.Price, &r.StartTime, &r.EndTime, &r.Contacts, &r.ContactsNumber, &r.Ctime); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}

	return
}

// SelBatchCodeName .
func (d *Dao) SelBatchCodeName(c context.Context, name string) (r *model.BatchCode, err error) {
	row := d.db.QueryRow(c, _selBatchCodeByNameSQL, name)
	r = new(model.BatchCode)
	if err = row.Scan(&r.ID, &r.BusinessID, &r.PoolID, &r.Status, &r.Type, &r.LimitDay, &r.MaxCount, &r.BatchName, &r.Reason, &r.Unit, &r.Count, &r.SurplusCount, &r.Price, &r.StartTime, &r.EndTime, &r.Contacts, &r.ContactsNumber, &r.Ctime); err != nil {
		if err == sql.ErrNoRows {
			r = nil
			err = nil
			fmt.Printf("this is  %+v", err)
			return
		}
		err = errors.WithStack(err)
		return
	}
	return
}

// SelBatchCodeID .
func (d *Dao) SelBatchCodeID(c context.Context, batchCodeID int64) (r *model.BatchCode, err error) {
	row := d.db.QueryRow(c, _selBatchCodeByIDSQL, batchCodeID)
	r = new(model.BatchCode)
	if err = row.Scan(&r.ID, &r.BusinessID, &r.PoolID, &r.Status, &r.Type, &r.LimitDay, &r.MaxCount, &r.BatchName, &r.Reason, &r.Unit, &r.Count, &r.SurplusCount, &r.Price, &r.StartTime, &r.EndTime, &r.Contacts, &r.ContactsNumber, &r.Ctime); err != nil {
		if err == sql.ErrNoRows {
			r = nil
			err = nil
			return
		}
		err = errors.WithStack(err)
		return
	}
	return
}

// SelCodeID .
func (d *Dao) SelCodeID(c context.Context, codeID int64) (r *model.ResourceCode, err error) {

	row := d.db.QueryRow(c, _selCodeByIDSQL, codeID)
	r = new(model.ResourceCode)
	if err = row.Scan(&r.ID, &r.BatchCodeID, &r.Status, &r.Code, &r.Mid, &r.UseTime, &r.Ctime); err != nil {
		if err == sql.ErrNoRows {
			r = nil
			err = nil
		}
		err = errors.WithStack(err)
		return
	}
	return
}

// TxAddBatchCode .
func (d *Dao) TxAddBatchCode(tx *sql.Tx, bc *model.BatchCode) (ID int64, err error) {
	var result xsql.Result
	if result, err = tx.Exec(_addBatchCodeSQL, bc.BusinessID, bc.PoolID, bc.Status, bc.BatchName, bc.Reason, bc.Unit, bc.Count, bc.SurplusCount, bc.Price, bc.StartTime, bc.EndTime, bc.Contacts, bc.ContactsNumber, bc.Type, bc.LimitDay, bc.MaxCount, bc.Operator); err != nil {
		err = errors.WithStack(err)
		return
	}
	if ID, err = result.LastInsertId(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// UpdateBatchCode .
func (d *Dao) UpdateBatchCode(c context.Context, bc *model.BatchCode) (eff int64, err error) {
	var result xsql.Result
	if result, err = d.db.Exec(c, _updateBatchCodeSQL, &bc.Status, &bc.BatchName, &bc.Reason, &bc.Price, &bc.Contacts, &bc.ContactsNumber, &bc.Type, &bc.LimitDay, &bc.MaxCount, &bc.Operator, &bc.ID); err != nil {
		err = errors.WithStack(err)
		return
	}
	if eff, err = result.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// UpdateCode .
func (d *Dao) UpdateCode(c context.Context, codeID int64, status int8) (eff int64, err error) {
	var result xsql.Result
	if result, err = d.db.Exec(c, _updateCodeSQL, status, codeID); err != nil {
		err = errors.WithStack(err)
		return
	}
	if eff, err = result.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// BatchAddCode .
func (d *Dao) BatchAddCode(tx *sql.Tx, codes []*model.ResourceCode) (err error) {
	values := make([]string, 0)
	for _, v := range codes {
		s := fmt.Sprintf("('%v','%v','%v')", v.BatchCodeID, v.Status, v.Code)
		values = append(values, s)
	}
	valuesStr := strings.Join(values, ",")
	if _, err = tx.Exec(_batchAddCodeSQL + valuesStr); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}
