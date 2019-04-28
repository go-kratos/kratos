package archive

import (
	"context"

	"go-common/app/job/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_inArchiveOperSQL = "INSERT INTO archive_oper(aid,uid,typeid,state,round,attribute,last_id,content) VALUES(?,399,?,?,?,?,?,?)"
	_arcPassedOperSQL = "SELECT id FROM archive_oper WHERE aid=? AND state>=? LIMIT 1"
)

// TxArchiveOper add archive oper log
func (d *Dao) TxArchiveOper(tx *sql.Tx, aid int64, typeID int16, state int8, round int8, attr int32, lastID int64, content string) (rows int64, err error) {
	res, err := tx.Exec(_inArchiveOperSQL, aid, typeID, state, round, attr, lastID, content)
	if err != nil {
		log.Error("tx.Exec error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// ArchiveOper add archive oper log
func (d *Dao) ArchiveOper(c context.Context, aid int64, typeID int16, state int8, round int8, attr int32, lastID int64, content string) (rows int64, err error) {
	res, err := d.db.Exec(c, _inArchiveOperSQL, aid, typeID, state, round, attr, lastID, content)
	if err != nil {
		log.Error("tx.Exec error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// PassedOper check archive passed
func (d *Dao) PassedOper(c context.Context, aid int64) (id int64, err error) {
	row := d.db.QueryRow(c, _arcPassedOperSQL, aid, archive.StateOpen)
	if err = row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}
