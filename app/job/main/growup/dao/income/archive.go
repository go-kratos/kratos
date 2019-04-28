package income

import (
	"context"
	"fmt"

	model "go-common/app/job/main/growup/model/income"

	"go-common/library/log"
)

const (
	_getArchiveByDateSQL   = "SELECT id, %s, mid, tag_id, income, date FROM %s WHERE id > ? AND date = ? ORDER BY id LIMIT ?"
	_getBgmIncomeByDateSQL = "SELECT id, sid, mid, income, date FROM bgm_income WHERE id > ? AND date = ? ORDER BY id LIMIT ?"
)

// GetArchiveByDate get archive by date
func (d *Dao) GetArchiveByDate(c context.Context, aid, table, date string, id int64, limit int) (archives []*model.ArchiveIncome, err error) {
	archives = make([]*model.ArchiveIncome, 0)
	rows, err := d.db.Query(c, fmt.Sprintf(_getArchiveByDateSQL, aid, table), id, date, limit)
	if err != nil {
		log.Error("GetArchiveByDate d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		archive := &model.ArchiveIncome{}
		err = rows.Scan(&archive.ID, &archive.AID, &archive.MID, &archive.TagID, &archive.Income, &archive.Date)
		if err != nil {
			log.Error("GetArchiveByDate rows.Scan error(%v)", err)
			return
		}
		archives = append(archives, archive)
	}
	return
}

// GetBgmIncomeByDate get bgm income by date
func (d *Dao) GetBgmIncomeByDate(c context.Context, date string, id int64, limit int) (archives []*model.ArchiveIncome, err error) {
	archives = make([]*model.ArchiveIncome, 0)
	rows, err := d.db.Query(c, _getBgmIncomeByDateSQL, id, date, limit)
	if err != nil {
		log.Error("GetArchiveByDate d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		archive := &model.ArchiveIncome{}
		err = rows.Scan(&archive.ID, &archive.AID, &archive.MID, &archive.Income, &archive.Date)
		if err != nil {
			log.Error("GetArchiveByDate rows.Scan error(%v)", err)
			return
		}
		archives = append(archives, archive)
	}
	return

}
