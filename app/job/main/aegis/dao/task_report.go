package dao

import (
	"context"
	"encoding/json"
	"sync"

	"go-common/app/job/main/aegis/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

var rmux sync.Mutex

//IncresByField .
func (d *Dao) IncresByField(c context.Context, bizid, flowid, uid int64, field string, value int64) (err error) {
	var (
		conn = d.redis.Get(c)
		hk   = model.PersonalHashKey(bizid, flowid, uid)
	)
	rmux.Lock()
	defer rmux.Unlock()
	defer conn.Close()

	if err = d.setSet(conn, hk); err != nil {
		return
	}
	if err = d.setHash(conn, hk, "ds"); err != nil {
		return
	}
	return d.setField(conn, hk, field, value)
}

//IncresTaskInOut 总进审量-出审量
func (d *Dao) IncresTaskInOut(c context.Context, bizid, flowid int64, inOrOut string) (err error) {
	var (
		conn = d.redis.Get(c)
		hk   = model.TotalHashKey(bizid, flowid)
	)
	rmux.Lock()
	defer rmux.Unlock()
	defer conn.Close()

	if err = d.setSet(conn, hk); err != nil {
		return
	}
	if err = d.setHash(conn, hk, inOrOut); err != nil {
		return
	}

	return d.setField(conn, hk, inOrOut, 1)
}

//FlushReport .
func (d *Dao) FlushReport(c context.Context) (data map[string][]byte, err error) {
	data = make(map[string][]byte)
	rmux.Lock()
	defer rmux.Unlock()

	conn := d.redis.Get(c)
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("SMEMBERS", model.SetKey))
	if err != nil {
		log.Error("SMEMBERS %s error(%v)", model.SetKey, err)
		return
	}
	if len(keys) == 0 {
		log.Info("FlushReport empty")
		return
	}

	for _, key := range keys {
		if err = conn.Send("HGETALL", key); err != nil {
			log.Error("HGETALL %s error(%v)", key, err)
			return
		}
	}
	conn.Flush()

	for _, key := range keys {
		var (
			bs []byte
			mp map[string]int64
		)
		if mp, err = redis.Int64Map(conn.Receive()); err != nil {
			log.Error("Receive error(%v)", err)
			return
		}
		if bs, err = json.Marshal(mp); err != nil {
			log.Error("Marshal mp(%+v) error(%v)", mp, err)
		}
		data[key] = bs
	}

	for _, key := range keys {
		conn.Do("DEL", key)
	}
	conn.Do("DEL", model.SetKey)
	return
}

//记录key
func (d *Dao) setSet(conn redis.Conn, hk string) (err error) {
	if _, err := conn.Do("SADD", model.SetKey, hk); err != nil {
		log.Error("setSet SADD(%s,%s) error(%v)", model.SetKey, hk, err)
	}
	return
}

//创建hash
func (d *Dao) setHash(conn redis.Conn, key string, defaultfield string) (err error) {
	var exist bool
	if exist, err = redis.Bool(conn.Do("EXISTS", key)); err != nil {
		log.Error("setHash EXISTS(%s) error(%v)", key, err)
		return
	}
	if !exist {
		if _, err = conn.Do("HMSET", key, defaultfield, 0); err != nil {
			log.Error("setHash HMSET(%s,%s,%d) error(%v)", key, defaultfield, 0, err)
		}
	}
	return
}

//每个field赋值
func (d *Dao) setField(conn redis.Conn, key string, field string, value int64) (err error) {
	var exist bool
	if exist, err = redis.Bool(conn.Do("HEXISTS", key, field)); err != nil {
		log.Error("setField HEXISTS(%s,%s,%s) error(%v)", key, field, err)
		return
	}
	if !exist {
		if _, err = conn.Do("HMSET", key, field, 0); err != nil {
			log.Error("setField HMSET(%s,%s,%d) error(%v)", key, field, 0, err)
		}
	}
	if _, err = conn.Do("HINCRBY", key, field, value); err != nil {
		log.Error("setField HINCRBY(%s,%s,%d) error(%v)", key, field, 1, err)
	}

	return nil
}
