package archive

import (
	"context"
	"database/sql"
	"fmt"

	"go-common/app/service/main/archive/api"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_maxAIDSQL    = "SELECT max(aid) FROM archive"
	_arcSQL       = "SELECT aid,mid,typeid,videos,copyright,title,cover,content,duration,attribute,state,access,pubtime,ctime,mission_id,order_id,redirect_url,forward,dynamic,cid,dimensions FROM archive WHERE aid=?"
	_arcsSQL      = "SELECT aid,mid,typeid,videos,copyright,title,cover,content,duration,attribute,state,access,pubtime,ctime,mission_id,order_id,redirect_url,forward,dynamic,cid,dimensions FROM archive WHERE aid IN (%s)"
	_arcStaffSQL  = "SELECT mid,title FROM archive_staff WHERE aid=?"
	_arcsStaffSQL = "SELECT aid,mid,title FROM archive_staff WHERE aid IN(%s)"
)

// MaxAID get max aid
func (d *Dao) MaxAID(c context.Context) (id int64, err error) {
	row := d.resultDB.QueryRow(c, _maxAIDSQL)
	if err = row.Scan(&id); err != nil {
		log.Error("row.Scan error(%v)", err)
		return
	}
	return
}

// archivePB get a archive by aid.
func (d *Dao) archive3(c context.Context, aid int64) (a *api.Arc, err error) {
	d.infoProm.Incr("archive3")
	row := d.resultDB.QueryRow(c, _arcSQL, aid)
	a = &api.Arc{}
	var dimension string
	if err = row.Scan(&a.Aid, &(a.Author.Mid), &a.TypeID, &a.Videos, &a.Copyright, &a.Title, &a.Pic, &a.Desc, &a.Duration,
		&a.Attribute, &a.State, &a.Access, &a.PubDate, &a.Ctime, &a.MissionID, &a.OrderID, &a.RedirectURL, &a.Forward, &a.Dynamic, &a.FirstCid, &dimension); err != nil {
		if err == sql.ErrNoRows {
			a = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
			d.errProm.Incr("result_db")
		}
		return
	}
	a.FillDimension(dimension)
	return
}

// archives3 multi get archives by avids.
func (d *Dao) archives3(c context.Context, aids []int64) (res map[int64]*api.Arc, err error) {
	d.infoProm.Incr("archives3")
	query := fmt.Sprintf(_arcsSQL, xstr.JoinInts(aids))
	rows, err := d.resultDB.Query(c, query)
	if err != nil {
		log.Error("db.Query(%s) error(%v)", query, err)
		d.errProm.Incr("result_db")
		return
	}
	defer rows.Close()
	res = make(map[int64]*api.Arc, len(aids))
	for rows.Next() {
		a := &api.Arc{}
		var dimension string
		if err = rows.Scan(&a.Aid, &(a.Author.Mid), &a.TypeID, &a.Videos, &a.Copyright, &a.Title, &a.Pic, &a.Desc, &a.Duration,
			&a.Attribute, &a.State, &a.Access, &a.PubDate, &a.Ctime, &a.MissionID, &a.OrderID, &a.RedirectURL, &a.Forward, &a.Dynamic, &a.FirstCid, &dimension); err != nil {
			log.Error("rows.Scan error(%v)", err)
			d.errProm.Incr("result_db")
			return
		}
		a.FillDimension(dimension)
		res[a.Aid] = a
	}
	err = rows.Err()
	return
}

// staff get archives staff by avid.
func (d *Dao) staff(c context.Context, aid int64) (res []*api.StaffInfo, err error) {
	d.infoProm.Incr("archive_staff")
	rows, err := d.resultDB.Query(c, _arcStaffSQL, aid)
	if err != nil {
		log.Error("d.resultDB.Query(%d) error(%v)", aid, err)
		d.errProm.Incr("result_db")
		return
	}
	defer rows.Close()
	for rows.Next() {
		as := &api.StaffInfo{}
		if err = rows.Scan(&as.Mid, &as.Title); err != nil {
			log.Error("rows.Scan error(%v)", err)
			d.errProm.Incr("result_db")
			return
		}
		res = append(res, as)
	}
	err = rows.Err()
	return
}

// staffs get archives staff by avids.
func (d *Dao) staffs(c context.Context, aids []int64) (res map[int64][]*api.StaffInfo, err error) {
	d.infoProm.Incr("archives_staff")
	query := fmt.Sprintf(_arcsStaffSQL, xstr.JoinInts(aids))
	rows, err := d.resultDB.Query(c, query)
	if err != nil {
		log.Error("d.resultDB.Query(%s) error(%v)", query, err)
		d.errProm.Incr("result_db")
		return
	}
	defer rows.Close()
	res = make(map[int64][]*api.StaffInfo)
	for rows.Next() {
		as := &api.StaffInfo{}
		var aid int64
		if err = rows.Scan(&aid, &as.Mid, &as.Title); err != nil {
			log.Error("rows.Scan error(%v)", err)
			d.errProm.Incr("result_db")
			return
		}
		res[aid] = append(res[aid], as)
	}
	err = rows.Err()
	return
}
