package like

import (
	"context"
	"encoding/json"
	"fmt"

	match "go-common/app/interface/main/activity/model/like"
	"go-common/library/cache/redis"
	"go-common/library/log"
	"go-common/library/time"
)

const (
	_keyMatch       = "mat_%d"
	_keyActMatch    = "am_%d"
	_keyObject      = "ob_%d"
	_keyObjects     = "os_%d"
	_keyUserLog     = "ugl_%d_%d"
	_keyMatchFollow = "mf_%d"
)

func keyMatch(id int64) string {
	return fmt.Sprintf(_keyMatch, id)
}

func keyActMatch(sid int64) string {
	return fmt.Sprintf(_keyActMatch, sid)
}

func keyObject(id int64) string {
	return fmt.Sprintf(_keyObject, id)
}

func keyObjects(sid int64) string {
	return fmt.Sprintf(_keyObjects, sid)
}

func keyUserLog(sid, mid int64) string {
	return fmt.Sprintf(_keyUserLog, sid, mid)
}

func keyMatchFollow(mid int64) string {
	return fmt.Sprintf(_keyMatchFollow, mid)
}

// MatchCache get match from cache.
func (d *Dao) MatchCache(c context.Context, id int64) (mat *match.Match, err error) {
	var (
		bs   []byte
		key  = keyMatch(id)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if bs, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
			mat = nil
		} else {
			log.Error("conn.Do(GET,%s) error(%v)", key, err)
		}
		return
	}
	mat = new(match.Match)
	if err = json.Unmarshal(bs, mat); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(bs), err)
	}
	return
}

// SetMatchCache set match to cache.
func (d *Dao) SetMatchCache(c context.Context, id int64, mat *match.Match) (err error) {
	var (
		bs   []byte
		key  = keyMatch(id)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if bs, err = json.Marshal(mat); err != nil {
		log.Error("json.Marshal() error(%v)", err)
		return
	}
	if err = conn.Send("SET", key, bs); err != nil {
		log.Error("conn.Send(SET,%s,%d) error(%v)", key, id, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.matchExpire); err != nil {
		log.Error("conn.Send(EXPIRE,%s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("add conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("add conn.Receive()%d error(%v)", i+1, err)
			return
		}
	}
	return

}

// ActMatchCache get match list from cache.
func (d *Dao) ActMatchCache(c context.Context, sid int64) (res []*match.Match, err error) {
	key := keyActMatch(sid)
	conn := d.redis.Get(c)
	defer conn.Close()
	values, err := redis.Values(conn.Do("ZRANGE", key, 0, -1, "WITHSCORES"))
	if err != nil {
		log.Error("conn.Do(ZRANGE, %s) error(%v)", key, err)
		return
	}
	if len(values) == 0 {
		return
	}
	var num int64
	for len(values) > 0 {
		bs := []byte{}
		if values, err = redis.Scan(values, &bs, &num); err != nil {
			log.Error("redis.Scan(%v) error(%v)", values, err)
			return
		}
		match := &match.Match{}
		if err = json.Unmarshal(bs, match); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", bs, err)
			return
		}
		res = append(res, match)
	}
	return
}

// SetActMatchCache set match list cache.
func (d *Dao) SetActMatchCache(c context.Context, sid int64, matchs []*match.Match) (err error) {
	key := keyActMatch(sid)
	conn := d.redis.Get(c)
	defer conn.Close()
	count := 0
	if err = conn.Send("DEL", key); err != nil {
		log.Error("conn.Send(DEL, %s) error(%v)", key, err)
		return
	}
	count++
	for _, match := range matchs {
		bs, _ := json.Marshal(match)
		if err = conn.Send("ZADD", key, match.Ctime, bs); err != nil {
			log.Error("conn.Send(ZADD, %s, %s) error(%v)", key, string(bs), err)
			return
		}
		count++
	}
	if err = conn.Send("EXPIRE", key, d.matchExpire); err != nil {
		log.Error("conn.Send(Expire, %s, %d) error(%v)", key, d.matchExpire, err)
		return
	}
	count++
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// ObjectCache get object from cache.
func (d *Dao) ObjectCache(c context.Context, id int64) (mat *match.Object, err error) {
	var (
		bs   []byte
		key  = keyObject(id)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if bs, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
			mat = nil
		} else {
			log.Error("conn.Do(GET,%s) error(%v)", key, err)
		}
		return
	}
	mat = new(match.Object)
	if err = json.Unmarshal(bs, mat); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(bs), err)
	}
	return
}

