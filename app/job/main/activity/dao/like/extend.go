package like

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

var (
	_batchUpExtSQL = "INSERT INTO like_extend (`lid`,`like`) VALUES %s ON DUPLICATE KEY UPDATE `like`=values(`like`)"
)

// AddExtend .
func (dao *Dao) AddExtend(c context.Context, query string) (res int64, err error) {
	rows, err := dao.db.Exec(c, fmt.Sprintf(_batchUpExtSQL, query))
	if err != nil {
		err = errors.Wrap(err, " dao.db.Exec()")
		return
	}
	return rows.RowsAffected()
}
