package dao

import (
	"context"
)

const (
	_selWhite = "SELECT mid FROM filter_ai_white WHERE state=0"
)

// AiWhite get AI white mids.
func (d *Dao) AiWhite(c context.Context) (res map[int64]int64, err error) {
	res = make(map[int64]int64)
	rows, err := d.mysql.Query(c, _selWhite)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			mid int64
		)
		if err = rows.Scan(&mid); err != nil {
			return
		}
		res[mid] = mid
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}
