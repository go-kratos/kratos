package dao

import (
	"context"
	"database/sql"
	"fmt"

	"go-common/app/service/main/rank/model"
	"go-common/library/log"
	xtime "go-common/library/time"
	"go-common/library/xstr"
)

const (
	_maxArchiveIDSQL      = `SELECT MAX(id) FROM archive`
	_archiveMetasSQL      = `SELECT id,typeid,pubtime FROM archive WHERE id>? ORDER BY id LIMIT ?`
	_archiveMetasMtimeSQL = `SELECT id,typeid,pubtime FROM archive WHERE id>? AND mtime BETWEEN ? AND ? ORDER BY mtime,id LIMIT ?`
	_archiveTypesSQL      = `SELECT id,pid FROM archive_type WHERE id in (%s)`
	_archiveStatsSQL      = `SELECT aid,click FROM archive_stat_%s WHERE aid in (%s)`
	_archiveStatsMtimeSQL = `SELECT id,aid,click FROM archive_stat_%s WHERE id > ? AND mtime BETWEEN ? AND ? ORDER BY mtime,id LIMIT ?`
	_archiveTVsSQL        = `SELECT aid,result,deleted,valid FROM ugc_archive WHERE aid in (%s)`
	_archiveTVsMtimeSQL   = `SELECT id,aid,result,deleted,valid FROM ugc_archive WHERE id > ? AND mtime BETWEEN ? AND ? ORDER BY mtime,id LIMIT ?`

	_archiveStatSharding = 100
)

// MaxOid .
func (d *Dao) MaxOid(c context.Context) (oid int64, err error) {
	row := d.dbArchive.QueryRow(c, _maxArchiveIDSQL)
	if err = row.Scan(&oid); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// ArchiveMetas .
func (d *Dao) ArchiveMetas(c context.Context, id int64, limit int) ([]*model.ArchiveMeta, error) {
	rows, err := d.dbArchive.Query(c, _archiveMetasSQL, id, limit)
	if err != nil {
		log.Error("d.dbArchive.Query(%s,%d,%d) error()", _archiveMetasSQL, id, limit, err)
		return nil, err
	}
	defer rows.Close()
	as := make([]*model.ArchiveMeta, 0)
	for rows.Next() {
		a := new(model.ArchiveMeta)
		if err = rows.Scan(&a.ID, &a.Typeid, &a.Pubtime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return nil, err
		}
		as = append(as, a)
	}
	return as, rows.Err()
}

// ArchiveMetasIncrs .
func (d *Dao) ArchiveMetasIncrs(c context.Context, aid int64, begin, end xtime.Time, limit int) ([]*model.ArchiveMeta, error) {
	rows, err := d.dbArchive.Query(c, _archiveMetasMtimeSQL, aid, begin, end, limit)
	if err != nil {
		log.Error("d.dbArchive.Query(%s,%d,%s,%s,%d) error()", _archiveMetasMtimeSQL, aid, begin, end, limit, err)
		return nil, err
	}
	defer rows.Close()
	as := make([]*model.ArchiveMeta, 0)
	for rows.Next() {
		a := new(model.ArchiveMeta)
		if err = rows.Scan(&a.ID, &a.Typeid, &a.Pubtime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return nil, err
		}
		as = append(as, a)
	}
	return as, rows.Err()
}

// ArchiveTypes .
func (d *Dao) ArchiveTypes(c context.Context, ids []int64) (map[int64]*model.ArchiveType, error) {
	idsStr := xstr.JoinInts(ids)
	rows, err := d.dbArchive.Query(c, fmt.Sprintf(_archiveTypesSQL, idsStr))
	if err != nil {
		log.Error("d.dbArchive.Query(%s) error(%v)", fmt.Sprintf(_archiveTypesSQL, idsStr), err)
		return nil, err
	}
	defer rows.Close()
	as := make(map[int64]*model.ArchiveType)
	for rows.Next() {
		a := new(model.ArchiveType)
		if err = rows.Scan(&a.ID, &a.Pid); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return nil, err
		}
		as[a.ID] = a
	}
	return as, rows.Err()
}

