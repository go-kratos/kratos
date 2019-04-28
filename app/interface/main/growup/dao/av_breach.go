package dao

import (
	"context"
	"fmt"

	"go-common/app/interface/main/growup/model"
	"go-common/library/log"
	xtime "go-common/library/time"
	"go-common/library/xstr"
)

const (
	// select
	_avBreachSQL       = "SELECT av_id, mid, cdate, money, ctype, reason FROM av_breach_record WHERE mid = ? AND cdate >= ? AND cdate <= ?"
	_avBreachByAvIDSQL = "SELECT av_id, mid, cdate, money, reason FROM av_breach_record WHERE av_id in (%s) AND ctype = ?"
	_avBreachByTypeSQL = "SELECT cdate, money FROM av_breach_record WHERE mid=? AND cdate >= ? AND cdate <= ? AND ctype in (%s)"
)

// ListAvBreach list av_breach_record by mid
func (d *Dao) ListAvBreach(c context.Context, mid int64, startTime, endTime string) (records []*model.AvBreach, err error) {
	records = make([]*model.AvBreach, 0)
	rows, err := d.db.Query(c, _avBreachSQL, mid, startTime, endTime)
	if err != nil {
		log.Error("ListAvBreach d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.AvBreach{}
		err = rows.Scan(&r.AvID, &r.MID, &r.CDate, &r.Money, &r.CType, &r.Reason)
		if err != nil {
			log.Error("ListAvBreach rows.Scan error(%v)", err)
			return
		}
		records = append(records, r)
	}

	err = rows.Err()
	return
}

// GetAvBreachByType get av breach record
func (d *Dao) GetAvBreachByType(c context.Context, mid int64, begin string, end string, typ []int64) (rs map[xtime.Time]int64, err error) {
	if len(typ) == 0 {
		typ = []int64{0, 1, 2, 3}
	}
	rows, err := d.db.Query(c, fmt.Sprintf(_avBreachByTypeSQL, xstr.JoinInts(typ)), mid, begin, end)
	if err != nil {
		log.Error("GetAvBreachByType d.db.Query error(%v)", err)
		return
	}
	rs = make(map[xtime.Time]int64)
	defer rows.Close()
	for rows.Next() {
		var cdate xtime.Time
		var money int64
		err = rows.Scan(&cdate, &money)
		if err != nil {
			log.Error("GetAvBreachByType rows.Scan error(%v)", err)
			return
		}
		rs[cdate] += money
	}
	return
}

// GetAvBreachs get av_breach map by avids
func (d *Dao) GetAvBreachs(c context.Context, avs []int64, ctype int) (breachs map[int64]*model.AvBreach, err error) {
	breachs = make(map[int64]*model.AvBreach)
	rows, err := d.db.Query(c, fmt.Sprintf(_avBreachByAvIDSQL, xstr.JoinInts(avs)), ctype)
	if err != nil {
		log.Error("ListAvBreach d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.AvBreach{}
		if err = rows.Scan(&r.AvID, &r.MID, &r.CDate, &r.Money, &r.Reason); err != nil {
			log.Error("ListAvBreach rows.Scan error(%v)", err)
			return
		}
		breachs[r.AvID] = r
	}
	err = rows.Err()
	return
}
