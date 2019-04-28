package dao

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-common/app/job/main/up-rating/model"

	"go-common/library/log"
)

const (
	// get base info from up_level_info
	_baseInfoSQL  = "SELECT id,mid,tag_id,inc_play,inc_coin,avs,maafans,mahfans,open_avs,lock_avs,cdate FROM up_level_info_%02d WHERE id > ? AND id <= ? LIMIT ?"
	_baseTotalSQL = "SELECT id,mid,total_fans,total_avs,total_coin,total_play FROM up_level_info_%02d WHERE id > ? AND cdate = '%s' LIMIT ?"

	// get up_level_info start & end
	_baseInfoStartSQL = "SELECT id FROM up_level_info_%02d WHERE cdate=? ORDER BY id LIMIT 1"
	_baseInfoEndSQL   = "SELECT id FROM up_level_info_%02d WHERE cdate=? ORDER BY id DESC LIMIT 1"
)

// GetBaseInfo get rating info
func (d *Dao) GetBaseInfo(c context.Context, month time.Month, start, end, limit int) (bs []*model.BaseInfo, id int, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_baseInfoSQL, month), start, end, limit)
	if err != nil {
		log.Error("d.db.Query Rating Base Info error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		b := &model.BaseInfo{}
		err = rows.Scan(&id, &b.MID, &b.TagID, &b.PlayIncr, &b.CoinIncr, &b.Avs, &b.MAAFans, &b.MAHFans, &b.OpenAvs, &b.LockedAvs, &b.Date)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		bs = append(bs, b)
	}
	return
}

// GetBaseTotal get total
func (d *Dao) GetBaseTotal(c context.Context, date time.Time, id, limit int64) (bs []*model.BaseInfo, err error) {
	bs = make([]*model.BaseInfo, 0)
	rows, err := d.db.Query(c, fmt.Sprintf(_baseTotalSQL, date.Month(), date.Format("2006-01-02")), id, limit)
	if err != nil {
		log.Error("d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		b := &model.BaseInfo{}
		err = rows.Scan(&b.ID, &b.MID, &b.TotalFans, &b.TotalAvs, &b.TotalCoin, &b.TotalPlay)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		bs = append(bs, b)
	}
	return
}

// BaseInfoStart get start id by date
func (d *Dao) BaseInfoStart(c context.Context, date time.Time) (start int, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_baseInfoStartSQL, date.Month()), date)
	if err = row.Scan(&start); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
	}
	return
}

// BaseInfoEnd get end id by date
func (d *Dao) BaseInfoEnd(c context.Context, date time.Time) (end int, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_baseInfoEndSQL, date.Month()), date)
	if err = row.Scan(&end); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
	}
	return
}
