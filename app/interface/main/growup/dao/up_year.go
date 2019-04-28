package dao

import (
	"context"
	"database/sql"
)

const (
	_upTagYearSQL = "SELECT tag1, tag2, tag3, tag4, tag5, tag6, total_income FROM up_tag_year WHERE mid = ?"
)

// GetUpYearTag get up year tag income
func (d *Dao) GetUpYearTag(c context.Context, mid int64) (tags []int64, income int64, err error) {
	tags = make([]int64, 6)
	err = d.db.QueryRow(c, _upTagYearSQL, mid).Scan(&tags[0], &tags[1], &tags[2], &tags[3], &tags[4], &tags[5], &income)
	if err == sql.ErrNoRows {
		err = nil
	}
	return
}
