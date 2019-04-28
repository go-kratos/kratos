package fav

import (
	"context"
	"encoding/json"
	"fmt"

	"go-common/app/job/main/favorite/model"
	favmdl "go-common/app/service/main/favorite/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	// _fid list fid => ([]model.Cover)
	_covers          = "fcs_"
	_folderKey       = "fi_%d_%d" // sortedset f_type_mid value:fid,score:ctime
	_oldRelationKey  = "r_%d_%d_%d"
	_allRelationKey  = "ar_%d_%d"
	_relationKey     = "r_%d_%d"  // sortedset r_mid_fid(mtime, oid)
	_relationOidsKey = "ro_%d_%d" // set ro_type_mid value:oids
	_cleanedKey      = "rc_%d_%d" // hash key:rc_type_mid field:fid value:timestamp
	// key fb_mid/100000  offset => mid%100000
	// bit value 1 mean unfaved; bit value 0 mean faved
	_favedBit = "fb_%d_%d"
	_bucket   = 100000
)

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

// allRelationKey return folder relation key.
func allRelationKey(mid, fid int64) string {
	return fmt.Sprintf(_allRelationKey, mid, fid)
}

// oldRelationKey return folder relation key.
func oldRelationKey(typ int8, mid, fid int64) string {
	return fmt.Sprintf(_oldRelationKey, typ, mid, fid)
}

// relationOidsKey return a user oids key.
func relationOidsKey(tp int8, mid int64) string {
	return fmt.Sprintf(_relationOidsKey, tp, mid)
}

// cleanKey return user whether cleaned key.
func cleanedKey(tp int8, mid int64) string {
	return fmt.Sprintf(_cleanedKey, tp, mid)
}

// redisKey make key for redis by prefix and mid
func coversKey(mid, fid int64) string {
	return fmt.Sprintf("%s%d_%d", _covers, mid, fid)
}

// PingRedis ping connection success.
func (d *Dao) pingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}

