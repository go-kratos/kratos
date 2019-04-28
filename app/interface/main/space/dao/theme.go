package dao

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-common/app/interface/main/space/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_themeKeyFmt  = `spc_them_%d`
	_themeSQL     = `SELECT sid,is_activated FROM dede_member_skin%d WHERE mid = ? AND expire > ?`
	_themeInfoSQL = `SELECT id,name,img_path,toutu,bgimg FROM dede_skin_mall WHERE is_disable = 0 AND id IN (%s)`
	_themeEditSQL = `UPDATE dede_member_skin%d SET is_activated = 1 WHERE mid = ? AND sid = ?`
	_themeUnSQL   = `UPDATE dede_member_skin%d SET is_activated = 0 WHERE mid = ?`
)

func themeHit(mid int64) int64 {
	return mid % 10
}

func themeKey(mid int64) string {
	return fmt.Sprintf(_themeKeyFmt, mid)
}

// ThemeInfoByMid get theme info by mid.
func (d *Dao) ThemeInfoByMid(c context.Context, mid int64) (res []*model.ThemeInfo, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_themeSQL, themeHit(mid)), mid, time.Now()); err != nil {
		log.Error("ThemeInfoByMid d.db.Query(%d) error(%v)", mid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.ThemeInfo)
		if err = rows.Scan(&r.SID, &r.IsActivated); err != nil {
			log.Error("ThemeInfoByMid row.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	return
}

// ThemeDetail get theme details.
func (d *Dao) ThemeDetail(c context.Context, themeIDs []int64) (res []*model.ThemeDetail, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_themeInfoSQL, xstr.JoinInts(themeIDs))); err != nil {
		log.Error("ThemeDetail d.db.Query(%v) error(%v)", themeIDs, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.ThemeDetail)
		if err = rows.Scan(&r.ID, &r.Name, &r.Icon, &r.TopPhoto, &r.BgImg); err != nil {
			log.Error("ThemeDetail row.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	return
}

// RawTheme get themes by mid.
func (d *Dao) RawTheme(c context.Context, mid int64) (res *model.ThemeDetails, err error) {
	var (
		themeInfo []*model.ThemeInfo
		themeIDs  []int64
		list      []*model.ThemeDetail
	)
	res = new(model.ThemeDetails)
	if themeInfo, err = d.ThemeInfoByMid(c, mid); err != nil || len(themeInfo) == 0 {
		return
	}
	themeInfoMap := make(map[int64]*model.ThemeInfo, len(themeInfo))
	for _, v := range themeInfo {
		themeIDs = append(themeIDs, v.SID)
		themeInfoMap[v.SID] = v
	}
	if list, err = d.ThemeDetail(c, themeIDs); err != nil {
		return
	}
	for _, v := range list {
		if theme, ok := themeInfoMap[v.ID]; ok && theme != nil {
			v.IsActivated = theme.IsActivated
		}
	}
	res.List = list
	return
}

// ThemeActive active theme.
func (d *Dao) ThemeActive(c context.Context, mid, themeID int64) (err error) {
	var (
		res sql.Result
		tx  *xsql.Tx
	)
	if tx, err = d.db.Begin(c); err != nil {
		log.Error("ThemeActive: d.db.Begin error(%v)", err)
		return
	}
	if res, err = tx.Exec(fmt.Sprintf(_themeUnSQL, themeHit(mid)), mid); err != nil {
		tx.Rollback()
		log.Error("ThemeActive: db.Exec(%d) error(%v)", mid, err)
		return
	}
	if _, err = tx.Exec(fmt.Sprintf(_themeEditSQL, themeHit(mid)), mid, themeID); err != nil {
		tx.Rollback()
		log.Error("ThemeActive: db.Exec(%d,%d) error(%v)", mid, themeID, err)
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("ThemeActive: tx.Commit error(%v)", err)
		return
	}
	_, err = res.RowsAffected()
	return
}
