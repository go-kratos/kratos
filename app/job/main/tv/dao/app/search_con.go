package app

import (
	"context"

	model "go-common/app/job/main/tv/model/pgc"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_PgcFromwhere = "FROM tv_ep_season WHERE `check` = 1 AND valid = 1 AND is_deleted = 0 "
	_PgcCont      = "SELECT id,category,cover,title,play_time,role,staff,`desc` " + _PgcFromwhere + "AND id > ? limit ?"
	_PgcContCount = " SELECT count(*) " + _PgcFromwhere
)

// PgcCont is used for getting valid pgc season data
func (d *Dao) PgcCont(ctx context.Context, id int, limit int) (res []*model.SearPgcCon, err error) {
	var (
		rows *sql.Rows
	)
	if rows, err = d.DB.Query(ctx, _PgcCont, id, limit); err != nil {
		log.Error("d.PgcCont.Query: %s error(%v)", _PgcCont, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r = &model.SearPgcCon{}
		if err = rows.Scan(&r.ID, &r.Category, &r.Cover, &r.Title, &r.PlayTime, &r.Role, &r.Staff, &r.Desc); err != nil {
			log.Error("PgcCont row.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("d.PgcCont.Query error(%v)", err)
	}
	return
}

// PgcContCount is used for getting valid data count
func (d *Dao) PgcContCount(ctx context.Context) (upCnt int, err error) {
	row := d.DB.QueryRow(ctx, _PgcContCount)
	if err = row.Scan(&upCnt); err != nil {
		log.Error("d.SeaContCount.Query: %s error(%v)", _PgcContCount, err)
	}
	return
}
