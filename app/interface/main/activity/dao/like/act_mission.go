package like

import (
	"context"
	"database/sql"
	"fmt"

	l "go-common/app/interface/main/activity/model/like"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_likeActMissionItemSQL    = "select id from like_act_mission_group where sid = ? and lid = ? and mid = ?"
	_addLikeActMissionSQL     = "insert into like_act_mission_group (`sid`,`mid`,`lid`,`action`,`ipv6`) values(?,?,?,?,?)"
	_likeActMissionFriendsSQL = "select id,sid,lid,mid,action from like_act_mission_group where sid = ? and lid = ? order by id asc limit 100 "
	_keyMissionLikeLimitFmt   = "go:m:l:limt:%d:%d"
	_keyMissionScoreFmt       = "go:l:a:m:score:%d"
	_keyMissionScoreStrFmt    = "go:l:a:m:score:str:%d:%d"
	_scoreMaxTime             = 2145888000
	_scoreMaxNum              = 1000000
	_maxMissionTop            = 200
)

// missionLikeLimitKey .
func missionLikeLimitKey(sid, mid int64) string {
	return fmt.Sprintf(_keyMissionLikeLimitFmt, sid, mid)
}

// likeActMissionScoreKey .
func likeActMissionScoreKey(sid int64) string {
	return fmt.Sprintf(_keyMissionScoreFmt, sid)
}

// likeMissionScoreStrKey .
func likeMissionScoreStrKey(sid, lid int64) string {
	return fmt.Sprintf(_keyMissionScoreStrFmt, sid, lid)
}

// scoreMaxTime .
func scoreMaxTime() int64 {
	return _scoreMaxTime
}

// scoreMaxNum only for num lower 1 million.
func scoreMaxNum() int64 {
	return _scoreMaxNum
}

// buildRankScore only for num lower 1 million.
func buildRankScore(num int64, ctime int64) int64 {
	var (
		maxTimeStamp = scoreMaxTime()
		timeScore    = maxTimeStamp - ctime
		rankScore    int64
	)
	rankScore = (num & 0xFFFFF) << 32
	rankScore |= timeScore & 0xFFFFFFFF
	return rankScore
}

// RawActMission .
func (d *Dao) RawActMission(c context.Context, sid, lid, mid int64) (res int64, err error) {
	act := &l.ActMissionGroup{}
	rows := d.db.QueryRow(c, _likeActMissionItemSQL, sid, lid, mid)
	if err = rows.Scan(&act.ID); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
		} else {
			err = errors.Wrap(err, "RawActMission:QueryRow")
			return
		}
	}
	res = act.ID
	return
}

// MissionLikeLimit .
func (d *Dao) MissionLikeLimit(c context.Context, sid, mid int64) (res int64, err error) {
	key := missionLikeLimitKey(sid, mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if res, err = redis.Int64(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			err = errors.Wrapf(err, "conn.Do(GET, %s)", key)
		}
	}
	return
}

