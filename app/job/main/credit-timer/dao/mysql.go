package dao

import (
	"context"
	"strconv"
	"time"

	"go-common/app/job/main/credit-timer/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_updateKPISQL         = "INSERT INTO blocked_kpi(mid,day,rate,rank,rank_per,rank_total) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE rate=?,rank=?,rank_per=?,rank_total=?"
	_updateKPIDataSQL     = "INSERT INTO blocked_kpi_data(mid,day,point,active_days,vote_total,vote_radio,blocked_total,opinion_num,opinion_likes,opinion_hates,vote_real_total) VALUES (?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE point=?,active_days=?,vote_total=?,vote_radio=?,blocked_total=?,opinion_num=?,opinion_likes=?,opinion_hates=?,vote_real_total=?"
	_updateKPIPointSQL    = "INSERT INTO blocked_kpi_point(mid,day,point,active_days,vote_total,vote_radio,blocked_total,opinion_num,opinion_likes,opinion_hates) VALUES (?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE point=?,active_days=?,vote_total=?,vote_radio=?,blocked_total=?,opinion_num=?,opinion_likes=?,opinion_hates=?"
	_updateCaseEndTimeSQL = "UPDATE blocked_case SET status=3 WHERE status = 1 AND end_time < ?"
	_updateCaseEndVoteSQL = "UPDATE blocked_case SET end_time=? WHERE status = 1 AND end_time > ? AND (vote_rule+vote_break+vote_delete) >= ?"
	_updateJurySQL        = "UPDATE blocked_jury SET status=2,invalid_reason=? WHERE status=1 AND expired < ?"
	_updateJuryExpiredSQL = "UPDATE blocked_jury SET status=1, expired=? WHERE mid =  ?"
	_updateVoteSQL        = "UPDATE blocked_case_vote SET vote=3 WHERE vote=0 AND expired > ? AND expired < ? "
	_selConfSQL           = "SELECT config_key,content FROM blocked_config"
	_selJurySQL           = "SELECT mid FROM blocked_jury WHERE status=1"
	_selJuryKPISQL        = "SELECT mid FROM blocked_jury WHERE expired BETWEEN ? AND ?"
	_countVoteTotalSQL    = "SELECT COUNT(*) FROM blocked_case_vote v INNER JOIN blocked_case c ON c.id=v.cid AND c.status=4 WHERE v.mid = ? AND v.vote IN (1,2,4) AND v.ctime BETWEEN ? AND ?"
	_countRightViolateSQL = "SELECT COUNT(*) FROM blocked_case_vote v INNER JOIN blocked_case c ON c.id=v.cid AND v.vote IN(1,4) AND c.judge_type = 1 AND c.status=4 WHERE v.mid = ? AND v.ctime BETWEEN ? AND ?"
	_countRightLegalSQL   = "SELECT COUNT(*) FROM blocked_case_vote v INNER JOIN blocked_case c ON c.id=v.cid AND v.vote = 2 AND c.judge_type = 2 AND c.status=4 WHERE v.mid = ? AND v.ctime BETWEEN ? AND ?"
	_CountBlockedSQL      = "SELECT COUNT(*) FROM blocked_info WHERE uid = ? AND status=0 AND ctime BETWEEN ? AND ?"
	_selKPIPointDaySQL    = "SELECT k.mid,k.day,k.point,k.active_days,k.vote_total,k.vote_radio,k.blocked_total,k.opinion_num,k.opinion_likes,k.opinion_hates,j.expired FROM blocked_kpi_point k INNER JOIN blocked_jury j ON k.mid = j.mid WHERE day = ? ORDER BY k.point desc"
	_selKPIPointSQL       = "SELECT mid,day,point,active_days,vote_total,vote_radio,blocked_total FROM blocked_kpi_point WHERE mid = ? AND day = ? ORDER BY point desc"
	_selKPISQL            = "SELECT mid,day,rate,rank,rank_per,rank_total FROM blocked_kpi WHERE mid = ?"
	_countActiveSQL       = "SELECT COUNT(*) FROM (SELECT DATE_FORMAT(ctime,'%Y-%m-%d') FROM blocked_case_vote WHERE vote IN(1,2,4) AND mid=? AND ctime BETWEEN ? AND ? GROUP BY DATE_FORMAT(ctime,'%Y-%m-%d')) t"
	_countOpinionSQL      = "SELECT COUNT(*) FROM blocked_opinion WHERE mid = ? AND state = 0 AND ctime BETWEEN ? AND ?"
	_opinionQualitySQL    = "SELECT COALESCE(SUM(likes),0),COALESCE(SUM(hates),0) FROM blocked_opinion WHERE mid = ? AND state = 0 AND ctime BETWEEN ? AND ?"
	_countVoteByTimeSQL   = "SELECT count(*) from blocked_case_vote where vote in(1,2,4) and mid=? and ctime between ? and ?"
)

