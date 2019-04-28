package dao

import (
	"context"

	"go-common/app/service/main/vip/model"
	"go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_allTips = "SELECT `id`,`platform`,`version`,`tip`,`link`,`start_time`,`end_time`,`level`,`judge_type`,`operator`,`deleted`, `position`,`ctime`,`mtime`FROM `vip_tips` WHERE  `deleted` = 0 AND `start_time` < ? AND `end_time`> ?;"
)

//AllTips all tips.
func (d *Dao) AllTips(c context.Context, now int64) (rs []*model.Tips, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _allTips, now, now); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.Tips)
		if err = rows.Scan(&r.ID, &r.Platform, &r.Version, &r.Tip, &r.Link, &r.StartTime, &r.EndTime, &r.Level, &r.JudgeType, &r.Operator, &r.Deleted, &r.Position, &r.Ctime, &r.Mtime); err != nil {
			rs = nil
			err = errors.WithStack(err)
		}
		rs = append(rs, r)
	}
	return
}