// ArchiveStats .
func (d *Dao) ArchiveStats(c context.Context, aids []int64) (map[int64]*model.ArchiveStat, error) {
	tableMap := make(map[int64][]int64)
	for _, aid := range aids {
		mod := aid % _archiveStatSharding
		tableMap[mod] = append(tableMap[mod], aid)
	}
	as := make(map[int64]*model.ArchiveStat)
	for tbl, aids := range tableMap {
		aidsStr := xstr.JoinInts(aids)
		rows, err := d.dbStat.Query(c, fmt.Sprintf(_archiveStatsSQL, fmt.Sprintf("%02d", tbl), aidsStr))
		if err != nil {
			log.Error("d.dbStat.Query(%s) error(%v)", fmt.Sprintf(_archiveStatsSQL, fmt.Sprintf("%02d", tbl), aidsStr), err)
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			a := new(model.ArchiveStat)
			if err = rows.Scan(&a.Aid, &a.Click); err != nil {
				log.Error("rows.Scan() error(%v)", err)
				return nil, err
			}
			as[a.Aid] = a
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}
	}
	return as, nil
}

// ArchiveStatsIncrs .
func (d *Dao) ArchiveStatsIncrs(c context.Context, tbl int, id int64, begin, end xtime.Time, limit int) ([]*model.ArchiveStat, error) {
	rows, err := d.dbStat.Query(c, fmt.Sprintf(_archiveStatsMtimeSQL, fmt.Sprintf("%02d", tbl)), id, begin, end, limit)
	if err != nil {
		log.Error("d.dbStat.Query(%s,%d,%s,%s,%d) error(%v)", fmt.Sprintf(_archiveStatsMtimeSQL, fmt.Sprintf("%02d", tbl)), id, begin, end, limit, err)
		return nil, err
	}
	defer rows.Close()
	as := make([]*model.ArchiveStat, 0)
	for rows.Next() {
		a := new(model.ArchiveStat)
		if err = rows.Scan(&a.ID, &a.Aid, &a.Click); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return nil, err
		}
		as = append(as, a)
	}
	return as, rows.Err()
}

// ArchiveTVs .
func (d *Dao) ArchiveTVs(c context.Context, aids []int64) (map[int64]*model.ArchiveTv, error) {
	aidsStr := xstr.JoinInts(aids)
	rows, err := d.dbTV.Query(c, fmt.Sprintf(_archiveTVsSQL, aidsStr))
	if err != nil {
		log.Error("d.dbTV.Query(%s) error(%v)", fmt.Sprintf(_archiveTVsSQL, aidsStr), err)
		return nil, err
	}
	defer rows.Close()
	as := make(map[int64]*model.ArchiveTv)
	for rows.Next() {
		a := new(model.ArchiveTv)
		if err = rows.Scan(&a.Aid, &a.Result, &a.Deleted, &a.Valid); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return nil, err
		}
		as[a.Aid] = a
	}
	return as, rows.Err()
}

// ArchiveTVsIncrs .
func (d *Dao) ArchiveTVsIncrs(c context.Context, id int64, begin, end xtime.Time, limit int) ([]*model.ArchiveTv, error) {
	rows, err := d.dbTV.Query(c, _archiveTVsMtimeSQL, id, begin, end, limit)
	if err != nil {
		log.Error("d.dbTV.Query(%s,%d,%s,%s,%d) error(%v)", _archiveTVsMtimeSQL, id, begin, end, limit, err)
		return nil, err
	}
	defer rows.Close()
	as := make([]*model.ArchiveTv, 0)
	for rows.Next() {
		a := new(model.ArchiveTv)
		if err = rows.Scan(&a.ID, &a.Aid, &a.Result, &a.Deleted, &a.Valid); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return nil, err
		}
		as = append(as, a)
	}
	return as, rows.Err()
}
