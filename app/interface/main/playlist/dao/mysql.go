package dao

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"go-common/app/interface/main/playlist/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_plArcSub              = 100
	_plArcSQL              = "SELECT aid,sort,`desc` FROM playlist_archive_%s WHERE pid = ? AND aid = ?"
	_plArcsSQL             = "SELECT aid,sort,`desc` FROM playlist_archive_%s WHERE pid = ? ORDER BY sort"
	_plArcAddSQL           = "INSERT INTO playlist_archive_%s (pid,aid,sort,`desc`) VALUES (?,?,?,?)"
	_plArcBatchAddSQL      = "INSERT INTO playlist_archive_%s (pid,aid,sort,`desc`) VALUES %s"
	_plArcDelSQL           = "DELETE FROM playlist_archive_%s WHERE pid = ? AND aid = ?"
	_plArcDelByPidSQL      = "DELETE FROM playlist_archive_%s WHERE pid = ?"
	_plArcBatchDelSQL      = "DELETE FROM playlist_archive_%s WHERE pid = ? AND aid in (%s)"
	_plArcDescEditSQL      = "UPDATE playlist_archive_%s SET `desc` = ? WHERE pid = ? AND aid = ?"
	_plArcSortEditSQL      = "UPDATE playlist_archive_%s SET sort = ? WHERE pid = ? AND aid = ?"
	_plArcSortBatchEditSQL = "UPDATE playlist_archive_%s SET sort = CASE %s END WHERE pid = ? AND aid IN (%s)"
	_plAddSQL              = "INSERT INTO playlist_stat (mid,fid,is_deleted,view,reply,fav,share) VALUES (?,?,0,0,0,0,0) ON DUPLICATE KEY UPDATE mid=?,fid=?,is_deleted=0,view=0,reply=0,fav=0,share=0"
	_plEditSQL             = "UPDATE playlist_stat SET mtime = ?  WHERE id = ?"
	_plDelSQL              = "UPDATE playlist_stat set is_deleted = 1 WHERE id = ?"
	_plByMidSQL            = "SELECT id,mid,fid,view,reply,fav,`share`,mtime  FROM playlist_stat WHERE is_deleted = 0 AND mid = ?"
	_plByPidsSQL           = "SELECT id,mid,fid,view,reply,fav,`share`,mtime  FROM playlist_stat WHERE id in (%s)"
	_plByPidSQL            = "SELECT id,mid,fid,view,reply,fav,`share`,mtime  FROM playlist_stat WHERE id = ?"
)

func plArcHit(pid int64) string {
	return fmt.Sprintf("%02d", pid%_plArcSub)
}

// Video get video by pid and aid
func (d *Dao) Video(c context.Context, pid, aid int64) (res *model.ArcSort, err error) {
	res = &model.ArcSort{}
	row := d.db.QueryRow(c, fmt.Sprintf(_plArcSQL, plArcHit(pid)), pid, aid)
	if err = row.Scan(&res.Aid, &res.Sort, &res.Desc); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("Video:row.Scan error(%v)", err)
		}
	}
	return
}

// Videos get playlist videos.
func (d *Dao) Videos(c context.Context, pid int64) (res []*model.ArcSort, err error) {
	var rows *xsql.Rows
	if rows, err = d.videosStmt[plArcHit(pid)].Query(c, pid); err != nil {
		log.Error("d.videosStmt[%s].Query(%d) error(%v)", plArcHit(pid), pid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.ArcSort)
		if err = rows.Scan(&r.Aid, &r.Sort, &r.Desc); err != nil {
			log.Error("row.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// AddArc add archive to playlist.
func (d *Dao) AddArc(c context.Context, pid, aid, sort int64, desc string) (lastID int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, fmt.Sprintf(_plArcAddSQL, plArcHit(pid)), pid, aid, desc, sort); err != nil {
		log.Error("AddArc: db.Exec(%d,%d,%d,%s) error(%v)", pid, aid, sort, desc, err)
		return
	}
	return res.LastInsertId()
}

// BatchAddArc add archives to playlist.
func (d *Dao) BatchAddArc(c context.Context, pid int64, arcSorts []*model.ArcSort) (lastID int64, err error) {
	var (
		res    sql.Result
		values []string
	)
	for _, v := range arcSorts {
		values = append(values, fmt.Sprintf("(%d,%d,%d,'%s')", pid, v.Aid, v.Sort, v.Desc))
	}
	if res, err = d.db.Exec(c, fmt.Sprintf(_plArcBatchAddSQL, plArcHit(pid), strings.Join(values, ","))); err != nil {
		log.Error("BatchAddArc: db.Exec(%d) error(%v)", pid, err)
		return
	}
	return res.LastInsertId()
}

// DelArc delete playlist archive.
func (d *Dao) DelArc(c context.Context, pid, aid int64) (affected int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, fmt.Sprintf(_plArcDelSQL, plArcHit(pid)), pid, aid); err != nil {
		log.Error("DelArc: db.Exec(%d,%d) error(%v)", pid, aid, err)
		return
	}
	return res.RowsAffected()
}

// BatchDelArc delete archives from  playlist.
func (d *Dao) BatchDelArc(c context.Context, pid int64, aids []int64) (affected int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, fmt.Sprintf(_plArcBatchDelSQL, plArcHit(pid), xstr.JoinInts(aids)), pid); err != nil {
		log.Error("BatchDelArc: db.Exec(%d,%v) error(%v)", pid, aids, err)
		return
	}
	return res.RowsAffected()
}

// UpdateArcDesc update playlist arc desc.
func (d *Dao) UpdateArcDesc(c context.Context, pid, aid int64, desc string) (affected int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, fmt.Sprintf(_plArcDescEditSQL, plArcHit(pid)), desc, pid, aid); err != nil {
		log.Error("UpdateArcDesc: db.Exec(%d,%d,%s) error(%v)", pid, aid, desc, err)
		return
	}
	return res.RowsAffected()
}

