package dao

import (
	"context"
	"fmt"

	"go-common/app/interface/main/playlist/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_plArcKey     = "pla_%d"
	_plArcDescKey = "plad_%d"
)

func keyPlArc(pid int64) string {
	return fmt.Sprintf(_plArcKey, pid)
}

func keyPlArcDesc(pid int64) string {
	return fmt.Sprintf(_plArcDescKey, pid)
}

// ArcsCache get playlist archives cache.
func (d *Dao) ArcsCache(c context.Context, pid int64, start, end int) (arcs []*model.ArcSort, err error) {
	var (
		plakey  = keyPlArc(pid)
		pladKey = keyPlArcDesc(pid)
		conn    = d.redis.Get(c)
		aids    []int64
		descs   []string
	)
	defer conn.Close()
	values, err := redis.Values(conn.Do("ZRANGE", plakey, start, end, "WITHSCORES"))
	if err != nil {
		log.Error("conn.Do(ZREVRANGE, %s) error(%v)", plakey, err)
		return
	}
	if len(values) == 0 {
		return
	}
	arcMap := make(map[int64]*model.ArcSort)
	args := redis.Args{}.Add(pladKey)
	for len(values) > 0 {
		arc := &model.ArcSort{}
		if values, err = redis.Scan(values, &arc.Aid, &arc.Sort); err != nil {
			log.Error("redis.Scan(%v) error(%v)", values, err)
			return
		}
		arcMap[arc.Aid] = arc
		aids = append(aids, arc.Aid)
		args = args.Add(arc.Aid)
	}
	if len(aids) > 0 {
		descs, err = redis.Strings(conn.Do("HMGET", args...))
		if err != nil {
			log.Error("conn.Do(HMGET %v) error(%v)", args, err)
			err = nil
		}
		descLen := len(descs)
		for k, aid := range aids {
			if arc, ok := arcMap[aid]; ok {
				if descLen >= k+1 {
					if desc := descs[k]; desc != "" {
						arc.Desc = desc
					}
				}
				arcs = append(arcs, arc)
			}
		}
	}
	return
}

