package dao

import (
	"context"
	"fmt"
	"time"

	model "go-common/app/interface/main/credit/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	ctime "go-common/library/time"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_insAnsLogSQL        = "INSERT INTO blocked_labour_answer_log(mid,score,content,start_time) VALUES(?,?,?,?)"
	_insQsSQL            = "INSERT INTO blocked_labour_question(question,ans,av_id,status,source) VALUES(?,?,?,?,?)"
	_updateQsSQL         = "UPDATE blocked_labour_question SET ans=?,status=? WHERE id=?"
	_delQsSQL            = "UPDATE blocked_labour_question SET isdel=? WHERE id=?"
	_selConfSQL          = "SELECT config_key,content from blocked_config"
	_selNoticeSQL        = "SELECT id,content,url FROM blocked_notice WHERE status=0 ORDER BY id DESC LIMIT 1"
	_selReasonSQL        = "SELECT id,reason,content FROM blocked_reason WHERE status=0"
	_selQuestionSQL      = "SELECT id,question,ans FROM blocked_labour_question WHERE id IN(%s) AND status=2 ORDER BY find_in_set(id,'%s') "
	_selAllQuestionSQL   = "SELECT id,question,ans,av_id,status,source,ctime,mtime FROM blocked_labour_question WHERE id IN(%s)  "
	_isAnsweredSQL       = "SELECT COUNT(*) FROM blocked_labour_answer_log WHERE mid=? AND score=100 AND mtime>=?"
	_noAuditQuestionSQL  = "SELECT id,av_id,question FROM blocked_labour_question WHERE status=2 ORDER BY id DESC LIMIT 20"
	_auditQuestionSQL    = "SELECT id,av_id,question FROM blocked_labour_question WHERE status=1 AND ans=0 ORDER BY id DESC LIMIT 20"
	_selKPISQL           = "SELECT k.mid,k.day,k.rate,k.rank,k.rank_per,k.rank_total,p.point,p.active_days,p.vote_total,p.opinion_likes,p.vote_real_total from blocked_kpi k inner join blocked_kpi_data p on k.mid=p.mid and k.day=p.day where k.mid = ?"
	_announcementInfoSQL = `SELECT id,title,sub_title,publish_status,stick_status,content,url,ctime,mtime FROM blocked_publish WHERE id = ? AND  publish_status = 1 AND status = 0`
	_announcementListSQL = `SELECT id,title,sub_title,publish_status,stick_status,content,url,ptype,ctime,mtime FROM blocked_publish WHERE publish_status = 1 AND status = 0 AND show_time <= ? ORDER BY  ctime desc`
	_publishsSQL         = "SELECT id,title,sub_title,stick_status,content,ctime FROM blocked_publish WHERE id IN (%s) AND publish_status = 1 AND status = 0"
	_selNewKPISQL        = "SELECT rate FROM blocked_kpi WHERE mid=? ORDER BY id DESC LIMIT 1"
)

// AddAnsLog add labour answer log.
func (d *Dao) AddAnsLog(c context.Context, mid int64, score int64, anstr string, stime ctime.Time) (affect int64, err error) {
	result, err := d.db.Exec(c, _insAnsLogSQL, mid, score, anstr, stime)
	if err != nil {
		log.Error("AddAnsLog: db.Exec(mid:%d,score:%d,ans:%s) error(%v)", mid, score, anstr, err)
		return
	}
	affect, err = result.LastInsertId()
	return
}

// AddQs add labour question log.
func (d *Dao) AddQs(c context.Context, qs *model.LabourQs) (err error) {
	if _, err = d.db.Exec(c, _insQsSQL, qs.Question, qs.Ans, qs.AvID, qs.Status, qs.Source); err != nil {
		log.Error("AddQs: db.Exec(as:%v) error(%v)", qs, err)
	}
	return
}

// SetQs set labour question field.
func (d *Dao) SetQs(c context.Context, id int64, ans int64, status int64) (err error) {
	if _, err = d.db.Exec(c, _updateQsSQL, ans, status, id); err != nil {
		log.Error("setQs: db.Exec(ans:%d status:%d) error(%v)", ans, status, err)
	}
	return
}

