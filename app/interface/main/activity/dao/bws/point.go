package bws

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	bwsmdl "go-common/app/interface/main/activity/model/bws"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_userPointKeyFmt = "bws_u_p_%d_%s"
	_pointsSQL       = "SELECT id,`name`,icon,fid,image,unlocked,lock_type,dic,rule,bid,lose_unlocked,other_ip,ower,ctime,mtime FROM act_bws_points WHERE del = 0 AND bid = ? ORDER BY ID"
	_userPointSQL    = "SELECT id,pid,points,ctime FROM act_bws_user_points WHERE bid = ? AND `key` = ? AND del = 0"
	_userPointAddSQL = "INSERT INTO act_bws_user_points(bid,pid,points,`key`) VALUES(?,?,?,?)"
)

func keyUserPoint(bid int64, key string) string {
	return fmt.Sprintf(_userPointKeyFmt, bid, key)
}

// RawPoints points list
func (d *Dao) RawPoints(c context.Context, bid int64) (res *bwsmdl.Points, err error) {
	var (
		rows *xsql.Rows
		rs   []*bwsmdl.Point
	)
	if rows, err = d.db.Query(c, _pointsSQL, bid); err != nil {
		log.Error("RawPoints: db.Exec(%d) error(%v)", bid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(bwsmdl.Point)
		if err = rows.Scan(&r.ID, &r.Name, &r.Icon, &r.Fid, &r.Image, &r.Unlocked, &r.LockType, &r.Dic, &r.Rule, &r.Bid, &r.LoseUnlocked, &r.OtherIp, &r.Ower, &r.Ctime, &r.Mtime); err != nil {
			log.Error("RawPoints:row.Scan() error(%v)", err)
			return
		}
		rs = append(rs, r)
	}
	if len(rs) > 0 {
		res = new(bwsmdl.Points)
		res.Points = rs
	}
	return
}

// RawUserPoints .
func (d *Dao) RawUserPoints(c context.Context, bid int64, key string) (rs []*bwsmdl.UserPoint, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, _userPointSQL, bid, key); err != nil {
		log.Error("RawUserPoints:db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(bwsmdl.UserPoint)
		if err = rows.Scan(&r.ID, &r.Pid, &r.Points, &r.Ctime); err != nil {
			log.Error("RawUserPoints:row.Scan() error(%v)", err)
			return
		}
		rs = append(rs, r)
	}
	return
}

// AddUserPoint .
func (d *Dao) AddUserPoint(c context.Context, bid, pid, points int64, key string) (lastID int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _userPointAddSQL, bid, pid, points, key); err != nil {
		log.Error("AddUserPoint error d.db.Exec(%d,%d,%d,%s) error(%v)", bid, pid, points, key, err)
		return
	}
	return res.LastInsertId()
}

// CacheUserPoints .
func (d *Dao) CacheUserPoints(c context.Context, bid int64, key string) (res []*bwsmdl.UserPoint, err error) {
	var (
		values   []interface{}
		cacheKey = keyUserPoint(bid, key)
		conn     = d.redis.Get(c)
	)
	defer conn.Close()
	if values, err = redis.Values(conn.Do("ZRANGE", cacheKey, 0, -1)); err != nil {
		log.Error("CacheUserAchieves conn.Do(ZRANGE, %s) error(%v)", cacheKey, err)
		return
	} else if len(values) == 0 {
		return
	}
	for len(values) > 0 {
		var bs []byte
		if values, err = redis.Scan(values, &bs); err != nil {
			log.Error("CacheUserAchieves redis.Scan(%v) error(%v)", values, err)
			return
		}
		item := new(bwsmdl.UserPoint)
		if err = json.Unmarshal(bs, &item); err != nil {
			log.Error("CacheUserAchieves conn.Do(ZRANGE, %s) error(%v)", cacheKey, err)
			continue
		}
		res = append(res, item)
	}
	return
}

// AddCacheUserPoints .
func (d *Dao) AddCacheUserPoints(c context.Context, bid int64, data []*bwsmdl.UserPoint, key string) (err error) {
	var bs []byte
	if len(data) == 0 {
		return
	}
	cacheKey := keyUserPoint(bid, key)
	conn := d.redis.Get(c)
	defer conn.Close()
	args := redis.Args{}.Add(cacheKey)
	for _, v := range data {
		if bs, err = json.Marshal(v); err != nil {
			log.Error("AddCacheUserPoints json.Marshal() error(%v)", err)
			return
		}
		args = args.Add(v.ID).Add(bs)
	}
	if err = conn.Send("ZADD", args...); err != nil {
		log.Error("AddCacheUserPoints conn.Send(ZADD, %s, %v) error(%v)", cacheKey, args, err)
		return
	}
	if err = conn.Send("EXPIRE", cacheKey, d.userPointExpire); err != nil {
		log.Error("AddCacheUserPoints conn.Send(Expire, %s) error(%v)", cacheKey, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("AddCacheUserPoints conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("AddCacheUserPoints conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// AppendUserPointsCache .
func (d *Dao) AppendUserPointsCache(c context.Context, bid int64, key string, point *bwsmdl.UserPoint) (err error) {
	var (
		bs       []byte
		ok       bool
		cacheKey = keyUserPoint(bid, key)
		conn     = d.redis.Get(c)
	)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("EXPIRE", cacheKey, d.userPointExpire)); err != nil || !ok {
		log.Error("AppendUserPointsCache conn.Do(EXPIRE %s) error(%v)", key, err)
		return
	}
	args := redis.Args{}.Add(cacheKey)
	if bs, err = json.Marshal(point); err != nil {
		log.Error("AppendUserPointsCache json.Marshal() error(%v)", err)
		return
	}
	args = args.Add(point.ID).Add(bs)
	if err = conn.Send("ZADD", args...); err != nil {
		log.Error("AppendUserPointsCache conn.Send(ZADD, %s, %v) error(%v)", cacheKey, args, err)
		return
	}
	if err = conn.Send("EXPIRE", cacheKey, d.userPointExpire); err != nil {
		log.Error("AppendUserPointsCache conn.Send(Expire, %s) error(%v)", cacheKey, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("AppendUserPointsCache conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("AppendUserPointsCache conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// DelCacheUserPoints .
func (d *Dao) DelCacheUserPoints(c context.Context, bid int64, key string) (err error) {
	cacheKey := keyUserPoint(bid, key)
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("DEL", cacheKey); err != nil {
		log.Error("DelCacheUserPoints conn.Do(DEL) key(%s) error(%v)", cacheKey, err)
	}
	return
}
