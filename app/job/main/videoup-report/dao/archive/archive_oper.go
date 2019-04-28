package archive

import (
	"context"
	"database/sql"

	"go-common/app/job/main/videoup-report/model/archive"
	"go-common/library/log"
)

const (
	_operLastRoundSQL    = "SELECT id,aid,uid,typeid,state,round,attribute,last_id,ctime,mtime FROM archive_oper WHERE aid = ? AND round = ? ORDER BY id DESC LIMIT 1"
	_operNextRoundSQL    = "SELECT id,aid,uid,typeid,state,round,attribute,last_id,ctime,mtime FROM archive_oper WHERE id > ? AND aid = ? AND round != ? ORDER BY id ASC LIMIT 1"
	_operInsertSQL       = "INSERT INTO archive_oper(aid,uid,typeid,state,content,round,attribute,last_id,remark) VALUES(?,399,?,?,?,?,?,?,?)"
	_lastArcOperSQL      = "SELECT aid,uid,typeid,content,round,attribute,last_id,remark FROM archive_oper WHERE aid=? AND uid!=399 ORDER BY mtime DESC LIMIT 1"
	_lastVideoOperUIDSQL = "SELECT uid FROM archive_video_oper WHERE vid=? AND uid != 399 ORDER BY mtime DESC LIMIT 1;"
	_lastVideoOperSQL    = "SELECT aid, uid, vid, status, content, attribute, last_id, remark FROM archive_video_oper WHERE vid=? AND uid != 399 ORDER BY mtime DESC LIMIT 1;"
)

// LastRoundOper get last archive round record.
func (d *Dao) LastRoundOper(c context.Context, aid int64, round int8) (oper *archive.Oper, err error) {
	row := d.db.QueryRow(c, _operLastRoundSQL, aid, round)
	oper = &archive.Oper{}
	if err = row.Scan(&oper.ID, &oper.AID, &oper.UID, &oper.TypeID, &oper.State, &oper.Round, &oper.Attribute, &oper.LastID, &oper.CTime, &oper.MTime); err != nil {
		log.Error("row.Scan error(%v)", err)
		return
	}
	return
}

// NextRoundOper get next archive round record.
func (d *Dao) NextRoundOper(c context.Context, id int64, aid int64, round int8) (oper *archive.Oper, err error) {
	row := d.db.QueryRow(c, _operNextRoundSQL, id, aid, round)
	oper = &archive.Oper{}
	if err = row.Scan(&oper.ID, &oper.AID, &oper.UID, &oper.TypeID, &oper.State, &oper.Round, &oper.Attribute, &oper.LastID, &oper.CTime, &oper.MTime); err != nil {
		log.Error("row.Scan error(%v)", err)
		return
	}
	return
}

// AddArchiveOper add archive operate log
func (d *Dao) AddArchiveOper(c context.Context, aid int64, attribute int32, typeid int16, state int, round int8, lastID int64, content, remark string) (id int64, err error) {
	var (
		res sql.Result
	)
	if res, err = d.db.Exec(c, _operInsertSQL, aid, typeid, state, content, round, attribute, lastID, remark); err != nil {
		log.Error("AddArchiveOper(%d,%d,%d,%d,%s,%d,%d,%d,%s) error(%v)", aid, typeid, state, content, round, attribute, lastID, remark, err)
		return
	}

	id, err = res.LastInsertId()
	return
}

//LastVideoOperUID get the last manual-operate operator id by vid
func (d *Dao) LastVideoOperUID(c context.Context, vid int64) (uid int64, err error) {
	if err = d.db.QueryRow(c, _lastVideoOperUIDSQL, vid).Scan(&uid); err != nil {
		log.Error("LastVideoOperUID db.row.Scan error(%v) vid(%d)", err, vid)
	}

	return
}

//LastVideoOper get the last manual-operate record by vid
func (d *Dao) LastVideoOper(c context.Context, vid int64) (oper *archive.VideoOper, err error) {
	oper = &archive.VideoOper{}
	if err = d.db.QueryRow(c, _lastVideoOperSQL, vid).Scan(&oper.AID, &oper.UID, &oper.VID, &oper.Status, &oper.Content, &oper.Attribute, &oper.LastID, &oper.Remark); err != nil {
		log.Error("LastVideoOper db.row.Scan error(%v), vid(%d)", err, vid)
	}
	return
}

// LastArcOper get a archive last history.
func (d *Dao) LastArcOper(c context.Context, aid int64) (re *archive.Oper, err error) {
	re = &archive.Oper{}
	if err = d.db.QueryRow(c, _lastArcOperSQL, aid).Scan(&re.AID, &re.UID, &re.TypeID, &re.Content, &re.Round, &re.Attribute, &re.LastID, &re.Remark); err != nil {
		log.Error(" LastArcOper db.row.Scan error(%v) aid(%d)", err, aid)
	}
	return
}