// UpdateKPI update KPI info.
func (d *Dao) UpdateKPI(c context.Context, r *model.Kpi) (err error) {
	if _, err = d.db.Exec(c, _updateKPISQL, r.Mid, r.Day, r.Rate, r.Rank, r.RankPer, r.RankTotal, r.Rate, r.Rank, r.RankPer, r.RankTotal); err != nil {
		log.Error("d.UpdateKPI err(%v)", err)
	}
	return
}

// UpdateKPIData update kpi_data info.
func (d *Dao) UpdateKPIData(c context.Context, r *model.KpiData) (err error) {
	if _, err = d.db.Exec(c, _updateKPIDataSQL, r.Mid, r.Day, r.Point, r.ActiveDays, r.VoteTotal, r.VoteRadio, r.BlockedTotal, r.OpinionNum, r.OpinionLikes, r.OpinionHates, r.VoteRealTotal, r.Point, r.ActiveDays, r.VoteTotal, r.VoteRadio, r.BlockedTotal, r.OpinionNum, r.OpinionLikes, r.OpinionHates, r.VoteRealTotal); err != nil {
		log.Error("d.UpdateKPIPoint err(%v)", err)
	}
	return
}

// UpdateKPIPoint update kpi point info.
func (d *Dao) UpdateKPIPoint(c context.Context, r *model.KpiPoint) (err error) {
	if _, err = d.db.Exec(c, _updateKPIPointSQL, r.Mid, r.Day, r.Point, r.ActiveDays, r.VoteTotal, r.VoteRadio, r.BlockedTotal, r.OpinionNum, r.OpinionLikes, r.OpinionHates, r.Point, r.ActiveDays, r.VoteTotal, r.VoteRadio, r.BlockedTotal, r.OpinionNum, r.OpinionLikes, r.OpinionHates); err != nil {
		log.Error("d.UpdateKPIPoint err(%v)", err)
	}
	return
}

// UpdateCaseEndTime update case status to CaseStatusDealing which expired time less than now.
func (d *Dao) UpdateCaseEndTime(c context.Context, now time.Time) (affect int64, err error) {
	rows, err := d.db.Exec(c, _updateCaseEndTimeSQL, now)
	if err != nil {
		log.Error("d.UpdateCaseEndTime err(%v)", err)
		return
	}
	return rows.RowsAffected()
}

// UpdateCaseEndVote update case status to CaseStatusDealing which vote total more than conf case vote total.
func (d *Dao) UpdateCaseEndVote(c context.Context, vt int64, ts time.Time) (affect int64, err error) {
	rows, err := d.db.Exec(c, _updateCaseEndVoteSQL, ts, ts, vt)
	if err != nil {
		log.Error("d.UpdateCaseEndVote err(%v)", err)
		return
	}
	return rows.RowsAffected()
}

// UpdateJury update jury status to expired which expired time less than ts.
func (d *Dao) UpdateJury(c context.Context, now time.Time) (affect int64, err error) {
	rows, err := d.db.Exec(c, _updateJurySQL, model.JuryExpire, now)
	if err != nil {
		log.Error("d.UpdateJury err(%v)", err)
		return
	}
	return rows.RowsAffected()
}

// UpdateJuryExpired update jury expired.
func (d *Dao) UpdateJuryExpired(c context.Context, mid int64, expired time.Time) (err error) {
	if _, err = d.db.Exec(c, _updateJuryExpiredSQL, expired, mid); err != nil {
		log.Error("d.UpdateJuryExpired err(%v)", err)
	}
	return
}

// UpdateVote update vote status to give up which do not vote and expired less than ts.
func (d *Dao) UpdateVote(c context.Context, now time.Time) (affect int64, err error) {
	rows, err := d.db.Exec(c, _updateVoteSQL, now.Add(-4*time.Hour), now)
	if err != nil {
		log.Error("d.updateVote err(%v)", err)
		return
	}
	return rows.RowsAffected()
}

// LoadConf load conf.
func (d *Dao) LoadConf(c context.Context) (vTotal int64, err error) {
	var (
		key   string
		value string
	)
	rows, err := d.db.Query(c, _selConfSQL)
	if err != nil {
		log.Error("d.loadConf err(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&key, &value); err != nil {
			log.Error("rows.Scan err(%v)", err)
			return
		}
		switch key {
		case "case_vote_max":
			if vTotal, err = strconv.ParseInt(value, 10, 64); err != nil {
				return
			}
		}
	}
	err = rows.Err()
	return
}

