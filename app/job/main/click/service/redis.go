package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go-common/library/cache/redis"
	"go-common/library/log"

	farm "github.com/dgryski/go-farm"
)

const (
	// aid_ip
	_anonymousePlayedKey = "nm:%d"
	// aid bvid
	_anonymouseBvIDKey = "nb:%d"
	// bvid's last played aid
	_bvIDLastPlayedKey = "bv:%d"
	// Mid last played key
	_midHashKey    = "mid:%d"
	_buvidToDidKey = "%d:bdid"
	// 改短成6分钟
	_hkeyExpire = 600
)

func (s *Service) midKey(mid int64) (key string) {
	key = fmt.Sprintf(_midHashKey, mid%s.c.HashNum)
	return
}

func (s *Service) bvKey(bvID string) (key string) {
	num := int64(farm.Hash32([]byte(bvID)))
	key = fmt.Sprintf(_bvIDLastPlayedKey, num%s.c.HashNum)
	return
}

func (s *Service) buvidToDidKey(buvid string, aid int64) (key string) {
	num := int64(farm.Hash32([]byte(fmt.Sprintf("%s_%d", buvid, aid))))
	key = fmt.Sprintf(_buvidToDidKey, num%s.c.HashNum)
	return
}

func (s *Service) anonymouseKey(aid, epID int64, ip string) (key string) {
	var str string
	if epID == 0 {
		str = strconv.Itoa(int(aid)) + ip
	} else {
		str = strconv.Itoa(int(aid)) + ip + strconv.Itoa(int(epID))
	}
	num := int64(farm.Hash32([]byte(str)))
	key = fmt.Sprintf(_anonymousePlayedKey, num%s.c.HashNum)
	return
}

func (s *Service) anonymouseBvIDKey(aid int64, bvid string) (key string) {
	str := strconv.Itoa(int(aid)) + bvid
	num := int64(farm.Hash32([]byte(str)))
	key = fmt.Sprintf(_anonymouseBvIDKey, num%s.c.HashNum)
	return
}

// canCount 同一个IP一分钟，同一个bvid五分钟
func (s *Service) canCount(c context.Context, aid, epID int64, ip string, stime int64, bvid string) (can bool) {
	var (
		err            error
		lastPlayTime   int64
		bvLastPlayTime int64
		hKey           = s.anonymouseKey(aid, epID, ip)
		hbvKey         = s.anonymouseBvIDKey(aid, bvid)
		conn           = s.redis.Get(c)
		pTime          = stime + s.c.CacheConf.NewAnonymousCacheTime
		now            = time.Now().Unix()
		ipCan          bool
		bvCan          bool
	)
	var field = aid
	if epID > 0 {
		field = epID
	}
	defer conn.Close()
	if lastPlayTime, err = redis.Int64(conn.Do("HGET", hKey, field)); err != nil {
		if err == redis.ErrNil {
			if _, err = conn.Do("HSET", hKey, field, pTime); err != nil {
				log.Error("conn.Do(HSET, %s, %d, %d) error(%v)", hKey, field, pTime, err)
				return
			}
			ipCan = true
		} else {
			log.Error("conn.Do(HGET, %s, %d) error(%v)", hKey, field, err)
			return
		}
	}
	if bvLastPlayTime, err = redis.Int64(conn.Do("HGET", hbvKey, field)); err != nil {
		if err == redis.ErrNil {
			if _, err = conn.Do("HSET", hbvKey, field, pTime); err != nil {
				log.Error("conn.Do(HSET, %s, %d, %d)", hbvKey, field, pTime)
				return
			}
			bvCan = true
		} else {
			log.Error("conn.Do(HGET, %s, %d) error(%v)", hbvKey, field, err)
			return
		}
	}
	if ipCan && bvCan {
		can = true
		return
	}
	if now > lastPlayTime && now > bvLastPlayTime {
		if err = conn.Send("HSET", hKey, field, now+s.c.CacheConf.NewAnonymousCacheTime); err != nil {
			log.Error("conn.Send(HSET, %s, %s, %d) error(%v)", hKey, field, now+s.c.CacheConf.NewAnonymousCacheTime, err)
			return
		}
		if err = conn.Send("EXPIRE", hKey, _hkeyExpire); err != nil {
			log.Error("conn.Do(EXPIRE, %s, %d) error(%v)", hKey, _hkeyExpire, err)
			return
		}
		if err = conn.Send("HSET", hbvKey, field, now+s.c.CacheConf.NewAnonymousBvCacheTime); err != nil {
			log.Error("conn.Send(HSET, %s, %d, %d) error(%v)", hbvKey, field, now+s.c.CacheConf.NewAnonymousBvCacheTime, err)
			return
		}
		if err = conn.Send("EXPIRE", hbvKey, _hkeyExpire); err != nil {
			log.Error("conn.Do(EXPIRE, %s, %d) error(%v)", hbvKey, _hkeyExpire, err)
			return
		}
		if err = conn.Flush(); err != nil {
			log.Error("conn.Flush error(%v)")
			return
		}
		for i := 0; i < 4; i++ {
			if _, err = conn.Receive(); err != nil {
				log.Error("conn.Receive error(%v)", err)
				return
			}
		}
		can = true
	}
	return
}

