package archive

import (
	"context"
	"database/sql"

	"go-common/app/admin/main/videoup/model/archive"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"
	"time"
)

const (
	_upDelaySQL      = "INSERT INTO archive_delay (aid,type,mid,state,dtime) VALUES(?,?,?,?,?) ON DUPLICATE KEY UPDATE mid=?,state=?,dtime=?,deleted_at='0000-00-00 00:00:00'"
	_upDelStateSQL   = "UPDATE archive_delay SET state=? WHERE aid=? AND type=?"
	_upDelayDtimeSQL = "UPDATE archive_delay SET dtime=? WHERE aid=? AND type=?"
	_delDelaySQL     = "UPDATE archive_delay SET deleted_at = ? WHERE aid=? AND type=?"
	_DelaySQL        = "SELECT aid,dtime,state,mid FROM archive_delay WHERE aid=? AND type=? AND deleted_at = 0"
)

// TxUpDelay update delay
func (d *Dao) TxUpDelay(tx *xsql.Tx, mid, aid int64, state, tp int8, dTime xtime.Time) (rows int64, err error) {
	res, err := tx.Exec(_upDelaySQL, aid, tp, mid, state, dTime, mid, state, dTime)
	if err != nil {
		log.Error("d.TxUpDelay.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpDelState update delay state
func (d *Dao) TxUpDelState(tx *xsql.Tx, aid int64, state, tp int8) (rows int64, err error) {
	res, err := tx.Exec(_upDelStateSQL, state, aid, tp)
	if err != nil {
		log.Error("d.TxUpDelState.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpDelayDtime update archive delaytime by aid.
func (d *Dao) TxUpDelayDtime(tx *xsql.Tx, aid int64, tp int8, dtime xtime.Time) (rows int64, err error) {
	res, err := tx.Exec(_upDelayDtimeSQL, dtime, aid, tp)
	if err != nil {
		log.Error("d.TxUpDelayDtime.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxDelDelay delete delay
func (d *Dao) TxDelDelay(tx *xsql.Tx, aid int64, tp int8) (rows int64, err error) {
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
	row := d.rddb.QueryRow(c, _DelaySQL, aid, tp)
	dl = &archive.Delay{}
	if err = row.Scan(&dl.Aid, &dl.DTime, &dl.State, &dl.Mid); err != nil {
		if err == sql.ErrNoRows {
			dl = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	return
}
