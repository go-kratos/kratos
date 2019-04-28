package dao

import (
	"context"
	"fmt"
	"strings"

	"go-common/library/xstr"
)

const (
	_addAuthoritySQL = "INSERT INTO up_white_list(mid) VALUES %s ON DUPLICATE KEY UPDATE is_deleted=0"
	_rmAuthoritySQL  = "UPDATE up_white_list SET is_deleted = 1 WHERE mid IN (%s)"
)

// AddAuthority .
func (d *Dao) AddAuthority(c context.Context, mids []int64) (int64, error) {
	if len(mids) == 0 {
		return 0, nil
	}
	s := make([]string, 0)
	for _, v := range mids {
		s = append(s, fmt.Sprintf("(%d)", v))
	}
	res, err := d.db.Exec(c, fmt.Sprintf(_addAuthoritySQL, strings.Join(s, ",")))
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// RmAuthority .
func (d *Dao) RmAuthority(c context.Context, mids []int64) (int64, error) {
	if len(mids) == 0 {
		return 0, nil
	}
	res, err := d.db.Exec(c, fmt.Sprintf(_rmAuthoritySQL, xstr.JoinInts(mids)))
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
