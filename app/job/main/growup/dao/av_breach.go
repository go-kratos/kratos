package dao

import (
	"context"
	"fmt"

	"go-common/app/job/main/growup/model"
	"go-common/library/log"
)

const (
	// list av_breach_record
	_avBreachRecordSQL = "SELECT id,av_id,mid,money,cdate,reason FROM av_breach_record WHERE cdate >= '%s' AND cdate <= '%s'"
	_avBreachPreSQL    = "SELECT aid,mid FROM av_breach_pre WHERE cdate = '%s' AND ctype = ? AND state = ?"

	// insert
	_inAvBreachPreSQL = "INSERT INTO av_breach_pre(aid, mid, cdate, ctype, state) VALUES(%s) ON DUPLICATE KEY UPDATE state=VALUES(state)"

	// update
	_upAvBreachPresSQL = "UPDATE av_breach_pre SET state = ? WHERE aid = ? AND ctype = ? AND cdate >= '%s' AND state = 1"
)

// GetAvBreach get av_breach by date
func (d *Dao) GetAvBreach(c context.Context, start, end string) (avs []*model.AvBreach, err error) {
	avs = make([]*model.AvBreach, 0)
	rows, err := d.db.Query(c, fmt.Sprintf(_avBreachRecordSQL, start, end))
	if err != nil {
		log.Error("GetAvBreach d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		list := &model.AvBreach{}
		err = rows.Scan(&list.ID, &list.AvID, &list.MID, &list.Money, &list.Date, &list.Reason)
		if err != nil {
			log.Error("GetAvBreach rows scan error(%v)", err)
			return
		}
		avs = append(avs, list)
	}

	err = rows.Err()
	return
}

// GetAvBreachPre get av breach pre by date and state
func (d *Dao) GetAvBreachPre(c context.Context, ctype, state int, cdate string) (avs []*model.AvBreach, err error) {
	avs = make([]*model.AvBreach, 0)
	rows, err := d.db.Query(c, fmt.Sprintf(_avBreachPreSQL, cdate), ctype, state)
	if err != nil {
		log.Error("GetAvBreachPre d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		list := &model.AvBreach{}
		err = rows.Scan(&list.AvID, &list.MID)
		if err != nil {
			log.Error("GetAvBreachPre rows scan error(%v)", err)
			return
		}
		avs = append(avs, list)
	}
	err = rows.Err()
	return
}

// InsertAvBreachPre insert into av_breach_pre
func (d *Dao) InsertAvBreachPre(c context.Context, val string) (rows int64, err error) {
	if val == "" {
		return
	}
	res, err := d.db.Exec(c, fmt.Sprintf(_inAvBreachPreSQL, val))
	if err != nil {
		log.Error("InsertAvBreachPre(%s) tx.Exec error(%v)", val, err)
		return
	}
	return res.RowsAffected()
}

// UpdateAvBreachPre update av breach pre
func (d *Dao) UpdateAvBreachPre(c context.Context, aid, ctype int64, date string, state int) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_upAvBreachPresSQL, date), state, aid, ctype)
	if err != nil {
		log.Error("UpdateAvBreachPre tx.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}