// AddArcCache add playlist archive cache.
func (d *Dao) AddArcCache(c context.Context, pid int64, arc *model.ArcSort) (err error) {
	var (
		plakey  = keyPlArc(pid)
		pladKey = keyPlArcDesc(pid)
		conn    = d.redis.Get(c)
		count   int
	)
	defer conn.Close()
	if _, err = redis.Bool(conn.Do("EXPIRE", plakey, d.plExpire)); err != nil {
		log.Error("conn.Do(EXPIRE %s) error(%v)", plakey, err)
		return
	}
	if _, err = redis.Bool(conn.Do("EXPIRE", pladKey, d.plExpire)); err != nil {
		log.Error("conn.Do(EXPIRE %s) error(%v)", pladKey, err)
		return
	}
	args1 := redis.Args{}.Add(plakey)
	args1 = args1.Add(arc.Sort).Add(arc.Aid)
	if err = conn.Send("ZADD", args1...); err != nil {
		log.Error("conn.Send(ZADD, %s, %v) error(%v)", plakey, args1, err)
		return
	}
	count++
	if arc.Desc != "" {
		args2 := redis.Args{}.Add(pladKey).Add(arc.Aid).Add(arc.Desc)
		if err = conn.Send("HSET", args2...); err != nil {
			log.Error("conn.Send(ZADD, %s, %v) error(%v)", plakey, args2, err)
			return
		}
		count++
		if err = conn.Send("EXPIRE", pladKey, d.plExpire); err != nil {
			log.Error("conn.Send(Expire, %s) error(%v)", pladKey, err)
			return
		}
		count++
	}
	if err = conn.Send("EXPIRE", plakey, d.plExpire); err != nil {
		log.Error("conn.Send(Expire, %s) error(%v)", pladKey, err)
		return
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

// SetArcsCache set playlist archives cache.
func (d *Dao) SetArcsCache(c context.Context, pid int64, arcs []*model.ArcSort) (err error) {
	var (
		plaKey  = keyPlArc(pid)
		pladKey = keyPlArcDesc(pid)
		conn    = d.redis.Get(c)
		addDesc bool
		count   int
	)
	defer conn.Close()
	if _, err = redis.Bool(conn.Do("EXPIRE", plaKey, d.plExpire)); err != nil {
		log.Error("conn.Do(EXPIRE %s) error(%v)", plaKey, err)
		return
	}
	if _, err = redis.Bool(conn.Do("EXPIRE", pladKey, d.plExpire)); err != nil {
		log.Error("conn.Do(EXPIRE %s) error(%v)", plaKey, err)
		return
	}
	args1 := redis.Args{}.Add(plaKey)
	args2 := redis.Args{}.Add(pladKey)
	for _, arc := range arcs {
		args1 = args1.Add(arc.Sort).Add(arc.Aid)
		if arc.Desc != "" {
			addDesc = true
			args2 = args2.Add(arc.Aid).Add(arc.Desc)
		}
	}
	if err = conn.Send("ZADD", args1...); err != nil {
		log.Error("conn.Send(ZADD, %s, %v) error(%v)", plaKey, args1, err)
		return
	}
	count++
	if addDesc {
		if err = conn.Send("HMSET", args2...); err != nil {
			log.Error("conn.Send(ZADD, %s, %v) error(%v)", pladKey, args2, err)
			return
		}
		count++
	}
	if err = conn.Send("EXPIRE", pladKey, d.plExpire); err != nil {
		log.Error("conn.Send(Expire, %s) error(%v)", pladKey, err)
		return
	}
	count++
	if err = conn.Send("EXPIRE", plaKey, d.plExpire); err != nil {
		log.Error("conn.Send(Expire, %s) error(%v)", plaKey, err)
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

// SetArcDescCache set playlist archive desc cache.
func (d *Dao) SetArcDescCache(c context.Context, pid, aid int64, desc string) (err error) {
	var (
		key  = keyPlArcDesc(pid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if _, err = redis.Bool(conn.Do("EXPIRE", key, d.plExpire)); err != nil {
		log.Error("conn.Do(EXPIRE %s) error(%v)", key, err)
		return
	}
	if err = conn.Send("HSET", key, aid, desc); err != nil {
		log.Error("conn.Send(HSET, %s, %d, %s) error(%v)", key, aid, desc, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.plExpire); err != nil {
		log.Error("conn.Send(Expire, %s) error(%v)", key, err)
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

// DelArcsCache delete  playlist archives cache.
func (d *Dao) DelArcsCache(c context.Context, pid int64, aids []int64) (err error) {
	var (
		plaKey  = keyPlArc(pid)
		pladKey = keyPlArcDesc(pid)
		conn    = d.redis.Get(c)
	)
	defer conn.Close()
	arg1 := redis.Args{}.Add(plaKey)
	arg2 := redis.Args{}.Add(pladKey)
	for _, aid := range aids {
		arg1 = arg1.Add(aid)
		arg2 = arg2.Add(aid)
	}
	if err = conn.Send("ZREM", arg1...); err != nil {
		log.Error("conn.Send(ZREM %s) error(%v)", plaKey, err)
		return
	}
	if err = conn.Send("HDEL", arg2...); err != nil {
		log.Error("conn.Send(HDEL %s) error(%v)", pladKey, err)
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

// DelCache del all cache .
func (d *Dao) DelCache(c context.Context, pid int64) (err error) {
	var (
		plaKey  = keyPlArc(pid)
		pladKey = keyPlArcDesc(pid)
		conn    = d.redis.Get(c)
	)
	defer conn.Close()
	if err = conn.Send("DEL", plaKey); err != nil {
		log.Error("conn.Send(DEL plaKey(%s) error(%v))", plaKey, err)
		return
	}
	if err = conn.Send("DEL", pladKey); err != nil {
		log.Error("conn.Send(DEL pladKey(%s) error(%v))", pladKey, err)
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
