package like

import (
	"context"
	"database/sql"
	"fmt"
	"go-common/app/admin/main/activity/model"

	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_likeActSumSQL  = "SELECT SUM(`action`) AS `like`,lid FROM like_action WHERE lid IN(%s) GROUP BY lid"
	_likeActListSQL = "SELECT id,mid FROM like_action WHERE lid = ? AND id > ? ORDER BY id LIMIT ?"
)

// BatchLikeActSum .
func (d *Dao) BatchLikeActSum(c context.Context, lids []int64) (res map[int64]int64, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_likeActSumSQL, xstr.JoinInts(lids)))
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			err = errors.Wrap(err, "d.db.Query()")
		}
		return
	}
	defer rows.Close()
	res = make(map[int64]int64)
	for rows.Next() {
		like := sql.NullInt64{}
		lid := sql.NullInt64{}
		if err = rows.Scan(&like, &lid); err != nil {
			err = errors.Wrap(err, "rows.Scan()")
			return
		}
		res[lid.Int64] = like.Int64
	}
	if err = rows.Err(); err != nil {
		err = errors.Wrap(err, "rows.Err()")
	}
	return
}

// LikeActList .
func (d *Dao) LikeActList(c context.Context, lid, minID, limit int64) (res []*model.LikeAction, err error) {
	rows, err := d.db.Query(c, _likeActListSQL, lid, minID, limit)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			err = errors.Wrap(err, "d.db.Query()")
		}
		return
	}
	defer rows.Close()
	for rows.Next() {
		action := new(model.LikeAction)
		if err = rows.Scan(&action.ID, &action.Mid); err != nil {
			err = errors.Wrap(err, "rows.Scan()")
			return
		}
		res = append(res, action)
	}
	if err = rows.Err(); err != nil {
		err = errors.Wrap(err, "rows.Err()")
	}
	return
}