// UpdateArcSort update playlist arc sort.
func (d *Dao) UpdateArcSort(c context.Context, pid, aid, sort int64) (affected int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, fmt.Sprintf(_plArcSortEditSQL, plArcHit(pid)), sort, pid, aid); err != nil {
		log.Error("UpdateArcSort: db.Exec(%d,%d,%d) error(%v)", pid, aid, sort, err)
		return
	}
	return res.RowsAffected()
}

// BatchUpdateArcSort batch update playlist arc sort.
func (d *Dao) BatchUpdateArcSort(c context.Context, pid int64, arcSorts []*model.ArcSort) (affected int64, err error) {
	var (
		caseStr string
		aids    []int64
		res     sql.Result
	)
	for _, v := range arcSorts {
		caseStr = fmt.Sprintf("%s WHEN aid = %d THEN %d", caseStr, v.Aid, v.Sort)
		aids = append(aids, v.Aid)
	}
	if res, err = d.db.Exec(c, fmt.Sprintf(_plArcSortBatchEditSQL, plArcHit(pid), caseStr, xstr.JoinInts(aids)), pid); err != nil {
		log.Error("BatchUpdateArcSort: db.Exec(%d,%s,%v) error(%v)", pid, caseStr, aids, err)
		return
	}
	return res.RowsAffected()
}

//Add playlist stat.
func (d *Dao) Add(c context.Context, mid, fid int64) (lastID int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _plAddSQL, mid, fid, mid, fid); err != nil {
		log.Error("Add:db.Exec(%d,%d) error(%v)", mid, fid, err)
		return
	}
	return res.LastInsertId()
}

// Del playlist stat.
func (d *Dao) Del(c context.Context, pid int64) (affected int64, err error) {
	var (
		res sql.Result
		tx  *xsql.Tx
	)
	if tx, err = d.db.Begin(c); err != nil {
		log.Error("d.db.Begin error(%v)", err)
		return
	}
	if res, err = tx.Exec(_plDelSQL, pid); err != nil {
		tx.Rollback()
		log.Error("DelPlaylist: db.Exec(%d) error(%v)", pid, err)
		return
	}
	if _, err = d.db.Exec(c, fmt.Sprintf(_plArcDelByPidSQL, plArcHit(pid)), pid); err != nil {
		tx.Rollback()
		log.Error("DelArc: db.Exec(%d) error(%v)", pid, err)
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// Update playlist stat.
func (d *Dao) Update(c context.Context, pid int64) (affected int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _plEditSQL, time.Now(), pid); err != nil {
		log.Error("Update mtime: db.Exec(%d) error(%v)", pid, err)
		return
	}
	return res.RowsAffected()
}

// PlsByMid get playlist by mid.
func (d *Dao) PlsByMid(c context.Context, mid int64) (res []*model.PlStat, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, _plByMidSQL, mid); err != nil {
		log.Error("PlsByMid:d.db.Query(%d) error(%v)", mid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.PlStat)
		if err = rows.Scan(&r.ID, &r.Mid, &r.Fid, &r.View, &r.Reply, &r.Fav, &r.Share, &r.MTime); err != nil {
			log.Error("PlsByMid:row.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// PlsByPid get playlist stat by pids.
func (d *Dao) PlsByPid(c context.Context, pids []int64) (res []*model.PlStat, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, fmt.Sprintf(_plByPidsSQL, xstr.JoinInts(pids))); err != nil {
		log.Error("PlsByPid: db.Exec(%s) error(%v)", xstr.JoinInts(pids), err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.PlStat)
		if err = rows.Scan(&r.ID, &r.Mid, &r.Fid, &r.View, &r.Reply, &r.Fav, &r.Share, &r.MTime); err != nil {
			log.Error("PlsByPid:row.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// PlByPid get playlist by pid.
func (d *Dao) PlByPid(c context.Context, pid int64) (res *model.PlStat, err error) {
	res = &model.PlStat{}
	row := d.db.QueryRow(c, _plByPidSQL, pid)
	if err = row.Scan(&res.ID, &res.Mid, &res.Fid, &res.View, &res.Reply, &res.Fav, &res.Share, &res.MTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("PlByPid:row.Scan error(%v)", err)
		}
	}
	return
}
