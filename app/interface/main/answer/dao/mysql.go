package dao

import (
	"context"
	xsql "database/sql"
	"fmt"
	"math/rand"
	"time"

	"go-common/app/interface/main/answer/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_shard = 10

	_getAnswerHistorySQL = "SELECT id,hid,mid,start_time,step_one_err_num,step_one_complete_time,step_two_start_time,complete_time,complete_result,score,is_pass_captcha,is_first_pass,passed_level,rank_id,step_extra_start_time,step_extra_complete_time,extra_score FROM answer_history_%d WHERE mid=? ORDER BY id DESC LIMIT 1"
	_getHistoryByHidSQL  = "SELECT id,hid,mid,start_time,step_one_err_num,step_one_complete_time,step_two_start_time,complete_time,complete_result,score,is_pass_captcha,is_first_pass,passed_level,rank_id,step_extra_start_time,step_extra_complete_time,extra_score FROM answer_history_%d WHERE hid=? ORDER BY id DESC LIMIT 1"
	_sharingIndexSQL     = "SELECT table_index FROM answer_history_mapping WHERE hid = ?"

	_addAnswerHistorySQL     = "INSERT INTO answer_history_%d (hid,mid,start_time,step_one_err_num,step_one_complete_time) VALUES (?,?,?,?,?)"
	_setAnswerHistorySQL     = "UPDATE answer_history_%d SET complete_result=?,complete_time=?,score=?,is_first_pass=?,passed_level=?,rank_id=? WHERE id=?"
	_updateAnswerLevelSQL    = "UPDATE answer_history_%d SET is_first_pass = ?,passed_level = ? WHERE id = ? "
	_updateCaptchaSQL        = "UPDATE answer_history_%d SET is_pass_captcha = ? WHERE id = ? "
	_updateStepTwoTimeSQL    = "UPDATE answer_history_%d SET step_two_start_time = ? WHERE id = ? "
	_updateExtraStartTimeSQL = "UPDATE answer_history_%d SET step_extra_start_time = ? WHERE id = ? "
	_updateExtraRetSQL       = "UPDATE answer_history_%d SET step_extra_complete_time = ?,extra_score = ? WHERE id = ? "

	_pendanHistorySQL    = "SELECT hid,status FROM answer_pendant_history WHERE mid = ?"
	_inPendantHistorySQL = "INSERT INTO answer_pendant_history (hid,mid) VALUES (?,?)"
	_upPendantHistorySQL = "UPDATE answer_pendant_history SET status = 1 WHERE mid = ? AND hid = ? AND status = 0"

	_allTypesSQL           = "SELECT id,typename,parentid,lablename FROM ans_v3_question_type ORDER BY parentid;"
	_questionByIdsSQL      = "SELECT id,type_id,question,ans1,ans2,ans3,ans4,mid FROM ans_v3_question WHERE id IN (%s)"
	_questionTypeSQL       = "SELECT id,type_id FROM ans_v3_question WHERE state = ?"
	_questionExtraByIdsSQL = "SELECT id,question,ans,status,origin_id,av_id,source,ctime,mtime FROM answer_extra_question WHERE state = 1 AND id IN (%s)"
	_questionExtraTypeSQL  = "SELECT id,ans FROM answer_extra_question WHERE isdel = 1 and state = 1 limit ?"
)

func hit(id int64) int64 {
	return id % _shard
}

// PendantHistory .
func (d *Dao) PendantHistory(c context.Context, mid int64) (hid int64, status int8, err error) {
	row := d.db.QueryRow(c, _pendanHistorySQL, mid)
	if err = row.Scan(&hid, &status); err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, nil
		}
		log.Error("PendantHistory(%d),error:%+v", mid, err)
		return 0, 0, errors.WithStack(err)
	}
	return hid, status, nil
}

// History get user's answer history by mid.
func (d *Dao) History(c context.Context, mid int64) (res *model.AnswerHistory, err error) {
	res = &model.AnswerHistory{}
	row := d.db.QueryRow(c, fmt.Sprintf(_getAnswerHistorySQL, hit(mid)), mid)
	if err = row.Scan(&res.ID, &res.Hid, &res.Mid, &res.StartTime, &res.StepOneErrTimes, &res.StepOneCompleteTime, &res.StepTwoStartTime, &res.CompleteTime,
		&res.CompleteResult, &res.Score, &res.IsPassCaptcha, &res.IsFirstPass, &res.PassedLevel, &res.RankID, &res.StepExtraStartTime, &res.StepExtraCompleteTime, &res.StepExtraScore); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("History(%d),error:%+v", mid, err)
		}
	}
	return
}

