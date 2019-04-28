package dao

import (
	"context"

	"go-common/app/admin/main/tag/model"
	"go-common/library/database/sql"
)

const (
	_inBussinessSQL      = "INSERT INTO business (type, name, appkey, remark, alias) VALUES (?,?,?,?,?)"
	_upBussinessSQL      = "UPDATE business SET name=?, appkey=?, remark=?, alias=? WHERE type=?"
	_upBussinessSteteSQL = "UPDATE business SET state=? WHERE type=?"
	_selBussinessSQL     = "SELECT type, name, appkey, remark, alias FROM business WHERE state=?"
	_selOneBusinessSQL   = "SELECT type, name, appkey, remark, alias FROM business WHERE type=?"
)

// InBusiness insert a business record
func (dao *Dao) InBusiness(c context.Context, tp int32, name, appkey, remark, alias string) (id int64, err error) {
	res, err := dao.db.Exec(c, _inBussinessSQL, tp, name, appkey, remark, alias)
	if err != nil {
		return
	}
	return res.LastInsertId()
}

// UpBusiness update business by type
func (dao *Dao) UpBusiness(c context.Context, name, appkey, remark, alias string, tp int32) (id int64, err error) {
	res, err := dao.db.Exec(c, _upBussinessSQL, name, appkey, remark, alias, tp)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// UpBusinessState logical delete a business record by set state to StateDelete
func (dao *Dao) UpBusinessState(c context.Context, state, tp int32) (id int64, err error) {
	res, err := dao.db.Exec(c, _upBussinessSteteSQL, state, tp)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// Business return one business instance by type
func (dao *Dao) Business(c context.Context, tp int32) (business *model.Business, err error) {
	row := dao.db.QueryRow(c, _selOneBusinessSQL, tp)
	business = new(model.Business)
	err = row.Scan(&business.Type, &business.Name, &business.Appkey, &business.Remark, &business.Alias)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
			business = nil
		}
	}
	return
}

// ListBusiness Gets gets all business records
func (dao *Dao) ListBusiness(c context.Context, state int32) (business []*model.Business, err error) {
	rows, err := dao.db.Query(c, _selBussinessSQL, state)
	if err != nil {
		return
	}
	defer rows.Close()
	business = make([]*model.Business, 0)
	for rows.Next() {
		b := new(model.Business)
		if err = rows.Scan(&b.Type, &b.Name, &b.Appkey, &b.Remark, &b.Alias); err != nil {
			return
		}
		business = append(business, b)
	}
	err = rows.Err()
	return
}
