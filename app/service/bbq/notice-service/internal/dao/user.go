package dao

import (
	"context"
)

const (
	_queryUserFansNum = "select `fan_total` from `user_statistics` where mid = ?;"
)

// FetchUserFansNum .
func (d *Dao) FetchUserFansNum(c context.Context, mid int64) (num int, err error) {
	row := d.db.QueryRow(c, _queryUserFansNum, mid)
	err = row.Scan(&num)
	return
}
