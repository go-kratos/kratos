package archive

import (
	"context"
	"database/sql"
	"time"

	"go-common/app/job/main/videoup/model/archive"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_upStateSQL      = "UPDATE archive SET state=? where id=?"
	_upAccessSQL     = "UPDATE archive SET access=? where id=?"
	_upRoundSQL      = "UPDATE archive SET round=? where id=?"
	_upAttrSQL       = "UPDATE archive SET attribute=attribute|? where id=?"
	_upCoverSQL      = "UPDATE archive SET cover=? where id=?"
	_upDuraSQL       = "UPDATE archive SET duration=? where id=?"
	_upAttrBitSQL    = "UPDATE archive SET attribute=attribute&(~(1<<?))|(?<<?) WHERE id=?"
	_upPTimeSQL      = "UPDATE archive SET pubtime=? WHERE id=?"
	_upDelayRoundSQL = "UPDATE archive SET round = ? WHERE state>=? AND round=? AND mtime >= ? AND mtime <= ?"
	// select
	_arcSQL         = "SELECT id,mid,typeid,copyright,author,title,cover,content,tag,duration,round,attribute,access,state,reject_reason,pubtime,ctime,mtime,forward FROM archive WHERE id=?"
	_upperArcStates = "SELECT id,state FROM archive WHERE mid=?"
)

// Archive get a archive by avid.
func (d *Dao) Archive(c context.Context, aid int64) (a *archive.Archive, err error) {
	var reason sql.NullString
	row := d.db.QueryRow(c, _arcSQL, aid)
	a = &archive.Archive{}
	if err = row.Scan(&a.Aid, &a.Mid, &a.TypeID, &a.Copyright, &a.Author, &a.Title, &a.Cover, &a.Desc, &a.Tag, &a.Duration,
		&a.Round, &a.Attribute, &a.Access, &a.State, &reason, &a.PTime, &a.CTime, &a.MTime, &a.Forward); err != nil {
		if err == sql.ErrNoRows {
			a = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	a.Reason = reason.String
	return
}

// UpperArcStateMap 获取UP主的稿件状态
func (d *Dao) UpperArcStateMap(c context.Context, mid int64) (sMap map[int64]int8, err error) {
	sMap = make(map[int64]int8)
	rows, err := d.rdb.Query(c, _upperArcStates, mid)
	if err != nil {
		log.Error("d.rdb.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := struct {
			ID    int64
			State int8
		}{}
		if err = rows.Scan(&a.ID, &a.State); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		sMap[a.ID] = a.State
	}
	return
}

// UpDelayRound update round to the end by conf
func (d *Dao) UpDelayRound(c context.Context, minTime, maxTime time.Time) (rows int64, err error) {
	res, err := d.db.Exec(c, _upDelayRoundSQL, archive.RoundEnd, archive.StateOpen, archive.RoundReviewFirstWaitTrigger, minTime, maxTime)
	if err != nil {
		log.Error("d.db.Exec error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpState update state of a archive by id.
func (d *Dao) TxUpState(tx *xsql.Tx, aid int64, state int8) (rows int64, err error) {
	res, err := tx.Exec(_upStateSQL, state, aid)
	if err != nil {
		log.Error("tx.Exec(%d, %d) error(%v)", state, aid, err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpAccess update access of a archive by id.
func (d *Dao) TxUpAccess(tx *xsql.Tx, aid int64, access int16) (rows int64, err error) {
	res, err := tx.Exec(_upAccessSQL, access, aid)
	if err != nil {
		log.Error("tx.Exec(%d, %d) error(%v)", access, aid, err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpRound update round of a archive by id.
func (d *Dao) TxUpRound(tx *xsql.Tx, aid int64, round int8) (rows int64, err error) {
	res, err := tx.Exec(_upRoundSQL, round, aid)
	if err != nil {
		log.Error("tx.Exec(%d, %d) error(%v)", round, aid, err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpAttr update attribute value of a archive by id.
func (d *Dao) TxUpAttr(tx *xsql.Tx, aid int64, attr archive.Attr) (rows int64, err error) {
	res, err := tx.Exec(_upAttrSQL, attr, aid)
	if err != nil {
		log.Error("tx.Exec(%d, %d) error(%v)", aid, attr, err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpCover update cover of a archive by id.
func (d *Dao) TxUpCover(tx *xsql.Tx, aid int64, cover string) (rows int64, err error) {
	res, err := tx.Exec(_upCoverSQL, cover, aid)
	if err != nil {
		log.Error("tx.Exec(%d, %s) error(%v)", aid, cover, err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// UpCover update cover of a archive by id.
func (d *Dao) UpCover(c context.Context, aid int64, cover string) (rows int64, err error) {
	res, err := d.db.Exec(c, _upCoverSQL, cover, aid)
	if err != nil {
		log.Error("tx.Exec(%d, %s) error(%v)", aid, cover, err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpArcDuration update duration of a archive by id.
func (d *Dao) TxUpArcDuration(tx *xsql.Tx, aid, duration int64) (rows int64, err error) {
	res, err := tx.Exec(_upDuraSQL, duration, aid)
	if err != nil {
		log.Error("tx.Exec(%d, %d) error(%v)", aid, duration, err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpAttrBit update attribute bit value of a archive by id.
func (d *Dao) TxUpAttrBit(tx *xsql.Tx, aid int64, v int32, bit uint) (rows int64, err error) {
	res, err := tx.Exec(_upAttrBitSQL, bit, v, bit, aid)
	if err != nil {
		log.Error("tx.Exec(%d, %d, %d) error(%v)", aid, v, bit, err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpPTime update ptime by aid
func (d *Dao) TxUpPTime(tx *xsql.Tx, aid int64, ptime time.Time) (rows int64, err error) {
	res, err := tx.Exec(_upPTimeSQL, ptime, aid)
	if err != nil {
		log.Error("tx.Exec(%s, %v, %v) error(%v)", _upPTimeSQL, ptime, aid, err)
		return
	}
	rows, err = res.RowsAffected()
	return
}
