package stat

import (
	"context"
	"strconv"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_ipBanKey    = "ipb_"
	_buvidBanKey = "bvb_"
)

func ipBanKey(id int64, ip string) (key string) {
	key = _ipBanKey + strconv.FormatInt(id, 10) + ip
	return
}

func buvidBanKey(id, mid int64, ip, buvid string) (key string) {
	key = _buvidBanKey + strconv.FormatInt(id, 10) + ip + buvid
	if mid != 0 {
		key += strconv.FormatInt(mid, 10)
	}
	return
}

// IPBan intercepts illegal views.
func (d *Dao) IPBan(c context.Context, id int64, ip string) (ban bool) {
	var (
		err   error
		exist bool
		key   = ipBanKey(id, ip)
		conn  = d.redis.Get(c)
	)
	defer conn.Close()
	if exist, err = redis.Bool(conn.Do("EXISTS", key)); err != nil {
		log.Error("conn.Do(EXISTS, %s) error(%+v)", key, err)
		return
	}
	if exist {
		ban = true
		return
	}
	if err = conn.Send("SET", key, "1"); err != nil {
		log.Error("conn.Send(EXPIRE, %s) error(%+v)", key, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.ipExpire); err != nil {
		log.Error("conn.Send(EXPIRE, %s) error(%+v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%+v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%+v)", err)
			return
		}
	}
	return
}

// BuvidBan intercepts illegal views.
func (d *Dao) BuvidBan(c context.Context, id, mid int64, ip, buvid string) (ban bool) {
	var (
		err   error
		exist bool
		key   = buvidBanKey(id, mid, ip, buvid)
		conn  = d.redis.Get(c)
	)
	defer conn.Close()
	if exist, err = redis.Bool(conn.Do("EXISTS", key)); err != nil {
		log.Error("conn.Do(EXISTS, %s) error(%+v)", key, err)
		return
	}
	if exist {
		ban = true
		return
	}
	if err = conn.Send("SET", key, "1"); err != nil {
		log.Error("conn.Send(EXPIRE, %s) error(%+v)", key, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.buvidExpire); err != nil {
		log.Error("conn.Send(EXPIRE, %s) error(%+v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%+v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%+v)", err)
			return
		}
	}
	return
}
