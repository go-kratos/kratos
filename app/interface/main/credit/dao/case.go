package dao

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	model "go-common/app/interface/main/credit/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_addBlockedCasesSQL       = "INSERT INTO blocked_case(status,mid,operator,origin_content,punish_result,origin_title,origin_type,origin_url,blocked_days,reason_type,relation_id,oper_id,business_time) VALUES %s"
	_insertVoteSQL            = "INSERT INTO blocked_case_vote(mid,cid,expired) VALUES(?,?,?)"
	_updateVoteSQL            = "INSERT INTO blocked_case_vote(mid,cid,vote) VALUES(?,?,?) ON DUPLICATE KEY UPDATE vote=?"
	_inBlockedCaseApplyLogSQL = "INSERT INTO blocked_case_apply_log(mid,case_id,apply_type,origin_reason,apply_reason) VALUES(?,?,?,?,?)"
	_updateCaseVoteTotalSQL   = "UPDATE blocked_case SET %s=%s+? WHERE id=?"
	_getCaseByIDSQL           = `SELECT id,mid,status,origin_content,punish_result,origin_title,origin_url,end_time,vote_rule,vote_break,vote_delete,origin_type,reason_type,judge_type,blocked_days,
	put_total,start_time,end_time,operator,ctime,mtime,relation_id,case_type FROM blocked_case WHERE id=? AND status IN (1,3,4,6)`
	_countCaseVoteSQL = "SELECT COUNT(*) FROM blocked_case_vote WHERE mid=? AND vote!=3"
	_isVoteByMIDSQL   = "SELECT id  FROM blocked_case_vote WHERE mid=? AND cid=? AND vote=0"
	_getVoteInfoSQL   = "SELECT id,cid,mid,vote,expired,mtime FROM blocked_case_vote WHERE mid=? AND cid=?"
	_loadMidVoteIDSQL = "SELECT v.cid,c.case_type FROM blocked_case_vote v INNER JOIN blocked_case c ON v.cid=c.id WHERE v.mid =? AND v.ctime >=? ORDER BY v.id DESC"
	_getCaseByIDsSQL  = `SELECT id,mid,status,origin_content,punish_result,origin_title,origin_url,end_time,vote_rule,vote_break,vote_delete,origin_type,reason_type,judge_type,blocked_days,
	put_total,start_time,end_time,operator,ctime,mtime,relation_id,case_type FROM blocked_case WHERE id IN(%s) AND status IN (1,3,4,6)`
	_caseRelationIDCountSQL = "SELECT COUNT(*) FROM blocked_case WHERE origin_type=? AND relation_id =?"
	_caseInfoIDsSQL         = `SELECT id,mid,status,origin_content,punish_result,origin_title,origin_url,end_time,vote_rule,vote_break,vote_delete,origin_type,reason_type,judge_type,blocked_days,
	put_total,start_time,end_time,operator,ctime,mtime,relation_id,case_type FROM blocked_case WHERE id IN (%s)`
	_caseVotesMIDSQL  = "SELECT id,cid,mid,vote,expired,mtime FROM blocked_case_vote WHERE id IN (%s)"
	_caseVoteIDMIDSQL = "SELECT id,cid FROM blocked_case_vote WHERE mid = ? ORDER BY mtime DESC LIMIT ?,?"
	_caseVoteIDTopSQL = "SELECT id,cid FROM blocked_case_vote WHERE mid = ? ORDER BY mtime DESC LIMIT 100"
)

// AddBlockedCases batch add blocked cases.
func (d *Dao) AddBlockedCases(c context.Context, bc []*model.ArgJudgeCase) (err error) {
	l := len(bc)
	valueStrings := make([]string, 0, l)
	valueArgs := make([]interface{}, 0, l*13)
	for _, b := range bc {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, strconv.FormatInt(2, 10))
		valueArgs = append(valueArgs, strconv.FormatInt(b.MID, 10))
		valueArgs = append(valueArgs, b.Operator)
		valueArgs = append(valueArgs, b.OContent)
		valueArgs = append(valueArgs, strconv.FormatInt(int64(b.PunishResult), 10))
		valueArgs = append(valueArgs, b.OTitle)
		valueArgs = append(valueArgs, strconv.FormatInt(int64(b.OType), 10))
		valueArgs = append(valueArgs, b.OURL)
		valueArgs = append(valueArgs, strconv.FormatInt(int64(b.BlockedDays), 10))
		valueArgs = append(valueArgs, strconv.FormatInt(int64(b.ReasonType), 10))
		valueArgs = append(valueArgs, b.RelationID)
		valueArgs = append(valueArgs, strconv.FormatInt(b.OperID, 10))
		if b.BCTime != 0 {
			valueArgs = append(valueArgs, b.BCTime)
		} else {
			valueArgs = append(valueArgs, "1979-12-31 16:00:00")
		}
	}
	stmt := fmt.Sprintf(_addBlockedCasesSQL, strings.Join(valueStrings, ","))
	_, err = d.db.Exec(c, stmt, valueArgs...)
	if err != nil {
		log.Error("AddBlockedCases: db.Exec(bc(%+v)) error(%v)", bc, err)
	}
	return
}

