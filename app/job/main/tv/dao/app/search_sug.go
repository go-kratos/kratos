package app

import (
	"context"

	model "go-common/app/job/main/tv/model/pgc"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_pgcSeaSug = "SELECT id,title FROM tv_ep_season WHERE `check` = ? AND valid = ? AND is_deleted = ?"
	_ugcSeaSug = "SELECT aid,title FROM ugc_archive WHERE result=1 AND valid=1 AND deleted=0 "
	_pgcType   = "pgc"
	_ugcType   = "ugc"
)

// PgcSeaSug is used for getting pgc search suggest content
func (d *Dao) PgcSeaSug(ctx context.Context) (res []*model.SearchSug, err error) {
	var (
		rows *sql.Rows
	)
	if rows, err = d.DB.Query(ctx, _pgcSeaSug, SeasonPassed, _CMSValid, _NotDeleted); err != nil {
		log.Error("d.PgcSeaSug.Query: %s error(%v)", _pgcSeaSug, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r = &model.SearchSug{
			Type: _pgcType,
		}
		if err = rows.Scan(&r.ID, &r.Term); err != nil {
			log.Error("PgcSeaSug row.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("d.PgcSeaSug.Query error(%v)", err)
	}
	return
}

// UgcSeaSug is used for getting ugc search suggest content
func (d *Dao) UgcSeaSug(ctx context.Context) (res []*model.SearchSug, err error) {
	var (
		rows *sql.Rows
	)
	if rows, err = d.DB.Query(ctx, _ugcSeaSug); err != nil {
		log.Error("d.UgcSeaSug.Query: %s error(%v)", _ugcSeaSug, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r = &model.SearchSug{
			Type: _ugcType,
		}
		if err = rows.Scan(&r.ID, &r.Term); err != nil {
			log.Error("UgcSeaSug row.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("d.UgcSeaSug.Query error(%v)", err)
	}
	return
}
