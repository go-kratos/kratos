package dao

import (
	"context"
	"strconv"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_qidByTypeID      = "v3_qus_tids_"
	_extraQidByTypeID = "v3_eq_t_"
)

func qusByType(tid int) string {
	return _qidByTypeID + strconv.FormatInt(int64(tid), 10)
}

func extraQidByType(tid int8) string {
	return _extraQidByTypeID + strconv.FormatInt(int64(tid), 10)
}

func (d *Dao) pingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}

// QidByType get question by type.
func (d *Dao) QidByType(c context.Context, tid int, num uint8) (ids []int64, err error) {
	key := qusByType(tid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if ids, err = redis.Int64s(conn.Do("SRANDMEMBER", key, num)); err != nil {
		log.Error("RandBaseQs conn.Send('SRANDMEMBER', %s, %d) error(%v)", key, num, err)
	}
	return
}

// SetQids set question ids.
func (d *Dao) SetQids(c context.Context, qs []int64, typeID int) (err error) {
	if len(qs) == 0 {
		return
	}
	key := qusByType(typeID)
	conn := d.redis.Get(c)
	defer conn.Close()
	args := make([]interface{}, 0, len(qs)+1)
	args = append(args, key)
	for _, q := range qs {
		args = append(args, q)
	}
	if _, err = conn.Do("SADD", args...); err != nil {
		log.Error("conn.Send(SADD, %v) error(%v)", args, err)
	}
	return
}

// SetExtraQids set extra question ids.
func (d *Dao) SetExtraQids(c context.Context, qs []int64, ans int8) (err error) {
	if len(qs) == 0 {
		return
	}
	key := extraQidByType(ans)
	conn := d.redis.Get(c)
	defer conn.Close()
	args := make([]interface{}, 0, len(qs)+1)
	args = append(args, key)
	for _, q := range qs {
		args = append(args, q)
	}
	if _, err = conn.Do("SADD", args...); err != nil {
		log.Error("conn.Send(SADD, %v) error(%v)", args, err)
	}
	return
}

// DelQidsCache del qids cahce.
func (d *Dao) DelQidsCache(c context.Context, typeID int) (err error) {
	key := qusByType(typeID)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("DEL", key); err != nil {
		log.Error("conn.Send(DEL, %s) error(%v)", key, err)
	}
	return
}

// DelExtraQidsCache del extra qids cahce.
func (d *Dao) DelExtraQidsCache(c context.Context, ans int8) (err error) {
	key := extraQidByType(ans)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("DEL", key); err != nil {
		log.Error("conn.Send(DEL, %s) error(%v)", key, err)
	}
	return
}

// ExtraQidByType extra qis by type.
func (d *Dao) ExtraQidByType(c context.Context, ans int8, num uint8) (ids []int64, err error) {
	key := extraQidByType(ans)
	conn := d.redis.Get(c)
	defer conn.Close()
	if ids, err = redis.Int64s(conn.Do("SRANDMEMBER", key, num)); err != nil {
		log.Error("ExtraQidByType conn.Send('SRANDMEMBER', %s, %d) error(%v)", key, num, err)
	}
	return
}
