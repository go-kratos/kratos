package like

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	l "go-common/app/interface/main/activity/model/like"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/stat/prom"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

// like_action sql and like state
const (
	_likeActInfosSQL    = "select id,lid from like_action where lid in (%s) and mid = ?"
	_likeActAddSQL      = "INSERT INTO like_action(lid,mid,action,sid,ipv6,ctime,mtime) VALUES(?,?,?,?,?,?,?)"
	_likeActSumSQL      = "select sum(action) likes, lid from like_action where sid = ? and lid in (%s) group by lid ORDER BY likes desc"
	_storyActSumSQL     = "select sum(action) likes from like_action where sid = ? and mid = ? and ctime >= ? and ctime <= ?"
	_storyEachActSumSQL = "select sum(action) likes from like_action where sid = ? and mid = ? and lid = ? and ctime >= ? and ctime <= ?"
	HasLike             = 1
	NoLike              = -1
	//Total number of activities set the old is bilibili-activity:like:%d
	_likeActScoreKeyFmt = "go:bl-a:l:%d"
	//Total number of comments for different types of manuscripts
	_likeActScoreTyoeKeyFmt = "go:bl:a:l:%d:%d"
	//liked key the old is bilibili-activity:like:%d:%d:%d
	_likeActKeyFmt = "go:bl-act:l:%d:%d:%d"
	//Total number of like the old is likes:oid:%d
	_likeLidKeyFmt = "go:ls:oid:%d"
	//Total number of activities like the old is sb:likes:count:%d
	_likeCountKeyFmt = "go:sb:ls:count:%d"
)

// likeActScoreKey .
func likeActScoreKey(sid int64) string {
	return fmt.Sprintf(_likeActScoreKeyFmt, sid)
}

// likeActScoreTypeKey .
func likeActScoreTypeKey(sid int64, ltype int) string {
	return fmt.Sprintf(_likeActScoreTyoeKeyFmt, ltype, sid)
}

func likeActKey(sid, lid, mid int64) string {
	return fmt.Sprintf(_likeActKeyFmt, sid, lid, mid)
}

// likeLidKey .
func likeLidKey(oid int64) string {
	return fmt.Sprintf(_likeLidKeyFmt, oid)
}

// likeCountKey .
func likeCountKey(sid int64) string {
	return fmt.Sprintf(_likeCountKeyFmt, sid)
}

// LikeActInfos get likesaction logs.
func (dao *Dao) LikeActInfos(c context.Context, lids []int64, mid int64) (likeActs map[int64]*l.Action, err error) {
	var rows *xsql.Rows
	if rows, err = dao.db.Query(c, fmt.Sprintf(_likeActInfosSQL, xstr.JoinInts(lids)), mid); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
		} else {
			err = errors.Wrapf(err, "LikeActInfos:Query(%s)", _likeActInfosSQL)
			return
		}
	}
	defer rows.Close()
	likeActs = make(map[int64]*l.Action, len(lids))
	for rows.Next() {
		a := &l.Action{}
		if err = rows.Scan(&a.ID, &a.Lid); err != nil {
			err = errors.Wrap(err, "LikeActInfos:scan()")
			return
		}
		likeActs[a.Lid] = a
	}
	if err = rows.Err(); err != nil {
		err = errors.Wrap(err, "LikeActInfos:rows.Err()")
	}
	return
}

// LikeActSums get like_action likes sum data .
func (dao *Dao) LikeActSums(c context.Context, sid int64, lids []int64) (res []*l.LidLikeSum, err error) {
	var rows *xsql.Rows
	if rows, err = dao.db.Query(c, fmt.Sprintf(_likeActSumSQL, xstr.JoinInts(lids)), sid); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
		} else {
			err = errors.Wrapf(err, "LikeActSums:Query(%s)", _likeActSumSQL)
			return
		}
	}
	defer rows.Close()
	res = make([]*l.LidLikeSum, 0, len(lids))
	for rows.Next() {
		a := &l.LidLikeSum{}
		if err = rows.Scan(&a.Likes, &a.Lid); err != nil {
			err = errors.Wrapf(err, "LikeActSums:Scan(%s)", _likeActSumSQL)
			return
		}
		res = append(res, a)
	}
	if err = rows.Err(); err != nil {
		err = errors.Wrapf(err, "LikeActSums:rows.Err(%s)", _likeActSumSQL)
	}
	return
}

// StoryLikeActSum .
func (dao *Dao) StoryLikeActSum(c context.Context, sid, mid int64, stime, etime string) (res int64, err error) {
	var tt sql.NullInt64
	row := dao.db.QueryRow(c, _storyActSumSQL, sid, mid, stime, etime)
	if err = row.Scan(&tt); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			err = errors.Wrap(err, "row.Scan()")
		}
	}
	res = tt.Int64
	return
}

