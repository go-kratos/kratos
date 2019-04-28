package ugc

import (
	"context"
	"fmt"

	"go-common/app/job/main/tv/model/common"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_passedArc = "SELECT aid, ctime FROM ugc_archive WHERE typeid IN (%s) AND result = 1 AND valid = 1 and deleted = 0"
)

// PassedArcs picks the passed Arc data for Idx Page
func (d *Dao) PassedArcs(c context.Context, tids []int64) (res []*common.IdxRank, err error) {
	if len(tids) == 0 {
		return
	}

	var (
		rows    *sql.Rows
		tidsStr = xstr.JoinInts(tids)
	)
	if rows, err = d.DB.Query(c, fmt.Sprintf(_passedArc, tidsStr)); err != nil {
		log.Error("d.PassedArcs.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r = &common.IdxRank{}
		if err = rows.Scan(&r.ID, &r.Ctime); err != nil {
			log.Error("PassedArcs row.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("d.PassedArcs.Query error(%v)", err)
	}
	return
}