// CacheMatchSubjects .
func (d *Dao) CacheMatchSubjects(c context.Context, ids []int64) (res map[int64]*match.Object, err error) {
	var (
		key  string
		args = redis.Args{}
		bss  [][]byte
	)
	for _, pid := range ids {
		key = keyObject(pid)
		args = args.Add(key)
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	if bss, err = redis.ByteSlices(conn.Do("MGET", args...)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("CacheMatchSubjects conn.Do(MGET,%s) error(%v)", key, err)
		}
		return
	}
	res = make(map[int64]*match.Object, len(ids))
	for _, bs := range bss {
		obj := new(match.Object)
		if bs == nil {
			continue
		}
		if err = json.Unmarshal(bs, obj); err != nil {
			log.Error("CacheMatchSubjects json.Unmarshal(%s) error(%v)", string(bs), err)
			err = nil
			continue
		}
		res[obj.ID] = obj
	}
	return
}

// SetObjectCache set object to cache.
func (d *Dao) SetObjectCache(c context.Context, id int64, object *match.Object) (err error) {
	var (
		bs   []byte
		key  = keyObject(id)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if bs, err = json.Marshal(object); err != nil {
		log.Error("json.Marshal() error(%v)", err)
		return
	}
	if err = conn.Send("SET", key, bs); err != nil {
		log.Error("conn.Send(HSET,%s,%d) error(%v)", key, id, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.matchExpire); err != nil {
		log.Error("conn.Send(EXPIRE,%s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("add conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("add conn.Receive()%d error(%v)", i+1, err)
			return
		}
	}
	return

}

// AddCacheMatchSubjects .
func (d *Dao) AddCacheMatchSubjects(c context.Context, data map[int64]*match.Object) (err error) {
	if len(data) == 0 {
		return
	}
	var (
		bs      []byte
		keyID   string
		keyIDs  []string
		argsPid = redis.Args{}
	)
	conn := d.redis.Get(c)
	defer conn.Close()
	for _, v := range data {
		if bs, err = json.Marshal(v); err != nil {
			log.Error("json.Marshal err(%v)", err)
			continue
		}
		keyID = keyObject(v.ID)
		keyIDs = append(keyIDs, keyID)
		argsPid = argsPid.Add(keyID).Add(string(bs))
	}
	if err = conn.Send("MSET", argsPid...); err != nil {
		log.Error("AddCacheMatchSubjects conn.Send(MSET) error(%v)", err)
		return
	}
	count := 1
	for _, v := range keyIDs {
		count++
		if err = conn.Send("EXPIRE", v, d.matchExpire); err != nil {
			log.Error("AddCacheMatchSubjects conn.Send(Expire, %s, %d) error(%v)", v, d.matchExpire, err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// ObjectsCache get object list from cache.
func (d *Dao) ObjectsCache(c context.Context, sid int64, start, end int) (res []*match.Object, total int, err error) {
	key := keyObjects(sid)
	conn := d.redis.Get(c)
	defer conn.Close()
	values, err := redis.Values(conn.Do("ZRANGE", key, start, end, "WITHSCORES"))
	if err != nil {
		log.Error("conn.Do(ZRANGE, %s) error(%v)", key, err)
		return
	}
	if len(values) == 0 {
		return
	}
	var num int64
	for len(values) > 0 {
		bs := []byte{}
		if values, err = redis.Scan(values, &bs, &num); err != nil {
			log.Error("redis.Scan(%v) error(%v)", values, err)
			return
		}
		object := &match.Object{}
		if err = json.Unmarshal(bs, object); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", bs, err)
			return
		}
		res = append(res, object)
	}
	total = from(num)
	return
}

// SetObjectsCache set object list cache.
func (d *Dao) SetObjectsCache(c context.Context, sid int64, objects []*match.Object, total int) (err error) {
	key := keyObjects(sid)
	conn := d.redis.Get(c)
	defer conn.Close()
	count := 0
	if err = conn.Send("DEL", key); err != nil {
		log.Error("conn.Send(DEL, %s) error(%v)", key, err)
		return
	}
	count++
	for _, object := range objects {
		bs, _ := json.Marshal(object)
		if err = conn.Send("ZADD", key, combine(object.GameStime, total), bs); err != nil {
			log.Error("conn.Send(ZADD, %s, %s) error(%v)", key, string(bs), err)
			return
		}
		count++
	}
	if err = conn.Send("EXPIRE", key, d.matchExpire); err != nil {
		log.Error("conn.Send(Expire, %s, %d) error(%v)", key, d.matchExpire, err)
		return
	}
	count++
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// UserLogCache get user log list from cache.
func (d *Dao) UserLogCache(c context.Context, sid, mid int64) (res []*match.UserLog, err error) {
	key := keyUserLog(sid, mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	values, err := redis.Values(conn.Do("ZREVRANGE", key, 0, -1, "WITHSCORES"))
	if err != nil {
		log.Error("conn.Do(ZREVRANGE, %s) error(%v)", key, err)
		return
	}
	if len(values) == 0 {
		return
	}
	var num int64
	for len(values) > 0 {
		bs := []byte{}
		if values, err = redis.Scan(values, &bs, &num); err != nil {
			log.Error("redis.Scan(%v) error(%v)", values, err)
			return
		}
		userLog := &match.UserLog{}
		if err = json.Unmarshal(bs, userLog); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", bs, err)
			return
		}
		res = append(res, userLog)
	}
	return
}

// SetUserLogCache set user log list cache.
func (d *Dao) SetUserLogCache(c context.Context, sid, mid int64, userLogs []*match.UserLog) (err error) {
	key := keyUserLog(sid, mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	count := 0
	if err = conn.Send("DEL", key); err != nil {
		log.Error("conn.Send(DEL, %s) error(%v)", key, err)
		return
	}
	count++
	for _, userLog := range userLogs {
		bs, _ := json.Marshal(userLog)
		if err = conn.Send("ZADD", key, userLog.Ctime, bs); err != nil {
			log.Error("conn.Send(ZADD, %s, %s) error(%v)", key, string(bs), err)
			return
		}
		count++
	}
	if err = conn.Send("EXPIRE", key, d.matchExpire); err != nil {
		log.Error("conn.Send(Expire, %s, %d) error(%v)", key, d.matchExpire, err)
		return
	}
	count++
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// DelUserLogCache delete user log cache.
func (d *Dao) DelUserLogCache(c context.Context, sid, mid int64) (err error) {
	key := keyUserLog(sid, mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	_, err = conn.Do("DEL", key)
	return
}

// DelActMatchCache del match   cache
func (d *Dao) DelActMatchCache(c context.Context, sid, matID int64) (err error) {
	matKey := keyMatch(matID)
	actMKey := keyActMatch(sid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("DEL", matKey); err != nil {
		log.Error("conn.Send(DEL, %s) error(%v)", actMKey, err)
		return
	}
	if err = conn.Send("DEL", actMKey); err != nil {
		log.Error("conn.Send(DEL, %s) error(%v)", actMKey, err)
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

// DelObjectCache del  object cache
func (d *Dao) DelObjectCache(c context.Context, objID, sid int64) (err error) {
	objKey := keyObject(objID)
	objsKey := keyObjects(sid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("DEL", objKey); err != nil {
		log.Error("conn.Send(DEL, %s) error(%v)", objKey, err)
		return
	}
	if err = conn.Send("DEL", objsKey); err != nil {
		log.Error("conn.Send(DEL, %s) error(%v)", objKey, err)
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

// AddFollow add follow teams
func (d *Dao) AddFollow(c context.Context, mid int64, teams []string) (err error) {
	var (
		bs   []byte
		key  = keyMatchFollow(mid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if bs, err = json.Marshal(teams); err != nil {
		log.Error("json.Marshal() error(%v)", err)
		return
	}
	if err = conn.Send("SET", key, bs); err != nil {
		log.Error("conn.Send(SET,%s,%d) error(%v)", key, mid, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.followExpire); err != nil {
		log.Error("conn.Send(EXPIRE,%s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("add conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("add conn.Receive()%d error(%v)", i+1, err)
			return
		}
	}
	return
}

// Follow get follow teams
func (d *Dao) Follow(c context.Context, mid int64) (res []string, err error) {
	var (
		bs   []byte
		key  = keyMatchFollow(mid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if bs, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(GET,%s) error(%v)", key, err)
		}
		return
	}
	if err = json.Unmarshal(bs, &res); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(bs), err)
	}
	return
}

func from(i int64) int {
	return int(i & 0xffff)
}

func combine(ctime time.Time, count int) int64 {
	return ctime.Time().Unix()<<16 | int64(count)
}