// HistoryByHid get user's answer history by mid.
func (d *Dao) HistoryByHid(c context.Context, hid int64) (res *model.AnswerHistory, err error) {
	res = &model.AnswerHistory{}
	row := d.db.QueryRow(c, fmt.Sprintf(_getHistoryByHidSQL, hit(hid)), hid)
	if err = row.Scan(&res.ID, &res.Hid, &res.Mid, &res.StartTime, &res.StepOneErrTimes, &res.StepOneCompleteTime, &res.StepTwoStartTime, &res.CompleteTime,
		&res.CompleteResult, &res.Score, &res.IsPassCaptcha, &res.IsFirstPass, &res.PassedLevel, &res.RankID, &res.StepExtraStartTime, &res.StepExtraCompleteTime, &res.StepExtraScore); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("HistoryByHid(%d),error:%+v", hid, err)
		}
	}
	return
}

// OldHistory get user's answer history by hid and sharing index.
func (d *Dao) OldHistory(c context.Context, hid, idx int64) (res *model.AnswerHistory, err error) {
	res = &model.AnswerHistory{}
	row := d.db.QueryRow(c, fmt.Sprintf(_getHistoryByHidSQL, idx), hid)
	if err = row.Scan(&res.ID, &res.Hid, &res.Mid, &res.StartTime, &res.StepOneErrTimes, &res.StepOneCompleteTime, &res.StepTwoStartTime, &res.CompleteTime,
		&res.CompleteResult, &res.Score, &res.IsPassCaptcha, &res.IsFirstPass, &res.PassedLevel, &res.RankID, &res.StepExtraStartTime, &res.StepExtraCompleteTime, &res.StepExtraScore); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("OldHistory(%d,%d),error:%+v", hid, idx, err)
		}
	}
	return
}

// SharingIndexByHid get old history sharingIndex by hid
func (d *Dao) SharingIndexByHid(c context.Context, hid int64) (res int64, err error) {
	row := d.db.QueryRow(c, _sharingIndexSQL, hid)
	if err = row.Scan(&res); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("OldHistory(%d),error:%+v", hid, err)
		}
	}
	return
}

// AddPendantHistory .
func (d *Dao) AddPendantHistory(c context.Context, mid, hid int64) (affected int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _inPendantHistorySQL, hid, mid); err != nil {
		log.Error("AddPendantHistory(%d,%d),error:%+v", mid, hid, err)
		return
	}
	return res.RowsAffected()
}

// SetHistory set user's answer history by id.
func (d *Dao) SetHistory(c context.Context, mid int64, his *model.AnswerHistory) (affected int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, fmt.Sprintf(_setAnswerHistorySQL, hit(mid)), his.CompleteResult, his.CompleteTime, his.Score, his.IsFirstPass, his.PassedLevel, his.RankID, his.ID); err != nil {
		log.Error("setAnswerHistory: db.Exec(%d, %v) error(%v)", mid, his, err)
		return
	}
	return res.RowsAffected()
}

// AddHistory add user's answer history by id.
func (d *Dao) AddHistory(c context.Context, mid int64, his *model.AnswerHistory) (affected int64, hid string, err error) {
	var res xsql.Result
	hid = fmt.Sprintf("%d%04d%02d", time.Now().Unix(), rand.Intn(9999), hit(mid))
	if res, err = d.db.Exec(c, fmt.Sprintf(_addAnswerHistorySQL, hit(mid)), hid, his.Mid, his.StartTime, his.StepOneErrTimes, his.StepOneCompleteTime); err != nil {
		log.Error("addAnswerHistory: db.Exec(%d, %v) error(%v)", mid, his, err)
		return
	}
	affected, err = res.RowsAffected()
	return
}

// UpdateLevel update answer history passedLevel and isFirstPass.
func (d *Dao) UpdateLevel(c context.Context, id int64, mid int64, isFirstPass, passedLevel int8) (ret int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, fmt.Sprintf(_updateAnswerLevelSQL, hit(mid)), isFirstPass, passedLevel, id); err != nil {
		log.Error("UpdateLevel: db.Exec(%d,%d,%d,%d) error(%v)", id, mid, isFirstPass, passedLevel, err)
		return
	}
	return res.RowsAffected()
}

