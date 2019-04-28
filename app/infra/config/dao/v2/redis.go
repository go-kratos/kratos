package v2

import (
	"context"
	"encoding/json"
	"time"

	"go-common/app/infra/config/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	expireDuration = 3 * time.Hour
)

// Hosts return service hosts from redis.
func (d *Dao) Hosts(c context.Context, svr string) (hosts []*model.Host, err error) {
	var (
		dels []string
		now  = time.Now()
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	res, err := redis.Strings(conn.Do("HGETALL", svr))
	if err != nil {
		log.Error("conn.Do(HGETALL, %s) error(%v)", svr, err)
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
		if _, err1 := conn.Do("HDEL", svr, dels); err1 != nil {
			log.Error("conn.Do(HDEL, %s, %v) error(%v)", svr, dels, err1)
		}
	}
	return
}

// SetHost add service host to redis.
func (d *Dao) SetHost(c context.Context, host *model.Host, svr string) (err error) {
	b, err := json.Marshal(host)
	if err != nil {
		log.Error("json.Marshal(%s) error(%v)", host, err)
		return
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("HSET", svr, host.Name, string(b)); err != nil {
		log.Error("conn.Do(SET, %s, %s, %v) error(%v)", svr, host.Name, host, err)
	}
	return
}

// ClearHost clear all hosts.
func (d *Dao) ClearHost(c context.Context, svr string) (err error) {
	var (
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if _, err = conn.Do("DEL", svr); err != nil {
		log.Error("conn.Do(DEL, %s) error(%v)", svr, err)
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
