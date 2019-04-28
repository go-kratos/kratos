package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"

	"go-common/app/job/main/aegis/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_hashexpire = 24 * 60 * 60
)

func personalKey(businessID, flowID int64, uid int64) string {
	return fmt.Sprintf("personal_%d_%d_%d", businessID, flowID, uid)
}

func publicKey(businessID, flowID int64) string {
	return fmt.Sprintf("{%d-%d}public_%d_%d", businessID, flowID, businessID, flowID)
}

func publicBackKey(businessID, flowID int64) string {
	return fmt.Sprintf("{%d-%d}publicBackup_%d_%d", businessID, flowID, businessID, flowID)
}

func delayKey(businessID, flowID int64, uid int64) string {
	return fmt.Sprintf("delay_%d_%d_%d", businessID, flowID, uid)
}

func haskKey(taskid int64) string {
	return fmt.Sprintf("task_%d", taskid)
}

func zsetKey(taskid int64) string {
	return fmt.Sprintf("%.11d", taskid)
}

// SetTask .
func (d *Dao) SetTask(c context.Context, task *model.Task) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()

	var bs []byte
	key := haskKey(task.ID)
	if bs, err = json.Marshal(task); err != nil {
		log.Error("json.Marshal(%+v) error(%v)", task, err)
		return
	}

	if err = conn.Send("SET", key, bs); err != nil {
		log.Error("HSET error(%v)", err)
		return
	}
	if err = conn.Send("EXPIRE", key, _hashexpire); err != nil {
		log.Error("EXPIRE error(%v)", err)
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

// GetTask .
func (d *Dao) GetTask(c context.Context, id int64) (task *model.Task, err error) {
	var bs []byte
	conn := d.redis.Get(c)
	defer conn.Close()

	key := haskKey(id)

	if bs, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		log.Error("conn.Get(%s) error(%v)", key, err)
		return
	}

	task = new(model.Task)
	if err = json.Unmarshal(bs, task); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(bs), err)
		return
	}

	return
}

// RemovePersonalTask 任务延迟或完成
func (d *Dao) RemovePersonalTask(c context.Context, businessID, flowID int64, uid, taskid int64) (err error) {
	key := personalKey(businessID, flowID, uid)
	return d.removeList(c, key, taskid)
}

// RemoveDelayTask 延迟任务完成
func (d *Dao) RemoveDelayTask(c context.Context, businessID, flowID int64, uid, taskid int64) (err error) {
	key := delayKey(businessID, flowID, uid)
	return d.removeList(c, key, taskid)
}

func (d *Dao) removeList(c context.Context, key string, id int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()

	if _, err = conn.Do("LREM", key, 0, id); err != nil {
		log.Error("LREM error(%v)", errors.WithStack(err))
	}
	return
}

// PushPersonalTask 放入本人任务池
func (d *Dao) PushPersonalTask(c context.Context, businessID, flowID int64, uid, taskid int64) (err error) {
	key := personalKey(businessID, flowID, uid)
	return d.pushList(c, key, taskid)
}

// PushDelayTask 延迟任务队列
func (d *Dao) PushDelayTask(c context.Context, businessID, flowID int64, uid, taskid int64) (err error) {
	key := delayKey(businessID, flowID, uid)
	return d.pushList(c, key, taskid)
}

func (d *Dao) pushList(c context.Context, key string, values ...interface{}) (err error) {
	var (
		conn = d.redis.Get(c)
	)
	defer conn.Close()

	args1 := []interface{}{key, 0}
	args1 = append(args1, values...)
	if _, err = conn.Do("LREM", args1...); err != nil {
		log.Error("conn.Do(RPUSH, %v) error(%v)", args1, err)
		return
	}
	args2 := []interface{}{key}
	args2 = append(args2, values...)
	if _, err = conn.Do("RPUSH", args2...); err != nil {
		log.Error("conn.Do(RPUSH, %v) error(%v)", args2, err)
	}
	return
}

