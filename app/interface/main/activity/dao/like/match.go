package like

import (
	"context"
	"database/sql"
	"fmt"

	match "go-common/app/interface/main/activity/model/like"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_matchSidSQL   = "SELECT id,sid,max_stake,stake,name,url,cover,ctime,mtime FROM act_matchs WHERE  status=0 AND sid = ? order by id"
	_matchSQL      = "SELECT id,sid,max_stake,stake,name,url,cover,ctime,mtime FROM act_matchs WHERE  status=0 AND id = ?"
	_objIDSQL      = "SELECT id,match_id,sid,home_name,home_logo,home_score,away_name,away_logo,away_score,result,ctime,mtime,stime,etime,game_stime FROM act_matchs_object WHERE status=0 AND id= ?"
	_objIDsSQL     = "SELECT id,match_id,sid,home_name,home_logo,home_score,away_name,away_logo,away_score,result,ctime,mtime,stime,etime,game_stime FROM act_matchs_object WHERE status=0 AND id IN (%s)"
	_objSIDSQL     = "SELECT o.id,o.match_id,o.sid,o.home_name,o.home_logo,o.home_score,o.away_name,o.away_logo,o.away_score,o.result,o.ctime,o.mtime,o.stime,o.etime,o.game_stime,m.name as match_name FROM act_matchs_object as o INNER JOIN act_matchs as m ON o.match_id=m.id  WHERE o.status=0  AND o.result=0 AND o.sid= ? ORDER BY o.game_stime ASC"
	_userAddSQL    = "INSERT INTO act_match_user_log (mid,match_id,m_o_id,sid,result,stake) VALUES (?,?,?,?,?,?)"
	_userResultSQL = "SELECT id,mid,match_id,m_o_id,sid,result,stake,ctime,mtime FROM act_match_user_log WHERE sid = ? AND mid = ? ORDER BY id DESC"
)

// Match get  Match
func (d *Dao) Match(c context.Context, id int64) (res *match.Match, err error) {
	res = &match.Match{}
	row := d.db.QueryRow(c, _matchSQL, id)
	if err = row.Scan(&res.ID, &res.Sid, &res.MaxStake, &res.Stake, &res.Name, &res.Url, &res.Cover, &res.Ctime, &res.Mtime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("Match:row.Scan error(%v)", err)
		}
	}
	return
}

// ActMatch get Activity Match
func (d *Dao) ActMatch(c context.Context, sid int64) (res []*match.Match, err error) {
	var (
		rows   *xsql.Rows
		matchs []*match.Match
	)
	if rows, err = d.db.Query(c, _matchSidSQL, sid); err != nil {
		log.Error("Match: db.Exec(%d) error(%v)", sid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(match.Match)
		if err = rows.Scan(&r.ID, &r.Sid, &r.MaxStake, &r.Stake, &r.Name, &r.Url, &r.Cover, &r.Ctime, &r.Mtime); err != nil {
			log.Error("Match:row.Scan() error(%v)", err)
			return
		}
		matchs = append(matchs, r)
	}
	res = matchs
	return
}

// Object get  object
func (d *Dao) Object(c context.Context, id int64) (res *match.Object, err error) {
	res = &match.Object{}
	row := d.db.QueryRow(c, _objIDSQL, id)
	if err = row.Scan(&res.ID, &res.MatchId, &res.Sid, &res.HomeName, &res.HomeLogo, &res.HomeScore, &res.AwayName, &res.AwayLogo, &res.AwayScore, &res.Result, &res.Ctime, &res.Mtime, &res.Stime, &res.Etime, &res.GameStime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("Match:row.Scan error(%v)", err)
		}
	}
	return
}

// RawMatchSubjects .
func (d *Dao) RawMatchSubjects(c context.Context, ids []int64) (res map[int64]*match.Object, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_objIDsSQL, xstr.JoinInts(ids))); err != nil {
		log.Error("RawMatchSubjects: d.db.Query(%v) error(%v)", ids, err)
		return
	}
	defer rows.Close()
	res = make(map[int64]*match.Object, len(ids))
	for rows.Next() {
		r := new(match.Object)
		if err = rows.Scan(&r.ID, &r.MatchId, &r.Sid, &r.HomeName, &r.HomeLogo, &r.HomeScore, &r.AwayName, &r.AwayLogo, &r.AwayScore, &r.Result, &r.Ctime, &r.Mtime, &r.Stime, &r.Etime, &r.GameStime); err != nil {
			log.Error("RawMatchSubjects:row.Scan() error(%v)", err)
			return
		}
		res[r.ID] = r
	}
	return
}

// ObjectsUnStart get unstart objects.
func (d *Dao) ObjectsUnStart(c context.Context, sid int64) (res []*match.Object, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, _objSIDSQL, sid); err != nil {
		log.Error("ObjectsUnStart: db.Exec(%d) error(%v)", sid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(match.Object)
		if err = rows.Scan(&r.ID, &r.MatchId, &r.Sid, &r.HomeName, &r.HomeLogo, &r.HomeScore, &r.AwayName, &r.AwayLogo, &r.AwayScore, &r.Result, &r.Ctime, &r.Mtime, &r.Stime, &r.Etime, &r.GameStime, &r.MatchName); err != nil {
			log.Error("ObjectsUnStart:row.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	return
}

// AddGuess add match user log
func (d *Dao) AddGuess(c context.Context, mid, matID, objID, sid, result, stake int64) (lastID int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _userAddSQL, mid, matID, objID, sid, result, stake); err != nil {
		log.Error("AddGuess: db.Exec(%d,%d,%d,%d,%d,%d) error(%v)", mid, matID, objID, sid, result, stake, err)
		return
	}
	return res.LastInsertId()
}

// ListGuess get match user log
func (d *Dao) ListGuess(c context.Context, sid, mid int64) (res []*match.UserLog, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _userResultSQL, sid, mid); err != nil {
		log.Error("ListGuess: db.Exec(%d,%d) error(%v)", sid, mid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(match.UserLog)
		if err = rows.Scan(&r.ID, &r.Mid, &r.MatchId, &r.MOId, &r.Sid, &r.Result, &r.Stake, &r.Ctime, &r.Mtime); err != nil {
			log.Error("ListGuess:row.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	return
}