// UpdateCaptcha update answer history captcha.
func (d *Dao) UpdateCaptcha(c context.Context, id int64, mid int64, isPassCaptcha int8) (ret int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, fmt.Sprintf(_updateCaptchaSQL, hit(mid)), isPassCaptcha, id); err != nil {
		log.Error("UpdateCaptcha: db.Exec(%d, %d, %d) error(%v)", id, mid, isPassCaptcha, err)
		return
	}
	return res.RowsAffected()
}

// UpdateStepTwoTime .
func (d *Dao) UpdateStepTwoTime(c context.Context, id int64, mid int64, t time.Time) (ret int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, fmt.Sprintf(_updateStepTwoTimeSQL, hit(mid)), t, id); err != nil {
		log.Error("UpdateCaptcha: db.Exec(%d,%d,%s) error(%v)", id, mid, t, err)
		return
	}
	return res.RowsAffected()
}

// UpdateExtraStartTime update extra start time.
func (d *Dao) UpdateExtraStartTime(c context.Context, id int64, mid int64, t time.Time) (ret int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, fmt.Sprintf(_updateExtraStartTimeSQL, hit(mid)), t, id); err != nil {
		log.Error("updateExtraStartTime: db.Exec(%d, %d, %s) error(%v)", id, mid, t, err)
		return
	}
	return res.RowsAffected()
}

// UpdateExtraRet update extra start time.
func (d *Dao) UpdateExtraRet(c context.Context, id int64, mid int64, t int64, extraScore int64) (ret int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, fmt.Sprintf(_updateExtraRetSQL, hit(mid)), t, extraScore, id); err != nil {
		log.Error("updateExtraRetSQL: db.Exec(%d, %d, %d, %d) error(%v)", id, mid, t, extraScore, err)
		return
	}
	return res.RowsAffected()
}

// UpPendantHistory update pendant history.
func (d *Dao) UpPendantHistory(c context.Context, mid, hid int64) (ret int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _upPendantHistorySQL, mid, hid); err != nil {
		log.Error("UpPendantHistory(%d,%d),error:%+v", mid, hid, err)
		return
	}
	return res.RowsAffected()
}

// QidsExtraByState get extra question ids by check
func (d *Dao) QidsExtraByState(c context.Context, size int) (res []*model.ExtraQst, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _questionExtraTypeSQL, size); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.ExtraQst)
		if err = rows.Scan(&r.ID, &r.Ans); err != nil {
			log.Error("QidsExtraByState(%d),error:%+v", size, err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// ExtraByIds get extra question in idstr
func (d *Dao) ExtraByIds(c context.Context, ids []int64) (res map[int64]*model.ExtraQst, err error) {
	var rows *sql.Rows
	res = make(map[int64]*model.ExtraQst, len(ids))
	idStr := xstr.JoinInts(ids)
	if rows, err = d.db.Query(c, fmt.Sprintf(_questionExtraByIdsSQL, idStr)); err != nil {
		log.Error("d.questionExtraByIds.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.ExtraQst)
		if err = rows.Scan(&r.ID, &r.Question, &r.Ans, &r.Status, &r.OriginID, &r.AvID, &r.Source, &r.Ctime, &r.Mtime); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res[r.ID] = r
	}
	err = rows.Err()
	return
}

// Types get all types
func (d *Dao) Types(c context.Context) (res []*model.TypeInfo, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _allTypesSQL); err != nil {
		log.Error("d.allTypesStmt.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.TypeInfo)
		if err = rows.Scan(&r.ID, &r.Name, &r.Parentid, &r.LabelName); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// ByIds get question in idStr
func (d *Dao) ByIds(c context.Context, ids []int64) (res map[int64]*model.Question, err error) {
	var rows *sql.Rows
	res = make(map[int64]*model.Question, len(ids))
	idStr := xstr.JoinInts(ids)
	if rows, err = d.db.Query(c, fmt.Sprintf(_questionByIdsSQL, idStr)); err != nil {
		log.Error("d.queryQuestionByIdsStmt.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.Question)
		ans := make([]string, 4)
		if err = rows.Scan(&r.ID, &r.TypeID, &r.Question, &ans[0], &ans[1], &ans[2], &ans[3], &r.Mid); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		r.Ans = ans
		res[r.ID] = r
	}
	err = rows.Err()
	return
}

// QidsByState get question ids by check
func (d *Dao) QidsByState(c context.Context, state int8) (res []*model.Question, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _questionTypeSQL, state); err != nil {
		log.Error("d.questionTypeStmt.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.Question)
		if err = rows.Scan(&r.ID, &r.TypeID); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}
