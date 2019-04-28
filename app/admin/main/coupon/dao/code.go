package dao

import (
	"bytes"
	"context"
	"fmt"

	"go-common/app/admin/main/coupon/model"
	"go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_codeCountSQL    = "SELECT COUNT(1) FROM bilibili_coupon.coupon_code WHERE 1=1 %s"
	_codeListSQL     = "SELECT id,batch_token,state,code,mid,coupon_token,coupon_type,ver,ctime,mtime FROM coupon_code WHERE 1=1 %s %s"
	_codeBlockSQL    = "UPDATE coupon_code SET state = ?, ver = ver +1 WHERE id = ? AND ver =?;"
	_codeByIDSQL     = "SELECT id,batch_token,state,code,mid,coupon_token,coupon_type,ver,ctime,mtime FROM coupon_code WHERE id = ?"
	_batchAddCodeSQL = "INSERT IGNORE INTO coupon_code(batch_token,state,code,coupon_type)VALUES "
)

// CountCode coupon code count.
func (d *Dao) CountCode(c context.Context, a *model.ArgCouponCode) (count int64, err error) {
	sql := fmt.Sprintf(_codeCountSQL, whereSQL(a))
	if err = d.db.QueryRow(c, sql).Scan(&count); err != nil {
		err = errors.Wrapf(err, "dao code list")
	}
	return
}

// CodeList code list.
func (d *Dao) CodeList(c context.Context, a *model.ArgCouponCode) (res []*model.CouponCode, err error) {
	listSQL := fmt.Sprintf(_codeListSQL, whereSQL(a), pageSQL(a.Pn, a.Ps))
	var rows *sql.Rows
	if rows, err = d.db.Query(c, listSQL); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.CouponCode)
		if err = rows.Scan(&r.ID, &r.BatchToken, &r.State, &r.Code, &r.Mid, &r.CouponToken, &r.CouponType, &r.Ver, &r.Ctime, &r.Mtime); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}

		res = append(res, r)
	}
	err = rows.Err()
	return
}

// UpdateCodeBlock update code block.
func (d *Dao) UpdateCodeBlock(c context.Context, a *model.CouponCode) (err error) {
	if _, err = d.db.Exec(c, _codeBlockSQL, a.State, a.ID, a.Ver); err != nil {
		err = errors.Wrapf(err, "dao update code block(%+v)", a)
	}
	return
}

// CodeByID code by id.
func (d *Dao) CodeByID(c context.Context, id int64) (r *model.CouponCode, err error) {
	r = new(model.CouponCode)
	if err = d.db.QueryRow(c, _codeByIDSQL, id).
		Scan(&r.ID, &r.BatchToken, &r.State, &r.Code, &r.Mid, &r.CouponToken, &r.CouponType, &r.Ver, &r.Ctime, &r.Mtime); err != nil {
		if err == sql.ErrNoRows {
			r = nil
			err = nil
			return
		}
		err = errors.Wrapf(err, "dao query code by id")
	}
	return
}

// BatchAddCode batch add code.
func (d *Dao) BatchAddCode(c context.Context, cs []*model.CouponCode) (err error) {
	var (
		buf bytes.Buffer
		sql string
	)
	buf.WriteString(_batchAddCodeSQL)
	for _, v := range cs {
		buf.WriteString("('")
		buf.WriteString(v.BatchToken)
		buf.WriteString("',")
		buf.WriteString(fmt.Sprintf("%d", v.State))
		buf.WriteString(",'")
		buf.WriteString(v.Code)
		buf.WriteString("',")
		buf.WriteString(fmt.Sprintf("%d", v.CouponType))
		buf.WriteString("),")
	}
	sql = buf.String()
	if _, err = d.db.Exec(c, sql[0:len(sql)-1]); err != nil {
		err = errors.Wrapf(err, "dao insert codes")
	}
	return
}

func whereSQL(a *model.ArgCouponCode) (sql string) {
	if a == nil {
		return
	}
	if a.Mid > 0 {
		sql += " AND mid = " + fmt.Sprintf("%d", a.Mid)
	}
	if a.Code != "" {
		sql += " AND code = '" + a.Code + "'"
	}
	if a.BatchToken != "" {
		sql += " AND batch_token = '" + a.BatchToken + "'"
	}
	return sql
}

func pageSQL(pn, ps int) (sql string) {
	if pn <= 0 {
		pn = 1
	}
	if ps <= 0 {
		ps = 20
	}
	return " ORDER BY ID DESC LIMIT " + fmt.Sprintf("%d,%d", (pn-1)*ps, ps)
}
