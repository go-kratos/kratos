package academy

import (
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/creative/model"

	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	// select
	_getArcsSQL = "SELECT id,oid,business FROM academy_archive WHERE state=? AND business=? AND id > ? order by id ASC limit ?"
)

//Archives get limit achives.
func (d *Dao) Archives(c context.Context, id int64, bs, limit int) (res []*model.OArchive, err error) {
	rows, err := d.db.Query(c, _getArcsSQL, 0, bs, id, limit)
	res = make([]*model.OArchive, 0)
	for rows.Next() {
		a := &model.OArchive{}
		if err = rows.Scan(&a.ID, &a.OID, &a.Business); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		res = append(res, a)
	}
	return
}

//UPHotByAIDs update hot by aids.
func (d *Dao) UPHotByAIDs(c context.Context, hots map[int64]int64) error {
	var oids []int64
	sql := "UPDATE academy_archive SET hot = CASE oid "
	for oid, hot := range hots {
		sql += fmt.Sprintf("WHEN %d THEN %d ", oid, hot)
		oids = append(oids, oid)
	}
	sql += fmt.Sprintf("END, mtime=? WHERE oid IN (%s)", xstr.JoinInts(oids))
	_, err := d.db.Exec(c, sql, time.Now())
	if err != nil {
		log.Error("d.db.Exec sql(%s) error(%v)", sql, err)
	}
	log.Info("d.db.Exec sql(%s) hots(%+v) error(%v)", sql, hots, err)
	return err
}
