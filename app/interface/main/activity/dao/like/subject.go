package like

import (
	"context"
	"fmt"
	"time"

	"go-common/app/interface/main/activity/model/like"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_selSubjectSQL    = "SELECT s.id,s.name,s.dic,s.cover,s.stime,f.interval,f.ltime,f.tlimit FROM act_subject s INNER JOIN act_time_config f ON s.id=f.sid WHERE s.id = ?"
	_votLogSQL        = "INSERT INTO act_online_vote_log(sid,aid,mid,stage,vote) VALUES(?,?,?,?,?)"
	_subjectNewestSQL = "SELECT id,ctime FROM act_subject WHERE state = 1 AND type IN (%s) AND stime <= ? AND etime >= ? ORDER BY ctime DESC LIMIT 1"
	_actSubjectSQL    = "SELECT id,name,dic,cover,stime,etime,flag,type,lstime,letime,act_url,uetime,ustime,level,h5_cover,rank,author,oid,state,ctime,mtime,like_limit,android_url,ios_url,daily_like_limit,daily_single_like_limit FROM act_subject WHERE id = ? and state = 1"
	_subjectInitSQL   = "SELECT id,name,dic,cover,stime,etime,flag,type,lstime,letime,act_url,uetime,ustime,level,h5_cover,rank,author,oid,state,ctime,mtime,like_limit,android_url,ios_url,daily_like_limit,daily_single_like_limit FROM act_subject WHERE id > ? order by id asc limit 1000"
	_subjectMaxIDSQL  = "SELECT id FROM act_subject order by id desc limit 1"
	//SubjectValidState act_subject valid state
	SubjectValidState = 1
)

// Subject Dao sql
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

// VoteLog Dao sql
func (dao *Dao) VoteLog(c context.Context, sid int64, aid int64, mid int64, stage int64, vote int64) (rows int64, err error) {
	rs, err := dao.voteLogStmt.Exec(c, sid, aid, mid, stage, vote)
	if err != nil {
		log.Error("d.VoteLog.Exec(%d, %d,%d, %d, %d) error(%v)", sid, aid, mid, stage, vote, err)
		return
	}
	rows, err = rs.RowsAffected()
	return
}

// NewestSubject get newest subject list.
func (dao *Dao) NewestSubject(c context.Context, typeIDs []int64) (res *like.SubItem, err error) {
	res = new(like.SubItem)
	now := time.Now()
	row := dao.db.QueryRow(c, fmt.Sprintf(_subjectNewestSQL, xstr.JoinInts(typeIDs)), now, now)
	if err = row.Scan(&res.ID, &res.Ctime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			err = errors.Wrap(err, "NewestPage:QueryRow")
		}
	}
	return
}

// RawActSubject get act_subject by id .
func (dao *Dao) RawActSubject(c context.Context, id int64) (res *like.SubjectItem, err error) {
	res = new(like.SubjectItem)
	row := dao.db.QueryRow(c, _actSubjectSQL, id)
	if err = row.Scan(&res.ID, &res.Name, &res.Dic, &res.Cover, &res.Stime, &res.Etime, &res.Flag, &res.Type, &res.Lstime, &res.Letime, &res.ActURL, &res.Uetime, &res.Ustime, &res.Level, &res.H5Cover, &res.Rank, &res.Author, &res.Oid, &res.State, &res.Ctime, &res.Mtime, &res.LikeLimit, &res.AndroidURL, &res.IosURL, &res.DailyLikeLimit, &res.DailySingleLikeLimit); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			err = errors.Wrap(err, "ActSubject:QueryRow")
		}
	}
	return
}

// SubjectListMoreSid get subject more sid .
func (dao *Dao) SubjectListMoreSid(c context.Context, minSid int64) (res []*like.SubjectItem, err error) {
	var rows *sql.Rows
	if rows, err = dao.db.Query(c, _subjectInitSQL, minSid); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			err = errors.Wrap(err, "SubjectInitialize:dao.db.Query()")
		}
		return
	}
	defer rows.Close()
	res = make([]*like.SubjectItem, 0, 1000)
	for rows.Next() {
		a := &like.SubjectItem{}
		if err = rows.Scan(&a.ID, &a.Name, &a.Dic, &a.Cover, &a.Stime, &a.Etime, &a.Flag, &a.Type, &a.Lstime, &a.Letime, &a.ActURL, &a.Uetime, &a.Ustime, &a.Level, &a.H5Cover, &a.Rank, &a.Author, &a.Oid, &a.State, &a.Ctime, &a.Mtime, &a.LikeLimit, &a.AndroidURL, &a.IosURL, &a.DailyLikeLimit, &a.DailySingleLikeLimit); err != nil {
			err = errors.Wrap(err, "SubjectInitialize:rows.Scan()")
			return
		}
		res = append(res, a)
	}
	if err = rows.Err(); err != nil {
		err = errors.Wrap(err, "SubjectInitialize:rows.Err()")
	}
	return
}

// SubjectMaxID get act_subject max id .
func (dao *Dao) SubjectMaxID(c context.Context) (res *like.SubjectItem, err error) {
	res = new(like.SubjectItem)
	row := dao.db.QueryRow(c, _subjectMaxIDSQL)
	if err = row.Scan(&res.ID); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			err = errors.Wrap(err, "SubjectMaxID:QueryRow")
		}
	}
	return
}
