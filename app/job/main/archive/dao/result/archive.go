package result

import (
	"context"
	"time"

	"go-common/app/job/main/archive/model/archive"
	"go-common/app/job/main/archive/model/result"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_arcSQL       = "SELECT aid,mid,typeid,videos,copyright,title,cover,content,duration,attribute,state,access,pubtime,ctime,mission_id,order_id,redirect_url,forward,dynamic FROM archive WHERE aid=?"
	_inArchiveSQL = `INSERT IGNORE INTO archive (aid,mid,typeid,videos,title,cover,content,duration,attribute,copyright,access,pubtime,state,mission_id,order_id,redirect_url,forward,dynamic,cid,dimensions)
			VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	_upArchiveSQL  = "UPDATE archive SET mid=?,typeid=?,videos=?,title=?,cover=?,content=?,duration=?,attribute=?,copyright=?,access=?,pubtime=?,state=?,mission_id=?,order_id=?,redirect_url=?,mtime=?,forward=?,dynamic=?,cid=?,dimensions=? WHERE aid=?"
	_delArchiveSQL = "UPDATE archive SET state=? WHERE aid=?"
	_upPassedSQL   = "SELECT aid FROM archive WHERE mid=? AND state>=0 ORDER BY pubtime DESC"
)

// UpPassed is
func (d *Dao) UpPassed(c context.Context, mid int64) (aids []int64, err error) {
	rows, err := d.db.Query(c, _upPassedSQL, mid)
	if err != nil {
		log.Error("d.db.Query(%s, %d) error(%v)", _upPassedSQL, mid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var aid int64
		if err = rows.Scan(&aid); err != nil {
			log.Error("rows.Scan(%d) error(%v)", aid, err)
			return
		}
		aids = append(aids, aid)
	}
	err = rows.Err()
	return
}

// Archive get a archive by aid.
func (d *Dao) Archive(c context.Context, aid int64) (a *result.Archive, err error) {
	row := d.db.QueryRow(c, _arcSQL, aid)
	a = &result.Archive{}
	if err = row.Scan(&a.AID, &a.Mid, &a.TypeID, &a.Videos, &a.Copyright, &a.Title, &a.Cover, &a.Content, &a.Duration,
		&a.Attribute, &a.State, &a.Access, &a.PubTime, &a.CTime, &a.MissionID, &a.OrderID, &a.RedirectURL, &a.Forward, &a.Dynamic); err != nil {
		if err == sql.ErrNoRows {
			a = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// TxAddArchive add archive result
func (d *Dao) TxAddArchive(c context.Context, tx *sql.Tx, a *archive.Archive, ad *archive.Addit, videoCnt int, firstCid int64, dimensions string) (rows int64, err error) {
	res, err := tx.Exec(_inArchiveSQL, a.ID, a.Mid, a.TypeID, videoCnt, a.Title, a.Cover, a.Content, a.Duration, a.Attribute, a.Copyright, a.Access, a.PubTime, a.State, ad.MissionID, ad.OrderID, ad.RedirectURL, a.Forward, ad.Dynamic, firstCid, dimensions)
	if err != nil {
		log.Error("tx.Exec(%s) error(%v)", _inArchiveSQL, err)
		return
	}
	return res.RowsAffected()
}

// TxUpArchive update archive result
func (d *Dao) TxUpArchive(c context.Context, tx *sql.Tx, a *archive.Archive, ad *archive.Addit, videoCnt int, firstCid int64, dimensions string) (rows int64, err error) {
	res, err := tx.Exec(_upArchiveSQL, a.Mid, a.TypeID, videoCnt, a.Title, a.Cover, a.Content, a.Duration, a.Attribute, a.Copyright, a.Access, a.PubTime, a.State, ad.MissionID, ad.OrderID, ad.RedirectURL, time.Now(), a.Forward, ad.Dynamic, firstCid, dimensions, a.ID)
	if err != nil {
		log.Error("tx.Exec(%s) error(%v)", _upArchiveSQL, err)
		return
	}
	return res.RowsAffected()
}

// TxDelArchive delete archive
func (d *Dao) TxDelArchive(c context.Context, tx *sql.Tx, aid int64) (rows int64, err error) {
	res, err := tx.Exec(_delArchiveSQL, archive.StateForbidUpDelete, aid)
	if err != nil {
		log.Error("tx.Execerror(%v)", err)
		return
	}
	return res.RowsAffected()
}
