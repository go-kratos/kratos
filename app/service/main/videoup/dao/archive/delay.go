package archive

import (
	"context"
	"time"

	"go-common/app/service/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"
)

const (
	//insert
	_inDelaySQL = "INSERT INTO archive_delay (mid,aid,state,type,dtime) VALUES (?,?,?,?,?)"
	//update
	_upDelaySQL = "INSERT INTO archive_delay (mid,aid,state,type,dtime,ctime) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE dtime=?,deleted_at='0000-00-00 00:00:00'"
	//delete
	_delDelaySQL = "UPDATE archive_delay SET deleted_at = ? WHERE aid=? AND type=?"
	//select
	_dTimeSQL = "SELECT aid,dtime,state FROM archive_delay WHERE aid=? AND type=? AND deleted_at = 0"
)

// TxAddDelay insert delay.
func (d *Dao) TxAddDelay(tx *sql.Tx, mid int64, aid int64, state, tp int8, dTime xtime.Time) (rows int64, err error) {
	res, err := tx.Exec(_inDelaySQL, mid, aid, state, tp, dTime)
	if err != nil {
		log.Error("d.inDelay.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpDelay update delay
func (d *Dao) TxUpDelay(tx *sql.Tx, mid, aid int64, state, tp int8, dTime xtime.Time) (rows int64, err error) {
	var now = time.Now()
	res, err := tx.Exec(_upDelaySQL, mid, aid, state, tp, dTime, now, dTime)
	if err != nil {
		log.Error("d.TxUpDelay.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxDelDelay delete delay
func (d *Dao) TxDelDelay(tx *sql.Tx, aid int64, tp int8) (rows int64, err error) {
	res, err := tx.Exec(_delDelaySQL, time.Now(), aid, tp)
	if err != nil {
		log.Error("d.TxDelDelay.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// Delay get a delay time by avid.
func (d *Dao) Delay(c context.Context, aid int64, tp int8) (dl *archive.Delay, err error) {
	row := d.rddb.QueryRow(c, _dTimeSQL, aid, tp)
	dl = &archive.Delay{}
	if err = row.Scan(&dl.Aid, &dl.DTime, &dl.State); err != nil {
		if err == sql.ErrNoRows {
			dl = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}