// DelNewCoverCache delete cover picture cache.
func (d *Dao) DelNewCoverCache(c context.Context, mid, fid int64) (err error) {
	var (
		key  = coversKey(mid, fid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if _, err = conn.Do("DEL", key); err != nil {
		log.Error("DEL %v failed error(%v)", key, err)
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

// FolderCache return a favorite folder from redis.
func (d *Dao) FolderCache(c context.Context, tp int8, mid, fid int64) (folder *favmdl.Folder, err error) {
	var (
		value []byte
		key   = folderKey(tp, mid)
		conn  = d.redis.Get(c)
	)
	defer conn.Close()
	if value, err = redis.Bytes(conn.Do("HGET", key, fid)); err != nil {
		if err == redis.ErrNil {
			err = nil
			folder = nil
		} else {
			log.Error("conn.Do(HGET, %v, %v) error(%v)", key, fid, err)
		}
		return
	}
	folder = &favmdl.Folder{}
	if err = json.Unmarshal(value, folder); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", value, err)
	}
	return
}

// DefaultFolderCache return default favorite folder from redis.
func (d *Dao) DefaultFolderCache(c context.Context, tp int8, mid int64) (folder *favmdl.Folder, err error) {
	var res map[int64]*favmdl.Folder
	if res, err = d.foldersCache(c, tp, mid); err != nil {
		return
	}
	if res == nil {
		return
	}
	for _, folder = range res {
		if folder.IsDefault() {
			return
		}
	}
	folder = nil
	return
}

// foldersCache return the user all folders from redis.
func (d *Dao) foldersCache(c context.Context, tp int8, mid int64) (res map[int64]*favmdl.Folder, err error) {
	var (
		values map[string]string
		key    = folderKey(tp, mid)
		conn   = d.redis.Get(c)
	)
	defer conn.Close()
	if values, err = redis.StringMap(conn.Do("HGETALL", key)); err != nil {
		if err == redis.ErrNil {
			return nil, nil
		}
		log.Error("conn.Do(HGETALL %s) error(%v)", key, err)
		return
	}
	res = make(map[int64]*favmdl.Folder, len(res))
	for _, data := range values {
		folder := &favmdl.Folder{}
		if err = json.Unmarshal([]byte(data), folder); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", data, err)
			return
		}
		res[folder.ID] = folder
	}
	return
}

// RelationCntCache return the folder all relation count from redis.
func (d *Dao) RelationCntCache(c context.Context, mid, fid int64) (cnt int, err error) {
	key := relationKey(mid, fid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if cnt, err = redis.Int(conn.Do("ZCARD", key)); err != nil {
		if err == redis.ErrNil {
			return model.CacheNotFound, nil
		}
		log.Error("conn.Do(ZCARD %s) error(%v)", key, err)
	}
	return
}

// MaxScore get the max score from sorted set
func (d *Dao) MaxScore(c context.Context, m *favmdl.Favorite) (score int64, err error) {
	key := allRelationKey(m.Mid, m.Fid)
	conn := d.redis.Get(c)
	defer conn.Close()
	values, err := redis.Values(conn.Do("ZREVRANGE", key, 0, 0, "WITHSCORES"))
	if err != nil {
		log.Error("conn.Do(ZREVRANGE, %s) error(%v)", key, err)
		return
	}
	if len(values) != 2 {
		err = fmt.Errorf("redis zrange items(%v) length not 2", values)
		return
	}
	var id int64
	redis.Scan(values, &id, &score)
	return
}

// AddAllRelationCache add a relation to redis.
func (d *Dao) AddAllRelationCache(c context.Context, m *favmdl.Favorite) (err error) {
	key := allRelationKey(m.Mid, m.Fid)
	score, err := d.MaxScore(c, m)
	if err != nil {
		return
	}
	if score <= 0 {
		log.Error("dao.AddAllRelationCache invalid score(%d)!%+v", d, *m)
		return
	}
	seq := int64(score/1e10) + 1
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("ZADD", key, seq*1e10+int64(m.MTime), m.Oid*100+int64(m.Type)); err != nil {
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

// AddRelationCache add a relation to redis.
func (d *Dao) AddRelationCache(c context.Context, m *favmdl.Favorite) (err error) {
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

// AddAllRelationsCache add a relation to redis.
func (d *Dao) AddAllRelationsCache(c context.Context, mid, fid int64, fs []*favmdl.Favorite) (err error) {
	key := allRelationKey(mid, fid)
	conn := d.redis.Get(c)
	defer conn.Close()
	for _, fav := range fs {
		if err = conn.Send("ZADD", key, fav.Sequence*1e10+uint64(fav.MTime), fav.Oid*100+int64(fav.Type)); err != nil {
			log.Error("conn.Send(ZADD %s,%d) error(%v)", key, fav.Oid, err)
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

// AddRelationsCache add a relation to redis.
func (d *Dao) AddRelationsCache(c context.Context, tp int8, mid, fid int64, fs []*favmdl.Favorite) (err error) {
	key := relationKey(mid, fid)
	conn := d.redis.Get(c)
	defer conn.Close()
	for _, fav := range fs {
		if err = conn.Send("ZADD", key, fav.MTime, fav.Oid); err != nil {
			log.Error("conn.Send(ZADD %s,%d) error(%v)", key, fav.Oid, err)
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

// DelAllRelationsCache delete the folder relation cache.
func (d *Dao) DelAllRelationsCache(c context.Context, mid, fid int64) (err error) {
	key := allRelationKey(mid, fid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("DEL", key); err != nil {
		log.Error("conn.Do(DEL %s) error(%v)", key, err)
	}
	return
}

// DelRelationCache delete one relation cache.
func (d *Dao) DelRelationCache(c context.Context, mid, fid, oid int64) (err error) {
	key := relationKey(mid, fid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("ZREM", key, oid); err != nil {
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

// DelAllRelationCache delete one relation cache.
func (d *Dao) DelAllRelationCache(c context.Context, mid, fid, oid int64, typ int8) (err error) {
	key := allRelationKey(mid, fid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("ZREM", key, oid*100+int64(typ)); err != nil {
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

// DelOldRelationsCache delete the folder relation cache. TODO:del at 2018.06.08
func (d *Dao) DelOldRelationsCache(c context.Context, typ int8, mid, fid int64) (err error) {
	key := oldRelationKey(typ, mid, fid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("DEL", key); err != nil {
		log.Error("conn.Do(DEL %s) error(%v)", key, err)
	}
	return
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

// RemRelationOidCache del favoured oid.
func (d *Dao) RemRelationOidCache(c context.Context, tp int8, mid, oid int64) (err error) {
	var (
		key  = relationOidsKey(tp, mid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if _, err = conn.Do("SREM", key, oid); err != nil {
		log.Error("conn.Do(%s,%d) error(%v)", key, oid, err)
	}
	return
}

// SetRelationOidsCache set favoured oids .
func (d *Dao) SetRelationOidsCache(c context.Context, tp int8, mid int64, oids []int64) (err error) {
	var (
		key  = relationOidsKey(tp, mid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	for _, oid := range oids {
		if err = conn.Send("SADD", key, oid); err != nil {
			log.Error("conn.Send error(%v)", err)
			return
		}
	}
	if err = conn.Send("EXPIRE", key, d.redisExpire); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < len(oids)+1; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// SetCleanedCache .
func (d *Dao) SetCleanedCache(c context.Context, typ int8, mid, fid, ftime, expire int64) (err error) {
	var (
		key  = cleanedKey(typ, mid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if err = conn.Send("HSET", key, fid, ftime); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
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
			log.Error("conn.Receive() error(%v)", err)
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
