package bws

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	bwsmdl "go-common/app/interface/main/activity/model/bws"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_userAchieveKeyFmt = "bws_u_a_%d_%s"
	_userAchieveCntFmt = "bws_a_c_%d_%s"
	_bwsLotteryKeyFmt  = "bws_l_%s_%d"
	_awardSQL          = "UPDATE act_bws_user_achievements SET award = 2 where `key`= ? AND aid = ?"
	_achievementsSQL   = "SELECT id,`name`,icon,dic,link_type,`unlock`,bid,icon_big,icon_active,icon_active_big,award,ctime,mtime,image,suit_id FROM act_bws_achievements WHERE del = 0  AND bid = ? ORDER BY ID"
	_userAchieveSQL    = "SELECT id,aid,award FROM act_bws_user_achievements WHERE bid = ? AND `key` = ? AND del = 0"
	_userAchieveAddSQL = "INSERT INTO act_bws_user_achievements(bid,aid,award,`key`) VALUES(?,?,?,?)"
	_countAchievesSQL  = "SELECT aid,COUNT(1) AS c FROM act_bws_user_achievements WHERE del = 0 AND bid = ? AND ctime BETWEEN ? AND ? GROUP BY aid HAVING c > 0"
	_nextDayHour       = 16
)

func keyUserAchieve(bid int64, key string) string {
	return fmt.Sprintf(_userAchieveKeyFmt, bid, key)
}

func keyAchieveCnt(bid int64, day string) string {
	return fmt.Sprintf(_userAchieveCntFmt, bid, day)
}

func keyLottery(aid int64, day string) string {
	if day == "" {
		day = time.Now().Format("20060102")
	}
	return fmt.Sprintf(_bwsLotteryKeyFmt, day, aid)
}

// Award  achievement award
func (d *Dao) Award(c context.Context, key string, aid int64) (err error) {
	if _, err = d.db.Exec(c, _awardSQL, key, aid); err != nil {
		log.Error("Award: db.Exec(%d,%s) error(%v)", aid, key, err)
	}
	return
}