func (s *Service) isReplay(c context.Context, mid, aid int64, bvID string, gapTime int64) (is bool) {
	var (
		hKey  string
		field string
		conn  = s.redis.Get(c)
		value int64
		err   error
		now   = time.Now().Unix()
	)
	defer conn.Close()
	if mid > 0 {
		hKey = s.midKey(mid)
		field = strconv.Itoa(int(mid))
	} else {
		hKey = s.bvKey(bvID)
		field = bvID
	}
	// aid << 32 | ptime
	if value, err = redis.Int64(conn.Do("HGET", hKey, field)); err != nil {
		if err != redis.ErrNil {
			log.Error("conn.Do(HGET, %s, %s) error(%s)", hKey, field, err)
			return
		}
		err = nil
	}
	if value != 0 {
		rOid := value >> 32
		rNow := value & 0xffffffff
		if rOid == aid && now-rNow < gapTime {
			is = true
			return
		}
	}
	value = aid<<32 | now
	if err = conn.Send("HSET", hKey, field, value); err != nil {
		log.Error("conn.Do(HSET, %s, %s, %s) error(%v)", hKey, field, value, err)
		return
	}
	if err = conn.Send("EXPIRE", hKey, _hkeyExpire); err != nil {
		log.Error("conn.Do(EXPIRE, %s, %d)", hKey, _hkeyExpire)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)")
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive error(%v)", err)
			return
		}
	}
	return
}

func (s *Service) getRealDid(c context.Context, buvid string, aid int64) (did string, err error) {
	var (
		conn = s.redis.Get(c)
		key  = s.buvidToDidKey(buvid, aid)
	)
	defer conn.Close()
	if did, err = redis.String(conn.Do("HGET", key, buvid)); err != nil {
		if err != redis.ErrNil {
			log.Error("redis.String(conn.Do(HGET, %s, %s)) error(%v)", key, buvid, err)
			return
		}
		err = nil
	}
	return
}

func (s *Service) setRealDid(c context.Context, buvid string, aid int64, did string) (err error) {
	var (
		conn = s.redis.Get(c)
		key  = s.buvidToDidKey(buvid, aid)
	)
	defer conn.Close()
	if err = conn.Send("HSET", key, buvid, did); err != nil {
		log.Error("conn.Do(HSET, %s, %s, %s) error(%v)", key, buvid, did, err)
		return
	}
	if err = conn.Send("EXPIRE", key, _hkeyExpire); err != nil {
		log.Error("conn.Send(EXPIRE, %s, %d) error(%v)", key, _hkeyExpire, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive error(%v)", err)
			return
		}
	}
	return
}
