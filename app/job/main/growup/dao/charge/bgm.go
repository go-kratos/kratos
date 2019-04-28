package charge

import (
	"context"
	"fmt"
	"time"

	model "go-common/app/job/main/growup/model/charge"

	"go-common/library/log"
)

const (
	_getBGMSQL    = "SELECT id,mid,sid,aid,cid,join_at,title FROM background_music WHERE id > ? ORDER BY id LIMIT ?"
	_bgmChargeSQL = "SELECT id,mid,sid,aid,cid,join_at,inc_charge,date FROM %s WHERE id > ? AND date = ? AND inc_charge > 0 ORDER BY id LIMIT ?"
	_bgmStatisSQL = "SELECT id,mid,sid,aid,cid,total_charge FROM bgm_charge_statis WHERE id > ? ORDER BY id LIMIT ?"

	_inBgmChargeTableSQL = "INSERT INTO %s(sid,aid,mid,cid,title,inc_charge,date,join_at) VALUES %s ON DUPLICATE KEY UPDATE inc_charge=VALUES(inc_charge)"
	_inBgmStatisSQL      = "INSERT INTO bgm_charge_statis(sid,aid,mid,cid,title,total_charge,join_at) VALUES %s ON DUPLICATE KEY UPDATE total_charge=VALUES(total_charge)"
)

// GetBgm get bgms
func (d *Dao) GetBgm(c context.Context, id int64, limit int64) (bs []*model.Bgm, last int64, err error) {
	rows, err := d.db.Query(c, _getBGMSQL, id, limit)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		b := &model.Bgm{}
		err = rows.Scan(&last, &b.MID, &b.SID, &b.AID, &b.CID, &b.JoinAt, &b.Title)
		if err != nil {
			return
		}
		bs = append(bs, b)
	}
	return
}

// BgmCharge get bgm charge by date
func (d *Dao) BgmCharge(c context.Context, date time.Time, id int64, limit int, table string) (bgms []*model.BgmCharge, err error) {
	bgms = make([]*model.BgmCharge, 0)
	rows, err := d.db.Query(c, fmt.Sprintf(_bgmChargeSQL, table), id, date, limit)
	if err != nil {
		log.Error("BgmCharge d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		bgm := &model.BgmCharge{}
		err = rows.Scan(&bgm.ID, &bgm.MID, &bgm.SID, &bgm.AID, &bgm.CID, &bgm.JoinAt, &bgm.IncCharge, &bgm.Date)
		if err != nil {
			log.Error("BgmCharge rows.Scan error(%v)", err)
			return
		}
		bgms = append(bgms, bgm)
	}
	return
}

// BgmStatis bgm statis
func (d *Dao) BgmStatis(c context.Context, id int64, limit int) (bgms []*model.BgmStatis, err error) {
	bgms = make([]*model.BgmStatis, 0)
	rows, err := d.db.Query(c, _bgmStatisSQL, id, limit)
	if err != nil {
		log.Error("BgmStatis d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		bgm := &model.BgmStatis{}
		err = rows.Scan(&bgm.ID, &bgm.MID, &bgm.SID, &bgm.AID, &bgm.CID, &bgm.TotalCharge)
		if err != nil {
			log.Error("BgmStatis rows.Scan error(%v)", err)
			return
		}
		bgms = append(bgms, bgm)
	}
	return
}

// InsertBgmChargeTable insert bgm charge
func (d *Dao) InsertBgmChargeTable(c context.Context, vals, table string) (rows int64, err error) {
	if vals == "" {
		return
	}
	res, err := d.db.Exec(c, fmt.Sprintf(_inBgmChargeTableSQL, table, vals))
	if err != nil {
		log.Error("InsertBgmChargeTable(%s) tx.Exec error(%v)", table, err)
		return
	}
	return res.RowsAffected()
}

// InsertBgmStatisBatch insert bgm statis
func (d *Dao) InsertBgmStatisBatch(c context.Context, vals string) (rows int64, err error) {
	if vals == "" {
		return
	}
	res, err := d.db.Exec(c, fmt.Sprintf(_inBgmStatisSQL, vals))
	if err != nil {
		log.Error("InsertBgmStatisBatch tx.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}
