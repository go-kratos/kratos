package archive

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-common/app/admin/main/videoup/model/archive"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"
	"go-common/library/xstr"
)

const (
	_upArcSQL       = "UPDATE archive SET title=?,content=?,copyright=?,cover=?,note=?,pubtime=?,mtime=? WHERE id=?"
	_upArcTpSQL     = "UPDATE archive SET typeid=? WHERE id=?"
	_upArcRound     = "UPDATE archive SET round=? WHERE id=?"
	_upArcState     = "UPDATE archive SET state=? WHERE id=?"
	_upAccessSQL    = "UPDATE archive SET access=? WHERE id=?"
	_upAuthorSQL    = "UPDATE archive SET mid=?,author=? WHERE id=?"
	_upPTimeSQL     = "UPDATE archive SET pubtime=? WHERE id=?"
	_upArcReasonSQL = "UPDATE archive SET reject_reason=?,forward=? WHERE id=?"
	_upArcAttrSQL   = "UPDATE archive SET attribute=attribute&(~(1<<?))|(?<<?) WHERE id=?"
	_upArcNote      = "UPDATE archive SET note=? WHERE id=?"
	_upArcCopyright = "UPDATE archive SET copyright=? WHERE id=?"
	_upArcMtime     = "UPDATE archive SET mtime=? WHERE id=?"
	_arcSQL         = "SELECT id,mid,title,access,attribute,reject_reason,tag,forward,round,state,copyright,cover,content,typeid,pubtime,ctime,mtime FROM archive WHERE id=?"
	_arcsSQL        = "SELECT id,mid,title,access,attribute,reject_reason,tag,forward,round,state,copyright,cover,content,typeid,pubtime,ctime,mtime FROM archive WHERE id in (%s)"
	_arcStatesSQL   = "SELECT id,state FROM archive WHERE id IN (%s)"
	_upArcTagSQL    = "UPDATE archive SET tag=? WHERE id=?"
)

