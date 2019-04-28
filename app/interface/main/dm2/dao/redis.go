package dao

import (
	"context"
	"fmt"

	"go-common/library/cache/redis"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	// dm xml list v1
	_prefixDM = "dm_v1_%d_%d" // dm_v1_type_oid

	_broadcastLimitFmt = "b_room_%d_%d" // b_room_type_oid
)

func keyDM(tp int32, oid int64) (key string) {
	return fmt.Sprintf(_prefixDM, tp, oid)
}

func keyBroadcastLimit(tp int32, oid int64) (key string) {
	return fmt.Sprintf(_broadcastLimitFmt, tp, oid)
}

// DMCache 获取redis列表中的弹幕.
func (d *Dao) DMCache(c context.Context, tp int32, oid int64) (res [][]byte, err error) {
	conn := d.dmRds.Get(c)
	key := keyDM(tp, oid)
	if res, err = redis.ByteSlices(conn.Do("ZRANGE", key, 0, -1)); err != nil {
		log.Error("conn.Do(ZRANGE %s) error(%v)", key, err)
	}
	conn.Close()
	return
}

// ExpireDMCache expire dm.
func (d *Dao) ExpireDMCache(c context.Context, tp int32, oid int64) (ok bool, err error) {
	key := keyDM(tp, oid)
	conn := d.dmRds.Get(c)
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.dmRdsExpire)); err != nil {
		log.Error("conn.Do(EXPIRE %s) error(%v)", key, err)
	}
	conn.Close()
	return
}

// TrimDMCache 从redis列表中pop掉count条弹幕.
func (d *Dao) TrimDMCache(c context.Context, tp int32, oid, count int64) (err error) {
	conn := d.dmRds.Get(c)
	key := keyDM(tp, oid)
	if _, err = conn.Do("ZREMRANGEBYRANK", key, 0, count-1); err != nil {
		log.Error("conn.Do(ZREMRANGEBYRANK %s) error(%v)", key, err)
	}
	conn.Close()
	return
}

// IncrPubCnt increase pub count of user.
func (d *Dao) IncrPubCnt(c context.Context, mid, color int64, mode, fontsize int32, ip, msg string) (err error) {
	conn := d.dmRds.Get(c)
	defer conn.Close()
	key := keyPubCntLock(mid, color, mode, fontsize, ip, msg)
	if err = conn.Send("INCR", key); err != nil {
		log.Error("conn.Send(INCR %s) error(%v)", key, err)
		return
	}
	if err = conn.Send("EXPIRE", key, 300); err != nil {
		log.Error("conn.Send(EXPIRE %s) error(%v)", key, err)
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

// PubCnt get dm pub count of user.
func (d *Dao) PubCnt(c context.Context, mid, color int64, mode, fontsize int32, ip, msg string) (count int64, err error) {
	conn := d.dmRds.Get(c)
	defer conn.Close()
	key := keyPubCntLock(mid, color, mode, fontsize, ip, msg)
	if count, err = redis.Int64(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(GET %s) error(%v)", key, err)
		}
		return
	}
	return
}

// IncrCharPubCnt increase character pub count of user.
func (d *Dao) IncrCharPubCnt(c context.Context, mid, oid int64) (err error) {
	conn := d.dmRds.Get(c)
	defer conn.Close()
	key := keyCharPubLock(mid, oid)
	if err = conn.Send("INCR", key); err != nil {
		log.Error("conn.Send(INCR %s) error(%v)", key, err)
		return
	}
	if err = conn.Send("EXPIRE", key, 60); err != nil {
		log.Error("conn.Send(EXPIRE %s) error(%v)", key, err)
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

// CharPubCnt get character pub count of user.
func (d *Dao) CharPubCnt(c context.Context, mid, oid int64) (count int64, err error) {
	conn := d.dmRds.Get(c)
	defer conn.Close()
	key := keyCharPubLock(mid, oid)
	if count, err = redis.Int64(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(GET %s) error(%v)", key, err)
		}
		return
	}
	return
}

// DelCharPubCnt delete char
func (d *Dao) DelCharPubCnt(c context.Context, mid, oid int64) (err error) {
	conn := d.dmRds.Get(c)
	key := keyCharPubLock(mid, oid)
	if _, err = conn.Do("DEL", key); err != nil {
		log.Error("conn.Do(DEL %s) error(%v)", key, err)
	}
	conn.Close()
	return
}

// BroadcastLimit .
func (d *Dao) BroadcastLimit(c context.Context, oid int64, tp int32, count, interval int) (err error) {
	conn := d.dmRds.Get(c)
	key := keyBroadcastLimit(tp, oid)
	defer conn.Close()
	incred, err := redis.Int64(conn.Do("INCR", key))
	if err != nil {
		return nil
	}
	if incred == 1 {
		conn.Do("EXPIRE", key, interval)
	}
	if incred > int64(count) {
		return ecode.LimitExceed
	}
	return
}
