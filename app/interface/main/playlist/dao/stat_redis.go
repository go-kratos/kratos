package dao

import (
	"context"
	"encoding/json"
	"fmt"

	"go-common/app/interface/main/playlist/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_statKey = "st_%d"
	_plKey   = "pl_%d"
)

func keyStat(mid int64) string {
	return fmt.Sprintf(_statKey, mid)
}

func keyPl(pid int64) string {
	return fmt.Sprintf(_plKey, pid)
}

// PlStatCache get stat from cache.
func (d *Dao) PlStatCache(c context.Context, mid, pid int64) (stat *model.PlStat, err error) {
	var (
		bs   []byte
		key  = keyStat(mid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if bs, err = redis.Bytes(conn.Do("HGET", key, pid)); err != nil {
		if err == redis.ErrNil {
			err = nil
			stat = nil
		} else {
			log.Error("conn.Do(HGET,%s,%d) error(%v)", key, pid, err)
		}
		return
	}
	stat = new(model.PlStat)
	if err = json.Unmarshal(bs, stat); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(bs), err)
	}
	return
}

// SetPlStatCache set playlist stat to cache.
func (d *Dao) SetPlStatCache(c context.Context, mid, pid int64, stat *model.PlStat) (err error) {
	var (
		bs     []byte
		ok     bool
		keyMid = keyStat(mid)
		keyPid = keyPl(pid)
		conn   = d.redis.Get(c)
	)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("EXPIRE", keyMid, d.statExpire)); err != nil {
		log.Error("conn.Do(EXPIRE %s) error(%v)", keyMid, err)
		return
	}
	if ok {
		if bs, err = json.Marshal(stat); err != nil {
			log.Error("json.Marshal() error(%v)", err)
			return
		}
		if err = conn.Send("HSET", keyMid, pid, bs); err != nil {
			log.Error("conn.Send(HSET,%s,%d) error(%v)", keyMid, pid, err)
			return
		}
		if err = conn.Send("EXPIRE", keyMid, d.statExpire); err != nil {
			log.Error("conn.Send(EXPIRE,%s) error(%v)", keyMid, err)
			return
		}
		if err = conn.Send("SET", keyPid, bs); err != nil {
			log.Error("conn.Send(SET,%s,%s) error(%v)", keyPid, string(bs), err)
			return
		}
		if err = conn.Send("EXPIRE", keyPid, d.plExpire); err != nil {
			log.Error("conn.Send(EXPIRE,%s) error(%v)", keyPid, err)
			return
		}
		if err = conn.Flush(); err != nil {
			log.Error("add conn.Flush error(%v)", err)
			return
		}
		for i := 0; i < 4; i++ {
			if _, err = conn.Receive(); err != nil {
				log.Error("add conn.Receive()%d error(%v)", i+1, err)
				return
			}
		}
	}
	return

}

