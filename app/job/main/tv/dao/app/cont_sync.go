package app

import (
	"context"
	dsql "database/sql"
	"fmt"
	"time"

	model "go-common/app/job/main/tv/model/pgc"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_pickDataSQL = "SELECT id, title, subtitle, `desc`, cover, epid, season_id, cid FROM tv_content " +
		"WHERE season_id = ? AND state = 1 AND audit_time < UNIX_TIMESTAMP(now()) AND is_deleted = 0"
	_pickSeasonSQL   = "SELECT id,origin_name,title,alias,category,`desc`,style,area,play_time,info,state,total_num,upinfo,staff,role,copyright,`check`,is_deleted,cover,producer,version,status FROM tv_ep_season WHERE id = ?"
	_pickEPSQL       = "SELECT a.id,a.season_id,a.title,a.long_title,a.cover,a.length,a.cid,b.pay_status FROM tv_ep_content AS a LEFT JOIN tv_content AS b ON a.id=b.epid WHERE a.id=? AND a.is_deleted=0"
	_waitCallSQL     = "UPDATE tv_content SET audit_time = ? WHERE epid = ? AND state = ? AND is_deleted = 0"
	_deleteEPSQL     = "UPDATE tv_ep_content SET is_deleted = 1 WHERE season_id = ? AND is_deleted = 0"
	_deleteContSQL   = "UPDATE tv_content SET is_deleted = 1 WHERE season_id = ? AND is_deleted = 0"
	_rejectContSQL   = "UPDATE tv_content SET audit_time = ? WHERE season_id = ? AND state = ? AND is_deleted = 0"
	_auditingContSQL = "UPDATE tv_content SET state = ? WHERE state = ? AND is_deleted = 0 AND epid IN (%s)"
	_removeContSQL   = "UPDATE tv_content SET state = ?,is_deleted = 1 WHERE state = ? AND is_deleted = 0 AND epid = ?"
	_ContSQL         = "SELECT id, title, subtitle, `desc`, cover, epid, season_id, cid FROM tv_content WHERE epid = ?"
	_readySns        = "SELECT DISTINCT a.season_id FROM tv_content a LEFT JOIN tv_ep_season b ON a.season_id = b.id " +
		"WHERE a.state = 1 AND a.is_deleted = 0 AND a.audit_time < UNIX_TIMESTAMP(now()) " +
		"AND b.`check` != 0 AND b.is_deleted = 0"
)

// RemoveCont is used to treat invalid ep data's content
func (d *Dao) RemoveCont(c context.Context, epid int) (nbRows int64, err error) {
	var (
		res dsql.Result
	)
	if res, err = d.DB.Exec(c, _removeContSQL, EPNotPass, EPToAudit, epid); err != nil {
		log.Error("_removeContSQL, failed to delay: (%v,%v), Error: %v", err)
		return
	}
	return res.RowsAffected()
}

