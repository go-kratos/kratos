package income

import (
	"context"
	"fmt"

	model "go-common/app/admin/main/growup/model/income"
	"go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"
)

const (
	// select
	_avDailyChargeSQL  = "SELECT av_id,mid,tag_id,upload_time,inc_charge,date FROM av_daily_charge_%02d WHERE av_id = ? LIMIT 31"
	_cmDailyChargeSQL  = "SELECT aid,mid,tag_id,upload_time,inc_charge,date FROM column_daily_charge WHERE aid = ?"
	_bgmDailyChargeSQL = "SELECT sid,aid,mid,join_at,inc_charge,date FROM bgm_daily_charge WHERE sid = ?"
	_upChargeRatioSQL  = "SELECT mid, ratio FROM up_charge_ratio LIMIT ?, ?"

	_archiveChargeStatisTableSQL = "SELECT avs,money_section,money_tips,charge,category_id,cdate FROM %s WHERE %s LIMIT ?,?"
	_archiveTotalChargeSQL       = "SELECT total_charge FROM %s WHERE %s LIMIT 1"
)

// GetAvDailyCharge get av_daily_charge by month
func (d *Dao) GetAvDailyCharge(c context.Context, month int, avID int64) (avs []*model.ArchiveCharge, err error) {
	if month < 1 || month > 12 {
		return nil, fmt.Errorf("error args month(%d)", month)
	}
	rows, err := d.db.Query(c, fmt.Sprintf(_avDailyChargeSQL, month), avID)
	if err != nil {
		log.Error("GetAvDailyCharge d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		list := &model.ArchiveCharge{}
		err = rows.Scan(&list.AID, &list.MID, &list.CategoryID, &list.UploadTime, &list.Charge, &list.Date)
		if err != nil {
			log.Error("GetAvDailyCharge rows scan error(%v)", err)
			return
		}
		avs = append(avs, list)
	}
	err = rows.Err()
	return
}

// GetColumnCharges get column daily charge
func (d *Dao) GetColumnCharges(c context.Context, aid int64) (cms []*model.ArchiveCharge, err error) {
	cms = make([]*model.ArchiveCharge, 0)
	rows, err := d.db.Query(c, _cmDailyChargeSQL, aid)
	if err != nil {
		log.Error("GetColumnCharge d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var uploadTime int64
		cm := &model.ArchiveCharge{}
		err = rows.Scan(&cm.AID, &cm.MID, &cm.CategoryID, &uploadTime, &cm.Charge, &cm.Date)
		if err != nil {
			log.Error("GetColumnCharge rows scan error(%v)", err)
			return
		}
		cm.UploadTime = xtime.Time(uploadTime)
		cms = append(cms, cm)
	}
	err = rows.Err()
	return
}

// GetBgmCharges get bgm daily charge
func (d *Dao) GetBgmCharges(c context.Context, aid int64) (bgms []*model.ArchiveCharge, err error) {
	bgms = make([]*model.ArchiveCharge, 0)
	rows, err := d.db.Query(c, _bgmDailyChargeSQL, aid)
	if err != nil {
		log.Error("GetBgmCharge d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		bgm := &model.ArchiveCharge{}
		err = rows.Scan(&bgm.AID, &bgm.AvID, &bgm.MID, &bgm.UploadTime, &bgm.Charge, &bgm.Date)
		if err != nil {
			log.Error("GetBgmCharge rows scan error(%v)", err)
			return
		}
		bgms = append(bgms, bgm)
	}
	err = rows.Err()
	return
}

// GetArchiveChargeStatis get archive charge statis from table and query
func (d *Dao) GetArchiveChargeStatis(c context.Context, table, query string, from, limit int) (archs []*model.ArchiveChargeStatis, err error) {
	archs = make([]*model.ArchiveChargeStatis, 0)
	if table == "" || query == "" {
		return nil, fmt.Errorf("error args table(%s), query(%s)", table, query)
	}
	rows, err := d.db.Query(c, fmt.Sprintf(_archiveChargeStatisTableSQL, table, query), from, limit)
	if err != nil {
		log.Error("GetArchiveChargeStatis d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		list := &model.ArchiveChargeStatis{}
		err = rows.Scan(&list.Avs, &list.MoneySection, &list.MoneyTips, &list.Charge, &list.CategroyID, &list.CDate)
		if err != nil {
			log.Error("GetArchiveChargeStatis rows scan error(%v)", err)
			return
		}
		archs = append(archs, list)
	}

	err = rows.Err()
	return
}

// GetTotalCharge get total charge by table and aid
func (d *Dao) GetTotalCharge(c context.Context, table, query string) (total int64, err error) {
	err = d.db.QueryRow(c, fmt.Sprintf(_archiveTotalChargeSQL, table, query)).Scan(&total)
	if err == sql.ErrNoRows {
		err = nil
	}
	return
}

// UpRatio get up charge ratio
func (d *Dao) UpRatio(c context.Context, from, limit int64) (ratio map[int64]int64, err error) {
	ratio = make(map[int64]int64)
	rows, err := d.db.Query(c, _upChargeRatioSQL, from, limit)
	if err != nil {
		log.Error("d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var mid, charge int64
		err = rows.Scan(&mid, &charge)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		ratio[mid] = charge
	}
	err = rows.Err()
	return
}
