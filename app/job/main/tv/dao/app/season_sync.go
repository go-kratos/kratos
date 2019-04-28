package app

import (
	"context"
	dsql "database/sql"
	"fmt"
	"time"

	model "go-common/app/job/main/tv/model/pgc"
	"go-common/library/database/sql"
	"go-common/library/log"
)

// season `check` meaning
const (
	SeasonRejected      = 0
	SeasonPassed        = 1
	SeasonToReAudit     = 2
	SeasonAuditing      = 3
	SeasonAlreadyPassed = 7
	EPToAudit           = 1
	EPAuditing          = 2
	EPPassed            = 3
	EPNotPass           = 4
	_CMSValid           = 1
	_NotDeleted         = 0
)

const (
	_modifiedSeason = "SELECT id,origin_name,title,alias,category,`desc`,style,area,play_time,info,state,total_num,upinfo,staff,role,copyright,`check`,is_deleted,cover,producer,version,status" +
		" FROM tv_ep_season WHERE `check` = ? AND audit_time < UNIX_TIMESTAMP(now()) AND is_deleted = 0 LIMIT 0,"
	_snEmpty        = "SELECT id FROM tv_content WHERE season_id = ? AND is_deleted = 0 LIMIT 1"
	_auditingSeason = "UPDATE tv_ep_season SET `check` = ? WHERE is_deleted = 0 AND id = ?"
	_delaySeason    = "UPDATE tv_ep_season SET audit_time = ? WHERE id = ?"
)

// SnEmpty determines whether the
func (d *Dao) SnEmpty(c context.Context, sid int64) (res bool, err error) {
	var epid int
	if err = d.DB.QueryRow(c, _snEmpty, sid).Scan(&epid); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			res = true
			return
		}
		log.Error("d.SnEmpty Error (%v)", err)
		return
	}
	res = false
	return
}

// ModSeason picks the modified season data to sync
func (d *Dao) ModSeason(c context.Context) (res []*model.TVEpSeason, err error) {
	var (
		rows    *sql.Rows
		isEmpty bool
	)
	if rows, err = d.DB.Query(c, _modifiedSeason+fmt.Sprintf("%d", d.conf.Sync.LConf.NbSeason), SeasonToReAudit); err != nil {
		log.Error("d._modifiedSeason.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r = &model.TVEpSeason{}
		if err = rows.Scan(&r.ID, &r.OriginName, &r.Title, &r.Alias, &r.Category, &r.Desc, &r.Style, &r.Area, &r.PlayTime, &r.Info,
			&r.State, &r.TotalNum, &r.Upinfo, &r.Staff, &r.Role, &r.Copyright, &r.Check, &r.IsDeleted, &r.Cover,
			&r.Producer, &r.Version, &r.Status); err != nil {
			log.Error("modSeason row.Scan() error(%v)", err)
			return
		}
		if isEmpty, err = d.SnEmpty(c, r.ID); err != nil {
			log.Error("modSeason SnEmpty Error (%v)", err)
			return
		}
		if !isEmpty { // we don't submit empty season to audit
			res = append(res, r)
		} else {
			d.DelaySeason(c, r.ID)
		}
	}
	if err = rows.Err(); err != nil {
		log.Error("d.ModSeason.Query error(%v)", err)
	}
	return
}

// AuditSeason updates the indicated season to Auditing Status
func (d *Dao) AuditSeason(c context.Context, sid int) (nbRows int64, err error) {
	var res dsql.Result
	if res, err = d.DB.Exec(c, _auditingSeason, SeasonAuditing, sid); err != nil {
		log.Error("_auditingSeason, failed to update to auditing: (%v), Error: %v", sid, err)
		return
	}
	return res.RowsAffected()
}

// DelaySeason postpones the season to sync in 30 minutes
func (d *Dao) DelaySeason(c context.Context, sid int64) (nbRows int64, err error) {
	var (
		res   dsql.Result
		delay = time.Now().Unix() + d.conf.Sync.Frequency.AuditDelay
	)
	if res, err = d.DB.Exec(c, _delaySeason, delay, sid); err != nil {
		log.Error("_delaySeason, failed to delay: (%v,%v), Error: %v", err)
		return
	}
	return res.RowsAffected()
}
