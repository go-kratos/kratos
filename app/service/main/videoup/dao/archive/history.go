package archive

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-common/app/service/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_inArcHistorySQL    = "INSERT INTO archive_edit_history(aid,mid,title,content,cover,tag) VALUE(?,?,?,?,?,?)"
	_inVideoHistorySQL  = "INSERT INTO archive_video_edit_history(aid,cid,hid,eptitle,description,filename) VALUE (?,?,?,?,?,?)"
	_upVideoHistorySQL  = "UPDATE archive_video_edit_history SET cid=? WHERE aid=? AND filename=? AND cid = 0"
	_arcHistorySQL      = "SELECT id,mid,aid,title,content,cover,tag,ctime FROM archive_edit_history WHERE id =?"
	_arcHistorysSQL     = "SELECT id,mid,aid,title,content,cover,tag,ctime FROM archive_edit_history WHERE aid =? and ctime >=? ORDER BY id DESC"
	_videoHistorySQL    = "SELECT cid,eptitle,description,filename FROM archive_video_edit_history WHERE hid =? ORDER BY id ASC"
	_inVideoHistorysSQL = "INSERT INTO archive_video_edit_history(aid,cid,hid,eptitle,description,filename) VALUES %s"
)

// TxAddArcHistory insert archive_edit_history.
func (d *Dao) TxAddArcHistory(tx *sql.Tx, aid, mid int64, title, content, cover, tag string) (hid int64, err error) {
	res, err := tx.Exec(_inArcHistorySQL, aid, mid, title, content, cover, tag)
	if err != nil {
		log.Error("d.inArcHistory.Exec() error(%v)", err)
		return
	}
	hid, err = res.LastInsertId()
	return
}

// TxAddVideoHistory insert archive_video_edit_history.
func (d *Dao) TxAddVideoHistory(tx *sql.Tx, hid int64, v *archive.Video) (rows int64, err error) {
	res, err := tx.Exec(_inVideoHistorySQL, v.Aid, v.Cid, hid, v.Title, v.Desc, v.Filename)
	if err != nil {
		log.Error("d.inVideoHistory.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpVideoHistory update cid to archive_video_edit_history
func (d *Dao) TxUpVideoHistory(tx *sql.Tx, aid, cid int64, filename string) (rows int64, err error) {
	res, err := tx.Exec(_upVideoHistorySQL, cid, aid, filename)
	if err != nil {
		log.Error("d.upVideoHistory.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// ArcHistory select archive edit history by hid.
func (d *Dao) ArcHistory(c context.Context, hid int64) (ah *archive.ArcHistory, err error) {
	row := d.rddb.QueryRow(c, _arcHistorySQL, hid)
	ah = &archive.ArcHistory{}
	if err = row.Scan(&ah.ID, &ah.Mid, &ah.Aid, &ah.Title, &ah.Content, &ah.Cover, &ah.Tag, &ah.CTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// ArcHistorys select archive edit history by aid.
func (d *Dao) ArcHistorys(c context.Context, aid int64, stime time.Time) (ahs []*archive.ArcHistory, err error) {
	rows, err := d.rddb.Query(c, _arcHistorysSQL, aid, stime)
	if err != nil {
		log.Error("d.arcHissStmt.Query(%d) error(%v)", aid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		ah := &archive.ArcHistory{}
		if err = rows.Scan(&ah.ID, &ah.Mid, &ah.Aid, &ah.Title, &ah.Content, &ah.Cover, &ah.Tag, &ah.CTime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		ahs = append(ahs, ah)
	}
	return
}

// VideoHistory select archive video edit history by hid.
func (d *Dao) VideoHistory(c context.Context, hid int64) (vhs []*archive.VideoHistory, err error) {
	rows, err := d.rddb.Query(c, _videoHistorySQL, hid)
	if err != nil {
		log.Error("d.videoHisStmt.Query(%d) error(%v)", hid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		vh := &archive.VideoHistory{}
		if err = rows.Scan(&vh.Cid, &vh.Title, &vh.Desc, &vh.Filename); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		vhs = append(vhs, vh)
	}
	return
}

// TxAddVideoHistorys batch add archive_video_history.
func (d *Dao) TxAddVideoHistorys(tx *sql.Tx, hid int64, vs []*archive.Video) (err error) {
	log.Info("info TxAddVideoHistorys: hid(%d)|vs(%+v)|cntVs(%d)", hid, vs, len(vs))
	l := len(vs)
	vStrs := make([]string, 0, l)
	vArgs := make([]interface{}, 0, l*6)
	for _, v := range vs {
		vStrs = append(vStrs, "(?, ?, ?, ?, ?, ?)")
		vArgs = append(vArgs, strconv.FormatInt(v.Aid, 10))
		vArgs = append(vArgs, strconv.FormatInt(v.Cid, 10))
		vArgs = append(vArgs, strconv.FormatInt(hid, 10))
		vArgs = append(vArgs, v.Title)
		vArgs = append(vArgs, v.Desc)
		vArgs = append(vArgs, v.Filename)
	}
	stmt := fmt.Sprintf(_inVideoHistorysSQL, strings.Join(vStrs, ","))
	_, err = tx.Exec(stmt, vArgs...)
	if err != nil {
		log.Error("TxAddVideoHistorys: tx.Exec(vs(%+v))|hid(%d) error(%v)", vs, hid, err)
	}
	return
}
