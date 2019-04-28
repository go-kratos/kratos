package dao

import (
	"context"
	"time"

	"go-common/app/service/main/filter/model"
	xsql "go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_areaRules = `SELECT b.id,a.typeid,b.mode,b.filter,b.level,a.level,b.source FROM filter_area AS a
	INNER JOIN filter_content AS b ON a.filterid=b.id AND a.area=? AND b.source=? AND b.state=0 AND a.is_delete=0 AND b.stime<? AND b.etime>?`
)

// FilterAreas get all filter by area and source.
func (d *Dao) FilterAreas(c context.Context, source int64, area string) (fs []*model.FilterAreaInfo, err error) {
	var (
		rows *xsql.Rows
		cur  = time.Now()
	)
	if rows, err = d.mysql.Query(c, _areaRules, area, source, cur, cur); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()

	var (
		filterMap = make(map[int64]*model.FilterAreaInfo, 1024)
	)
	for rows.Next() {
		var (
			fi               = &model.FilterAreaInfo{}
			tpid             int64
			level, areaLevel int8
		)
		if err = rows.Scan(&fi.ID, &tpid, &fi.Mode, &fi.Filter, &level, &areaLevel, &fi.Source); err != nil {
			err = errors.WithStack(err)
			return
		}
		if f, ok := filterMap[fi.ID]; ok {
			f.TpIDs = append(f.TpIDs, tpid)
		} else {
			filterMap[fi.ID] = fi
			fi.Area = area
			fi.TpIDs = []int64{tpid}
			fi.SetLevel(level, areaLevel)
			fs = append(fs, fi)
		}
	}
	if err = rows.Err(); err != nil {
		err = errors.WithStack(err)
	}
	return
}
