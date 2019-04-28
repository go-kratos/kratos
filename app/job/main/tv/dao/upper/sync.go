package upper

import (
	"context"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_toSyncUps     = "SELECT mid FROM ugc_uploader WHERE submit = 1 AND deleted = 0"
	_finishSyncUps = "UPDATE ugc_uploader SET submit = 0 WHERE mid = ? AND deleted = 0"
)

// TosyncUps gets all the uppers that need to sync ( with archives )
func (d *Dao) TosyncUps(c context.Context) (mids []int64, err error) {
	var rows *xsql.Rows
	if rows, err = d.DB.Query(c, _toSyncUps); err != nil {
		log.Error("d.TosyncUps.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r int64
		if err = rows.Scan(&r); err != nil {
			log.Error("TosyncUps row.Scan() error(%v)", err)
			return
		}
		mids = append(mids, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("d.ToSyncUps.Query error(%v)", err)
	}
	return
}

// FinsyncUps is used to update the upper's submit from 1 to 0, after the sync is finished
func (d *Dao) FinsyncUps(c context.Context, mid int64) (err error) {
	if _, err = d.DB.Exec(c, _finishSyncUps, mid); err != nil {
		log.Error("FinsyncUps Error: %v", mid, err)
	}
	return
}
