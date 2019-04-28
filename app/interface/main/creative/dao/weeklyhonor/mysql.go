package weeklyhonor

import (
	"context"

	model "go-common/app/interface/main/creative/model/weeklyhonor"
	"go-common/library/log"
)

const (
	_honorLogsSQL        = `SELECT id,mid,hid,count,ctime,mtime FROM weeklyhonor WHERE mid=?`
	_upsertCountSQL      = `INSERT INTO weeklyhonor (mid,hid,count) VALUES (?,?,?) ON DUPLICATE KEY UPDATE count=count+1`
	_upsertClickCountSQL = `INSERT INTO weeklyhonor_click (mid,count) VALUES (?,?) ON DUPLICATE KEY UPDATE count=count+1`
)

// pingMySQL check mysql connection.
func (d *Dao) pingMySQL(c context.Context) error {
	return d.db.Ping(c)
}

// HonorLogs .
func (d *Dao) HonorLogs(c context.Context, mid int64) (hls map[int]*model.HonorLog, err error) {
	rows, err := d.db.Query(c, _honorLogsSQL, mid)
	if err != nil {
		log.Error("d.db.Query(%s,%d) error(%v)", _honorLogsSQL, mid, err)
		return
	}
	defer rows.Close()
	hls = make(map[int]*model.HonorLog)
	for rows.Next() {
		h := new(model.HonorLog)
		if err = rows.Scan(&h.ID, &h.MID, &h.HID, &h.Count, &h.CTime, &h.MTime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		hls[h.HID] = h
	}
	err = rows.Err()
	return
}

// UpsertCount .
func (d *Dao) UpsertCount(c context.Context, mid int64, hid int) (err error) {
	if _, err = d.db.Exec(c, _upsertCountSQL, mid, hid, 1); err != nil {
		log.Error("d.db.Exec(%s,%d,%d) error(%v)", _upsertCountSQL, mid, hid, err)
		return
	}
	return
}

// UpsertClickCount log weeklyhonor click count
func (d *Dao) UpsertClickCount(c context.Context, mid int64) (err error) {
	if _, err = d.db.Exec(c, _upsertClickCountSQL, mid, 1); err != nil {
		log.Error("d.db.Exec(%s,%d) error(%v)", _upsertClickCountSQL, mid, err)
		return
	}
	return
}
