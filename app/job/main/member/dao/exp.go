package dao

import (
	"context"
	"fmt"

	"go-common/app/job/main/member/model"
	"go-common/library/database/sql"
	"go-common/library/time"
)

const (
	_initExp       = "INSERT IGNORE INTO user_exp_%02d (mid) VALUES (?)"
	_updateExpAped = `UPDATE user_exp_%02d SET exp=exp+?,flag=flag|? WHERE mid=? AND flag&?=0`
	_updateExpFlag = `UPDATE user_exp_%02d SET addtime=?,flag=? WHERE mid=? AND addtime<?`
	_SelExp        = "SELECT mid,exp,flag,addtime,mtime FROM user_exp_%02d where mid=?"
)

// InitExp init user exp
func (d *Dao) InitExp(c context.Context, mid int64) (err error) {
	_, err = d.db.Exec(c, fmt.Sprintf(_initExp, hit(mid)), mid)
	return
}

// SelExp get new user exp by mid
func (d *Dao) SelExp(c context.Context, mid int64) (exp *model.NewExp, err error) {
	exp = &model.NewExp{}
	row := d.db.QueryRow(c, fmt.Sprintf(_SelExp, hit(mid)), mid)
	if err = row.Scan(&exp.Mid, &exp.Exp, &exp.Flag, &exp.Addtime, &exp.Mtime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}
	return
}

// UpdateExpAped update user exp append mode.
func (d *Dao) UpdateExpAped(c context.Context, mid, exp int64, flag int32) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_updateExpAped, hit(mid)), exp, flag, mid, flag)
	if err != nil {
		return
	}
	rows, _ = res.RowsAffected()
	return
}

// UpdateExpFlag update user exp.
func (d *Dao) UpdateExpFlag(c context.Context, mid int64, flag int32, addtime time.Time) (err error) {
	_, err = d.db.Exec(c, fmt.Sprintf(_updateExpFlag, hit(mid)), addtime, flag, mid, addtime)
	return
}
