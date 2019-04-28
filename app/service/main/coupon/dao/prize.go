package dao

import (
	"context"
	xsql "database/sql"

	"go-common/app/service/main/coupon/model"
	"go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_userCardSQL       = "SELECT card_type,state,batch_token,coupon_token,act_id FROM coupon_user_card WHERE mid=? AND act_id=? AND card_type=?"
	_userCardsSQL      = "SELECT card_type,state,batch_token,coupon_token,act_id FROM coupon_user_card WHERE mid=? AND act_id=?"
	_addUserCardSQL    = "INSERT INTO coupon_user_card(mid,card_type,state,batch_token,coupon_token,act_id) VALUES (?,?,?,?,?,?)"
	_updateUserCardSQL = "UPDATE coupon_user_card SET state=? WHERE mid=? AND coupon_token=?"
)

// UserCard .
func (d *Dao) UserCard(c context.Context, mid, actID int64, cardType int8) (r *model.CouponUserCard, err error) {
	var row *sql.Row
	r = &model.CouponUserCard{}
	row = d.db.QueryRow(c, _userCardSQL, mid, actID, cardType)
	if err = row.Scan(&r.CardType, &r.State, &r.BatchToken, &r.CouponToken, &r.ActID); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
			return
		}
		err = errors.WithStack(err)
		return
	}
	return
}

// UserCards .
func (d *Dao) UserCards(c context.Context, mid, actID int64) (res []*model.CouponUserCard, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _userCardsSQL, mid, actID); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.CouponUserCard{}
		if err = rows.Scan(&r.CardType, &r.State, &r.BatchToken, &r.CouponToken, &r.ActID); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// AddUserCard .
func (d *Dao) AddUserCard(c context.Context, tx *sql.Tx, uc *model.CouponUserCard) (a int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_addUserCardSQL, uc.MID, uc.CardType, uc.State, uc.BatchToken, uc.CouponToken, uc.ActID); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// UpdateUserCard .
func (d *Dao) UpdateUserCard(c context.Context, mid int64, state int8, couponToken string) (a int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _updateUserCardSQL, state, mid, couponToken); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}
