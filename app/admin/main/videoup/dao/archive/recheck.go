package archive

import (
	"context"
	"fmt"

	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_recheckByAid           = "SELECT id,type,aid,uid,state,ctime,mtime FROM archive_recheck WHERE aid =? and type = ?"
	_recheckBatchIDByAid    = "SELECT id,aid FROM archive_recheck WHERE aid IN (%s) AND type=?"
	_recheckBatchStateByAid = "SELECT aid,state FROM archive_recheck WHERE aid IN (%s) AND type=?"
	_upRecheckState         = "UPDATE archive_recheck SET state=? WHERE aid =? and type = ?"
)

// TxUpRecheckState update recheck state
func (d *Dao) TxUpRecheckState(tx *sql.Tx, tp int, aid int64, state int8) (err error) {
	if _, err = tx.Exec(_upRecheckState, state, aid, tp); err != nil {
		log.Error("TxUpRecheckState Exec(%d,%d,%d) error(%v)", state, tp, aid, err)
		return
	}
	return
}

// RecheckByAid find archive recheck
func (d *Dao) RecheckByAid(c context.Context, tp int, aid int64) (recheck *archive.Recheck, err error) {
	row := d.db.QueryRow(c, _recheckByAid, aid, tp)
	recheck = &archive.Recheck{}
	if err = row.Scan(&recheck.ID, &recheck.Type, &recheck.AID, &recheck.UID, &recheck.State, &recheck.CTime, &recheck.MTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			recheck = nil
		} else {
			log.Error("RecheckByAid row.Scan(%d,%d) error(%v)", tp, aid, err)
		}
		return
	}
	return
}

//RecheckIDByAID find states by ids
func (d *Dao) RecheckIDByAID(c context.Context, tp int, aids []int64) (ids []int64, existAID []int64, err error) {
	var (
		rows    *sql.Rows
		id, aid int64
	)
	aidstr := xstr.JoinInts(aids)
	ids = []int64{}
	existAID = []int64{}
	if rows, err = d.db.Query(c, fmt.Sprintf(_recheckBatchIDByAid, aidstr), tp); err != nil {
		log.Error("RecheckIDByAID d.db.Query error(%v) type(%d) aids(%s)", err, tp, aidstr)
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&id, &aid); err != nil {
			log.Error("RecheckIDByAID rows.Scan error(%v) type(%d) aids(%s)", err, tp, aidstr)
			return
		}

		ids = append(ids, id)
		existAID = append(existAID, aid)
	}
	return
}

func (d *Dao) RecheckStateMap(c context.Context, tp int, aids []int64) (m map[int64]int8, err error) {
	var (
		rows  *sql.Rows
		aid   int64
		state int8
	)
	m = make(map[int64]int8)
	if len(aids) == 0 {
		return
	}
	str := xstr.JoinInts(aids)
	if rows, err = d.db.Query(c, fmt.Sprintf(_recheckBatchStateByAid, str), tp); err != nil {
		log.Error("RecheckStateMap d.db.Query error(%v) type(%d) aids(%s)", err, tp, str)
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&aid, &state); err != nil {
			log.Error("RecheckStateMap rows.Scan error(%v) type(%d) aids(%s)", err, tp, str)
			return
		}
		m[aid] = state
	}
	return
}
