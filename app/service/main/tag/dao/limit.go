package dao

import (
	"context"

	"go-common/app/service/main/tag/model"
	"go-common/library/log"
)

var (
	_limitResourceSQL = "SELECT oid,`type`,`operation` FROM limit_resource WHERE type=? ;"
)

// LimitRes .
func (d *Dao) LimitRes(c context.Context, tye int32) (res []*model.ResourceLimit, err error) {
	rows, err := d.db.Query(c, _limitResourceSQL, tye)
	if err != nil {
		log.Error("d.db.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.ResourceLimit{}
		if err = rows.Scan(&r.Oid, &r.Type, &r.Attr); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	return
}

var (
	_limitUserSQL = "SELECT mid FROM limit_user"
)

// WhiteUser .
func (d *Dao) WhiteUser(c context.Context) (midm map[int64]struct{}, err error) {
	rows, err := d.db.Query(c, _limitUserSQL)
	if err != nil {
		log.Error("d.db.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	midm = make(map[int64]struct{})
	for rows.Next() {
		var mid int64
		if err = rows.Scan(&mid); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		midm[mid] = struct{}{}
	}
	return
}