// DelQs del labour question.
func (d *Dao) DelQs(c context.Context, id int64, isDel int64) (err error) {
	if _, err = d.db.Exec(c, _delQsSQL, isDel, id); err != nil {
		log.Error("setQs: db.Exec(id:%d isDel:%d) error(%v)", id, isDel, err)
	}
	return
}

// LoadConf load conf.
func (d *Dao) LoadConf(c context.Context) (cf map[string]string, err error) {
	cf = make(map[string]string)
	rows, err := d.db.Query(c, _selConfSQL)
	if err != nil {
		log.Error("d.loadConf err(%v)", err)
		return
	}
	defer rows.Close()
	var (
		key   string
		value string
	)
	for rows.Next() {
		if err = rows.Scan(&key, &value); err != nil {
			log.Error("rows.Scan err(%v)", err)
			return
		}
		cf[key] = value
	}
	err = rows.Err()
	return
}

// Notice get notice
func (d *Dao) Notice(c context.Context) (n *model.Notice, err error) {
	row := d.db.QueryRow(c, _selNoticeSQL)
	if err != nil {
		log.Error("Notice: d.QueryRow error(%v)", err)
		return
	}
	n = &model.Notice{}
	if err = row.Scan(&n.ID, &n.Content, &n.URL); err != nil {
		log.Error("row.Scan() error(%v)", err)
	}
	return
}