// StoryEachLikeAct .
func (dao *Dao) StoryEachLikeAct(c context.Context, sid, mid, lid int64, stime, etime string) (res int64, err error) {
	var tt sql.NullInt64
	row := dao.db.QueryRow(c, _storyEachActSumSQL, sid, mid, lid, stime, etime)
	if err = row.Scan(&tt); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			err = errors.Wrap(err, "row.Scan()")
		}
	}
	res = tt.Int64
	return
}

// SetRedisCache .
func (dao *Dao) SetRedisCache(c context.Context, sid, lid, score int64, likeType int) (err error) {
	var (
		conn        = dao.redis.Get(c)
		key         = likeActScoreKey(sid)
		lidKey      = likeLidKey(lid)
		lidCountKey = likeCountKey(sid)
		max         = 3
	)
	defer conn.Close()
	if err = conn.Send("ZINCRBY", key, score, lid); err != nil {
		err = errors.Wrap(err, "conn.Send(ZINCRBY) likeActScoreKey")
		return
	}
	if err = conn.Send("INCRBY", lidKey, score); err != nil {
		err = errors.Wrap(err, "conn.Send(INCR) likeLidKey")
		return
	}
	if likeType != 0 {
		max++
		typeKey := likeActScoreTypeKey(sid, likeType)
		if err = conn.Send("ZINCRBY", typeKey, score, lid); err != nil {
			err = errors.Wrap(err, "conn.Send(ZINCRBY) likeActScoreTypeKey")
			return
		}
	}
	if err = conn.Send("INCRBY", lidCountKey, score); err != nil {
		err = errors.Wrap(err, "conn.Send(INCR) likeLidKey")
		return
	}
	if err = conn.Flush(); err != nil {
		err = errors.Wrap(err, " conn.Set()")
		return
	}
	for i := 0; i < max; i++ {
		if _, err = conn.Receive(); err != nil {
			err = errors.Wrap(err, fmt.Sprintf("conn.Receive()%d", i+1))
			return
		}
	}
	return
}

// RedisCache get cache order by like .
func (dao *Dao) RedisCache(c context.Context, sid int64, start, end int) (res []*l.LidLikeRes, err error) {
	var (
		conn = dao.redis.Get(c)
		key  = likeActScoreKey(sid)
	)
	defer conn.Close()
	values, err := redis.Values(conn.Do("ZREVRANGE", key, start, end, "WITHSCORES"))
	if err != nil {
		err = errors.Wrap(err, "conn.Do(ZREVRANGE)")
		return
	}
	if len(values) == 0 {
		return
	}
	res = make([]*l.LidLikeRes, 0, len(values))
	for len(values) > 0 {
		t := &l.LidLikeRes{}
		if values, err = redis.Scan(values, &t.Lid, &t.Score); err != nil {
			err = errors.Wrap(err, "redis.Scan")
			return
		}
		res = append(res, t)
	}
	return
}

// LikeActZscore .
func (dao *Dao) LikeActZscore(c context.Context, sid, lid int64) (res int64, err error) {
	var (
		conn = dao.redis.Get(c)
		key  = likeActScoreKey(sid)
	)
	defer conn.Close()
	if res, err = redis.Int64(conn.Do("ZSCORE", key, lid)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			err = errors.Wrap(err, "conn.Do(ZSCORE)")
		}
	}
	return
}

// SetInitializeLikeCache initialize like_action like data .
func (dao *Dao) SetInitializeLikeCache(c context.Context, sid int64, lidLikeAct map[int64]int64, typeLike map[int64]int) (err error) {
	var (
		conn = dao.redis.Get(c)
		max  = 0
		key  = likeActScoreKey(sid)
		args = redis.Args{}.Add(key)
	)
	defer conn.Close()
	for k, val := range lidLikeAct {
		args = args.Add(val).Add(k)
		if typeLike[k] != 0 {
			keyType := likeActScoreTypeKey(sid, typeLike[k])
			argsType := redis.Args{}.Add(keyType).Add(val).Add(k)
			if err = conn.Send("ZADD", argsType...); err != nil {
				log.Error("SetInitializeLikeCache:conn.Send(zadd) args(%v) error(%v)", argsType, err)
				return
			}
			max++
		}
	}
	if err = conn.Send("ZADD", args...); err != nil {
		log.Error("SetInitializeLikeCache:conn.Send(zadd) args(%v) error(%v)", args, err)
		return
	}
	max++
	if err = conn.Flush(); err != nil {
		err = errors.Wrap(err, "SetInitializeLikeCache:conn.Set()")
		return
	}
	for i := 0; i < max; i++ {
		if _, err = conn.Receive(); err != nil {
			err = errors.Wrapf(err, "SetInitializeLikeCache:conn.Receive()%d", i+1)
		}
	}
	return
}

