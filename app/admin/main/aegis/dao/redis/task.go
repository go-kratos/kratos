package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-common/app/admin/main/aegis/model/common"
	modtask "go-common/app/admin/main/aegis/model/task"
	"go-common/library/cache/redis"
	"go-common/library/log"

	"github.com/pkg/errors"
)

/*
 四个缓存
1. personalPool 存储每个用户的任务池 类型为List []taskid
2. publicPool 存储还未被领取的任务 类型为Sorted Set weight-taskid
3. delayPool 存储每个用户的延迟任务 类型为List
4. hashTask 存储所有任务的所有其他字段信息
*/

const (
	_hashexpire = 24 * 60 * 60
)

func personalKey(businessID, flowID int64, uid int64) string {
	return fmt.Sprintf("personal_%d_%d_%d", businessID, flowID, uid)
}

func publicKey(businessID, flowID int64) string {
	return fmt.Sprintf("{%d-%d}public_%d_%d", businessID, flowID, businessID, flowID)
}

func delayKey(businessID, flowID int64, uid int64) string {
	return fmt.Sprintf("delay_%d_%d_%d", businessID, flowID, uid)
}

func haskKey(taskid int64) string {
	return fmt.Sprintf("task_%d", taskid)
}

// formatID 在sorted set里面，id要扩展出来，否则排序不对
func formatID(taskid int64) string {
	return fmt.Sprintf("%.11d", taskid)
}

// CountPersonalTask 统计个数
func (d *Dao) CountPersonalTask(c context.Context, opt *common.BaseOptions) (count int64, err error) {
	conn := d.cluster.Get(c)
	defer conn.Close()
	key := personalKey(opt.BusinessID, opt.FlowID, opt.UID)
	if count, err = redis.Int64(conn.Do("LLEN", key)); err != nil {
		log.Error("conn.Do(LLEN,%s) error(%v)", key, err)
	}
	return
}

// RangePersonalTask 从本人的任务池取
func (d *Dao) RangePersonalTask(c context.Context, opt *modtask.ListOptions) (tasks map[int64]*modtask.Task, count int64, hitids, missids []int64, err error) {
	tasks, count, hitids, missids, err = d.rangefuncCluster(c, "personal", opt)
	return
}

// RangeDealyTask .
func (d *Dao) RangeDealyTask(c context.Context, opt *modtask.ListOptions) (tasks map[int64]*modtask.Task, count int64, hitids, missids []int64, err error) {
	return d.rangefuncCluster(c, "delay", opt)
}

// RangePublicTask .
func (d *Dao) RangePublicTask(c context.Context, opt *modtask.ListOptions) (tasks map[int64]*modtask.Task, count int64, hitids, missids []int64, err error) {
	return d.rangefuncCluster(c, "public", opt)
}

// PushPersonalTask 放入本人任务池
func (d *Dao) PushPersonalTask(c context.Context, opt *common.BaseOptions, ids ...interface{}) (err error) {
	key := personalKey(opt.BusinessID, opt.FlowID, opt.UID)
	return d.pushList(c, key, ids...)
}

// RemovePersonalTask 任务延迟或完成
func (d *Dao) RemovePersonalTask(c context.Context, opt *common.BaseOptions, ids ...interface{}) (err error) {
	key := personalKey(opt.BusinessID, opt.FlowID, opt.UID)
	return d.removeList(c, key, ids...)
}

// RemoveDelayTask 延迟任务完成
func (d *Dao) RemoveDelayTask(c context.Context, opt *common.BaseOptions, ids ...interface{}) (err error) {
	key := delayKey(opt.BusinessID, opt.FlowID, opt.UID)
	return d.removeList(c, key, ids...)
}

func (d *Dao) removeList(c context.Context, key string, ids ...interface{}) (err error) {
	conn := d.cluster.Get(c)
	defer conn.Close()
	for _, id := range ids {
		if err = conn.Send("LREM", key, 1, id); err != nil {
			log.Error("LREM error(%v)", errors.WithStack(err))
			continue
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("Flush error(%v)", errors.WithStack(err))
		return
	}
	for i := 0; i < len(ids); i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("Receive error(%v)", errors.WithStack(err))
			return
		}
	}
	return
}

// Release 清空个人任务 直接删除key
func (d *Dao) Release(c context.Context, opt *common.BaseOptions, delay bool) (err error) {
	var conn = d.cluster.Get(c)
	defer conn.Close()
	key := personalKey(opt.BusinessID, opt.FlowID, opt.UID)

	log.Info("Redis Release(%+v) delay(%v)", opt, delay)
	if delay {
		debug := func(msg string) {
			tasks, count, hitids, missids, err1 := d.rangefuncCluster(c, "personal", &modtask.ListOptions{
				BaseOptions: *opt,
				Pager: common.Pager{
					Pn: 1,
					Ps: 10,
				},
			})

			for id, task := range tasks {
				log.Info(msg+" Release task(%d)(%+v)", id, task)
			}
			log.Info(msg+" Release count(%d)", count)
			log.Info(msg+" Release hitids(%+v)", hitids)
			log.Info(msg+" Release missids(%+v)", missids)
			log.Info(msg+" Release err(%+v)", err1)
		}
		debug("Before")

		if _, err = conn.Do("LTRIM", key, 0, 0); err != nil {
			log.Error("LTRIM ReleasePersonalTask(%s), error(%v)", key, err)
		}

		debug("Middle")

		time.AfterFunc(5*time.Minute, func() {
			d.Release(context.Background(), opt, false)
		})

		debug("After")
	} else {
		log.Info("Redis DEL Release(%+v) delay(%v)", opt, delay)
		if _, err = conn.Do("DEL", key); err != nil {
			log.Error("DEL ReleasePersonalTask(%s), error(%v)", key, err)
		}
	}

	return
}

