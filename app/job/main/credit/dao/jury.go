package dao

import (
	"context"
	"time"

	"go-common/app/job/main/credit/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_addBlockInfoSQL = `INSERT INTO blocked_info (uid,origin_title,blocked_remark,origin_url,origin_content,origin_content_modify,origin_type,
		punish_time,punish_type,moral_num,blocked_days,publish_status,reason_type,operator_name,blocked_forever,blocked_type,case_id,oper_id)
		VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	_updateKPISQL          = "UPDATE blocked_kpi set rate=1,rank_per=1 where id=?"
	_updateKPIPendentSQL   = "UPDATE blocked_kpi set pendent_status=1 where id=?"
	_updateKPIHandlerSQL   = "UPDATE blocked_kpi set handler_status=1 where id=?"
	_updateCaseSQL         = "UPDATE blocked_case SET status=?,judge_type=? WHERE id=?"
	_invalidJurySQL        = "UPDATE blocked_jury SET status=2,invalid_reason=? where mid=?"
	_updateVoteRightSQL    = "UPDATE blocked_jury SET vote_total = vote_total + 1, vote_right = vote_right + 1 WHERE mid = ?"
	_updateVoteTotalSQL    = "UPDATE blocked_jury SET vote_total = vote_total + 1 WHERE mid = ?"
	_updatePunishResultSQL = "UPDATE blocked_case SET punish_result=? WHERE id=?"
	_selKPISQL             = "SELECT id,mid,rate from blocked_kpi where rate in(1,2,3) and day=?"
	_selKPIInfoSQL         = "SELECT id,mid,handler_status from blocked_kpi where id=?"
	_selCaseByIDSQL        = "SELECT id,mid,status,judge_type,relation_id from blocked_case where id=?"
	_countKPIRateSQL       = "SELECT COUNT(*) AS num FROM blocked_kpi WHERE mid=? AND rate<=4"
)

// AddBlockInfo add user block info.
func (d *Dao) AddBlockInfo(c context.Context, b *model.BlockedInfo, ts time.Time) (id int64, err error) {
	res, err := d.db.Exec(c, _addBlockInfoSQL, b.UID, b.OriginTitle, b.BlockedRemark, b.OriginURL, b.OriginContent, b.OriginContentModify, b.OriginType,
		ts, b.PunishType, b.MoralNum, b.BlockedDays, b.PublishStatus, b.ReasonType, b.OperatorName, b.BlockedForever, b.BlockedType, b.CaseID, b.OPID)
	if err != nil {
		log.Error("d.AddBlockInfo err(%v)", err)
	}
	return res.LastInsertId()
}

// UpdateKPI update kpi status to st.
func (d *Dao) UpdateKPI(c context.Context, id int64) (err error) {
	if _, err = d.db.Exec(c, _updateKPISQL, id); err != nil {
		log.Error("d.UpdateKPI err(%v)", err)
	}
	return
}

// UpdateKPIPendentStatus update blocked_kpi status to st.
func (d *Dao) UpdateKPIPendentStatus(c context.Context, id int64) (err error) {
	if _, err = d.db.Exec(c, _updateKPIPendentSQL, id); err != nil {
		log.Error("d.UpdatePendentStatus err(%v)", err)
	}
	return
}

// UpdateKPIHandlerStatus update blocked_kpi handler status.
func (d *Dao) UpdateKPIHandlerStatus(c context.Context, id int64) (err error) {
	if _, err = d.db.Exec(c, _updateKPIHandlerSQL, id); err != nil {
		log.Error("d.UpdatePendentStatus err(%v)", err)
	}
	return
}

// UpdateCase update case status to st.
func (d *Dao) UpdateCase(c context.Context, st, jt, id int64) (err error) {
	if _, err = d.db.Exec(c, _updateCaseSQL, st, jt, id); err != nil {
		log.Error("d.UpdateCase err(%v)", err)
	}
	return
}

// InvalidJury set jury invalid.
func (d *Dao) InvalidJury(c context.Context, reason, mid int64) (err error) {
	if _, err = d.db.Exec(c, _invalidJurySQL, reason, mid); err != nil {
		log.Error("d.InvalidJury err(%v)", err)
	}
	return
}

// UpdateVoteRight update vote total and vote right.
func (d *Dao) UpdateVoteRight(c context.Context, mid int64) (err error) {
	if _, err = d.db.Exec(c, _updateVoteRightSQL, mid); err != nil {
		log.Error("d.UpdateVoteRight err(%v)", err)
	}
	return
}

// UpdateVoteTotal update vote total.
func (d *Dao) UpdateVoteTotal(c context.Context, mid int64) (err error) {
	if _, err = d.db.Exec(c, _updateVoteTotalSQL, mid); err != nil {
		log.Error("d.UpdateVoteTotal err(%v)", err)
	}
	return
}

// UpdatePunishResult update table blocked_case punish_result field =0
func (d *Dao) UpdatePunishResult(c context.Context, id int64, punishResult int8) (err error) {
	if _, err = d.db.Exec(c, _updatePunishResultSQL, punishResult, id); err != nil {
		log.Error("d.UpdatePunishResult err(%v)", err)
	}
	return
}

// KPIList get kpi list.
func (d *Dao) KPIList(c context.Context, day string) (res []model.Kpi, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _selKPISQL, day); err != nil {
		log.Error("dao.KPIList error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := model.Kpi{}
		if err = rows.Scan(&r.ID, &r.Mid, &r.Rate); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// KPIInfo get KPI info.
func (d *Dao) KPIInfo(c context.Context, id int64) (res model.Kpi, err error) {
	row := d.db.QueryRow(c, _selKPIInfoSQL, id)
	if err = row.Scan(&res.ID, &res.Mid, &res.HandlerStatus); err != nil {
		log.Error("d.KPIInfo err(%v)", err)
	}
	return
}

// CaseByID get case info by id.
func (d *Dao) CaseByID(c context.Context, id int64) (res model.Case, err error) {
	row := d.db.QueryRow(c, _selCaseByIDSQL, id)
	if err = row.Scan(&res.ID, &res.Mid, &res.Status, &res.JudgeType, &res.RelationID); err != nil {
		log.Error("d.BlockCount err(%v)", err)
	}
	return
}

// CountKPIRate count KPI rate<=4(C).
func (d *Dao) CountKPIRate(c context.Context, mid int64) (count int, err error) {
	row := d.db.QueryRow(c, _countKPIRateSQL, mid)
	if err = row.Scan(&count); err != nil {
		log.Error("d.CountKPIRate(mid:%d) err(%v)", mid, err)
	}
	return
}