// InsVote insert user vote.
func (d *Dao) InsVote(c context.Context, mid int64, cid int64, t int64) (err error) {
	if _, err = d.db.Exec(c, _insertVoteSQL, mid, cid, time.Now().Add(time.Duration(t)*time.Minute)); err != nil {
		log.Error("InsVote: db.Exec(%d,%d) error(%v)", mid, cid, err)
	}
	return
}

// Setvote set user vote.
func (d *Dao) Setvote(c context.Context, mid, cid, vote int64) (err error) {
	if _, err = d.db.Exec(c, _updateVoteSQL, mid, cid, vote, vote); err != nil {
		log.Error("Setvote: db.Exec(%d,%d,%d) error(%v)", mid, cid, vote, err)
	}
	return
}

// SetVoteTx set vote info by tx.
func (d *Dao) SetVoteTx(tx *sql.Tx, mid, cid int64, vote int8) (affect int64, err error) {
	row, err := tx.Exec(_updateVoteSQL, mid, cid, vote, vote)
	if err != nil {
		log.Error("SetVoteTx err(%v)", err)
		return
	}
	return row.LastInsertId()
}

// AddCaseReasonApply add case reason apply log.
func (d *Dao) AddCaseReasonApply(c context.Context, mid, cid int64, applyType, originReason, applyReason int8) (err error) {
	_, err = d.db.Exec(c, _inBlockedCaseApplyLogSQL, mid, cid, applyType, originReason, applyReason)
	if err != nil {
		log.Error("AddCaseReasonApply err(%v)", err)
		return
	}
	return
}

// AddCaseVoteTotal add case vote total.
func (d *Dao) AddCaseVoteTotal(c context.Context, field string, cid int64, voteNum int8) (err error) {
	sql := fmt.Sprintf(_updateCaseVoteTotalSQL, field, field)
	if _, err = d.db.Exec(c, sql, voteNum, cid); err != nil {
		log.Error("AddCaseVoteTotal: db.Exec(%d) error(%v)", cid, err)
	}
	return
}

// CaseInfo jury get case info.
func (d *Dao) CaseInfo(c context.Context, cid int64) (r *model.BlockedCase, err error) {
	row := d.db.QueryRow(c, _getCaseByIDSQL, cid)
	r = &model.BlockedCase{}
	if err = row.Scan(&r.ID, &r.MID, &r.Status, &r.OriginContent, &r.PunishResult, &r.OriginTitle, &r.OriginURL, &r.EndTime, &r.VoteRule, &r.VoteBreak, &r.VoteDelete, &r.OriginType, &r.ReasonType, &r.JudgeType, &r.BlockedDays, &r.PutTotal, &r.StartTime, &r.EndTime, &r.Operator, &r.CTime, &r.MTime, &r.RelationID, &r.CaseType); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
			return
		}
	}
	return
}

// CountCaseVote jury count case vote total.
func (d *Dao) CountCaseVote(c context.Context, mid int64) (r int64, err error) {
	row := d.db.QueryRow(c, _countCaseVoteSQL, mid)
	if err = row.Scan(&r); err != nil {
		log.Error("row.Scan() error(%v)", err)
	}
	return
}

// IsVote jury  user is vote.
func (d *Dao) IsVote(c context.Context, mid int64, cid int64) (r int64, err error) {
	row := d.db.QueryRow(c, _isVoteByMIDSQL, mid, cid)
	if err = row.Scan(&r); err != nil {
		if err == sql.ErrNoRows {
			r = 0
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
	}
	return
}

// VoteInfo jury user get vote info.
func (d *Dao) VoteInfo(c context.Context, mid int64, cid int64) (r *model.VoteInfo, err error) {
	row := d.db.QueryRow(c, _getVoteInfoSQL, mid, cid)
	r = &model.VoteInfo{}
	if err = row.Scan(&r.ID, &r.CID, &r.MID, &r.Vote, &r.Expired, &r.Mtime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
		}
	}
	return
}

// LoadVoteIDsMid  load user vote case ids.
func (d *Dao) LoadVoteIDsMid(c context.Context, mid int64, day int) (cases map[int64]*model.SimCase, err error) {
	now := time.Now()
	t := time.Date(now.Year(), now.Month(), now.Day()-day, 0, 0, 0, 0, now.Location())
	rows, err := d.db.Query(c, _loadMidVoteIDSQL, mid, t)
	if err != nil {
		log.Error("d.db.Query(%d %v) error(%v)", mid, t, err)
		return
	}
	defer rows.Close()
	cases = make(map[int64]*model.SimCase)
	for rows.Next() {
		mcase := &model.SimCase{}
		if err = rows.Scan(&mcase.ID, &mcase.CaseType); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		cases[mcase.ID] = mcase
	}
	return
}

