package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	modtask "go-common/app/admin/main/aegis/model/task"
	"go-common/library/cache/redis"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

func lockKey(businessID, flowID int64) string {
	return fmt.Sprintf("aegis_lock_%d_%d", businessID, flowID)
}

func (d *Dao) getlock(c context.Context, bizid, flowid int64) (ok bool) {
	var (
		conn = d.cluster.Get(c)
		key  = lockKey(bizid, flowid)
		err  error
	)
	defer conn.Close()

	if ok, err = redis.Bool(conn.Do("SETNX", key, "1")); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(SETNX(%s)) error(%v)", key, err)
			return
		}
	}

	if ok {
		conn.Do("EXPIRE", key, 3)
	}
	return
}

//SeizeTask .
func (d *Dao) SeizeTask(c context.Context, businessID, flowID, uid, count int64) (hitids []int64, missids []int64, others map[int64]int64, err error) {
	var (
		lock   bool
		pubkey = publicKey(businessID, flowID)
		ids    []int64
	)

	// 1. 抢占分布式锁
	for lc := 0; lc < 3; lc++ {
		if lock = d.getlock(c, businessID, flowID); lock {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	if !lock {
		log.Error("getlock getlock fail(%d,%d,%d)", businessID, flowID, uid)
		err = ecode.AegisTaskBusy
		return
	}

	conn := d.cluster.Get(c)
	defer conn.Close()
	defer conn.Do("DEL", lockKey(businessID, flowID))

	var (
		head, tail = int64(0), int64(count)
	)
	// 2. 从 public 按权重从高到低取出一批来
	for {
		if ids, err = redis.Int64s(conn.Do("ZRANGE", pubkey, head, tail)); err != nil {
			log.Error("redis (ZRANGE,%s,%d,%d) error(%v)", pubkey, 0, count, err)
			return
		}
		head += count
		tail += count
		if len(ids) == 0 {
			break
		}

		for _, id := range ids {
			if err = conn.Send("GET", haskKey(id)); err != nil {
				log.Error("redis (GET,%s) error(%v)", haskKey(id), err)
				return
			}
		}
		conn.Flush()

		var enough bool

		for _, id := range ids {
			var (
				bs []byte
				e  error
			)
			bs, e = redis.Bytes(conn.Receive())

			if e != nil {
				log.Error("Receive Weight(%d) error(%v)", id, errors.WithStack(e))
				missids = append(missids, id)
				continue
			}

			task := &modtask.Task{}
			if e = json.Unmarshal(bs, task); err != nil {
				log.Error("json.Unmarshal error(%v)", errors.WithStack(e))
				missids = append(missids, id)
				continue
			}
			if task.ID != id {
				log.Error("id(%d-%d)不匹配", task.ID, id)
				missids = append(missids, id)
				continue
			}

			if task.UID != 0 && task.UID != uid {
				log.Info("id(%d) 任务已经指派给(%d)", task.ID, task.UID)
				missids = append(missids, id)
				continue
			}

			hitids = append(hitids, id)

			if len(hitids) >= int(count) {
				enough = true
				break
			}
		}
		if enough {
			break
		}
	}

	personKey := personalKey(businessID, flowID, uid)
	for _, id := range hitids {
		conn.Send("ZREM", pubkey, formatID(id))
		conn.Send("LREM", personKey, 0, id)
		conn.Send("RPUSH", personKey, id)
	}
	conn.Flush()
	for i := 0; i < len(hitids)*3; i++ {
		conn.Receive()
	}

	log.Info("rangefunc count(%d) hitids(%v) missids(%v)", count, hitids, missids)
	return
}

/*
遍历personal,delay,public。
在缓存中进行状态校验，public还要补充缓存权重
*/
func (d *Dao) rangefuncCluster(c context.Context, listtype string, opt *modtask.ListOptions) (tasks map[int64]*modtask.Task, count int64, hitids, missids []int64, err error) {
	var (
		key              string
		LENCMD, RANGECMD = "LLEN", "LRANGE"
		ids              []int64
	)

	conn := d.cluster.Get(c)
	defer conn.Close()

	switch listtype {
	case "public":
		LENCMD, RANGECMD = "ZCARD", "ZRANGE"
		key = publicKey(opt.BusinessID, opt.FlowID)
	case "personal":
		key = personalKey(opt.BusinessID, opt.FlowID, opt.UID)
	case "delay":
		key = delayKey(opt.BusinessID, opt.FlowID, opt.UID)
	}

	// 1. 长度
	if count, err = redis.Int64(conn.Do(LENCMD, key)); err != nil {
		log.Error("redis (%s,%s) error(%v)", LENCMD, key, err)
		return
	}
	if count == 0 {
		return
	}

	if ids, err = redis.Int64s(conn.Do(RANGECMD, key, (opt.Pn-1)*opt.Ps, opt.Pn*opt.Ps-1)); err != nil {
		log.Error("redis (%s,%s,%d,%d) error(%v)", LENCMD, key, (opt.Pn-1)*opt.Ps, opt.Pn*opt.Ps, err)
		return
	}

	for _, id := range ids {
		if err = conn.Send("GET", haskKey(id)); err != nil {
			log.Error("redis (GET,%s) error(%v)", haskKey(id), err)
			return
		}
		if listtype == "public" {
			if err = conn.Send("ZSCORE", key, formatID(id)); err != nil {
				log.Error("redis (ZSCORE,%s,%s) error(%v)", key, formatID(id), err)
				return
			}
		}
	}
	conn.Flush()

	tasks = make(map[int64]*modtask.Task)
	for _, id := range ids {
		var (
			bs []byte
			e  error
			wt int64
		)
		bs, e = redis.Bytes(conn.Receive())
		if listtype == "public" {
			wt, _ = redis.Int64(conn.Receive())
			wt = -wt
		}

		if e != nil {
			log.Error("Receive Weight(%d) error(%v)", id, errors.WithStack(e))
			missids = append(missids, id)
			continue
		}

		task := &modtask.Task{}
		if e = json.Unmarshal(bs, task); err != nil {
			log.Error("json.Unmarshal error(%v)", errors.WithStack(e))
			missids = append(missids, id)
			continue
		}
		if task.ID != id {
			log.Error("id(%d-%d)不匹配", task.ID, id)
			missids = append(missids, id)
			continue
		}
		// 缓存里状态同步不实时，不能用作校验

		tasks[task.ID] = task
		hitids = append(hitids, id)
	}

	log.Info("rangefunc count(%d) hitids(%v) missids(%v)", count, hitids, missids)
	return
}
