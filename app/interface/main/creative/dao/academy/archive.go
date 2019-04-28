package academy

import (
	"context"
	"fmt"

	"database/sql"
	sqlx "go-common/library/database/sql"
	"go-common/library/xstr"

	"go-common/app/interface/main/creative/dao/tool"
	"go-common/app/interface/main/creative/model/academy"
	"go-common/library/log"
)

const (
	// select
	_getArcByOIDAndBusSQL = "SELECT id, oid, business, state, ctime, mtime FROM academy_archive WHERE oid=? AND business=?"
	_getSearByBusSQL      = "SELECT a.oid, a.business, GROUP_CONCAT(t.tid SEPARATOR ',') AS tidstr FROM academy_archive AS a LEFT JOIN academy_archive_tag as t on t.oid = a.oid"
	_getArcCountSQL       = "SELECT count(DISTINCT a.oid) FROM (SELECT oid FROM academy_archive  WHERE state=0 AND business=?) AS a LEFT JOIN academy_archive_tag as t on t.oid=a.oid"
	_getTagByOidsSQL      = "SELECT oid, tid FROM academy_archive_tag WHERE state=0 AND oid IN (%s)"
)

//Archive get one achive.
func (d *Dao) Archive(c context.Context, oid int64, bs int) (a *academy.Archive, err error) {
	row := d.db.QueryRow(c, _getArcByOIDAndBusSQL, oid, bs)
	a = &academy.Archive{}
	if err = row.Scan(&a.ID, &a.OID, &a.Business, &a.State, &a.CTime, &a.MTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("row.Scan error(%v)", err)
	}
	return
}

//ArchiveCount get all achive count.
func (d *Dao) ArchiveCount(c context.Context, tids []int64, bs int) (count int, err error) {
	sqlStr := _getArcCountSQL
	if len(tids) > 0 {
		sqlStr += fmt.Sprintf(" WHERE t.tid IN (%s)", xstr.JoinInts(tids))
	}
	if err = d.db.QueryRow(c, sqlStr, bs).Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("d.db.QueryRow error(%v)", err)
	}
	return
}

//SearchArchive get all oid & tid.
func (d *Dao) SearchArchive(c context.Context, tidsMap map[int][]int64, bs int) (res []*academy.Archive, err error) {
	var (
		rows *sqlx.Rows
		tids []int64
	)
	for _, v := range tidsMap {
		tids = append(tids, v...)
	}
	total := len(tids)
	sqlStr := _getSearByBusSQL
	if total > 0 {
		sqlStr += fmt.Sprintf(" WHERE a.state=0 AND a.business=? AND t.tid IN (%s) GROUP BY a.oid  ORDER BY a.mtime DESC", xstr.JoinInts(tids))
	} else {
		sqlStr += " WHERE a.state=0 AND a.business=? GROUP BY a.oid  ORDER BY a.mtime DESC"
	}
	rows, err = d.db.Query(c, sqlStr, bs)
	if err != nil {
		log.Error("d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	res = make([]*academy.Archive, 0)
	origin := make([]*academy.Archive, 0)
	var tidStr string
	for rows.Next() {
		a := &academy.Archive{}
		if err = rows.Scan(&a.OID, &a.Business, &tidStr); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		a.TIDs, _ = xstr.SplitInts(tidStr)
		origin = append(origin, a)
	}
	if total > 0 {
		var cts, ots, clts, acts []int64
		if v, ok := tidsMap[academy.Course]; ok {
			cts = v
		}
		if v, ok := tidsMap[academy.Operation]; ok {
			ots = v
		}
		if v, ok := tidsMap[academy.Classify]; ok {
			clts = v
		}
		if v, ok := tidsMap[academy.ArticleClass]; ok {
			acts = v
		}
		for _, v := range origin {
			log.Info("search tag 课程级别(%+v)|运营标签(%+v)|分类标签(%+v)|专栏分类(%+v)|当前稿件标签(%+v)", cts, ots, clts, acts, v.TIDs)
			if tool.ContainAtLeastOne(cts, v.TIDs) &&
				tool.ContainAtLeastOne(ots, v.TIDs) &&
				tool.ContainAtLeastOne(clts, v.TIDs) &&
				tool.ContainAtLeastOne(acts, v.TIDs) {
				res = append(res, v)
			}
		}
	} else {
		res = origin
	}
	return
}

//ArchiveTagsByOids get all tids by oids.
func (d *Dao) ArchiveTagsByOids(c context.Context, oids []int64) (res map[int64][]int64, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_getTagByOidsSQL, xstr.JoinInts(oids)))
	if err != nil {
		log.Error("d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	res = make(map[int64][]int64)
	for rows.Next() {
		a := &academy.ArchiveTag{}
		if err = rows.Scan(&a.OID, &a.TID); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		res[a.OID] = append(res[a.OID], a.TID)
	}
	return
}