// SetStatsCache set playlist stat list  to cache.
func (d *Dao) SetStatsCache(c context.Context, mid int64, plStats []*model.PlStat) (err error) {
	var (
		bs      []byte
		keyPid  string
		keyPids []string
		argsPid = redis.Args{}
	)
	keyMid := keyStat(mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = redis.Bool(conn.Do("EXPIRE", keyMid, d.statExpire)); err != nil {
		log.Error("conn.Do(EXPIRE %s) error(%v)", keyMid, err)
		return
	}
	argsMid := redis.Args{}.Add(keyMid)
	for _, v := range plStats {
		if bs, err = json.Marshal(v); err != nil {
			log.Error("json.Marshal err(%v)", err)
			continue
		}
		argsMid = argsMid.Add(v.ID).Add(string(bs))
		keyPid = keyPl(v.ID)
		keyPids = append(keyPids, keyPid)
		argsPid = argsPid.Add(keyPid).Add(string(bs))
	}
	if err = conn.Send("HMSET", argsMid...); err != nil {
		log.Error("conn.Send(HMSET, %s) error(%v)", keyMid, err)
		return
	}
	if err = conn.Send("EXPIRE", keyMid, d.statExpire); err != nil {
		log.Error("conn.Send(Expire, %s, %d) error(%v)", keyMid, d.statExpire, err)
		return
	}
	if err = conn.Send("MSET", argsPid...); err != nil {
		log.Error("conn.Send(MSET) error(%v)", err)
		return
	}
	count := 3
	for _, v := range keyPids {
		count++
		if err = conn.Send("EXPIRE", v, d.plExpire); err != nil {
			log.Error("conn.Send(Expire, %s, %d) error(%v)", v, d.plExpire, err)
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

// PlsCache get playlist by pids from cache.
func (d *Dao) PlsCache(c context.Context, pids []int64) (res []*model.PlStat, err error) {
	var (
		key  string
		args = redis.Args{}
	)
	for _, pid := range pids {
		key = keyPl(pid)
		args = args.Add(key)
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	var (
		bss [][]byte
	)
	if bss, err = redis.ByteSlices(conn.Do("MGET", args...)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("PlsCache conn.Do(MGET,%s) error(%v)", key, err)
		}
		return
	}
	for _, bs := range bss {
		stat := new(model.PlStat)
		if bs == nil {
			continue
		}
		if err = json.Unmarshal(bs, stat); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", string(bs), err)
			err = nil
			continue
		}
		res = append(res, stat)
	}
	return
}

// SetPlCache set playlist to cache.
func (d *Dao) SetPlCache(c context.Context, plStats []*model.PlStat) (err error) {
	var (
		bs      []byte
		keyPid  string
		keyPids []string
		argsPid = redis.Args{}
	)
	conn := d.redis.Get(c)
	defer conn.Close()
	for _, v := range plStats {
		if bs, err = json.Marshal(v); err != nil {
			log.Error("json.Marshal err(%v)", err)
			continue
		}
		keyPid = keyPl(v.ID)
		keyPids = append(keyPids, keyPid)
		argsPid = argsPid.Add(keyPid).Add(string(bs))
	}
	if err = conn.Send("MSET", argsPid...); err != nil {
		log.Error("conn.Send(MSET) error(%v)", err)
		return
	}
	count := 1
	for _, v := range keyPids {
		count++
		if err = conn.Send("EXPIRE", v, d.plExpire); err != nil {
			log.Error("conn.Send(Expire, %s, %d) error(%v)", v, d.plExpire, err)
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

// DelPlCache delete playlist from  redis.
func (d *Dao) DelPlCache(c context.Context, mid, pid int64) (err error) {
	var (
		key     = keyPl(pid)
		plaKey  = keyPlArc(pid)
		pladKey = keyPlArcDesc(pid)
		keyStat = keyStat(mid)
		conn    = d.redis.Get(c)
	)
	defer conn.Close()
	if err = conn.Send("DEL", key); err != nil {
		log.Error("conn.Send(DEL %s) error(%v)", key, err)
		return
	}
	if err = conn.Send("DEL", plaKey); err != nil {
		log.Error("conn.Send(DEL %s) error(%v)", plaKey, err)
		return
	}
	if err = conn.Send("DEL", pladKey); err != nil {
		log.Error("conn.Send(DEL %s) error(%v)", pladKey, err)
		return
	}
	if err = conn.Send("HDEL", keyStat, pid); err != nil {
		log.Error("conn.Send(HDEL,%s,%d) error(%v)", keyStat, pid, err)
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	for i := 0; i < 4; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// StatsCache get playlist  stats from cache.
func (d *Dao) StatsCache(c context.Context, mid int64) (res []*model.PlStat, err error) {
	key := keyStat(mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	var (
		bss [][]byte
	)
	if bss, err = redis.ByteSlices(conn.Do("HGETALL", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("StatCache conn.Do(HGETALL,%s) error(%v)", key, err)
		}
		return
	}
	for i := 1; i <= len(bss); i += 2 {
		stat := new(model.PlStat)
		if err = json.Unmarshal(bss[i], stat); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", string(bss[i]), err)
			continue
		}
		res = append(res, stat)
	}
	return
}
