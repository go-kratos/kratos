package databus

import (
	"context"
	"encoding/json"

	"go-common/app/admin/main/videoup/model/archive"
	"go-common/app/admin/main/videoup/model/message"
	"go-common/library/cache/redis"
	"go-common/library/conf/env"
	"go-common/library/log"
)

const (
	_multSyncList  = "m_sync_list"
	_prefixMsgInfo = "videoup_admin_msg"
)

// PopMsgCache get databus message from redis
func (d *Dao) PopMsgCache(c context.Context) (msg *message.Videoup, err error) {
	var (
		conn = d.redis.Get(c)
		bs   []byte
	)
	defer conn.Close()
	if bs, err = redis.Bytes(conn.Do("LPOP", fixRedisList(_prefixMsgInfo))); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(LPOP, %s) error(%v)", fixRedisList(_prefixMsgInfo), err)
		}
		return
	}
	msg = &message.Videoup{}
	if err = json.Unmarshal(bs, msg); err != nil {
		log.Error("json.Unmarshal error(%v)", err)
	}
	return
}

// PushMsgCache add message into redis.
func (d *Dao) PushMsgCache(c context.Context, msg *message.Videoup) (err error) {
	var (
		bs   []byte
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if bs, err = json.Marshal(msg); err != nil {
		log.Error("json.Marshal(%s) error(%v)", bs, err)
		return
	}
	if _, err = conn.Do("RPUSH", fixRedisList(_prefixMsgInfo), bs); err != nil {
		log.Error("conn.Do(RPUSH, %s) error(%v)", bs, err)
	}
	return
}

func fixRedisList(list string) (target string) {
	if env.DeployEnv == env.DeployEnvPre {
		target = "pre_" + list
	} else {
		target = list
	}
	return
}

// PushMultSync rpush stuct item to redis
func (d *Dao) PushMultSync(c context.Context, sync *archive.MultSyncParam) (ok bool, err error) {
	var (
		conn = d.redis.Get(c)
		bs   []byte
	)
	defer conn.Close()
	if bs, err = json.Marshal(sync); err != nil {
		log.Error("json.Marshal(%v) error(%v)", sync, err)
		return
	}
	if err = conn.Send("SADD", fixRedisList(_multSyncList), bs); err != nil {
		log.Error("conn.Send(SADD, %s, %s) error(%v)", fixRedisList(_multSyncList), bs, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	if ok, err = redis.Bool(conn.Receive()); err != nil {
		log.Error("conn.Receive error(%v)", err)
	}
	return
}

// PopMultSync lpop stuct item from redis
func (d *Dao) PopMultSync(c context.Context) (res *archive.MultSyncParam, err error) {
	var (
		conn = d.redis.Get(c)
		bs   []byte
		sync = &archive.MultSyncParam{}
	)
	defer conn.Close()
	if bs, err = redis.Bytes(conn.Do("SPOP", fixRedisList(_multSyncList))); err != nil && err != redis.ErrNil {
		log.Error("redis.Bytes(conn.Do(SPOP, %s)) error(%v)", fixRedisList(_multSyncList), err)
		return
	}
	if len(bs) == 0 {
		return
	}
	if err = json.Unmarshal(bs, sync); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", sync, err)
		return
	}
	res = sync
	return
}
