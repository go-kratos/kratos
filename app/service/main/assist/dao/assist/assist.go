package assist

import (
	"context"

	"go-common/app/service/main/assist/model/assist"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	// insert
	_inAssSQL = "INSERT INTO assist (mid,assist_mid) VALUE (?,?) ON DUPLICATE KEY UPDATE state=0"
	// update
	_delAssSQL = "UPDATE assist SET state=1 WHERE mid=? AND assist_mid=?"
	// select
	_assSQL       = "SELECT mid,assist_mid,ctime,mtime FROM assist WHERE mid=? AND assist_mid=? AND state=0"
	_asssSQL      = "SELECT mid,assist_mid,ctime,mtime FROM assist WHERE mid=? AND state=0"
	_assCntSQL    = "SELECT count(*) FROM assist WHERE mid=? AND state=0"
	_assUpsSQL    = "SELECT mid, ctime FROM assist WHERE assist_mid=? AND state=0 limit ?,?"
	_assUpsCntSQL = "SELECT count(*) as total FROM assist WHERE assist_mid=? AND state=0"
)

// AddAssist add assist
func (d *Dao) AddAssist(c context.Context, mid, assistMid int64) (id int64, err error) {
	res, err := d.db.Exec(c, _inAssSQL, mid, assistMid)
	if err != nil {
		log.Error("d.inAss.Exec error(%v)", err)
		return
	}
	id, err = res.LastInsertId()
	return
}

// DelAssist del assist
func (d *Dao) DelAssist(c context.Context, mid, assistMid int64) (row int64, err error) {
	res, err := d.db.Exec(c, _delAssSQL, mid, assistMid)
	if err != nil {
		log.Error("d.delAss.Exec error(%v)", err)
		return
	}
	row, err = res.RowsAffected()
	return
}

// Assist get one Assist from assist database.
func (d *Dao) Assist(c context.Context, mid, assistMid int64) (a *assist.Assist, err error) {
	row := d.db.QueryRow(c, _assSQL, mid, assistMid)
	a = &assist.Assist{}
	if err = row.Scan(&a.Mid, &a.AssistMid, &a.CTime, &a.MTime); err != nil {
		if err == sql.ErrNoRows {
			a = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// Assists get all Assists from assist database.
func (d *Dao) Assists(c context.Context, mid int64) (as []*assist.Assist, err error) {
	as = make([]*assist.Assist, 0)
	rows, err := d.db.Query(c, _asssSQL, mid)
	if err != nil {
		log.Error("db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := &assist.Assist{}
		if err = rows.Scan(&a.Mid, &a.AssistMid, &a.CTime, &a.MTime); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		as = append(as, a)
	}
	return
}

// AssistCnt get assist count.
func (d *Dao) AssistCnt(c context.Context, mid int64) (count int, err error) {
	row := d.db.QueryRow(c, _assCntSQL, mid)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			count = 0
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// Ups get ups who already sign me as assist.
func (d *Dao) Ups(c context.Context, assistMid, pn, ps int64) (mids []int64, ups map[int64]*assist.Up, total int64, err error) {
	mids = make([]int64, 0)
	ups = make(map[int64]*assist.Up, ps)
	// default is empty json array
	rows, err := d.db.Query(c, _assUpsSQL, assistMid, (pn-1)*ps, pn*ps)
	if err != nil {
		log.Error("db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		up := &assist.Up{}
		if err = rows.Scan(&up.Mid, &up.CTime); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		ups[up.Mid] = up
		mids = append(mids, up.Mid)
	}
	row := d.db.QueryRow(c, _assUpsCntSQL, assistMid)
	if err = row.Scan(&total); err != nil {
		if err == sql.ErrNoRows {
			total = 0
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}
