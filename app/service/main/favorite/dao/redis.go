package dao

import (
	"context"
	"encoding/json"
	"fmt"

	"go-common/app/service/main/favorite/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
	xtime "go-common/library/time"
)

const (
	_folderKey       = "fi_%d_%d" // sortedset fi_type_mid value:fid,score:ctime
	_relationKey     = "r_%d_%d"  // sortedset r_mid_fid(mtime, oid)
	_allRelationKey  = "ar_%d_%d" // sortedset ar_mid_fid(mtime, oid)
	_relationOidsKey = "ro_%d_%d" // set ro_type_mid value:oids
	_cleanedKey      = "rc_%d_%d" // hash key:rc_type_mid field:fid value:timestamp

	// key fb_mid/100000  offset => mid%100000
	// bit value 1 mean unfaved; bit value 0 mean faved
	_favedBit = "fb_%d_%d"
	_bucket   = 100000
)

// isFavedKey return user's fav flag key
func favedBitKey(tp int8, mid int64) string {
	return fmt.Sprintf(_favedBit, tp, mid/_bucket)
}

// folderKey return a user folder key.
func folderKey(tp int8, mid int64) string {
	return fmt.Sprintf(_folderKey, tp, mid)
}

// relationKey return folder relation key.
func relationKey(mid, fid int64) string {
	return fmt.Sprintf(_relationKey, mid, fid)
}

// relationKey return folder relation key.
func allRelationKey(mid, fid int64) string {
	return fmt.Sprintf(_allRelationKey, mid, fid)
}

// relationOidsKey return a user oids key.
func relationOidsKey(tp int8, mid int64) string {
	return fmt.Sprintf(_relationOidsKey, tp, mid)
}

// isCleanedKey return user whether cleaned key.
func isCleanedKey(tp int8, mid int64) string {
	return fmt.Sprintf(_cleanedKey, tp, mid)
}

// pingRedis check redis connection
func (d *Dao) pingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}

// ExpireRelations expire folder relations cache.
func (d *Dao) ExpireRelations(c context.Context, mid, fid int64) (ok bool, err error) {
	key := relationKey(mid, fid)
	conn := d.redis.Get(c)
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.redisExpire)); err != nil {
		log.Error("conn.Do(EXPIRE %s) error(%v)", key, err)
	}
	conn.Close()
	return
}

// ExpireAllRelations expire folder relations cache.
func (d *Dao) ExpireAllRelations(c context.Context, mid, fid int64) (ok bool, err error) {
	key := allRelationKey(mid, fid)
	conn := d.redis.Get(c)
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.redisExpire)); err != nil {
		log.Error("conn.Do(EXPIRE %s) error(%v)", key, err)
	}
	conn.Close()
	return
}

// ExpireFolder expire folder cache.
func (d *Dao) ExpireFolder(c context.Context, tp int8, mid int64) (ok bool, err error) {
	key := folderKey(tp, mid)
	conn := d.redis.Get(c)
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.redisExpire)); err != nil {
		log.Error("conn.Do(EXPIRE %s) error(%v)", key, err)
	}
	conn.Close()
	return
}

