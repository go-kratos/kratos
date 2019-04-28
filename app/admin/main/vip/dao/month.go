package dao

import (
	"context"
	xsql "database/sql"

	"go-common/app/admin/main/vip/model"
	"go-common/library/database/sql"

	"github.com/pkg/errors"
)

// GetMonth .
func (d *Dao) GetMonth(c context.Context, id int64) (m *model.VipMonth, err error) {
	res := d.db.QueryRow(c, _getMonth, id)
	m = new(model.VipMonth)
	if err = res.Scan(&m.ID, &m.Month, &m.MonthType, &m.Operator, &m.Status, &m.Mtime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			m = nil
		} else {
			err = errors.WithStack(err)
		}
	}
	return
}

// MonthList get all month.
func (d *Dao) MonthList(c context.Context) (res []*model.VipMonth, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _allMonth); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.VipMonth{}
		if err = rows.Scan(&r.ID, &r.Month, &r.MonthType, &r.Operator, &r.Status, &r.Mtime); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}

// MonthEdit .
func (d *Dao) MonthEdit(c context.Context, id int64, status int8, op string) (eff int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _updateMonthStatus, status, op, id); err != nil {
		err = errors.WithStack(err)
		return
	}
	if eff, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// GetPrice .
func (d *Dao) GetPrice(c context.Context, id int64) (r *model.VipMonthPrice, err error) {
	row := d.db.QueryRow(c, _monthPriceSQL, id)
	r = new(model.VipMonthPrice)
	if err = row.Scan(&r.ID, &r.MonthID, &r.MonthType, &r.Money, &r.Selected, &r.FirstDiscountMoney, &r.DiscountMoney, &r.StartTime, &r.EndTime, &r.Remark, &r.Operator); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
		} else {
			err = errors.WithStack(err)
		}
	}
	return
}

// PriceList .
func (d *Dao) PriceList(c context.Context, mID int64) (res []*model.VipMonthPrice, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _allMonthPrice, mID); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.VipMonthPrice{}
		if err = rows.Scan(&r.ID, &r.MonthID, &r.MonthType, &r.Money, &r.Selected, &r.FirstDiscountMoney, &r.DiscountMoney, &r.StartTime, &r.EndTime, &r.Remark, &r.Operator); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}

// PriceAdd .
func (d *Dao) PriceAdd(c context.Context, mp *model.VipMonthPrice) (eff int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _addMonthPrice, mp.MonthID, mp.MonthType, mp.Money, mp.FirstDiscountMoney, mp.DiscountMoney, mp.StartTime, mp.EndTime, mp.Remark, mp.Operator); err != nil {
		err = errors.WithStack(err)
		return
	}
	if eff, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// PriceEdit .
func (d *Dao) PriceEdit(c context.Context, mp *model.VipMonthPrice) (eff int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _editMonthPrice, mp.MonthType, mp.Money, mp.FirstDiscountMoney, mp.DiscountMoney, mp.StartTime, mp.EndTime, mp.Remark, mp.Operator, mp.ID); err != nil {
		err = errors.WithStack(err)
		return
	}
	if eff, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}