// LikeActAdd add like_action .
func (dao *Dao) LikeActAdd(c context.Context, likeAct *l.Action) (id int64, err error) {
	var res sql.Result
	if res, err = dao.db.Exec(c, _likeActAddSQL, likeAct.Lid, likeAct.Mid, likeAct.Action, likeAct.Sid, likeAct.IPv6, time.Now(), time.Now()); err != nil {
		err = errors.Wrapf(err, "d.db.Exec(%s)", _likeActAddSQL)
		return
	}
	return res.LastInsertId()
}

// LikeActLidCounts get lid score.
func (dao *Dao) LikeActLidCounts(c context.Context, lids []int64) (res map[int64]int64, err error) {
	var (
		conn = dao.redis.Get(c)
		args = redis.Args{}
		ss   []int64
	)
	defer conn.Close()
	for _, lid := range lids {
		args = args.Add(likeLidKey(lid))
	}
	if ss, err = redis.Int64s(conn.Do("MGET", args...)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			err = errors.Wrapf(err, "redis.Ints(conn.Do(HMGET,%v)", args)
		}
		return
	}
	res = make(map[int64]int64, len(lids))
	for key, val := range ss {
		res[lids[key]] = val
	}
	return
}

// LikeActs get data from cache if miss will call source method, then add to cache.
func (dao *Dao) LikeActs(c context.Context, sid, mid int64, lids []int64) (res map[int64]int, err error) {
	var (
		miss         []int64
		likeActInfos map[int64]*l.Action
		missVal      map[int64]int
	)
	if len(lids) == 0 {
		return
	}
	addCache := true
	res, err = dao.CacheLikeActs(c, sid, mid, lids)
	if err != nil {
		addCache = false
		res = nil
		err = nil
	}
	for _, key := range lids {
		if (res == nil) || (res[key] == 0) {
			miss = append(miss, key)
		}
	}
	prom.CacheHit.Add("LikeActs", int64(len(lids)-len(miss)))
	if len(miss) == 0 {
		return
	}
	if likeActInfos, err = dao.LikeActInfos(c, miss, mid); err != nil {
		err = errors.Wrapf(err, "dao.LikeActInfos(%v) error(%v)", miss, err)
		return
	}
	if res == nil {
		res = make(map[int64]int)
	}
	missVal = make(map[int64]int, len(miss))
	for _, mcLid := range miss {
		if _, ok := likeActInfos[mcLid]; ok {
			res[mcLid] = HasLike
		} else {
			res[mcLid] = NoLike
		}
		missVal[mcLid] = res[mcLid]
	}
	if !addCache {
		return
	}
	dao.AddCacheLikeActs(c, sid, mid, missVal)
	return
}

// CacheLikeActs res value val -1:no like 1:has like 0:no value.
func (dao *Dao) CacheLikeActs(c context.Context, sid, mid int64, lids []int64) (res map[int64]int, err error) {
	l := len(lids)
	if l == 0 {
		return
	}
	keysMap := make(map[string]int64, l)
	keys := make([]string, 0, l)
	for _, id := range lids {
		key := likeActKey(sid, id, mid)
		keysMap[key] = id
		keys = append(keys, key)
	}
	conn := dao.mc.Get(c)
	defer conn.Close()
	replies, err := conn.GetMulti(keys)
	if err != nil {
		prom.BusinessErrCount.Incr("mc:CacheLikeActs")
		log.Errorv(c, log.KV("CacheLikeActs", fmt.Sprintf("%+v", err)), log.KV("keys", keys))
		return
	}
	for key, reply := range replies {
		var v string
		err = conn.Scan(reply, &v)
		if err != nil {
			prom.BusinessErrCount.Incr("mc:CacheLikeActs")
			log.Errorv(c, log.KV("CacheLikeActs", fmt.Sprintf("%+v", err)), log.KV("key", key))
			return
		}
		r, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			prom.BusinessErrCount.Incr("mc:CacheLikeActs")
			log.Errorv(c, log.KV("CacheLikeActs", fmt.Sprintf("%+v", err)), log.KV("key", key))
			return res, err
		}
		if res == nil {
			res = make(map[int64]int, len(keys))
		}
		res[keysMap[key]] = int(r)
	}
	return
}

// AddCacheLikeActs Set data to mc
func (dao *Dao) AddCacheLikeActs(c context.Context, sid, mid int64, values map[int64]int) (err error) {
	if len(values) == 0 {
		return
	}
	conn := dao.mc.Get(c)
	defer conn.Close()
	for id, val := range values {
		key := likeActKey(sid, id, mid)
		bs := []byte(strconv.FormatInt(int64(val), 10))
		item := &memcache.Item{Key: key, Value: bs, Expiration: dao.mcPerpetualExpire, Flags: memcache.FlagRAW}
		if err = conn.Set(item); err != nil {
			prom.BusinessErrCount.Incr("mc:AddCacheLikeActs")
			log.Errorv(c, log.KV("AddCacheLikeActs", fmt.Sprintf("%+v", err)), log.KV("key", key))
			return
		}
	}
	return
}
