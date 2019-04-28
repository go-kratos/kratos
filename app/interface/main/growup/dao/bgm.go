package dao

import (
	"context"
	"fmt"

	"go-common/app/interface/main/growup/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_bgmCountSQL = "SELECT count(distinct sid) FROM background_music WHERE mid=?"
	_bgmTitleSQL = "SELECT sid,title FROM background_music WHERE sid IN (%s)"

	_bgmWhiteSQL = "SELECT COUNT(*) FROM bgm_white_list WHERE mid = ?"

	_bgmIncomeByMIDSQL = "SELECT sid,income,total_income,date FROM bgm_income WHERE mid = ? AND date >= ? AND date <= ?"
	_bgmIncomeBySIDSQL = "SELECT aid,income,date FROM bgm_income WHERE sid = ? AND date <= ?"
)

// BgmWhiteList bgm_white_list
func (d *Dao) BgmWhiteList(c context.Context, mid int64) (count int, err error) {
	row := d.db.QueryRow(c, _bgmWhiteSQL, mid)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row scan error(%v)", err)
		}
	}
	return
}

// BGMCount get up bgm count
func (d *Dao) BGMCount(c context.Context, mid int64) (count int, err error) {
	row := d.db.QueryRow(c, _bgmCountSQL, mid)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row scan error(%v)", err)
		}
	}
	return
}

// GetBgmTitle get bgm titles
func (d *Dao) GetBgmTitle(c context.Context, ids []int64) (titles map[int64]string, err error) {
	titles = make(map[int64]string)
	rows, err := d.db.Query(c, fmt.Sprintf(_bgmTitleSQL, xstr.JoinInts(ids)))
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var sid int64
		var title string
		err = rows.Scan(&sid, &title)
		if err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		titles[sid] = title
	}
	return
}

// ListBgmIncome list bgm income
func (d *Dao) ListBgmIncome(c context.Context, mid int64, startTime, endTime string) (bgms []*model.ArchiveIncome, err error) {
	bgms = make([]*model.ArchiveIncome, 0)
	rows, err := d.db.Query(c, _bgmIncomeByMIDSQL, mid, startTime, endTime)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		bgm := &model.ArchiveIncome{}
		err = rows.Scan(&bgm.ArchiveID, &bgm.Income, &bgm.TotalIncome, &bgm.Date)
		if err != nil {
			log.Error("ListColumnIncome rows.Scan error(%v)", err)
			return
		}
		bgms = append(bgms, bgm)
	}
	err = rows.Err()
	return
}

// ListBgmIncomeByID list bgm_income by sid
func (d *Dao) ListBgmIncomeByID(c context.Context, id int64, endTime string) (bgms []*model.ArchiveIncome, err error) {
	bgms = make([]*model.ArchiveIncome, 0)
	rows, err := d.db.Query(c, _bgmIncomeBySIDSQL, id, endTime)
	if err != nil {
		log.Error("ListBgmIncomeByID d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		bgm := &model.ArchiveIncome{}
		err = rows.Scan(&bgm.ArchiveID, &bgm.Income, &bgm.Date)
		if err != nil {
			log.Error("ListBgmIncomeByID rows.Scan error(%v)", err)
			return
		}
		bgms = append(bgms, bgm)
	}
	err = rows.Err()
	return
}
