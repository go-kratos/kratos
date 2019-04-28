package dao

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-common/app/job/main/credit/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"
	"go-common/library/xstr"
)

const (
	_inBlockedCasesSQL         = "INSERT INTO blocked_case(mid,status,origin_content,punish_result,origin_title,origin_type,origin_url,blocked_days,reason_type,relation_id,oper_id,business_time) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)"
	_inBlockedCaseModifyLogSQL = "INSERT INTO blocked_case_modify_log(case_id,modify_type,origin_reason,modify_reason) VALUES (?,?,?,?)"
	_upBlockedCaseReasonSQL    = "UPDATE blocked_case SET reason_type = ?  WHERE id = ?"
	_upBlockedCaseStatusSQL    = "UPDATE blocked_case SET status = ?  WHERE id = ?"
	_upGrantCaseSQL            = "UPDATE blocked_case SET status = 1, start_time = ?, end_time = ? WHERE id IN (%s)"
	_grantCaseSQL              = "SELECT id,mid,start_time,end_time,vote_rule,vote_break,vote_delete,case_type FROM blocked_case WHERE status = 8 ORDER BY mtime ASC limit ?"
	_caseVoteSQL               = "SELECT mid,vote,expired FROM blocked_case_vote WHERE id = ?"
	_caseRelationIDCountSQL    = "SELECT COUNT(*) FROM blocked_case WHERE origin_type=? AND relation_id =?"
	_caseVotesCIDSQL           = "SELECT mid,vote FROM blocked_case_vote WHERE cid = ?"
	_countCaseMIDSQL           = "SELECT COUNT(*) FROM blocked_case WHERE mid = ? AND origin_type = ? AND ctime >= ?"
	_caseApplyReasonsSQL       = "SELECT DISTINCT(apply_reason) FROM blocked_case_apply_log WHERE case_id = ? AND apply_type = 1"
	_caseApplyReasonNumSQL     = "SELECT COUNT(*) AS num, case_id, apply_type, apply_reason, origin_reason FROM `blocked_case_apply_log` WHERE case_id = ? AND apply_type = 1 AND apply_reason IN (%s)  GROUP BY apply_reason ORDER BY mtime ASC"
	_casesStatusSQL            = "SELECT id,status FROM blocked_case WHERE id IN (%s)"
)

// AddBlockedCase add blocked case.
func (d *Dao) AddBlockedCase(c context.Context, ca *model.Case) (err error) {
	if _, err = d.db.Exec(c, _inBlockedCasesSQL, ca.Mid, ca.Status, ca.OriginContent, ca.PunishResult, ca.OriginTitle, ca.OriginType, ca.OriginURL, ca.BlockedDay, ca.ReasonType, ca.RelationID, ca.OPID, ca.BCtime); err != nil {
		log.Error("d.AddBlockedCase(%+v) err(%v)", ca, err)
	}
	return
}

// AddBlockedCaseModifyLog add blocked case modify log.
func (d *Dao) AddBlockedCaseModifyLog(c context.Context, cid int64, mType, oReason, mReason int8) (err error) {
	if _, err = d.db.Exec(c, _inBlockedCaseModifyLogSQL, cid, mType, oReason, mReason); err != nil {
		log.Error("d.AddBlockedCaseModifyLog(%d,%d,%d,%d) err(%v)", cid, mType, oReason, mReason, err)
	}
	return
}

// UpBlockedCaseReason update blocked case reason.
func (d *Dao) UpBlockedCaseReason(c context.Context, cid int64, reasonType int8) (affect int64, err error) {
	var result sql.Result
	if result, err = d.db.Exec(c, _upBlockedCaseReasonSQL, reasonType, cid); err != nil {
		log.Error("d.UpBlockedCaseReason(%d,%d) err(%v)", cid, reasonType, err)
		return
	}
	return result.RowsAffected()
}

// UpBlockedCaseStatus update blocked case status.
func (d *Dao) UpBlockedCaseStatus(c context.Context, cid int64, status int8) (err error) {
	if _, err = d.db.Exec(c, _upBlockedCaseStatusSQL, status, cid); err != nil {
		log.Error("d.UpBlockedCaseStatus(%d,%d) err(%v)", cid, status, err)
	}
	return
}

// UpGrantCase update blocked_case from que to granting.
func (d *Dao) UpGrantCase(c context.Context, ids []int64, stime xtime.Time, etime xtime.Time) (err error) {
	if _, err = d.db.Exec(c, fmt.Sprintf(_upGrantCaseSQL, xstr.JoinInts(ids)), stime, etime); err != nil {
		log.Error("d.UpGrantCase(%s,%v,%v) err(%v)", xstr.JoinInts(ids), stime, etime, err)
	}
	return
}