// PickData picks the source content data to sync
func (d *Dao) PickData(c context.Context, sid int64) (res [][]*model.Content, err error) {
	var (
		rows   *sql.Rows
		nbData = d.conf.Sync.LConf.SizeMsg
		conts  []*model.Content
	)
	if rows, err = d.DB.Query(c, _pickDataSQL, sid); err != nil {
		log.Error("d._pickDataSQL error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r = &model.Content{}
		if err = rows.Scan(&r.ID, &r.Title, &r.Subtitle, &r.Desc, &r.Cover, &r.EPID, &r.SeasonID, &r.CID); err != nil {
			log.Error("Conts row.Scan() error(%v)", err)
			return
		}
		conts = append(conts, r)
		if len(conts) >= nbData {
			res = append(res, conts)
			conts = []*model.Content{}
		}
	}
	if err = rows.Err(); err != nil {
		log.Error("d.PickData.Query error(%v)", err)
	}
	if len(conts) > 0 {
		res = append(res, conts)
	}
	return
}

// Season gets one season data
func (d *Dao) Season(c context.Context, sid int) (r *model.TVEpSeason, err error) {
	row := d.DB.QueryRow(c, _pickSeasonSQL, sid)
	r = &model.TVEpSeason{}
	if err = row.Scan(&r.ID, &r.OriginName, &r.Title, &r.Alias, &r.Category, &r.Desc, &r.Style, &r.Area, &r.PlayTime, &r.Info,
		&r.State, &r.TotalNum, &r.Upinfo, &r.Staff, &r.Role, &r.Copyright, &r.Check, &r.IsDeleted, &r.Cover, &r.Producer, &r.Version, &r.Status); err != nil {
		return
	}
	return
}

// EP gets one EP data
func (d *Dao) EP(c context.Context, epid int) (r *model.TVEpContent, err error) {
	row := d.DB.QueryRow(c, _pickEPSQL, epid)
	r = &model.TVEpContent{}
	if err = row.Scan(&r.ID, &r.SeasonID, &r.Title, &r.LongTitle, &r.Cover, &r.Length, &r.CID, &r.PayStatus); err != nil {
		if err != sql.ErrNoRows {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	return
}

// Cont picks the content data
func (d *Dao) Cont(c context.Context, epid int) (res *model.Content, err error) {
	row := d.DB.QueryRow(c, _ContSQL, epid)
	res = &model.Content{}
	if err = row.Scan(&res.ID, &res.Title, &res.Subtitle, &res.Desc, &res.Cover, &res.EPID, &res.SeasonID, &res.CID); err != nil {
		if err != sql.ErrNoRows {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	return
}

// WaitCall updates the audit time ( wait for caliing also )
func (d *Dao) WaitCall(c context.Context, epid int) (nbRows int64, err error) {
	var (
		res   dsql.Result
		delay = time.Now().Unix() + int64(d.conf.Sync.Frequency.WaitCall)
	)
	if res, err = d.DB.Exec(c, _waitCallSQL, delay, epid, EPToAudit); err != nil {
		log.Error("_waitCallSQL, failed to delay: (%v,%v), Error: %v", err)
		return
	}
	return res.RowsAffected()
}

// DeleteEP deletes ep data of a deleted season
func (d *Dao) DeleteEP(c context.Context, sid int) (nbRows int64, err error) {
	var res dsql.Result
	if res, err = d.DB.Exec(c, _deleteEPSQL, sid); err != nil {
		log.Error("_deleteEPSQL, failed to delete: (%v), Error: %v", sid, err)
		return
	}
	return res.RowsAffected()
}

// DeleteCont deletes cont data of a deleted season
func (d *Dao) DeleteCont(c context.Context, sid int) (nbRows int64, err error) {
	var res dsql.Result
	if res, err = d.DB.Exec(c, _deleteContSQL, sid); err != nil {
		log.Error("_deleteContSQL, failed to delete: (%v), Error: %v", sid, err)
		return
	}
	return res.RowsAffected()
}

// RejectCont postpones its content in 1 day
func (d *Dao) RejectCont(c context.Context, sid int) (nbRows int64, err error) {
	var (
		res   dsql.Result
		delay = time.Now().Unix() + int64(d.conf.Sync.Frequency.RejectWait)
	)
	if res, err = d.DB.Exec(c, _rejectContSQL, delay, sid, EPToAudit); err != nil {
		log.Error("_rejectContSQL, failed to reject: (%v), Error: %v", sid, err)
		return
	}
	return res.RowsAffected()
}

// AuditingCont updates the content state from 1 ( auditing ) to 2
func (d *Dao) AuditingCont(c context.Context, conts []*model.Content) (nbRows int64, err error) {
	var (
		res   dsql.Result
		epids []int64
	)
	for _, v := range conts {
		epids = append(epids, int64(v.EPID))
	}
	if res, err = d.DB.Exec(c, fmt.Sprintf(_auditingContSQL, xstr.JoinInts(epids)), EPAuditing, EPToAudit); err != nil {
		log.Error("_auditingContSQL, failed to update to auditing: (%v), Error: %v", epids, err)
		return
	}
	return res.RowsAffected()
}

// ReadySns picks ready to sync seasons
func (d *Dao) ReadySns(c context.Context) (res []int64, err error) {
	var rows *sql.Rows
	if rows, err = d.DB.Query(c, _readySns); err != nil {
		log.Error("d.ReadySns.Query: %s error(%v)", _readySns, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r int64
		if err = rows.Scan(&r); err != nil {
			log.Error("ReadySns row.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("ReadySns rows.Err() error(%v)", err)
	}
	return
}
