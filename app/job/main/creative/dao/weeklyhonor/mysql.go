package weeklyhonor

import (
	"context"
	"fmt"

	model "go-common/app/interface/main/creative/model/weeklyhonor"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_honorLogsSQL       = `SELECT id,mid,hid,count,ctime,mtime FROM weeklyhonor WHERE mid=?`
	_honorLatestLogsSQL = `SELECT id,mid,hid,mtime FROM weeklyhonor WHERE mid in (%s) and mtime >= ?`
	_honorClicksSQL     = `SELECT mid,count FROM weeklyhonor_click WHERE mid in (%s) and mtime >= ?`
	_upsertCountSQL     = `INSERT INTO weeklyhonor (mid,hid,count) VALUES (?,?,?) ON DUPLICATE KEY UPDATE count=count+1`
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

// LatestHonorLogs list latest honor logs by mids
func (d *Dao) LatestHonorLogs(c context.Context, mids []int64) (hls []*model.HonorLog, err error) {
	latestSun := model.LatestSunday()
	midsStr := xstr.JoinInts(mids)
	if midsStr == "" {
		return
	}
	sql := fmt.Sprintf(_honorLatestLogsSQL, midsStr)
	rows, err := d.db.Query(c, sql, latestSun)
	if err != nil {
		log.Error("d.db.Query(%s,%q,%v) error(%v)", _honorLatestLogsSQL, mids, latestSun, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		h := new(model.HonorLog)
		if err = rows.Scan(&h.ID, &h.MID, &h.HID, &h.MTime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		hls = append(hls, h)
	}
	err = rows.Err()
	return
}

// UpsertCount .
func (d *Dao) UpsertCount(c context.Context, mid int64, hid int) (affected int64, err error) {
	res, err := d.db.Exec(c, _upsertCountSQL, mid, hid, 1)
	if err != nil {
		log.Error("d.db.Exec(%s,%d,%d,%d) error(%v)", _upsertCountSQL, mid, hid, 1, err)
		return
	}
	return res.RowsAffected()
}

// ClickCounts  honor click count map
func (d *Dao) ClickCounts(c context.Context, mids []int64) (res map[int64]int32, err error) {
	res = make(map[int64]int32)
	twoWeekAgo := model.LatestSunday().AddDate(0, 0, -14)
	midsStr := xstr.JoinInts(mids)
	if midsStr == "" {
		return
	}
	sql := fmt.Sprintf(_honorClicksSQL, midsStr)
	rows, err := d.db.Query(c, sql, twoWeekAgo)
	if err != nil {
		log.Error("d.db.Query(%s,%v) error(%v)", sql, twoWeekAgo, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var mid, count int64
		if err = rows.Scan(&mid, &count); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		if count > 0 {
			res[mid] = int32(count)
		}
	}
	return
}
