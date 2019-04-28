package resource

import (
	"context"
	"database/sql"

	"go-common/app/interface/main/web-show/model/resource"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"
)

const (
	_inVdoActSQL       = "INSERT IGNORE INTO video_ads_active (name,aid,cid,url,skipable,strategy,mtime) VALUES(?,?,?,?,?,?,?)"
	_selAllVdoActSQL   = "SELECT name,aid,cid,url,skipable,strategy FROM video_ads_active ORDER BY id"
	_selVdoMTCntActSQL = `SELECT FROM_UNIXTIME(ROUND(AVG(UNIX_TIMESTAMP(mtime)))) FROM video_ads_active`
	_delAllVdoActSQL   = "DELETE FROM video_ads_active"
)

// initActive init
func (dao *Dao) initActive() {
	dao.selAllVdoActStmt = dao.videodb.Prepared(_selAllVdoActSQL)
	dao.selVdoActMTCntStmt = dao.videodb.Prepared(_selVdoMTCntActSQL)
	dao.delAllVdoActStmt = dao.videodb.Prepared(_delAllVdoActSQL)
}

// TxInsertVideo dao
func (dao *Dao) TxInsertVideo(tx *xsql.Tx, vad resource.VideoAD) (err error) {
	if _, err = tx.Exec(_inVdoActSQL, vad.Name, vad.Aid, vad.Cid, vad.URL, vad.Skipable, vad.Strategy, vad.MTime); err != nil {
		log.Error("tx.Stmt(dao.inStmt).Exec(), err (%v)", err)
	}
	return
}

//VideoAds dao
func (dao *Dao) VideoAds(c context.Context) (vads map[int64][]*resource.VideoAD, err error) {
	rows, err := dao.selAllVdoActStmt.Query(c)
	if err != nil {
		log.Error("dao.selAllVdoStmt query error (%v)", err)
		return
	}
	defer rows.Close()
	vads = map[int64][]*resource.VideoAD{}
	for rows.Next() {
		vad := &resource.VideoAD{}
		if err = rows.Scan(&vad.Name, &vad.Aid, &vad.Cid, &vad.URL, &vad.Skipable, &vad.Strategy); err != nil {
			log.Error("rows.Scan err (%v)", err)
			vads = nil
			return
		}
		vads[vad.Aid] = append(vads[vad.Aid], vad)
	}
	return
}

// ActVideoMTimeCount dao
func (dao *Dao) ActVideoMTimeCount(c context.Context) (mtime xtime.Time, err error) {
	row := dao.selVdoActMTCntStmt.QueryRow(c)
	if err = row.Scan(&mtime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan(), err (%v)", err)
		}
	}
	return
}

// DelAllVideo dao
func (dao *Dao) DelAllVideo(c context.Context) (err error) {
	if _, err = dao.delAllVdoActStmt.Exec(c); err != nil {
		log.Error("mysqlDB.Exec error(%v)", err)
	}
	return
}
