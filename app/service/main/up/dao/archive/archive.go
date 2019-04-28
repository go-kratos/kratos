package archive

import (
	"context"
	"fmt"

	"go-common/app/service/main/up/model"
	"go-common/library/database/sql"
	"go-common/library/time"
	"go-common/library/xstr"
)

const (
	_arcsAidsSQL = "SELECT aid,pubtime,copyright FROM archive WHERE aid IN (%s) AND state>=0 ORDER BY pubtime DESC"
)

// ArcsAids get archives by aids.
func (d *Dao) ArcsAids(c context.Context, ids []int64) (aids []int64, ptimes []time.Time, copyrights []int8, aptm map[int64]*model.AidPubTime, err error) {
	d.infoProm.Incr("ArcAids")
	rows, err := d.resultDB.Query(c, fmt.Sprintf(_arcsAidsSQL, xstr.JoinInts(ids)))
	if err != nil {
		return
	}
	defer rows.Close()
	aptm = make(map[int64]*model.AidPubTime, len(ids))
	for rows.Next() {
		var (
			aid       int64
			ptime     time.Time
			copyright int8
		)
		if err = rows.Scan(&aid, &ptime, &copyright); err != nil {
			if err == sql.ErrNoRows {
				err = nil
				return
			}
			return
		}
		aptm[aid] = &model.AidPubTime{Aid: aid, PubDate: ptime, Copyright: copyright}
		aids = append(aids, aid)
		ptimes = append(ptimes, ptime)
		copyrights = append(copyrights, copyright)
	}
	return
}