// RemovePublicTask 移除
func (d *Dao) RemovePublicTask(c context.Context, businessID, flowID int64, taskid int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := publicKey(businessID, flowID)
	args := []interface{}{key}
	args = append(args, zsetKey(taskid))
	if _, err = conn.Do("ZREM", args...); err != nil {
		log.Error("(ZREM,%v) error(%v)", args, errors.WithStack(err))
	}
	return err
}

// PushPublicTask 放入实时任务池
func (d *Dao) PushPublicTask(c context.Context, tasks ...*model.Task) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()

	for _, task := range tasks {
		key := publicKey(task.BusinessID, task.FlowID)
		if _, err = conn.Do("ZADD", key, -task.Weight, zsetKey(task.ID)); err != nil {
			log.Error("conn.Do(ZADD,%s) error(%v)", key, errors.WithStack(err))
		}
	}

	return
}

// SetWeight set weight
func (d *Dao) SetWeight(c context.Context, businessID, flowID int64, id, weight int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()

	var (
		ow  int64
		key = publicKey(businessID, flowID)
		zid = zsetKey(id)
	)

	if ow, err = redis.Int64(conn.Do("ZSCORE", key, zid)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("redis (ZSCORE,%s,%s) error(%v)", key, zid, err)
		}
		return
	}

	// 为了从大到小排序，weight取负值
	nw := -(weight + ow)
	if _, err = conn.Do("ZINCRBY", key, nw, zid); err != nil {
		log.Error("redis (ZINCRBY,%s,%s,%d) error(%v)", key, nw, zid, err)
		return
	}

	return
}

// GetWeight get Weight
func (d *Dao) GetWeight(c context.Context, businessID, flowID int64, id int64) (weight int64, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()

	key := publicKey(businessID, flowID)

	weight, err = redis.Int64(conn.Do("ZSCORE", key, zsetKey(id)))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(ZSCORE %s %s) error(%v)", key, zsetKey(id), errors.WithStack(err))
		}
	}
	weight = -weight
	return
}

// TopWeights .
func (d *Dao) TopWeights(c context.Context, businessID, flowID int64, toplen int64) (wis []*model.WeightItem, err error) {
	key := publicKey(businessID, flowID)
	return d.zrange(c, key, 0, toplen)
}

// CreateUnionSet  创建分身集合
func (d *Dao) CreateUnionSet(c context.Context, businessID, flowID int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()

	key := publicKey(businessID, flowID)
	backKey := publicBackKey(businessID, flowID)

	if _, err = conn.Do("ZUNIONSTORE", backKey, 1, key); err != nil {
		log.Error("conn.Do(ZUNIONSTORE,%s) error(%v)", key, errors.WithStack(err))
	}

	return
}

// RangeUinonSet  批次取出
func (d *Dao) RangeUinonSet(c context.Context, businessID, flowID int64, start, stop int64) (wis []*model.WeightItem, err error) {
	key := publicBackKey(businessID, flowID)
	return d.zrange(c, key, start, stop)
}

// DeleteUinonSet  清空分身
func (d *Dao) DeleteUinonSet(c context.Context, businessID, flowID int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()

	backKey := publicBackKey(businessID, flowID)

	if _, err = conn.Do("DEL", backKey); err != nil {
		log.Error("conn.Do(ZUNIODELNSTORE,%s) error(%v)", backKey, errors.WithStack(err))
	}

	return
}

func (d *Dao) zrange(c context.Context, key string, start, stop int64) (wis []*model.WeightItem, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()

	reply, err := redis.Int64s(conn.Do("ZRANGE", key, start, stop, "WITHSCORES"))
	if err != nil {
		log.Error("conn.Do(ZADD,%s) error(%v)", key, errors.WithStack(err))
		return
	}
	// 单数是id,双数是weight
	for i := 0; i < len(reply); i += 2 {
		wi := &model.WeightItem{}
		wi.ID = reply[i]
		wi.Weight = -reply[i+1]
		wis = append(wis, wi)
	}

	return
}
