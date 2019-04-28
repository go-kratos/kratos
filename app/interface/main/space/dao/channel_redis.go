package dao

import (
	"context"
	"encoding/json"
	"fmt"

	"go-common/app/interface/main/space/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_keyCl    = "cl_%d"
	_keyClArc = "cla_%d_%d"
)

func keyCl(mid int64) string {
	return fmt.Sprintf(_keyCl, mid)
}

func keyClArc(mid, cid int64) string {
	return fmt.Sprintf(_keyClArc, mid, cid)
}

func keyClArcSort(mid, cid int64) string {
	return keyClArc(mid, cid) + "_s"
}

// ChannelCache get channel cache.
func (d *Dao) ChannelCache(c context.Context, mid, cid int64) (channel *model.Channel, err error) {
	var (
		bs   []byte
		key  = keyCl(mid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if bs, err = redis.Bytes(conn.Do("HGET", key, cid)); err != nil {
		if err == redis.ErrNil {
			err = nil
			channel = nil
		} else {
			log.Error("conn.Do(HGET,%s,%d) error(%v)", key, cid, err)
		}
		return
	}
	channel = new(model.Channel)
	if err = json.Unmarshal(bs, channel); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(bs), err)
	}
	return
}

// SetChannelCache add channel data cache.
func (d *Dao) SetChannelCache(c context.Context, mid, cid int64, channel *model.Channel) (err error) {
	var (
		bs   []byte
		ok   bool
		key  = keyCl(mid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if bs, err = json.Marshal(channel); err != nil {
		log.Error("json.Marshal() error(%v)", err)
		return
	}
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.clExpire)); err != nil || !ok {
		log.Error("conn.Do(EXPIRE %s) error(%v)", key, err)
		return
	}
	if err = conn.Send("HSET", key, cid, bs); err != nil {
		log.Error("conn.Send(HSET,%s,%d) error(%v)", key, cid, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.clExpire); err != nil {
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

// DelChannelCache delete channel cache from list.
func (d *Dao) DelChannelCache(c context.Context, mid, cid int64) (err error) {
	var (
		key     = keyCl(mid)
		arcsKey = keyClArc(mid, cid)
		sortKey = keyClArcSort(mid, cid)
		conn    = d.redis.Get(c)
	)
	defer conn.Close()
	if err = conn.Send("HDEL", key, cid); err != nil {
		log.Error("conn.Send(HDEL,%s,%d) error(%v)", key, cid, err)
		return
	}
	if err = conn.Send("DEL", arcsKey); err != nil {
		log.Error("conn.Send(DEL,%s) error(%v)", arcsKey, err)
	}
	if err = conn.Send("DEL", sortKey); err != nil {
		log.Error("conn.Send(DEL,%s) error(%v)", sortKey, err)
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	for i := 0; i < 3; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("add conn.Receive()%d error(%v)", i+1, err)
			return
		}
	}
	return
}

// SetChannelListCache add channel data cache.
func (d *Dao) SetChannelListCache(c context.Context, mid int64, channelList []*model.Channel) (err error) {
	var (
		bs   []byte
		key  = keyCl(mid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if err = conn.Send("DEL", key); err != nil {
		log.Error("conn.Send(DEL, %s) error(%v)", key, err)
		return
	}
	args := redis.Args{}.Add(key)
	for _, channel := range channelList {
		if bs, err = json.Marshal(channel); err != nil {
			log.Error("json.Marshal() error(%v)", err)
			continue
		} else {
			args = args.Add(channel.Cid).Add(string(bs))
		}
	}
	if err = conn.Send("HMSET", args...); err != nil {
		log.Error("conn.Send(HMSET, %s) error(%v)", key, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.clExpire); err != nil {
		log.Error("conn.Send(Expire, %s, %d) error(%v)", key, d.clExpire, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 3; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// ChannelListCache get channel list cache.
func (d *Dao) ChannelListCache(c context.Context, mid int64) (channels []*model.Channel, err error) {
	var (
		bss  [][]byte
		key  = keyCl(mid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if bss, err = redis.ByteSlices(conn.Do("HGETALL", key)); err != nil {
		log.Error("conn.Do(HGETALL,%s) error(%v)", key, err)
		return
	}
	for i := 1; i <= len(bss); i += 2 {
		channel := new(model.Channel)
		if err = json.Unmarshal(bss[i], channel); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", string(bss[i]), err)
			continue
		}
		channels = append(channels, channel)
	}
	return
}

// ChannelArcsCache get channel archives cache.
func (d *Dao) ChannelArcsCache(c context.Context, mid, cid int64, start, end int, order bool) (arcs []*model.ChannelArc, err error) {
	var (
		bss     [][]byte
		values  []interface{}
		key     = keyClArc(mid, cid)
		sortKey = keyClArcSort(mid, cid)
		conn    = d.redis.Get(c)
		cmd     = "ZREVRANGE"
	)
	defer conn.Close()
	if order {
		cmd = "ZRANGE"
	}
	if values, err = redis.Values(conn.Do(cmd, sortKey, start, end, "WITHSCORES")); err != nil {
		log.Error("conn.Do(ZREVRANGE, %s) error(%v)", sortKey, err)
		return
	} else if len(values) == 0 {
		return
	}
	arg := redis.Args{}.Add(key)
	for len(values) > 0 {
		arcSort := new(model.ChannelArcSort)
		if values, err = redis.Scan(values, &arcSort.Aid, &arcSort.OrderNum); err != nil {
			log.Error("redis.Scan(%v) error(%v)", values, err)
			return
		}
		if arcSort.Aid > 0 {
			arg = arg.Add(arcSort.Aid)
		}
	}
	if bss, err = redis.ByteSlices(conn.Do("HMGET", arg...)); err != nil {
		log.Error("conn.Do(HMGET,%s) error(%v)", key, err)
		return
	}
	for _, bs := range bss {
		if len(bs) == 0 {
			continue
		}
		if len(bs) > 0 {
			arc := new(model.ChannelArc)
			if err = json.Unmarshal(bs, arc); err != nil {
				log.Error("json.Unmarshal(%s) mid(%d) cid(%d) error(%v)", string(bs), mid, cid, err)
				err = nil
				continue
			}
			arcs = append(arcs, arc)
		}
	}
	return
}

// AddChannelArcCache add channel archives cache.
func (d *Dao) AddChannelArcCache(c context.Context, mid, cid int64, arcs []*model.ChannelArc) (err error) {
	var (
		bs      []byte
		ok      bool
		key     = keyClArc(mid, cid)
		sortKey = keyClArcSort(mid, cid)
		conn    = d.redis.Get(c)
	)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.clExpire)); err != nil && ok {
		log.Error("conn.Do(EXPIRE %s) error(%v)", key, err)
		return
	}
	if ok, err = redis.Bool(conn.Do("EXPIRE", sortKey, d.clExpire)); err != nil && ok {
		log.Error("conn.Do(EXPIRE %s) error(%v)", key, err)
		return
	}
	args1 := redis.Args{}.Add(key)
	args2 := redis.Args{}.Add(sortKey)
	for _, arc := range arcs {
		if bs, err = json.Marshal(arc); err != nil {
			log.Error("json.Marshal() error(%v)", err)
			return
		}
		args1 = args1.Add(arc.Aid).Add(string(bs))
		args2 = args2.Add(arc.OrderNum).Add(arc.Aid)
	}
	if err = conn.Send("HMSET", args1...); err != nil {
		log.Error("conn.Send(HMSET, %s, %v) error(%v)", key, args1, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.clExpire); err != nil {
		log.Error("conn.Send(Expire, %s) error(%v)", key, err)
		return
	}
	if err = conn.Send("ZADD", args2...); err != nil {
		log.Error("conn.Send(ZADD, %s, %v) error(%v)", sortKey, args2, err)
		return
	}
	if err = conn.Send("EXPIRE", sortKey, d.clExpire); err != nil {
		log.Error("conn.Send(Expire, %s) error(%v)", sortKey, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
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

// SetChannelArcSortCache set channel archives sort cache
func (d *Dao) SetChannelArcSortCache(c context.Context, mid, cid int64, sort []*model.ChannelArcSort) (err error) {
	var (
		key     = keyClArc(mid, cid)
		sortKey = keyClArcSort(mid, cid)
		conn    = d.redis.Get(c)
	)
	defer conn.Close()
	if err = conn.Send("DEL", sortKey); err != nil {
		log.Error("conn.Send(DEL, %s) error(%v)", sortKey, err)
		return
	}
	args := redis.Args{}.Add(sortKey)
	for _, v := range sort {
		args = args.Add(v.OrderNum).Add(v.Aid)
	}
	if err = conn.Send("ZADD", args...); err != nil {
		log.Error("conn.Send(ZADD, %s) error(%v)", sortKey, err)
		return
	}
	if err = conn.Send("EXPIRE", sortKey, d.clExpire); err != nil {
		log.Error("conn.Send(Expire, %s, %d) error(%v)", sortKey, d.clExpire, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.clExpire); err != nil {
		log.Error("conn.Send(Expire, %s, %d) error(%v)", key, d.clExpire, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
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

// DelChannelArcCache delete channel archive cache from cache list.
func (d *Dao) DelChannelArcCache(c context.Context, mid, cid, aid int64) (err error) {
	key := keyClArc(mid, cid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("HDEL", key, aid); err != nil {
		log.Error("conn.Send(ZREM,%s,%d) error(%v)", key, aid, err)
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		log.Error("conn.Receive() error(%v)", err)
	}
	return
}

// DelChannelArcsCache delete all channel arcs cache when delete channel
func (d *Dao) DelChannelArcsCache(c context.Context, mid, cid int64) (err error) {
	key := keyClArc(mid, cid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("DEL", key); err != nil {
		log.Error("conn.Send(DEL,%s,%d,%d) error(%v)", key, mid, cid, err)
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		log.Error("conn.Receive() error(%v)", err)
	}
	return
}

// SetChannelArcsCache add channel archive cache.
func (d *Dao) SetChannelArcsCache(c context.Context, mid, cid int64, arcs []*model.ChannelArc) (err error) {
	var (
		bs   []byte
		key1 = keyClArc(mid, cid)
		key2 = keyClArcSort(mid, cid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if err = conn.Send("DEL", key1); err != nil {
		log.Error("conn.Send(DEL, %s) error(%v)", key1, err)
		return
	}
	if err = conn.Send("DEL", key2); err != nil {
		log.Error("conn.Send(DEL, %s) error(%v)", key2, err)
		return
	}
	args1 := redis.Args{}.Add(key1)
	args2 := redis.Args{}.Add(key2)
	for _, arc := range arcs {
		if bs, err = json.Marshal(arc); err != nil {
			log.Error("json.Marshal() error(%v)", err)
			continue
		} else {
			args1 = args1.Add(arc.Aid).Add(string(bs))
		}
		args2 = args2.Add(arc.OrderNum).Add(arc.Aid)
	}
	if err = conn.Send("HMSET", args1...); err != nil {
		log.Error("conn.Send(HMSET, %s) error(%v)", key1, err)
		return
	}
	if err = conn.Send("ZADD", args2...); err != nil {
		log.Error("conn.Send(ZADD, %s) error(%v)", key2, err)
		return
	}
	if err = conn.Send("EXPIRE", key1, d.clExpire); err != nil {
		log.Error("conn.Send(Expire, %s, %d) error(%v)", key1, d.clExpire, err)
		return
	}
	if err = conn.Send("EXPIRE", key2, d.clExpire); err != nil {
		log.Error("conn.Send(Expire, %s, %d) error(%v)", key2, d.clExpire, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 6; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}
