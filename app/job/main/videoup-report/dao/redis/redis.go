package redis

import (
	"context"
	"fmt"
	"go-common/library/log"
	"time"
)

const (
	_videoJamTime = "va_v_jam_time"
)

// SetVideoJam set video traffic jam time
func (d *Dao) SetVideoJam(c context.Context, jamTime int) (err error) {

	var conn = d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("SET", _videoJamTime, jamTime); err != nil {
		log.Error("conn.Do(SET, %s, %d) error(%v)", _videoJamTime, jamTime, err)
	}
	return
}

// TrackAddRedis add track redis
func (d *Dao) TrackAddRedis(c context.Context, key, value string) (err error) {
	if key == "" {
		log.Warn("TrackAddRedis add empty key(%s) value(%v)", key, value)
		return fmt.Errorf("empty key")
	}
	var conn = d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("ZADD", key, time.Now().Unix(), value); err != nil {
		log.Error("conn.Do(ZADD, %s, %s) error(%v)", key, value, err)
	}
	return
}

// TrackAddVideosRedis add track video redis
func (d *Dao) TrackAddVideosRedis(c context.Context, keys []string, value string) (err error) {
	var conn = d.redis.Get(c)
	defer conn.Close()
	for _, key := range keys {
		if err = conn.Send("ZADD", key, time.Now().Unix(), value); err != nil {
			log.Error("conn.Send(ZADD, %s, %s) error(%v)", key, value, err)
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("add conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < len(keys); i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("add conn.Receive(%d) error(%v)", i+1, err)
			return
		}
	}
	return
}

// TrackRemRedis remove track redis
func (d *Dao) TrackRemRedis(c context.Context, key, value string) (err error) {
	var conn = d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("ZREM", key, value); err != nil {
		log.Error("conn.Do(ZREM, %s, %s) error(%v)", key, value, err)
	}
	return
}

// TrackRemVideosRedis remove video track redis
func (d *Dao) TrackRemVideosRedis(c context.Context, keys []string, value string) (err error) {
	var conn = d.redis.Get(c)
	defer conn.Close()
	for _, key := range keys {
		if err = conn.Send("ZREM", key, value); err != nil {
			log.Error("conn.Send(ZREM, %s, %s) error(%v)", key, value, err)
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("add conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < len(keys); i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("add conn.Receive(%d) error(%v)", i+1, err)
			return
		}
	}
	return
}