// RemFidsRedis del user's fids in redis
func (d *Dao) RemFidsRedis(c context.Context, typ int8, mid int64, fs ...*model.Folder) (err error) {
	var (
		key  = folderKey(typ, mid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	args := []interface{}{key}
	for _, f := range fs {
		args = append(args, f.ID)
	}
	if err = conn.Send("ZREM", args...); err != nil {
		log.Error("conn.Send(ZREM %s,%v) error(%v)", key, args, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.redisExpire); err != nil {
		log.Error("conn.Send(EXPIRE) error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	for i := 0; i < len(fs)+1; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// AddFidsRedis set user's fids to redis
func (d *Dao) AddFidsRedis(c context.Context, typ int8, mid int64, fs ...*model.Folder) (err error) {
	var (
		key  = folderKey(typ, mid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	for _, f := range fs {
		if err = conn.Send("ZADD", key, f.CTime, f.ID); err != nil {
			log.Error("conn.Send(ZADD %s,%s,%d) error(%v)", key, f.CTime, f.ID, err)
			return
		}
	}
	if err = conn.Send("EXPIRE", key, d.redisExpire); err != nil {
		log.Error("conn.Send(EXPIRE) error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	for i := 0; i < len(fs)+1; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// FidsRedis return the user's all fids from redis.
func (d *Dao) FidsRedis(c context.Context, tp int8, mid int64) (fids []int64, err error) {
	var (
		key  = folderKey(tp, mid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if fids, err = redis.Int64s(conn.Do("ZRANGE", key, 0, -1)); err != nil {
		log.Error("conn.Do(ZRANGE, %s) error(%v)", key, err)
		return
	}
	return
}

// DelFidsRedis delete user's all fids from redis.
func (d *Dao) DelFidsRedis(c context.Context, typ int8, mid int64) (err error) {
	var (
		key  = folderKey(typ, mid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if _, err = conn.Do("DEL", key); err != nil {
		log.Error("conn.Do(DEL %s) error(%v)", key, err)
	}
	return
}

// AddFoldersCache add the user all folders to redis.
func (d *Dao) AddFoldersCache(c context.Context, tp int8, mid int64, folders []*model.Folder) (err error) {
	var (
		folder *model.Folder
		value  []byte
		key    = folderKey(tp, mid)
		conn   = d.redis.Get(c)
	)
	defer conn.Close()
	for _, folder = range folders {
		if value, err = json.Marshal(folder); err != nil {
			return
		}
		if err = conn.Send("HSET", key, folder.ID, value); err != nil {
			log.Error("conn.Send(HSET %s,%d,%s) error(%v)", key, folder.ID, value)
			return
		}
	}
	if err = conn.Send("EXPIRE", key, d.redisExpire); err != nil {
		log.Error("conn.Send(EXPIRE) error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < len(folders)+1; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive error(%v)", err)
			return
		}
	}
	return
}

// FolderRelationsCache return the folder all relations from redis.
func (d *Dao) FolderRelationsCache(c context.Context, typ int8, mid, fid int64, start, end int) (res []*model.Favorite, err error) {
	conn := d.redis.Get(c)
	key := relationKey(mid, fid)
	defer conn.Close()
	values, err := redis.Values(conn.Do("ZREVRANGE", key, start, end, "WITHSCORES"))
	if err != nil {
		log.Error("conn.Do(ZREVRANGE,%s,%d,%d) error(%v)", key, start, end, err)
		return
	}
	if len(values) == 0 {
		return
	}
	var oid, t int64
	for len(values) > 0 {
		if values, err = redis.Scan(values, &oid, &t); err != nil {
			log.Error("redis.Scan() error(%v)", err)
			return
		}
		res = append(res, &model.Favorite{Oid: oid, Mid: mid, Fid: fid, Type: typ, MTime: xtime.Time(t)})
	}
	return
}

// CntRelationsCache return the folder all relation count from redis.
func (d *Dao) CntRelationsCache(c context.Context, mid, fid int64) (cnt int, err error) {
	var conn = d.redis.Get(c)
	key := relationKey(mid, fid)
	defer conn.Close()
	if cnt, err = redis.Int(conn.Do("ZCARD", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		}
		log.Error("conn.Do(ZCARD %s) error(%v)", key, err)
		return
	}
	return
}

// CntAllRelationsCache return the folder all relation count from redis.
func (d *Dao) CntAllRelationsCache(c context.Context, mid, fid int64) (cnt int, err error) {
	var conn = d.redis.Get(c)
	key := allRelationKey(mid, fid)
	defer conn.Close()
	if cnt, err = redis.Int(conn.Do("ZCARD", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		}
		log.Error("conn.Do(ZCARD %s) error(%v)", key, err)
		return
	}
	return
}

// FolderRelationsCache return the folder all relations from redis.
func (d *Dao) FolderAllRelationsCache(c context.Context, typ int8, mid, fid int64, start, end int) (res []*model.Favorite, err error) {
	conn := d.redis.Get(c)
	key := allRelationKey(mid, fid)
	defer conn.Close()
	values, err := redis.Values(conn.Do("ZREVRANGE", key, start, end, "WITHSCORES"))
	if err != nil {
		log.Error("conn.Do(ZREVRANGE,%s,%d,%d) error(%v)", key, start, end, err)
		return
	}
	if len(values) == 0 {
		return
	}
	var oid, score int64
	for len(values) > 0 {
		if values, err = redis.Scan(values, &oid, &score); err != nil {
			log.Error("redis.Scan() error(%v)", err)
			return
		}
		res = append(res, &model.Favorite{Oid: oid / 100, Mid: mid, Fid: fid, Type: int8(oid % 100), MTime: xtime.Time(score % 1e10)})
	}
	return
}

// MultiExpireRelations expire folders's relations cache.
func (d *Dao) MultiExpireAllRelations(c context.Context, mid int64, fids []int64) (map[int64]bool, error) {
	okMap := make(map[int64]bool, len(fids))
	conn := d.redis.Get(c)
	defer conn.Close()
	for _, fid := range fids {
		key := allRelationKey(mid, fid)
		if err := conn.Send("EXPIRE", key, d.redisExpire); err != nil {
			log.Error("redis.send(EXPIRE %s) error(%v)", key, err)
			return nil, err
		}
	}
	if err := conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return nil, err
	}
	for _, fid := range fids {
		ok, err := redis.Bool(conn.Receive())
		if err != nil {
			log.Error("conn.Receive(%d) error(%v)", fid, err)
			return nil, err
		}
		okMap[fid] = ok
	}
	return okMap, nil
}

// MultiExpireRelations expire folders's relations cache.
func (d *Dao) MultiExpireRelations(c context.Context, mid int64, fids []int64) (map[int64]bool, error) {
	okMap := make(map[int64]bool, len(fids))
	conn := d.redis.Get(c)
	defer conn.Close()
	for _, fid := range fids {
		key := relationKey(mid, fid)
		if err := conn.Send("EXPIRE", key, d.redisExpire); err != nil {
			log.Error("redis.send(EXPIRE %s) error(%v)", key, err)
			return nil, err
		}
	}
	if err := conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return nil, err
	}
	for _, fid := range fids {
		ok, err := redis.Bool(conn.Receive())
		if err != nil {
			log.Error("conn.Receive(%d) error(%v)", fid, err)
			return nil, err
		}
		okMap[fid] = ok
	}
	return okMap, nil
}

// ExpireRelationOids set expire for faved oids.
func (d *Dao) ExpireRelationOids(c context.Context, tp int8, mid int64) (ok bool, err error) {
	key := relationOidsKey(tp, mid)
	var conn = d.redis.Get(c)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.redisExpire)); err != nil {
		log.Error("conn.Do(EXPIRE, %s) error(%v)", key, err)
	}
	return
}

// AddRelationOidCache add favoured oid.
func (d *Dao) AddRelationOidCache(c context.Context, tp int8, mid, oid int64) (err error) {
	var (
		key  = relationOidsKey(tp, mid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if err = conn.Send("SADD", key, oid); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.redisExpire); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
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

// RemRelationOidCache del faved oid.
func (d *Dao) RemRelationOidCache(c context.Context, tp int8, mid, oid int64) (err error) {
	var (
		key  = relationOidsKey(tp, mid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if err = conn.Send("SREM", key, oid); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.redisExpire); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
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

// IsFavedCache return true or false to judge object whether faved by user.
func (d *Dao) IsFavedCache(c context.Context, tp int8, mid int64, oid int64) (isFaved bool, err error) {
	var (
		key  = relationOidsKey(tp, mid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if isFaved, err = redis.Bool(conn.Do("SISMEMBER", key, oid)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("HGET %v %v error(%v)", key, oid, err)
		}
	}
	return
}

// IsFavedsCache return true or false to judge object whether faved by user.
func (d *Dao) IsFavedsCache(c context.Context, tp int8, mid int64, oids []int64) (favoreds map[int64]bool, err error) {
	var (
		key  = relationOidsKey(tp, mid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	for _, oid := range oids {
		if err = conn.Send("SISMEMBER", key, oid); err != nil {
			log.Error("conn.Send(SISMEMBER %s,%d) error(%v)", key, oid, err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	favoreds = make(map[int64]bool, len(oids))
	for _, oid := range oids {
		faved, err := redis.Bool(conn.Receive())
		if err != nil {
			log.Error("conn.Receive() error(%v)", err)
			continue
		}
		favoreds[oid] = faved
	}
	return
}

// SetFavedBit set unfaved user bit to 0
func (d *Dao) SetFavedBit(c context.Context, tp int8, mid int64) (err error) {
	key := favedBitKey(tp, mid)
	offset := mid % _bucket
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("SETBIT", key, offset, 0); err != nil {
		log.Error("conn.DO(SETBIT) key(%s) offset(%d) err(%v)", key, offset, err)
	}
	return
}

// SetUnFavedBit set unfaved user bit to 1
func (d *Dao) SetUnFavedBit(c context.Context, tp int8, mid int64) (err error) {
	key := favedBitKey(tp, mid)
	offset := mid % _bucket
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("SETBIT", key, offset, 1); err != nil {
		log.Error("conn.DO(SETBIT) key(%s) offset(%d) err(%v)", key, offset, err)
	}
	return
}

// FavedBit check if user had fav video by bit value
func (d *Dao) FavedBit(c context.Context, tp int8, mid int64) (unfaved bool, err error) {
	key := favedBitKey(tp, mid)
	offset := mid % _bucket
	conn := d.redis.Get(c)
	defer conn.Close()
	if unfaved, err = redis.Bool(conn.Do("GETBIT", key, offset)); err != nil {
		log.Error("conn.DO(GETBIT) key(%s) offset(%d) err(%v)", key, offset, err)
	}
	return
}

// DelRelationsCache delete the folder relation cache.
func (d *Dao) DelAllRelationsCache(c context.Context, mid, fid int64) (err error) {
	key := allRelationKey(mid, fid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("DEL", key); err != nil {
		log.Error("conn.Do(DEL %s) error(%v)", key, err)
	}
	return
}

// DelRelationsCache delete the folder relation cache.
func (d *Dao) DelRelationsCache(c context.Context, mid, fid int64) (err error) {
	key := relationKey(mid, fid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("DEL", key); err != nil {
		log.Error("conn.Do(DEL %s) error(%v)", key, err)
	}
	return
}

// AddRelationCache add a relation to redis.
func (d *Dao) AddRelationCache(c context.Context, m *model.Favorite) (err error) {
	key := relationKey(m.Mid, m.Fid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("ZADD", key, m.MTime, m.Oid); err != nil {
		log.Error("conn.Send(ZADD %s,%d) error(%v)", key, m.Oid, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.redisExpire); err != nil {
		log.Error("conn.Send(EXPIRE) error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
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

// RecentOidsCache return three recent oids of all user's folders from redis.
func (d *Dao) RecentOidsCache(c context.Context, typ int8, mid int64, fids []int64) (rctFidsMap map[int64][]int64, missFids []int64, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	for _, fid := range fids {
		key := relationKey(mid, fid)
		if err = conn.Send("ZREVRANGE", key, 0, 2); err != nil {
			log.Error("conn.Do(ZREVRANGE, %s,%d) error(%v)", key, fid, err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	rctFidsMap = make(map[int64][]int64, len(fids))
	for _, fid := range fids {
		oids, err := redis.Int64s(conn.Receive())
		if err != nil {
			log.Error("redis.Strings()err(%v)", err)
			return nil, fids, err
		}
		if len(oids) == 0 {
			missFids = append(missFids, fid)
		}
		rctFidsMap[fid] = oids
	}
	return
}

// BatchOidsRedis return the user's 1000 oids from redis.
func (d *Dao) BatchOidsRedis(c context.Context, tp int8, mid int64, limit int) (oids []int64, err error) {
	var (
		key  = relationOidsKey(tp, mid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if oids, err = redis.Int64s(conn.Do("SRANDMEMBER", key, limit)); err != nil {
		log.Error("conn.Do(SRANDMEMBER, %s) error(%v)", key, err)
		return
	}
	return
}

// IsCleaned check if user had do clean action
func (d *Dao) IsCleaned(c context.Context, typ int8, mid, fid int64) (cleanedTime int64, err error) {
	var (
		key  = isCleanedKey(typ, mid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if cleanedTime, err = redis.Int64(conn.Do("HGET", key, fid)); err != nil {
		if err == redis.ErrNil {
			err = nil
			cleanedTime = 0
		} else {
			log.Error("conn.Do(HGET, %v, %v) error(%v)", key, fid, err)
		}
	}
	return
}

// SetCleanedCache set cleand flag.
func (d *Dao) SetCleanedCache(c context.Context, typ int8, mid, fid, ftime, expire int64) (err error) {
	var (
		key  = isCleanedKey(typ, mid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if err = conn.Send("HSET", key, fid, ftime); err != nil {
		log.Error("conn.Send error(%v)", err)
	}
	if err = conn.Send("EXPIRE", key, expire); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive error(%v)", err)
			return
		}
	}
	return
}

// DelRelationOidsCache .
func (d *Dao) DelRelationOidsCache(c context.Context, typ int8, mid int64) (err error) {
	key := relationOidsKey(typ, mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("DEL", key); err != nil {
		log.Error("conn.Do(DEL %s) error(%v)", key, err)
	}
	return
}
