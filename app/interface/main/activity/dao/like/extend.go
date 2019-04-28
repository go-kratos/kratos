package like

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

const (
	_addLikeExtendSQL = "insert into like_extend (lid,`like`) values %s ON DUPLICATE KEY UPDATE `like` = values(`like`);"
)

// AddExtend .
func (dao *Dao) AddExtend(c context.Context, query string) (res int64, err error) {
	rows, err := dao.db.Exec(c, fmt.Sprintf(_addLikeExtendSQL, query))
	if err != nil {
		err = errors.Wrap(err, " dao.db.Exec()")
		return
	}
	return rows.RowsAffected()
}