// ReasonList get reason list
func (d *Dao) ReasonList(c context.Context) (res []*model.Reason, err error) {
	rows, err := d.db.Query(c, _selReasonSQL)
	if err != nil {
		log.Error("reasonList: d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.Reason)
		if err = rows.Scan(&r.ID, &r.Reason, &r.Content); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// QsList get question list.
func (d *Dao) QsList(c context.Context, idStr string) (res []*model.LabourQs, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_selQuestionSQL, idStr, idStr)); err != nil {
		log.Error("d.QuestionList.Query(%s) error(%v)", idStr, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.LabourQs)
		if err = rows.Scan(&r.ID, &r.Question, &r.Ans); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// QsAllList get question list.
func (d *Dao) QsAllList(c context.Context, idStr string) (mlab map[int64]*model.LabourQs, labs []*model.LabourQs, avIDs []int64, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_selAllQuestionSQL, idStr)); err != nil {
		return
	}
	defer rows.Close()
	mlab = make(map[int64]*model.LabourQs, 40)
	for rows.Next() {
		r := new(model.LabourQs)
		if err = rows.Scan(&r.ID, &r.Question, &r.Ans, &r.AvID, &r.Status, &r.Source, &r.Ctime, &r.Mtime); err != nil {
			if err == sql.ErrNoRows {
				err = nil
			}
			mlab = nil
			labs = nil
			avIDs = nil
			return
		}
		mlab[r.ID] = r
		labs = append(labs, r)
		if r.AvID != 0 {
			avIDs = append(avIDs, r.AvID)
		}
	}
	err = rows.Err()
	return
}

// AnswerStatus get blocked user answer status.
func (d *Dao) AnswerStatus(c context.Context, mid int64, ts time.Time) (status bool, err error) {
	row := d.db.QueryRow(c, _isAnsweredSQL, mid, ts)
	var count int64
	if err = row.Scan(&count); err != nil {
		if err != sql.ErrNoRows {
			err = errors.Wrap(err, "AnswerStatus")
			return
		}
		count = 0
		err = nil
	}
	status = count > 0
	return
}

// LastNoAuditQuestion get new no audit question data.
func (d *Dao) LastNoAuditQuestion(c context.Context) (res []*model.LabourQs, avIDs []int64, err error) {
	rows, err := d.db.Query(c, _noAuditQuestionSQL)
	if err != nil {
		log.Error("d.db.Query err(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		labourQs := &model.LabourQs{}
		if err = rows.Scan(&labourQs.ID, &labourQs.AvID, &labourQs.Question); err != nil {
			if err == sql.ErrNoRows {
				err = nil
			} else {
				log.Error("rows.Scan err(%v)", err)
			}
			res = []*model.LabourQs{}
			avIDs = []int64{}
			return
		}
		res = append(res, labourQs)
		if labourQs.AvID != 0 {
			avIDs = append(avIDs, labourQs.AvID)
		}
	}
	err = rows.Err()
	return
}

// LastAuditQuestion get new  audit question data.
func (d *Dao) LastAuditQuestion(c context.Context) (res []*model.LabourQs, avIDs []int64, err error) {
	rows, err := d.db.Query(c, _auditQuestionSQL)
	if err != nil {
		log.Error("d.db.Query err(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		labourQs := &model.LabourQs{}
		if err = rows.Scan(&labourQs.ID, &labourQs.AvID, &labourQs.Question); err != nil {
			if err == sql.ErrNoRows {
				err = nil
			} else {
				log.Error("rows.Scan err(%v)", err)
			}
			res = []*model.LabourQs{}
			avIDs = []int64{}
			return
		}
		res = append(res, labourQs)
		if labourQs.AvID != 0 {
			avIDs = append(avIDs, labourQs.AvID)
		}
	}
	err = rows.Err()
	return
}

// KPIList get kpi list.
func (d *Dao) KPIList(c context.Context, mid int64) (res []*model.KPIData, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _selKPISQL, mid); err != nil {
		log.Error("KpiList: d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.KPIData)
		if err = rows.Scan(&r.Mid, &r.Day, &r.Rate, &r.RankPer, &r.RankPer, &r.RankTotal, &r.Point, &r.ActiveDays, &r.VoteTotal, &r.OpinionLikes, &r.VoteRealTotal); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// AnnouncementInfo get announcement detail.
func (d *Dao) AnnouncementInfo(c context.Context, aid int64) (r *model.BlockedAnnouncement, err error) {
	row := d.db.QueryRow(c, _announcementInfoSQL, aid)
	if err = row.Scan(&r.ID, &r.Title, &r.SubTitle, &r.Content, &r.PublishStatus, &r.StickStatus, &r.URL, &r.Ptype, &r.CTime, &r.MTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		err = errors.Wrap(err, "AnnouncementInfo")
	}
	return
}

// AnnouncementList get accnoucement list.
func (d *Dao) AnnouncementList(c context.Context) (res []*model.BlockedAnnouncement, err error) {
	rows, err := d.db.Query(c, _announcementListSQL, time.Now())
	if err != nil {
		err = errors.Wrap(err, "AnnouncementList")
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.BlockedAnnouncement)
		if err = rows.Scan(&r.ID, &r.Title, &r.SubTitle, &r.PublishStatus, &r.StickStatus, &r.Content, &r.URL, &r.Ptype, &r.CTime, &r.MTime); err != nil {
			if err == sql.ErrNoRows {
				err = nil
				return
			}
			err = errors.Wrap(err, "AnnouncementList")
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// BatchPublishs get publish info list.
func (d *Dao) BatchPublishs(c context.Context, ids []int64) (res map[int64]*model.BlockedAnnouncement, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_publishsSQL, xstr.JoinInts(ids)))
	if err != nil {
		err = errors.Wrap(err, "BatchPublishs")
		return
	}
	res = make(map[int64]*model.BlockedAnnouncement)
	defer rows.Close()
	for rows.Next() {
		r := new(model.BlockedAnnouncement)
		if err = rows.Scan(&r.ID, &r.Title, &r.SubTitle, &r.StickStatus, &r.Content, &r.CTime); err != nil {
			if err == sql.ErrNoRows {
				err = nil
				return
			}
			err = errors.Wrap(err, "BatchPublishs")
			return
		}
		res[r.ID] = r
	}
	err = rows.Err()
	return
}

// NewKPI return user newest KPI rate.
func (d *Dao) NewKPI(c context.Context, mid int64) (rate int8, err error) {
	row := d.db.QueryRow(c, _selNewKPISQL, mid)
	if err != nil {
		log.Error("NewKPI: d.QueryRow error(%v)", err)
		return
	}
	if err = row.Scan(&rate); err != nil {
		if err == sql.ErrNoRows {
			rate = model.KPILevelC
			err = nil
			return
		}
		log.Error("row.Scan() error(%v)", err)
	}
	return
}
