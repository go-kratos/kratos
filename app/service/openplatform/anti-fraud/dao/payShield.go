package dao

import (
	"context"
	"go-common/app/service/openplatform/anti-fraud/model"
	"go-common/library/log"
)

const (
	_ipList    = "select ip,count(1) as num from shield_ip_log where mtime >= ? group by ip order by num desc limit 50"
	_ipDetail  = "select ip,uid from shield_ip_log where ip = ? and mtime >= ? and mtime <= ?"
	_uidList   = "select uid,count(1) as num from shield_user_log where mtime >= ? group by uid order by num desc limit 50"
	_uidDetail = "select ip,uid from shield_user_log where uid = ? and  mtime >= ? and mtime <= ?"
)

// ShieldIPList .
func (d *Dao) ShieldIPList(c context.Context, mtime string) (res []*model.IPListDetail, err error) {
	res = make([]*model.IPListDetail, 0)

	rows, err := d.payShieldDb.Query(c, _ipList, mtime)
	if err != nil {
		log.Warn("select err %s %v", _ipList, err)
		return
	}

	for rows.Next() {
		r := new(model.IPListDetail)
		if err = rows.Scan(&r.IP, &r.Num); err != nil {
			log.Warn("scan err %v", err)
			return
		}
		res = append(res, r)
	}

	return
}

// ShieldIPDetail .
func (d *Dao) ShieldIPDetail(c context.Context, ip, stime, etime string) (res []*model.ListDetail, err error) {
	res = make([]*model.ListDetail, 0)

	rows, err := d.payShieldDb.Query(c, _ipDetail, ip, stime, etime)
	if err != nil {
		log.Warn("select err %s %v", _ipDetail, err)
		return
	}

	for rows.Next() {
		r := new(model.ListDetail)
		if err = rows.Scan(&r.IP, &r.UID); err != nil {
			log.Warn("scan err %v", err)
			return
		}
		res = append(res, r)
	}

	return
}

// ShieldUIDList .
func (d *Dao) ShieldUIDList(c context.Context, mtime string) (res []*model.UIDListDetail, err error) {
	res = make([]*model.UIDListDetail, 0)

	rows, err := d.payShieldDb.Query(c, _uidList, mtime)
	if err != nil {
		log.Warn("select err %s %v", _uidList, err)
		return
	}

	for rows.Next() {
		r := new(model.UIDListDetail)
		if err = rows.Scan(&r.UID, &r.Num); err != nil {
			log.Warn("scan err %v", err)
			return
		}
		res = append(res, r)
	}

	return
}

// ShieldUIDDetail .
func (d *Dao) ShieldUIDDetail(c context.Context, uid, stime, etime string) (res []*model.ListDetail, err error) {
	res = make([]*model.ListDetail, 0)

	rows, err := d.payShieldDb.Query(c, _uidDetail, uid, stime, etime)
	if err != nil {
		log.Warn("select err %s %v", _uidDetail, err)
		return
	}

	for rows.Next() {
		r := new(model.ListDetail)
		if err = rows.Scan(&r.IP, &r.UID); err != nil {
			log.Warn("scan err %v", err)
			return
		}
		res = append(res, r)
	}

	return
}
