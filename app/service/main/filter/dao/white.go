package dao

import (
	"context"

	"go-common/app/service/main/filter/model"
	xsql "go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_areaWhites = "SELECT fwc.id,fwc.content,fwc.mode,fwa.tpid FROM filter_white_area AS fwa INNER JOIN filter_white_content AS fwc ON fwa.content_id=fwc.id WHERE fwa.area=? AND fwa.state=0"
)

// WhiteAreas get all whites by area
func (d *Dao) WhiteAreas(c context.Context, area string) (fs []*model.WhiteAreaInfo, err error) {
	var rows *xsql.Rows
	if rows, err = d.mysql.Query(c, _areaWhites, area); err != nil {
		return
	}
	defer rows.Close()
	var (
		fsMap = make(map[int64]*model.WhiteAreaInfo)
	)
	for rows.Next() {
		var (
			f    = &model.WhiteAreaInfo{}
			tpID int64
		)
		if err = rows.Scan(&f.ID, &f.Content, &f.Mode, &tpID); err != nil {
			err = errors.WithStack(err)
			return
		}
		if w, ok := fsMap[f.ID]; ok {
			w.TpIDs = append(w.TpIDs, tpID)
		} else {
			fsMap[f.ID] = f
			f.TpIDs = []int64{tpID}
			fs = append(fs, f)
		}
	}
	if err = rows.Err(); err != nil {
		err = errors.WithStack(err)
	}
	return
}