// JuryList get jury list.
func (d *Dao) JuryList(c context.Context) (res []int64, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Prepared(_selJurySQL).Query(c); err != nil {
		log.Error("dao.JuryList error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r int64
		if err = rows.Scan(&r); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// JuryKPI get jury list.
func (d *Dao) JuryKPI(c context.Context, begin, end string) (res []int64, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Prepared(_selJuryKPISQL).Query(c, begin, end); err != nil {
		log.Error("dao.JuryKPI error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r int64
		if err = rows.Scan(&r); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// CountVoteTotal get vote total.
func (d *Dao) CountVoteTotal(c context.Context, mid int64, begin, end string) (count int64, err error) {
	row := d.db.QueryRow(c, _countVoteTotalSQL, mid, begin, end)
	if err = row.Scan(&count); err != nil {
		log.Error("d.CountVoteTotal err(%v)", err)
	}
	return
}

// CountVoteRightViolate get vote right violate count.
func (d *Dao) CountVoteRightViolate(c context.Context, mid int64, begin, end string) (count int64, err error) {
	row := d.db.QueryRow(c, _countRightViolateSQL, mid, begin, end)
	if err = row.Scan(&count); err != nil {
		log.Error("d.CountVoteRightViolate err(%v)", err)
	}
	return
}

// CountVoteRightLegal get vote right legal count.
func (d *Dao) CountVoteRightLegal(c context.Context, mid int64, begin, end string) (count int64, err error) {
	row := d.db.QueryRow(c, _countRightLegalSQL, mid, begin, end)
	if err = row.Scan(&count); err != nil {
		log.Error("d.CountVoteRightLegal err(%v)", err)
	}
	return
}

// CountBlocked get user block count ofter ts.
func (d *Dao) CountBlocked(c context.Context, mid int64, begin, end string) (count int64, err error) {
	row := d.db.QueryRow(c, _CountBlockedSQL, mid, begin, end)
	if err = row.Scan(&count); err != nil {
		log.Error("d.CountBlocked err(%v)", err)
	}
	return
}

// KPIPointDay get KPI point day list.
func (d *Dao) KPIPointDay(c context.Context, day string) (res []model.KpiPoint, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _selKPIPointDaySQL, day); err != nil {
		log.Error("dao.JuryKpi error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := model.KpiPoint{}
		if err = rows.Scan(&r.Mid, &r.Day, &r.Point, &r.ActiveDays, &r.VoteTotal, &r.VoteRadio, &r.BlockedTotal, &r.OpinionNum, &r.OpinionLikes, &r.OpinionHates, &r.Expired); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// KPIPoint get kpi point.
func (d *Dao) KPIPoint(c context.Context, mid int64, day string) (r model.KpiPoint, err error) {
	row := d.db.QueryRow(c, _selKPIPointSQL, mid, day)
	if err = row.Scan(&r.Mid, &r.Day, &r.Point, &r.ActiveDays, &r.VoteTotal, &r.VoteRadio, &r.BlockedTotal); err != nil {
		log.Error("d.KPIPoint err(%v)", err)
	}
	return
}

// KPIList get kpi list.
func (d *Dao) KPIList(c context.Context, mid int64) (res []model.Kpi, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _selKPISQL, mid); err != nil {
		log.Error("dao.KPIList error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := model.Kpi{}
		if err = rows.Scan(&r.Mid, &r.Day, &r.Rate, &r.Rank, &r.RankPer, &r.RankTotal); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// CountVoteActive get vote active days count.
func (d *Dao) CountVoteActive(c context.Context, mid int64, begin, end string) (count int64, err error) {
	row := d.db.QueryRow(c, _countActiveSQL, mid, begin, end)
	if err = row.Scan(&count); err != nil {
		log.Error("d.CountVoteActive err(%v)", err)
	}
	return
}

// CountOpinion count user opinion in 30 days.
func (d *Dao) CountOpinion(c context.Context, mid int64, begin, end string) (count int64, err error) {
	row := d.db.QueryRow(c, _countOpinionSQL, mid, begin, end)
	if err = row.Scan(&count); err != nil {
		log.Error("d.CountOpinion(mid:%d begin:%s end:%s) err(%v)", mid, begin, end, err)
	}
	return
}

// OpinionQuality count user opinion quality(fields likes - hates) in 30days.
func (d *Dao) OpinionQuality(c context.Context, mid int64, begin, end string) (likes, hates int64, err error) {
	row := d.db.QueryRow(c, _opinionQualitySQL, mid, begin, end)
	if err = row.Scan(&likes, &hates); err != nil {
		if err != sql.ErrNoRows {
			log.Error("d.OpinionQuality(mid:%d begin:%s end:%s) err(%v)", mid, begin, end, err)
			return
		}
		err = nil
	}
	return
}

// CountVoteByTime count vote by time.
func (d *Dao) CountVoteByTime(c context.Context, mid int64, begin, end time.Time) (count int64, err error) {
	var row *sql.Row
	if row = d.db.QueryRow(c, _countVoteByTimeSQL, mid, begin, end); err != nil {
		log.Error("d.CountVoteByTime.Query(%d) error(%v)", mid, err)
		return
	}
	if err = row.Scan(&count); err != nil {
		log.Error("row.Scan() error(%v)", err)
	}
	return
}
