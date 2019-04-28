package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-common/app/infra/config/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	expireDuration = 3 * time.Hour
	_hostKey       = "%s_%s"
)

// Hostkey host cache key
func hostKey(svr, env string) string {
	return fmt.Sprintf(_hostKey, svr, env)
}

// Hosts return service hosts from redis.
func (d *Dao) Hosts(c context.Context, svr, env string) (hosts []*model.Host, err error) {
	var (
		dels    []string
		now     = time.Now()
		hostkey = hostKey(svr, env)
		conn    = d.redis.Get(c)
	)
	defer conn.Close()
	res, err := redis.Strings(conn.Do("HGETALL", hostkey))
	if err != nil {
		log.Error("conn.Do(HGETALL, %s) error(%v)", hostkey, err)
		return
	}
	for i, r := range res {
		if i%2 == 0 {
			continue
		}
		h := &model.Host{}
		if err = json.Unmarshal([]byte(r), h); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", r, err)
			return
		}
		if now.Sub(h.HeartbeatTime.Time()) <= d.expire+5 {
			h.State = model.HostOnline
			hosts = append(hosts, h)
		} else if now.Sub(h.HeartbeatTime.Time()) >= expireDuration {
			dels = append(dels, h.Name)
		} else {
			h.State = model.HostOffline
			hosts = append(hosts, h)
		}
	}
	if len(dels) > 0 {
		if _, err1 := conn.Do("HDEL", hostkey, dels); err1 != nil {
			log.Error("conn.Do(HDEL, %s, %v) error(%v)", hostkey, dels, err1)
		}
	}
	return
}

// SetHost add service host to redis.
func (d *Dao) SetHost(c context.Context, host *model.Host, svr, env string) (err error) {
	hostkey := hostKey(svr, env)
	b, err := json.Marshal(host)
	if err != nil {
		log.Error("json.Marshal(%s) error(%v)", host, err)
		return
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("HSET", hostkey, host.Name, string(b)); err != nil {
		log.Error("conn.Do(SET, %s, %s, %v) error(%v)", hostkey, host.Name, host, err)
	}
	return
}

// ClearHost clear all hosts.
func (d *Dao) ClearHost(c context.Context, svr, env string) (err error) {
	var (
		hostkey = hostKey(svr, env)
		conn    = d.redis.Get(c)
	)
	defer conn.Close()
	if _, err = conn.Do("DEL", hostkey); err != nil {
		log.Error("conn.Do(DEL, %s) error(%v)", hostkey, err)
	}
	return
}

// Ping check Redis connection
func (d *Dao) pingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}
