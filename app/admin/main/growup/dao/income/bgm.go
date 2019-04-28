package income

import (
	"context"
	"fmt"

	model "go-common/app/admin/main/growup/model/income"

	"go-common/library/log"
)

const (
	// select
	_getAvByBgmIncome = "SELECT aid FROM bgm_income WHERE sid = ? AND date >= ? AND date <= ?"
	_bgmIncomeSQL     = "SELECT id,sid,mid,join_at,total_income,income,tax_money,date FROM bgm_income WHERE id > ? AND date >= ? AND date <= ? AND %s is_deleted = 0 ORDER BY id LIMIT ?"
)

// GetAvByBgm get av_id by bgm id
func (d *Dao) GetAvByBgm(c context.Context, sid int64, from, to string) (avs map[int64]struct{}, err error) {
	avs = make(map[int64]struct{})
	rows, err := d.db.Query(c, _getAvByBgmIncome, sid, from, to)
	if err != nil {
		log.Error("GetAvByBgm d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var avID int64
		err = rows.Scan(&avID)
		if err != nil {
			log.Error("GetAvByBgm rows scan error(%v)", err)
			return
		}
		avs[avID] = struct{}{}
	}
	err = rows.Err()
	return
}

// GetBgmIncome get bgm income by query
func (d *Dao) GetBgmIncome(c context.Context, id int64, query string, from, to string, limit int, typ int) (bgms []*model.ArchiveIncome, err error) {
	bgms = make([]*model.ArchiveIncome, 0)
	if query != "" {
		query += " AND"
	}
	rows, err := d.db.Query(c, fmt.Sprintf(_bgmIncomeSQL, query), id, from, to, limit)
	if err != nil {
		log.Error("GetBgmIncome d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		list := &model.ArchiveIncome{}
		err = rows.Scan(&list.ID, &list.AvID, &list.MID, &list.UploadTime, &list.TotalIncome, &list.Income, &list.TaxMoney, &list.Date)
		if err != nil {
			log.Error("GetBgmIncome rows scan error(%v)", err)
			return
		}
		list.Type = typ
		bgms = append(bgms, list)
	}

	err = rows.Err()
	return
}
