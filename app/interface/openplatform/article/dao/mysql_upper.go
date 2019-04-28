package dao

import (
	"context"
	"fmt"

	"go-common/app/interface/openplatform/article/model"
	"go-common/library/log"
	"go-common/library/xstr"
)

// UpperPassed upper passed articles
func (d *Dao) UpperPassed(c context.Context, mid int64) (aids [][2]int64, err error) {
	rows, err := d.upPassedStmt.Query(c, mid)
	if err != nil {
		PromError("db:up文章列表")
		log.Error("getUpPasStmt.Query(%d) error(%+v)", mid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			aid, ptime int64
			attributes int32
		)
		if err = rows.Scan(&aid, &ptime, &attributes); err != nil {
			log.Error("rows.Scan error(%+v)", err)
			return
		}
		if !model.NoDistributeAttr(attributes) {
			aids = append(aids, [2]int64{aid, ptime})
		}
	}
	err = rows.Err()
	promErrorCheck(err)
	return
}

// UppersPassed uppers passed articles
func (d *Dao) UppersPassed(c context.Context, mids []int64) (aidm map[int64][][2]int64, err error) {
	rows, err := d.articleDB.Query(c, fmt.Sprintf(_uppersPassedSQL, xstr.JoinInts(mids)))
	if err != nil {
		PromError("db:批量查询up文章列表")
		log.Error("UpsPassed error(%+v)", err)
		return
	}
	defer rows.Close()
	aidm = make(map[int64][][2]int64, len(mids))
	for rows.Next() {
		var (
			aid, mid, ptime int64
			attributes      int32
		)
		if err = rows.Scan(&aid, &mid, &ptime, &attributes); err != nil {
			log.Error("rows.Scan error(%+v)", err)
			return
		}
		if !model.NoDistributeAttr(attributes) {
			aidm[mid] = append(aidm[mid], [2]int64{aid, ptime})
		}
	}
	for _, mid := range mids {
		if aidm[mid] == nil {
			aidm[mid] = [][2]int64{}
		}
	}
	err = rows.Err()
	promErrorCheck(err)
	return
}