// RawAchievements  achievements list
func (d *Dao) RawAchievements(c context.Context, bid int64) (res *bwsmdl.Achievements, err error) {
	var (
		rows *xsql.Rows
		rs   []*bwsmdl.Achievement
	)
	if rows, err = d.db.Query(c, _achievementsSQL, bid); err != nil {
		log.Error("RawAchievements: db.Exec(%d) error(%v)", bid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(bwsmdl.Achievement)
		if err = rows.Scan(&r.ID, &r.Name, &r.Icon, &r.Dic, &r.LockType, &r.Unlock, &r.Bid, &r.IconBig, &r.IconActive, &r.IconActiveBig, &r.Award, &r.Ctime, &r.Mtime, &r.Image, &r.SuitID); err != nil {
			log.Error("RawAchievements:row.Scan() error(%v)", err)
			return
		}
		rs = append(rs, r)
	}
	if len(rs) > 0 {
		res = new(bwsmdl.Achievements)
		res.Achievements = rs
	}
	return
}

// RawUserAchieves get user achievements from db.
func (d *Dao) RawUserAchieves(c context.Context, bid int64, key string) (rs []*bwsmdl.UserAchieve, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, _userAchieveSQL, bid, key); err != nil {
		log.Error("RawUserAchieves: db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(bwsmdl.UserAchieve)
		if err = rows.Scan(&r.ID, &r.Aid, &r.Award); err != nil {
			log.Error("RawUserAchieves:row.Scan() error(%v)", err)
			return
		}
		rs = append(rs, r)
	}
	return
}

// RawAchieveCounts  achievements user count
func (d *Dao) RawAchieveCounts(c context.Context, bid int64, day string) (res []*bwsmdl.CountAchieves, err error) {
	var (
		rows *xsql.Rows
	)
	start, _ := time.Parse("20060102-15:04:05", day+"-00:00:00")
	end, _ := time.Parse("20060102-15:04:05", day+"-23:59:59")
	if rows, err = d.db.Query(c, _countAchievesSQL, bid, start, end); err != nil {
		log.Error("RawCountAchieves: db.Exec(%d) error(%v)", bid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(bwsmdl.CountAchieves)
		if err = rows.Scan(&r.Aid, &r.Count); err != nil {
			log.Error("RawCountAchieves:row.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	return
}

// AddUserAchieve .
func (d *Dao) AddUserAchieve(c context.Context, bid, aid, award int64, key string) (lastID int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _userAchieveAddSQL, bid, aid, award, key); err != nil {
		log.Error("AddUserAchieve error d.db.Exec(%d,%d,%s) error(%v)", bid, aid, key, err)
		return
	}
	return res.LastInsertId()
}

// CacheUserAchieves .
func (d *Dao) CacheUserAchieves(c context.Context, bid int64, key string) (res []*bwsmdl.UserAchieve, err error) {
	var (
		values   []interface{}
		cacheKey = keyUserAchieve(bid, key)
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
		item := new(bwsmdl.UserAchieve)
		if err = json.Unmarshal(bs, &item); err != nil {
			log.Error("CacheUserAchieves conn.Do(ZRANGE, %s) error(%v)", cacheKey, err)
			continue
		}
		res = append(res, item)
	}
	return
}

// AddCacheUserAchieves .
func (d *Dao) AddCacheUserAchieves(c context.Context, bid int64, data []*bwsmdl.UserAchieve, key string) (err error) {
	var bs []byte
	if len(data) == 0 {
		return
	}
	cacheKey := keyUserAchieve(bid, key)
	conn := d.redis.Get(c)
	defer conn.Close()
	args := redis.Args{}.Add(cacheKey)
	for _, v := range data {
		if bs, err = json.Marshal(v); err != nil {
			log.Error("AddCacheUserAchieves json.Marshal() error(%v)", err)
			return
		}
		args = args.Add(v.ID).Add(bs)
	}
	if err = conn.Send("ZADD", args...); err != nil {
		log.Error("AddCacheUserAchieves conn.Send(ZADD, %s, %v) error(%v)", cacheKey, args, err)
		return
	}
	if err = conn.Send("EXPIRE", cacheKey, d.userAchExpire); err != nil {
		log.Error("AddCacheUserAchieves conn.Send(Expire, %s) error(%v)", cacheKey, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("AddCacheUserAchieves conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("AddCacheUserAchieves conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// DelCacheUserAchieves .
func (d *Dao) DelCacheUserAchieves(c context.Context, bid int64, key string) (err error) {
	cacheKey := keyUserAchieve(bid, key)
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("DEL", cacheKey); err != nil {
		log.Error("DelCacheUserAchieves conn.Do(DEL) key(%s) error(%v)", cacheKey, err)
	}
	return
}

// AppendUserAchievesCache .
func (d *Dao) AppendUserAchievesCache(c context.Context, bid int64, key string, achieve *bwsmdl.UserAchieve) (err error) {
	var (
		bs       []byte
		ok       bool
		cacheKey = keyUserAchieve(bid, key)
		conn     = d.redis.Get(c)
	)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("EXPIRE", cacheKey, d.userAchExpire)); err != nil || !ok {
		log.Error("AppendUserAchievesCache conn.Do(EXPIRE %s) error(%v)", key, err)
		return
	}
	args := redis.Args{}.Add(cacheKey)
	if bs, err = json.Marshal(achieve); err != nil {
		log.Error("AppendUserAchievesCache json.Marshal() error(%v)", err)
		return
	}
	args = args.Add(achieve.ID).Add(bs)
	if err = conn.Send("ZADD", args...); err != nil {
		log.Error("AddCacheUserAchieves conn.Send(ZADD, %s, %v) error(%v)", cacheKey, args, err)
		return
	}
	if err = conn.Send("EXPIRE", cacheKey, d.userAchExpire); err != nil {
		log.Error("AddCacheUserAchieves conn.Send(Expire, %s) error(%v)", cacheKey, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("AddCacheUserAchieves conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("AddCacheUserAchieves conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// CacheAchieveCounts  get achieve counts from cache
func (d *Dao) CacheAchieveCounts(c context.Context, bid int64, day string) (res []*bwsmdl.CountAchieves, err error) {
	var (
		bss  []int64
		key  = keyAchieveCnt(bid, day)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if bss, err = redis.Int64s(conn.Do("HGETALL", key)); err != nil {
		log.Error("CacheAchieveCounts conn.Do(HGETALL,%s) error(%v)", key, err)
		return
	}
	for i := 1; i < len(bss); i += 2 {
		item := &bwsmdl.CountAchieves{Aid: bss[i-1], Count: bss[i]}
		res = append(res, item)
	}
	return
}

// AddCacheAchieveCounts set achieve counts  to cache
func (d *Dao) AddCacheAchieveCounts(c context.Context, bid int64, res []*bwsmdl.CountAchieves, day string) (err error) {
	if len(res) == 0 {
		return
	}
	key := keyAchieveCnt(bid, day)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("DEL", key); err != nil {
		log.Error("AddCacheAchieveCounts conn.Send(DEL, %s) error(%v)", key, err)
		return
	}
	args := redis.Args{}.Add(key)
	for _, v := range res {
		args = args.Add(v.Aid).Add(v.Count)
	}
	if err = conn.Send("HMSET", args...); err != nil {
		log.Error("AddCacheAchieveCounts conn.Send(HMSET, %s) error(%v)", key, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.achCntExpire); err != nil {
		log.Error("AddCacheAchieveCounts conn.Send(Expire, %s, %d) error(%v)", key, d.achCntExpire, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("AddCacheAchieveCounts conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// IncrCacheAchieveCounts incr achieve counts  to cache
func (d *Dao) IncrCacheAchieveCounts(c context.Context, bid, aid int64, day string) (err error) {
	var (
		key  = keyAchieveCnt(bid, day)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if err = conn.Send("EXPIRE", key, d.achCntExpire); err != nil {
		log.Error("IncrCacheAchieveCounts conn.Send(Expire, %s, %d) error(%v)", key, d.achCntExpire, err)
		return
	}
	if err = conn.Send("HINCRBY", key, aid, 1); err != nil {
		log.Error("IncrCacheAchieveCounts conn.Send(HMSET, %s, %d) error(%v)", key, aid, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("IncrCacheAchieveCounts conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// DelCacheAchieveCounts delete achieve cnt cache.
func (d *Dao) DelCacheAchieveCounts(c context.Context, bid int64, day string) (err error) {
	cacheKey := keyAchieveCnt(bid, day)
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("DEL", cacheKey); err != nil {
		log.Error("DelCacheAchieveCounts conn.Do(DEL) key(%s) error(%v)", cacheKey, err)
	}
	return
}

// AddLotteryMidCache add lottery mid cache.
func (d *Dao) AddLotteryMidCache(c context.Context, aid, mid int64) (err error) {
	now := time.Now()
	hour := now.Hour()
	dayInt, _ := strconv.ParseInt(now.Format("20060102"), 10, 64)
	if hour >= _nextDayHour {
		dayInt = dayInt + 1
	}
	cacheKey := keyLottery(aid, strconv.FormatInt(dayInt, 10))
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("SADD", cacheKey, mid); err != nil {
		log.Error("AddLotteryCache conn.Do(LPUSH, %s, %d) error(%v)", cacheKey, mid, err)
	}
	return
}

// CacheLotteryMid .
func (d *Dao) CacheLotteryMid(c context.Context, aid int64, day string) (mid int64, err error) {
	var (
		cacheKey = keyLottery(aid, day)
		conn     = d.redis.Get(c)
	)
	defer conn.Close()
	if mid, err = redis.Int64(conn.Do("SPOP", cacheKey)); err != nil && err != redis.ErrNil {
		log.Error("LotteryMidCache SPOP key(%s) error(%v)", cacheKey, err)
	}
	return
}

// CacheLotteryMids .
func (d *Dao) CacheLotteryMids(c context.Context, aid int64, day string) (mids []int64, err error) {
	var cacheKey string
	conn := d.redis.Get(c)
	defer conn.Close()
	cacheKey = keyLottery(aid, day)
	if mids, err = redis.Int64s(conn.Do("SMEMBERS", cacheKey)); err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		log.Error("redis.Int64s(conn.Do(SMEMEBERS,%s)) error(%v)", cacheKey, err)
	}
	return
}