// InrcMissionLikeLimit .
func (d *Dao) InrcMissionLikeLimit(c context.Context, sid, mid int64, val int64) (res bool, err error) {
	key := missionLikeLimitKey(sid, mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if res, err = redis.Bool(conn.Do("INCRBY", key, val)); err != nil {
		log.Error("InrcMissionLikeLimit:conn.Do(INCR) %s %d", key, val)
	}
	return
}

// SetMissionTop .
func (d *Dao) SetMissionTop(c context.Context, sid, lid int64, score int64, ctime int64) (count int64, err error) {
	var (
		conn      = d.redis.Get(c)
		key       = likeActMissionScoreKey(sid)
		strKey    = likeMissionScoreStrKey(sid, lid)
		maxNum    = scoreMaxNum()
		rankScore int64
	)
	defer conn.Close()
	if count, err = redis.Int64(conn.Do("INCRBY", strKey, score)); err != nil {
		err = errors.Wrap(err, "conn.Do(INCRBY)")
		return
	}
	// the score could not more than 1 million
	if count > maxNum {
		log.Warn("SetMissionTop over max score (%d) error", count)
		return
	}
	rankScore = buildRankScore(count, ctime)
	// set top list
	if _, err = redis.Bool(conn.Do("ZADD", key, rankScore, lid)); err != nil {
		err = errors.Wrapf(err, "conn.Do(ZADD) key(%s)", key)
		return
	}
	if _, err = conn.Do("ZREMRANGEBYRANK", key, 0, -(_maxMissionTop + 1)); err != nil {
		err = errors.Wrapf(err, "conn.Do(ZREMRANGEBYRANK) key(%s)", key)
	}
	return
}

// MissionLidScore .
func (d *Dao) MissionLidScore(c context.Context, sid, lid int64) (score int64, err error) {
	var (
		conn = d.redis.Get(c)
		key  = likeMissionScoreStrKey(sid, lid)
	)
	defer conn.Close()
	if score, err = redis.Int64(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			err = errors.Wrapf(err, "conn.Do(GET) %s", key)
		}
	}
	return
}

// MissionLidRank .
func (d *Dao) MissionLidRank(c context.Context, sid, lid int64) (rank int64, err error) {
	var (
		conn = d.redis.Get(c)
		key  = likeActMissionScoreKey(sid)
	)
	defer conn.Close()
	if rank, err = redis.Int64(conn.Do("ZREVRANK", key, lid)); err != nil {
		if err == redis.ErrNil {
			err = nil
			rank = -1
		} else {
			err = errors.Wrapf(err, "conn.Do(GET) %s", key)
		}
	}
	return
}

// MissionScoreList score list .
func (d *Dao) MissionScoreList(c context.Context, sid int64, start, end int) (data []int64, err error) {
	var (
		key  = likeActMissionScoreKey(sid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	vs, err := redis.Values(conn.Do("ZREVRANGE", key, start, end, "WITHSCORES"))
	if err != nil {
		PromError("redis:获取MissionScoreList", "conn.Do(ZREVRANGE,%s,%d,%d) error(%v)", key, start, end, err)
		return
	}
	data = make([]int64, 0, len(vs))
	for len(vs) > 0 {
		var (
			lid   int64
			score int64
		)
		if vs, err = redis.Scan(vs, &lid, &score); err != nil {
			PromError("redis:获取获取MissionScoreList", "redis.Scan(%v) error(%v)", vs, err)
			return
		}
		if lid != 0 {
			data = append(data, lid)
		}
	}
	return
}

// AddActMission .
func (d *Dao) AddActMission(c context.Context, act *l.ActMissionGroup) (actID int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _addLikeActMissionSQL, act.Sid, act.Mid, act.Lid, act.Action, act.IPv6); err != nil {
		err = errors.Wrapf(err, "d.db.Exec(%s)", _addLikeActMissionSQL)
		return
	}
	return res.LastInsertId()
}

// RawActMissionFriends .
func (d *Dao) RawActMissionFriends(c context.Context, sid, lid int64) (res *l.ActMissionGroups, err error) {
	rows, err := d.db.Query(c, _likeActMissionFriendsSQL, sid, lid)
	if err != nil {
		err = errors.Wrapf(err, " d.db.Query(%s)", _likeActMissionFriendsSQL)
		return
	}
	list := make([]*l.ActMissionGroup, 0, 100)
	for rows.Next() {
		n := &l.ActMissionGroup{}
		if err = rows.Scan(&n.ID, &n.Sid, &n.Lid, &n.Mid, &n.Action); err != nil {
			err = errors.Wrapf(err, " d.db.rows.Scan(%s)", _likeActMissionFriendsSQL)
			return
		}
		list = append(list, n)
	}
	res = &l.ActMissionGroups{ActMissionGroups: list}
	if err = rows.Err(); err != nil {
		err = errors.Wrap(err, "rows.Err()")
	}
	return
}
