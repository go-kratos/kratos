package dao

import (
	"context"

	"go-common/app/service/main/vip/model"
	"go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_jointlysSQL = "SELECT id,title,content,operator,start_time,end_time,link,is_hot,ctime,mtime FROM vip_jointly WHERE end_time > ? AND deleted = 0;"
)

// Jointlys jointly list.
func (d *Dao) Jointlys(c context.Context, now int64) (res []*model.Jointly, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _jointlysSQL, now); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.Jointly)
		if err = rows.Scan(&r.ID, &r.Title, &r.Content, &r.Operator, &r.StartTime, &r.EndTime, &r.Link, &r.IsHot, &r.CTime, &r.MTime); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}
