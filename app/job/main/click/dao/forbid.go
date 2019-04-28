package dao

import (
	"context"

	"go-common/app/job/main/click/model"
	"go-common/library/log"
)

const (
	_forbidSQL        = "SELECT aid,plat,lv,locked FROM archive_click_forbid WHERE locked = ?"
	_upForbidSQL      = "INSERT INTO archive_click_forbid(aid,plat,lv,locked) VALUES(?,?,?,?) ON DUPLICATE KEY UPDATE lv=?,locked=?"
	_allForbidMidsSQL = "SELECT mid FROM archive_mid_forbid WHERE status=1"
	_upForbidMidSQL   = "INSERT INTO archive_mid_forbid(mid,status) VALUE(?,?) ON DUPLICATE KEY UPDATE status=?"
)

// ForbidMids is
func (d *Dao) ForbidMids(c context.Context) (mids map[int64]struct{}, err error) {
	rows, err := d.db.Query(c, _allForbidMidsSQL)
	if err != nil {
		log.Error("d.db.Query(%s) error(%v)", _allForbidMidsSQL, err)
		return
	}
	defer rows.Close()
	mids = make(map[int64]struct{})
	for rows.Next() {
		var mid int64
		if err = rows.Scan(&mid); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		mids[mid] = struct{}{}
	}
	err = rows.Err()
	return
}

// UpMidForbidStatus is
func (d *Dao) UpMidForbidStatus(c context.Context, mid int64, status int8) (err error) {
	_, err = d.db.Exec(c, _upForbidMidSQL, mid, status, status)
	return
}

// Forbids is
func (d *Dao) Forbids(c context.Context) (forbids map[int64]map[int8]*model.Forbid, err error) {
	rows, err := d.db.Query(c, _forbidSQL, model.ValueForLocked)
	if err != nil {
		log.Error("d.db.Query(%s) error(%v)", _forbidSQL, model.ValueForLocked, err)
		return
	}
	defer rows.Close()
	forbids = make(map[int64]map[int8]*model.Forbid)
	for rows.Next() {
		var f = &model.Forbid{}
		if err = rows.Scan(&f.AID, &f.Plat, &f.Lv, &f.Locked); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		if _, ok := forbids[f.AID]; !ok {
			forbids[f.AID] = make(map[int8]*model.Forbid)
		}
		forbids[f.AID][f.Plat] = f
	}
	err = rows.Err()
	return
}

// UpForbid is
func (d *Dao) UpForbid(c context.Context, aid int64, plat, lock, lv int8) (rows int64, err error) {
	res, err := d.db.Exec(c, _upForbidSQL, aid, plat, lv, lock, lv, lock)
	if err != nil {
		log.Error("d.db.Exec(%s) error(%v)", _upForbidSQL, err)
		return
	}
	return res.RowsAffected()
}
