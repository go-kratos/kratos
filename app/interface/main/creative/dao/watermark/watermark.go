package watermark

import (
	"context"
	"database/sql"
	"go-common/app/interface/main/creative/model/watermark"
	"go-common/library/log"
)

const (
	// insert
	_inWmSQL = "INSERT IGNORE INTO watermark (mid,uname,state,type,position,url,info,md5,ctime,mtime) VALUES (?,?,?,?,?,?,?,?,?,?)"
	// update
	_upWmSQL = "UPDATE watermark SET mid=?,uname=?,state=?,type=?,position=?,url=?,info=?,md5=?,mtime=? WHERE id=?"
	// select
	_getWmSQL = "SELECT id,mid,uname,state,type,position,url,info,md5,ctime,mtime FROM watermark WHERE mid=?"
)

// AddWaterMark  create watermark.
func (d *Dao) AddWaterMark(c context.Context, w *watermark.Watermark) (id int64, err error) {
	res, err := d.db.Exec(c, _inWmSQL, w.MID, w.Uname, w.State, w.Ty, w.Pos, w.URL, w.Info, w.MD5, w.CTime, w.MTime)
	if err != nil {
		log.Error("d.db.Exec(%d) error(%v)", w.MID, err)
		return
	}
	id, err = res.LastInsertId()
	return
}

// UpWaterMark update watermark info.
func (d *Dao) UpWaterMark(c context.Context, w *watermark.Watermark) (rows int64, err error) {
	res, err := d.db.Exec(c, _upWmSQL, w.MID, w.Uname, w.State, w.Ty, w.Pos, w.URL, w.Info, w.MD5, w.MTime, w.ID)
	if err != nil {
		log.Error("d.db.Exec(%d) error(%v)", w.MID, err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// WaterMark get watermark.
func (d *Dao) WaterMark(c context.Context, mid int64) (w *watermark.Watermark, err error) {
	row := d.db.QueryRow(c, _getWmSQL, mid)
	w = &watermark.Watermark{}
	if err = row.Scan(&w.ID, &w.MID, &w.Uname, &w.State, &w.Ty, &w.Pos, &w.URL, &w.Info, &w.MD5, &w.CTime, &w.MTime); err != nil {
		if err == sql.ErrNoRows {
			w = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	return
}
