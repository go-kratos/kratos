package share

import (
	"context"
	"fmt"
	"time"

	shamdl "go-common/app/interface/main/web-goblin/model/share"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_shaSub  = 100
	_shasSQL = "SELECT id,mid,day_count,cycle,share_date,ctime,mtime  FROM gb_share_%s WHERE mid = ? and cycle = ?"
)

func shaHit(mid int64) string {
	return fmt.Sprintf("%02d", mid%_shaSub)
}

// Shares get shares.
func (d *Dao) Shares(c context.Context, mid int64) (res []*shamdl.Share, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, fmt.Sprintf(_shasSQL, shaHit(mid)), mid, time.Now().Format("200601")); err != nil {
		log.Error("Shares d.db.Query(%d) error(%v)", mid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(shamdl.Share)
		if err = rows.Scan(&r.ID, &r.Mid, &r.DayCount, &r.Cycle, &r.ShareDate, &r.Ctime, &r.Mtime); err != nil {
			log.Error("Shares:row.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}