// Grantcase  get case from state of CaseStatusQueueing.
func (d *Dao) Grantcase(c context.Context, limit int) (mcases map[int64]*model.SimCase, err error) {
	mcases = make(map[int64]*model.SimCase)
	rows, err := d.db.Query(c, _grantCaseSQL, limit)
	if err != nil {
		log.Error("d.Grantcase err(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		ca := &model.SimCase{}
		if err = rows.Scan(&ca.ID, &ca.Mid, &ca.Stime, &ca.Etime, &ca.VoteRule, &ca.VoteBreak, &ca.VoteDelete, &ca.CaseType); err != nil {
			log.Error("rows.Scan err(%v)", err)
			return
		}
		mcases[ca.ID] = ca
	}
	err = rows.Err()
	return
}

// CaseVote get blocked case vote info.
func (d *Dao) CaseVote(c context.Context, id int64) (res *model.CaseVote, err error) {
	res = &model.CaseVote{}
	row := d.db.QueryRow(c, _caseVoteSQL, id)
	if err = row.Scan(&res.MID, &res.Vote, &res.Expired); err != nil {
		log.Error("d.CaseVote err(%v)", err)
	}
	return
}

// CaseRelationIDCount get case relation_id count.
func (d *Dao) CaseRelationIDCount(c context.Context, tp int8, relationID string) (count int64, err error) {
	row := d.db.QueryRow(c, _caseRelationIDCountSQL, tp, relationID)
	if err = row.Scan(&count); err != nil {
		log.Error("d.caseRelationIDCount err(%v)", err)
	}
	return
}

// CaseVotesCID is blocked case vote list by cid.
func (d *Dao) CaseVotesCID(c context.Context, cid int64) (res []*model.CaseVote, err error) {
	rows, err := d.db.Query(c, _caseVotesCIDSQL, cid)
	if err != nil {
		log.Error("d.CaseVoteCID(%d) err(%v)", cid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		ca := &model.CaseVote{}
		if err = rows.Scan(&ca.MID, &ca.Vote); err != nil {
			log.Error("rows.Scan err(%v)", err)
			return
		}
		res = append(res, ca)
	}
	err = rows.Err()
	return
}

// CountCaseMID get count case by mid.
func (d *Dao) CountCaseMID(c context.Context, mid int64, tp int8) (count int64, err error) {
	row := d.db.QueryRow(c, _countCaseMIDSQL, mid, tp, time.Now().AddDate(0, 0, -2))
	if err = row.Scan(&count); err != nil {
		log.Error("d.CountCaseMID err(%v)", err)
	}
	return
}

// CaseApplyReasons case apply reasons.
func (d *Dao) CaseApplyReasons(c context.Context, cid int64) (aReasons []int64, err error) {
	rows, err := d.db.Query(c, _caseApplyReasonsSQL, cid)
	if err != nil {
		log.Error("d.CaseApplyReasons(%d) err(%v)", cid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var aReason int64
		if err = rows.Scan(&aReason); err != nil {
			log.Error("rows.Scan err(%v)", err)
			return
		}
		aReasons = append(aReasons, aReason)
	}
	err = rows.Err()
	return
}

// CaseApplyReasonNum case group apply reasons num.
func (d *Dao) CaseApplyReasonNum(c context.Context, cid int64, aReasons []int64) (cas []*model.CaseApplyModifyLog, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_caseApplyReasonNumSQL, xstr.JoinInts(aReasons)), cid)
	if err != nil {
		log.Error("d.CaseApplyReasonNum(%s) err(%v)", xstr.JoinInts(aReasons), err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		ca := &model.CaseApplyModifyLog{}
		if err = rows.Scan(&ca.Num, &ca.CID, &ca.AType, &ca.AReason, &ca.OReason); err != nil {
			log.Error("rows.Scan err(%v)", err)
			return
		}
		cas = append(cas, ca)
	}
	err = rows.Err()
	return
}

// CasesStatus get cases status.
func (d *Dao) CasesStatus(c context.Context, cids []int64) (map[int64]int8, error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_casesStatusSQL, xstr.JoinInts(cids)))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cs := make(map[int64]int8, len(cids))
	for rows.Next() {
		var (
			cid    int64
			status int8
		)
		if err = rows.Scan(&cid, &status); err != nil {
			if err == xsql.ErrNoRows {
				err = nil
			}
			return nil, err
		}
		cs[cid] = status
	}
	err = rows.Err()
	return cs, err
}
