package archive

import (
	"context"
	sql2 "database/sql"

	"go-common/app/job/main/videoup-report/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_recheckByAid      = "SELECT id,type,aid,uid,state,ctime,mtime FROM archive_recheck WHERE aid =? and type = ?"
	_inRecheck         = "INSERT INTO archive_recheck (type,aid) VALUES (?,?)"
	_upRecheckState    = "UPDATE archive_recheck SET state=? WHERE aid =? and type = ?"
	_upMidRecheckState = "UPDATE archive_recheck as r LEFT JOIN archive AS a ON a.id = r.aid SET r.state=? WHERE a.mid =? and r.type = ? AND r.state = ?"
)

// AddRecheckAids add recheck aids
func (d *Dao) AddRecheckAids(c context.Context, tp int, aids []int64, ignRechecked bool) (err error) {
	for _, aid := range aids {
		recheck, _ := d.RecheckByAid(c, tp, aid)
		if recheck != nil {
			if ignRechecked && recheck.State != archive.RecheckStateIgnore {
				continue
			}
			if recheck.State == archive.RecheckStateWait {
				log.Info("d.AddRecheckAids(%d) already in recheck", aid)
				continue
			}
			if err = d.UpdateRecheckState(c, tp, aid, archive.RecheckStateWait); err != nil {
				log.Error("d.UpdateRecheckState error(%v)", err)
				continue
			}
		} else if _, err = d.db.Exec(c, _inRecheck, tp, aid); err != nil {
			log.Error("d.AddRecheckAids.Exec error(%v)", err)
			continue
		}
		a, err := d.ArchiveByAid(c, aid)
		if err != nil {
			log.Error("d.ArchiveByAid error(%v)", err)
			err = nil
			continue
		}
		tpStr := archive.RecheckType(tp)
		if tpStr != "" {
			d.AddArchiveOper(c, aid, a.Attribute, a.TypeID, a.State, a.Round, 0, "", "待"+tpStr)
		}
	}
	return
}

// UpdateMidRecheckState 设置某个UP主的未回查稿件回查状态
func (d *Dao) UpdateMidRecheckState(c context.Context, tp int, mid int64, state int8) (err error) {
	if _, err = d.db.Exec(c, _upMidRecheckState, state, mid, tp, archive.RecheckStateWait); err != nil {
		log.Error("d.updateRecheckState.Exec error(%v)", err)
		return
	}
	return
}

// TxAddRecheckAID add recheck aid to db
func (d *Dao) TxAddRecheckAID(tx *sql.Tx, tp int, aid int64) (id int64, err error) {
	var (
		res sql2.Result
	)
	if res, err = tx.Exec(_inRecheck, tp, aid); err != nil {
		log.Error("TxAddRecheckAID error(%v) type(%d) aid(%d)", err, tp, aid)
		return
	}

	id, err = res.LastInsertId()
	return
}

// UpdateRecheckState update recheck state
func (d *Dao) UpdateRecheckState(c context.Context, tp int, aid int64, state int8) (err error) {
	if _, err = d.db.Exec(c, _upRecheckState, state, aid, tp); err != nil {
		log.Error("d.updateRecheckState.Exec error(%v)", err)
		return
	}
	return
}

//TxUpRecheckState update recheck state
func (d *Dao) TxUpRecheckState(tx *sql.Tx, tp int, aid int64, state int8) (row int64, err error) {
	var res sql2.Result
	if res, err = tx.Exec(_upRecheckState, state, aid, tp); err != nil {
		log.Error("d.TxUpRecheckState.Exec error(%v)", err)
		return
	}

	row, err = res.RowsAffected()
	return
}

// RecheckByAid find archive recheck
func (d *Dao) RecheckByAid(c context.Context, tp int, aid int64) (recheck *archive.Recheck, err error) {
	row := d.db.QueryRow(c, _recheckByAid, aid, tp)
	recheck = &archive.Recheck{}
	if err = row.Scan(&recheck.ID, &recheck.Type, &recheck.Aid, &recheck.UID, &recheck.State, &recheck.CTime, &recheck.MTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		recheck = nil
		return
	}
	return
}
