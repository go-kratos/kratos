package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-common/app/service/main/history/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const _deleteDuration = 3600 * 12

// keyIndex return history index key.
func keyIndex(business string, mid int64) string {
	return fmt.Sprintf("i_%d_%s", mid, business)
}

// keyHistory return history key.
func keyHistory(business string, mid int64) string {
	return fmt.Sprintf("h_%d_%s", mid, business)
}

// HistoriesCache return the user histories from redis.
func (d *Dao) HistoriesCache(c context.Context, merges []*model.Merge) (res []*model.History, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	for _, merge := range merges {
		key := keyHistory(d.BusinessesMap[merge.Bid].Name, merge.Mid)
		if err = conn.Send("HGET", key, merge.Kid); err != nil {
			log.Error("conn.Send(HGET %v %v) error(%v)", key, merge.Kid, err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	for i := 0; i < len(merges); i++ {
		var value []byte
		if value, err = redis.Bytes(conn.Receive()); err != nil {
			if err == redis.ErrNil {
				err = nil
				continue
			}
			log.Error("conn.Receive error(%v)", err)
			return
		}
		h := &model.History{}
		if err = json.Unmarshal(value, h); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", value, err)
			err = nil
			continue
		}
		h.BusinessID = d.BusinessNames[h.Business].ID
		res = append(res, h)
	}
	return
}

// TrimCache trim history.
func (d *Dao) TrimCache(c context.Context, business string, mid int64, limit int) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	aids, err := redis.Int64s(conn.Do("ZRANGE", keyIndex(business, mid), 0, -limit-1))
	if err != nil {
		log.Error("conn.Do(ZRANGE %v) error(%v)", keyIndex(business, mid), err)
		return
	}
	if len(aids) == 0 {
		return
	}
	return d.DelCache(c, business, mid, aids)
}

// DelCache delete the history redis.
func (d *Dao) DelCache(c context.Context, business string, mid int64, aids []int64) (err error) {
	var (
		key1  = keyIndex(business, mid)
		key2  = keyHistory(business, mid)
		args1 = []interface{}{key1}
		args2 = []interface{}{key2}
	)
	for _, aid := range aids {
		args1 = append(args1, aid)
		args2 = append(args2, aid)
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("ZREM", args1...); err != nil {
		log.Error("conn.Send(ZREM %s,%v) error(%v)", key1, aids, err)
		return
	}
	if err = conn.Send("HDEL", args2...); err != nil {
		log.Error("conn.Send(HDEL %s,%v) error(%v)", key2, aids, err)
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

// DelLock delete proc lock
func (d *Dao) DelLock(c context.Context) (ok bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := "his_job_del_proc"
	if err = conn.Send("SETNX", key, time.Now().Unix()); err != nil {
		log.Error("DelLock conn.SETNX() error(%v)", err)
		return
	}
	if err = conn.Send("EXPIRE", key, _deleteDuration); err != nil {
		log.Error("DelLock conn.Expire() error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("DelLock conn.Flush() error(%v)", err)
		return
	}
	if ok, err = redis.Bool(conn.Receive()); err != nil {
		log.Error("conn.Receive() error(%v)", err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		log.Error("conn.Receive() error(%v)", err)
	}
	return
}
