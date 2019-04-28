package like

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"go-common/app/job/main/activity/model/like"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_selSubjectSQL      = "SELECT s.id,s.name,s.dic,s.cover,s.stime,f.interval,f.ltime,f.tlimit FROM act_subject s INNER JOIN act_time_config f ON s.id=f.sid WHERE s.id=?"
	_inOnlineLogSQL     = "INSERT INTO act_online_vote_end_log(sid,aid,stage,yes,no) VALUES(?,?,?,?,?)"
	_subjectsSQL        = "SELECT id,name,dic,cover,stime,etime FROM act_subject WHERE state = 1 AND type IN (%s) AND stime <= ? AND etime>= ?"
	_addLotteryTimesURI = "/matsuri/api/add/times"
)

// Subject subject
func (dao *Dao) Subject(c context.Context, sid int64) (n *like.Subject, err error) {
	rows := dao.subjectStmt.QueryRow(c, sid)
	n = &like.Subject{}
	if err = rows.Scan(&n.ID, &n.Name, &n.Dic, &n.Cover, &n.Stime, &n.Interval, &n.Ltime, &n.Tlimit); err != nil {
		if err == sql.ErrNoRows {
			n = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	return
}

// InOnlinelog InOnlinelog
func (dao *Dao) InOnlinelog(c context.Context, sid, aid, stage, yes, no int64) (rows int64, err error) {
	rs, err := dao.inOnlineLog.Exec(c, sid, aid, stage, yes, no)
	if err != nil {
		log.Error("d.InOnlinelog.Exec(%d, %d, %d, %d, %d) error(%v)", sid, aid, stage, yes, no, err)
		return
	}
	return rs.RowsAffected()
}

// SubjectList get online subject list by type.
func (dao *Dao) SubjectList(c context.Context, types []int64, ts time.Time) (res []*like.Subject, err error) {
	rows, err := dao.db.Query(c, fmt.Sprintf(_subjectsSQL, xstr.JoinInts(types)), ts, ts)
	if err != nil {
		err = errors.Wrapf(err, "SubjectList:d.db.Query(%v,%d)", types, ts.Unix())
		return
	}
	defer rows.Close()
	for rows.Next() {
		n := new(like.Subject)
		if err = rows.Scan(&n.ID, &n.Name, &n.Dic, &n.Cover, &n.Stime, &n.Etime); err != nil {
			err = errors.Wrapf(err, "SubjectList:row.Scan row (%v,%d)", types, ts.Unix())
			return
		}
		res = append(res, n)
	}
	if err = rows.Err(); err != nil {
		err = errors.Wrapf(err, "SubjectList:rowsErr(%v,%d)", types, ts.Unix())
	}
	return
}

// SubjectTotalStat total stat.
func (dao *Dao) SubjectTotalStat(c context.Context, sid int64) (rs *like.SubjectTotalStat, err error) {
	req := dao.es.NewRequest(_activity).Index(_activity).WhereEq("state", 1).WhereEq("sid", sid).Sum("click").Sum("likes").Sum("fav").Sum("coin")
	res := new(struct {
		Result struct {
			SumCoin []struct {
				Value float64 `json:"value"`
			} `json:"sum_coin"`
			SumFav []struct {
				Value float64 `json:"value"`
			} `json:"sum_fav"`
			SumLikes []struct {
				Value float64 `json:"value"`
			} `json:"sum_likes"`
			SumClick []struct {
				Value float64 `json:"value"`
			} `json:"sum_click"`
		}
		Page struct {
			Total int `json:"total"`
		}
	})
	if err = req.Scan(c, &res); err != nil || res == nil {
		log.Error("SearchArc req.Scan error(%v)", err)
		return
	}
	rs = &like.SubjectTotalStat{
		SumCoin: int64(res.Result.SumCoin[0].Value),
		SumFav:  int64(res.Result.SumFav[0].Value),
		SumLike: int64(res.Result.SumLikes[0].Value),
		SumView: int64(res.Result.SumClick[0].Value),
		Count:   res.Page.Total,
	}
	return
}

// AddLotteryTimes .
func (dao *Dao) AddLotteryTimes(c context.Context, sid, mid int64) (err error) {
	params := url.Values{}
	params.Set("act_id", strconv.FormatInt(sid, 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
	}
	if err = dao.httpClient.Get(c, dao.addLotteryTimesURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		err = errors.Wrapf(err, "dao.client.Get(%s)", dao.addLotteryTimesURL+"?"+params.Encode())
		return
	}
	if res.Code != ecode.OK.Code() {
		err = ecode.Int(res.Code)
	}
	return
}