// CaseVoteIDs get user's vote info by ids.
func (d *Dao) CaseVoteIDs(c context.Context, ids []int64) (mbc map[int64]*model.BlockedCase, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_getCaseByIDsSQL, xstr.JoinInts(ids)))
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	mbc = make(map[int64]*model.BlockedCase, len(ids))
	for rows.Next() {
		r := new(model.BlockedCase)
		if err = rows.Scan(&r.ID, &r.MID, &r.Status, &r.OriginContent, &r.PunishResult, &r.OriginTitle, &r.OriginURL, &r.EndTime, &r.VoteRule, &r.VoteBreak, &r.VoteDelete, &r.OriginType, &r.ReasonType, &r.JudgeType, &r.BlockedDays, &r.PutTotal, &r.StartTime, &r.EndTime, &r.Operator, &r.CTime, &r.MTime, &r.RelationID, &r.CaseType); err != nil {
			if err == sql.ErrNoRows {
				mbc = nil
				err = nil
				return
			}
			err = errors.WithStack(err)
			return
		}
		mbc[r.ID] = r
	}
	err = rows.Err()
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

// CaseInfoIDs  get case info by ids.
func (d *Dao) CaseInfoIDs(c context.Context, ids []int64) (cases map[int64]*model.BlockedCase, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_caseInfoIDsSQL, xstr.JoinInts(ids)))
	if err != nil {
		log.Error("d.CaseInfoIDs err(%v)", err)
		return
	}
	defer rows.Close()
	cases = make(map[int64]*model.BlockedCase, len(ids))
	for rows.Next() {
		ca := &model.BlockedCase{}
		if err = rows.Scan(&ca.ID, &ca.MID, &ca.Status, &ca.OriginContent, &ca.PunishResult, &ca.OriginTitle, &ca.OriginURL, &ca.EndTime, &ca.VoteRule,
			&ca.VoteBreak, &ca.VoteDelete, &ca.OriginType, &ca.ReasonType, &ca.JudgeType, &ca.BlockedDays, &ca.PutTotal, &ca.StartTime, &ca.EndTime,
			&ca.Operator, &ca.CTime, &ca.MTime, &ca.RelationID, &ca.CaseType); err != nil {
			log.Error("row.Scan err(%v)", err)
			return
		}
		cases[ca.ID] = ca
	}
	return
}

// CaseVotesMID get user's vote case ids.
func (d *Dao) CaseVotesMID(c context.Context, ids []int64) (mvo map[int64]*model.VoteInfo, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_caseVotesMIDSQL, xstr.JoinInts(ids)))
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	mvo = make(map[int64]*model.VoteInfo, len(ids))
	for rows.Next() {
		vo := new(model.VoteInfo)
		if err = rows.Scan(&vo.ID, &vo.CID, &vo.MID, &vo.Vote, &vo.Expired, &vo.Mtime); err != nil {
			if err == sql.ErrNoRows {
				mvo = nil
				err = nil
				return
			}
			err = errors.WithStack(err)
			return
		}
		mvo[vo.CID] = vo
	}
	err = rows.Err()
	return
}

// CaseVoteIDMID get user's vote case ids and cids.
func (d *Dao) CaseVoteIDMID(c context.Context, mid, pn, ps int64) (vids []int64, cids []int64, err error) {
	rows, err := d.db.Query(c, _caseVoteIDMIDSQL, mid, (pn-1)*ps, ps)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var vid, cid int64
		if err = rows.Scan(&vid, &cid); err != nil {
			if err == sql.ErrNoRows {
				vids = nil
				cids = nil
				err = nil
				return
			}
			err = errors.WithStack(err)
			return
		}
		vids = append(vids, vid)
		cids = append(cids, cid)
	}
	err = rows.Err()
	return
}

// CaseVoteIDTop get user's vote case ids and cids by top 100.
func (d *Dao) CaseVoteIDTop(c context.Context, mid int64) (vids []int64, cids []int64, err error) {
	rows, err := d.db.Query(c, _caseVoteIDTopSQL, mid)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var vid, cid int64
		if err = rows.Scan(&vid, &cid); err != nil {
			if err == sql.ErrNoRows {
				vids = nil
				cids = nil
				err = nil
				return
			}
			err = errors.WithStack(err)
			return
		}
		vids = append(vids, vid)
		cids = append(cids, cid)
	}
	err = rows.Err()
	return
}
