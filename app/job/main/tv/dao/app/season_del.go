package app

import (
	"context"
	dsql "database/sql"
	"fmt"

	model "go-common/app/job/main/tv/model/pgc"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_deletedSeasonSQL = "SELECT id,origin_name,title,alias,category,`desc`,style,area,play_time,info,state,total_num,upinfo,staff,role,copyright,`check`,is_deleted " +
		"FROM tv_ep_season WHERE is_deleted = 1 AND `check` = ? AND audit_time < UNIX_TIMESTAMP(now()) LIMIT 0,"
	_delSyncSeasonSQL = "UPDATE tv_ep_season SET `check` = ? WHERE is_deleted = 1 AND id = ?"
)

// DelSeason picks the modified season data to sync
func (d *Dao) DelSeason(c context.Context) (res []*model.TVEpSeason, err error) {
	var (
		rows *sql.Rows
	)
	if rows, err = d.DB.Query(c, _deletedSeasonSQL+fmt.Sprintf("%d", d.conf.Sync.LConf.NbSeason), SeasonToReAudit); err != nil {
		log.Error("d._deletedSeasonSQL.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r = &model.TVEpSeason{}
		if err = rows.Scan(&r.ID, &r.OriginName, &r.Title, &r.Alias, &r.Category, &r.Desc, &r.Style, &r.Area, &r.PlayTime, &r.Info,
			&r.State, &r.TotalNum, &r.Upinfo, &r.Staff, &r.Role, &r.Copyright, &r.Check, &r.IsDeleted); err != nil {
			log.Error("modSeason row.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("d.DelSeason.Query error(%v)", err)
	}
	return
}

// RejectSeason updates the indicated season to rejected status
func (d *Dao) RejectSeason(c context.Context, sid int) (nbRows int64, err error) {
	var res dsql.Result
	if res, err = d.DB.Exec(c, _delSyncSeasonSQL, SeasonRejected, sid); err != nil {
		log.Error("_delSyncSeason, failed to update to auditing: (%v), Error: %v", sid, err)
		return
	}
	return res.RowsAffected()
}
