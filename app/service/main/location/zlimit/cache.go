package zlimit

import (
	"context"
	"strconv"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

func keyZlimit(aid int64) (key string) {
	key = _prefixBlackList + strconv.FormatInt(aid, 10)
	return
}

// existsRule if existes ruls in redis
func (s *Service) existsRule(c context.Context, aid int64) (ok bool, err error) {
	conn := s.redis.Get(c)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("EXISTS", keyZlimit(aid))); err != nil {
		log.Error("conn.DO(HEXISTS) error(%v)", err)
	}
	return
}

// rule get zone rule from redis
func (s *Service) rule(c context.Context, aid int64, zoneids []int64) (res []int64, err error) {
	var playauth int64
	key := keyZlimit(aid)
	conn := s.redis.Get(c)
	defer conn.Close()
	for _, v := range zoneids {
		if err = conn.Send("HGET", key, v); err != nil {
			log.Error("rule conn.Send(HGET, %s, %d) error(%v)", key, v, err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("rule conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < len(zoneids); i++ {
		if playauth, err = redis.Int64(conn.Receive()); err != nil {
			if err != redis.ErrNil {
				log.Error("rule conn.Receive()%d error(%v)", i+1, err)
				return
			}
			err = nil
		}
		res = append(res, playauth)
	}
	return
}

// addRule add zone rule from redis
func (s *Service) addRule(c context.Context, zoneids map[int64]map[int64]int64) (err error) {
	var key string
	conn := s.redis.Get(c)
	defer conn.Close()
	count := 0
	for aid, zids := range zoneids {
		if key == "" {
			key = keyZlimit(aid)
		}
		for zid, auth := range zids {
			if err = conn.Send("HSET", key, zid, auth); err != nil {
				log.Error("add conn.Send error(%v)", err)
				return
			}
			count++
		}
	}
	if err = conn.Send("EXPIRE", key, s.expire); err != nil {
		log.Error("add conn.Send error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("add conn.Flush error(%v)", err)
		return
	}
	for i := 0; i <= count; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("add conn.Receive()%d error(%v)", i+1, err)
			return
		}
	}
	return
}