// TxUpArchive update archive by aid.
func (d *Dao) TxUpArchive(tx *xsql.Tx, aid int64, title, content, cover, note string, copyright int8, pTime xtime.Time) (rows int64, err error) {
	res, err := tx.Exec(_upArcSQL, title, content, copyright, cover, note, pTime, time.Now(), aid)
	if err != nil {
		log.Error("d.TxUpArchive.tx.Exec error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpArcTypeID update archive type_id by aid
func (d *Dao) TxUpArcTypeID(tx *xsql.Tx, aid int64, typeID int16) (rows int64, err error) {
	res, err := tx.Exec(_upArcTpSQL, typeID, aid)
	if err != nil {
		log.Error("d.TxUpArcTypeID.tx.Exec error(%v) ", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpArcRound update archive round by aid
func (d *Dao) TxUpArcRound(tx *xsql.Tx, aid int64, round int8) (rows int64, err error) {
	res, err := tx.Exec(_upArcRound, round, aid)
	if err != nil {
		log.Error("d.TxUpArcRound.tx.Exec error(%v) ", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpArcState update archive state by aid
func (d *Dao) TxUpArcState(tx *xsql.Tx, aid int64, state int8) (rows int64, err error) {
	res, err := tx.Exec(_upArcState, state, aid)
	if err != nil {
		log.Error("d.TxUpArcState.tx.Exec error(%v)", err)
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpArcAccess update archive by aid.
func (d *Dao) TxUpArcAccess(tx *xsql.Tx, aid int64, access int16) (rows int64, err error) {
	res, err := tx.Exec(_upAccessSQL, access, aid)
	if err != nil {
		log.Error("d.TxUpArcAccess.tx.Exec error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpArcAuthor update archive mid && author  by aid.
func (d *Dao) TxUpArcAuthor(tx *xsql.Tx, aid, mid int64, author string) (rows int64, err error) {
	res, err := tx.Exec(_upAuthorSQL, mid, author, aid)
	if err != nil {
		log.Error("d.TxUpArcAuthor.tx.Exec error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpArcPTime update ptime by aid
func (d *Dao) TxUpArcPTime(tx *xsql.Tx, aid int64, pTime xtime.Time) (rows int64, err error) {
	res, err := tx.Exec(_upPTimeSQL, pTime, aid)
	if err != nil {
		log.Error("tx.Exec(%s, %d, %v) error(%v)", _upPTimeSQL, pTime, aid, err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpArcReason update archive reject_reason && forward_id  by aid
func (d *Dao) TxUpArcReason(tx *xsql.Tx, aid, forward int64, reason string) (rows int64, err error) {
	res, err := tx.Exec(_upArcReasonSQL, reason, forward, aid)
	if err != nil {
		log.Error("d.TxUpArcReason.tx.Exec error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpArcAttr update attribute by aid.
func (d *Dao) TxUpArcAttr(tx *xsql.Tx, aid int64, bit uint, val int32) (rows int64, err error) {
	res, err := tx.Exec(_upArcAttrSQL, bit, val, bit, aid)
	if err != nil {
		log.Error("d.upArcAttr.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpArcNote update note by aid.
func (d *Dao) TxUpArcNote(tx *xsql.Tx, aid int64, note string) (rows int64, err error) {
	res, err := tx.Exec(_upArcNote, note, aid)
	if err != nil {
		log.Error("d.upArcNote.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpArcCopyRight update copyright by aid.
func (d *Dao) TxUpArcCopyRight(tx *xsql.Tx, aid int64, copyright int8) (rows int64, err error) {
	res, err := tx.Exec(_upArcCopyright, copyright, aid)
	if err != nil {
		log.Error("d.TxUpArcCopyRight.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpArcMtime update mtime by aid.
func (d *Dao) TxUpArcMtime(tx *xsql.Tx, aid int64) (rows int64, err error) {
	res, err := tx.Exec(_upArcMtime, time.Now(), aid)
	if err != nil {
		log.Error("d.upArcNote.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// Archive get archive by aid
func (d *Dao) Archive(c context.Context, aid int64) (a *archive.Archive, err error) {
	var (
		row         = d.rddb.QueryRow(c, _arcSQL, aid)
		reason, tag sql.NullString
	)
	a = &archive.Archive{}
	if err = row.Scan(&a.Aid, &a.Mid, &a.Title, &a.Access, &a.Attribute, &reason, &tag, &a.Forward, &a.Round, &a.State,
		&a.Copyright, &a.Cover, &a.Desc, &a.TypeID, &a.PTime, &a.CTime, &a.MTime); err != nil {
		if err == xsql.ErrNoRows {
			a = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	a.RejectReason = reason.String
	a.Tag = tag.String
	return
}

// Archives get archives by aids
func (d *Dao) Archives(c context.Context, aids []int64) (am map[int64]*archive.Archive, err error) {
	am = make(map[int64]*archive.Archive)
	if len(aids) == 0 {
		return
	}
	var reason, tag sql.NullString
	rows, err := d.rddb.Query(c, fmt.Sprintf(_arcsSQL, xstr.JoinInts(aids)))
	if err != nil {
		log.Error("db.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := &archive.Archive{}
		if err = rows.Scan(&a.Aid, &a.Mid, &a.Title, &a.Access, &a.Attribute, &reason, &tag, &a.Forward, &a.Round, &a.State,
			&a.Copyright, &a.Cover, &a.Desc, &a.TypeID, &a.PTime, &a.CTime, &a.MTime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		a.RejectReason = reason.String
		a.Tag = tag.String
		am[a.Aid] = a
	}
	return
}

// ArcStateMap get archive id and state map
func (d *Dao) ArcStateMap(c context.Context, aids []int64) (sMap map[int64]int, err error) {
	sMap = make(map[int64]int)
	if len(aids) == 0 {
		return
	}
	rows, err := d.rddb.Query(c, fmt.Sprintf(_arcStatesSQL, xstr.JoinInts(aids)))
	if err != nil {
		log.Error("db.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := struct {
			ID    int64
			State int
		}{}
		if err = rows.Scan(&a.ID, &a.State); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		sMap[a.ID] = a.State
	}
	return
}

//TxUpTag update archive tag
func (d *Dao) TxUpTag(tx *xsql.Tx, aid int64, tags string) (id int64, err error) {
	res, err := tx.Exec(_upArcTagSQL, tags, aid)
	if err != nil {
		log.Error("TxUpTag tx.Exec error(%v) aid(%d) tags(%s)", err, aid, tags)
		return
	}

	return res.LastInsertId()
}
