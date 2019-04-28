package archive

import (
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_delaysSQL        = "SELECT id,aid,dtime,type,state FROM archive_delay WHERE aid=? AND deleted_at = 0 ORDER BY dtime DESC LIMIT 1"
	_getNowDelaysSQL  = "SELECT id,aid,dtime,type,state FROM archive_delay WHERE dtime<=? AND deleted_at = 0"
	_delAdminDelaySQL = "UPDATE archive_delay SET deleted_at = ? WHERE aid=? AND type=1"
	_delDelayByIdsSQL = "UPDATE archive_delay SET deleted_at = ? WHERE id IN(%s)"
)

// Delay get delay by aid
func (d *Dao) Delay(c context.Context, aid int64) (delay *archive.Delay, err error) {
	rows := d.db.QueryRow(c, _delaysSQL, aid)
	delay = &archive.Delay{}
	if err = rows.Scan(&delay.ID, &delay.Aid, &delay.DTime, &delay.Type, &delay.State); err != nil {
		if err == sql.ErrNoRows {
			delay = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// NowDelays get current minute dalay archive
func (d *Dao) NowDelays(c context.Context, dtime time.Time) (delays []*archive.Delay, err error) {
	rows, err := d.db.Query(c, _getNowDelaysSQL, dtime)
	if err != nil {
		log.Error("d.db.Query(%s) error(%v)", dtime, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		v := &archive.Delay{}
		if err = rows.Scan(&v.ID, &v.Aid, &v.DTime, &v.Type, &v.State); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		delays = append(delays, v)
	}
	return
}

// DelAdminDelay delete admin delay by aid
func (d *Dao) DelAdminDelay(c context.Context, aid int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _delAdminDelaySQL, time.Now(), aid)
	if err != nil {
		log.Error("d.db.Exec(%d) error(%v)", aid, err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// DelDelayByIds delete delays by ids
func (d *Dao) DelDelayByIds(c context.Context, ids []int64) (rows int64, err error) {
	if len(ids) == 0 {
		return
	}
	res, err := d.db.Exec(c, fmt.Sprintf(_delDelayByIdsSQL, xstr.JoinInts(ids)), time.Now())
	if err != nil {
		log.Error("d.DelDelayByIds() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}
