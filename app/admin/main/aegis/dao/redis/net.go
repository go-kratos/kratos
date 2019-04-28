package redis

import (
	"context"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

//SetMulti setex key expire val, kvs as multi key val
func (d *Dao) SetMulti(c context.Context, kvs map[string]string) (err error) {
	if len(kvs) == 0 {
		return
	}

	conn := d.cluster.Get(c)
	defer conn.Close()

	//拼接参数
	for key, val := range kvs {
		if err = conn.Send("SETEX", key, d.netExpire, val); err != nil {
			log.Error("SetMulti conn.send(SETEX) error(%v)", err)
			return
		}
	}

	if err = conn.Flush(); err != nil {
		log.Error("SetMulti conn.Flush error(%v)", err)
	}
	return
}

//MGet get multi key value
func (d *Dao) MGet(c context.Context, keys ...string) (dest []string, err error) {
	//检测参数
	if len(keys) == 0 {
		return
	}

	//拼接查询参数+ redis请求
	args := redis.Args{}
	for _, item := range keys {
		args = args.Add(item)
	}
	conn := d.cluster.Get(c)
	defer conn.Close()
	if dest, err = redis.Strings(conn.Do("MGET", args...)); err != nil {
		log.Error("MGet conn.Do(mget) error(%v) args(%+v)", err, args)
	}
	return
}

//DelMulti del redis keys
func (d *Dao) DelMulti(c context.Context, keys ...string) (err error) {
	if len(keys) == 0 {
		return
	}

	conn := d.cluster.Get(c)
	defer conn.Close()

	args := redis.Args{}
	for _, k := range keys {
		args = args.Add(k)
	}

	if _, err = redis.Int(conn.Do("DEL", args...)); err != nil {
		log.Error("DelMulti conn.Do(del) error(%v) args(%+v)", err, args)
	}
	return
}
