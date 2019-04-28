package share

import (
	"context"
	"fmt"
	"time"

	"go-common/library/cache/redis"
	"go-common/library/log"
	xip "go-common/library/net/ip"
)

const (
	_prefixShare = "as"
	// the first shared every day
	_prefixFirst = "afs"
)

func redisKey(aid int64) string {
	return fmt.Sprintf("%s_%d", _prefixShare, aid)
}

func redisFKey(mid int64, date string) string {
	return fmt.Sprintf("%s_%s_%d", _prefixFirst, date, mid%10000)
}

// AddShare add a share with mid and ip to redis.
func (d *Dao) AddShare(c context.Context, mid, aid int64, ip string) (ok bool, err error) {
	var (
		key   = redisKey(aid)
		value = (mid << 32) | int64(xip.InetAtoN(ip))
		conn  = d.rds.Get(c)
	)
	defer conn.Close()
	if err = conn.Send("SADD", key, value); err != nil {
		log.Error("conn.Send(SADD, %s, %d) error(%v)", key, value, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.expire); err != nil {
		log.Error("conn.Send(EXPIRE, %s, %d) error(%v)", key, value, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
	}
	if ok, err = redis.Bool(conn.Receive()); err != nil {
		log.Error("conn.Receive error(%v)", err)
	}
	if _, err = conn.Receive(); err != nil {
		log.Error("conn.Receive error(%v)", err)
	}
	return
}

// HadFirstShare if key not exist need set expire to one day.
func (d *Dao) HadFirstShare(c context.Context, mid, aid int64, ip string) (had bool, err error) {
	var (
		date  = time.Now().Format("0102")
		key   = redisFKey(mid, date)
		addOk bool
		conn  = d.rds.Get(c)
	)
	defer conn.Close()
	if addOk, err = redis.Bool(conn.Do("SADD", key, mid)); err != nil {
		log.Error("conn.Send(SADD, %s, %d) error(%v)", key, mid, err)
		return
	}
	if addOk {
		if _, err = conn.Do("EXPIRE", key, 24*60*60); err != nil {
			log.Error("conn.Send(EXPIRE, %s, 24*60*60) error(%v)", key, err)
		}
	}
	had = !addOk
	return
}