// PushDelayTask 延迟任务队列
func (d *Dao) PushDelayTask(c context.Context, opt *common.BaseOptions, ids ...interface{}) (err error) {
	key := delayKey(opt.BusinessID, opt.FlowID, opt.UID)
	return d.pushList(c, key, ids...)
}

func (d *Dao) pushList(c context.Context, key string, values ...interface{}) (err error) {
	var (
		conn = d.cluster.Get(c)
	)
	defer conn.Close()

	for _, id := range values {
		if _, err = conn.Do("LREM", key, 0, id); err != nil {
			log.Error("conn.Do(LREM, %v, %v, %v) error(%v)", key, 0, id, err)
			return
		}
	}

	args2 := []interface{}{key}
	args2 = append(args2, values...)
	if _, err = conn.Do("RPUSH", args2...); err != nil {
		log.Error("conn.Do(RPUSH, %v) error(%v)", args2, err)
		return
	}
	return
}

// RemovePublicTask 移除
func (d *Dao) RemovePublicTask(c context.Context, opt *common.BaseOptions, ids ...interface{}) (err error) {
	conn := d.cluster.Get(c)
	defer conn.Close()
	key := publicKey(opt.BusinessID, opt.FlowID)
	args := []interface{}{key}
	nids := []interface{}{}
	for _, id := range ids {
		nids = append(nids, fmt.Sprintf("%.11d", id))
	}

	args = append(args, nids...)
	if _, err = conn.Do("ZREM", args...); err != nil {
		log.Error("(ZREM,%v) error(%v)", args, errors.WithStack(err))
	}
	return err
}

// PopPublicTask 从实时任务池取出来
func (d *Dao) PopPublicTask(c context.Context, businessID, flowID, count int64) (ids []int64, err error) {
	conn := d.cluster.Get(c)
	defer conn.Close()
	key := publicKey(businessID, flowID)
	var tempids []int64
	if tempids, err = redis.Int64s(conn.Do("ZRANGE", key, 0, count)); err != nil {
		log.Error("conn.Do(ZADD,%s) error(%v)", key, errors.WithStack(err))
	}

	ids = append(ids, tempids...)
	return
}

// PushPublicTask 放入实时任务池
func (d *Dao) PushPublicTask(c context.Context, tasks ...*modtask.Task) (err error) {
	conn := d.cluster.Get(c)
	defer conn.Close()

	for _, task := range tasks {
		key := publicKey(task.BusinessID, task.FlowID)
		id := formatID(task.ID)
		fmt.Println("id:", id)
		if _, err = conn.Do("ZADD", key, -task.Weight, id); err != nil {
			log.Error("conn.Do(ZADD,%s) error(%v)", key, errors.WithStack(err))
		}
	}

	return
}

// SetTask .
func (d *Dao) SetTask(c context.Context, tasks ...*modtask.Task) (err error) {
	conn := d.cluster.Get(c)
	defer conn.Close()

	for _, task := range tasks {
		var bs []byte
		key := haskKey(task.ID)
		if bs, err = json.Marshal(task); err != nil {
			log.Error("json.Marshal(%+v) error(%v)", task, err)
			continue
		}

		if err = conn.Send("SET", key, bs); err != nil {
			log.Error("SET error(%v)", err)
			continue
		}
		if err = conn.Send("EXPIRE", key, _hashexpire); err != nil {
			log.Error("EXPIRE error(%v)", err)
			continue
		}
	}

	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}

	for i := 0; i < 2*len(tasks); i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive error(%v)", err)
			return
		}
	}
	return
}

// GetTask .
func (d *Dao) GetTask(c context.Context, ids []int64) (tasks []*modtask.Task, err error) {
	conn := d.cluster.Get(c)
	defer conn.Close()

	for _, id := range ids {
		key := haskKey(id)
		conn.Send("GET", key)
	}
	conn.Flush()
	var data []byte
	for i := 0; i < len(ids); i++ {
		if data, err = redis.Bytes(conn.Receive()); err != nil {
			log.Errorv(c, log.KV("event", "redis_get"), log.KV("key", haskKey(ids[i])), log.KV("error", err))
			return
		}
		task := new(modtask.Task)
		if err = json.Unmarshal(data, task); err != nil {
			log.Errorv(c, log.KV("event", "json.Unmarshal"), log.KV("task", string(data)), log.KV("error", err))
			return
		}
		tasks = append(tasks, task)
	}
	return
}
