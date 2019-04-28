package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	tmod "go-common/app/job/main/videoup-report/model/task"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_twexpire = 24 * 60 * 60 // 1 day
)

func key(id int64) string {
	return fmt.Sprintf("tw_%d", id)
}

//SetWeight 设置权重配置
func (d *Dao) SetWeight(c context.Context, mcases map[int64]*tmod.WeightParams) (err error) {
	conn := d.secondary.Get(c)
	defer conn.Close()

	log.Info("SetWeight before len(%d) time(%v)", len(mcases), time.Now())
	for tid, mcase := range mcases {
		var bs []byte
		key := key(tid)
		if bs, err = json.Marshal(mcase); err != nil {
			log.Error("json.Marshal(%+v) error(%v)", mcase, err)
			continue
		}

		if err = conn.Send("SET", key, bs); err != nil {
			log.Error("SET error(%v)", err)
			continue
		}
		if err = conn.Send("EXPIRE", key, _twexpire); err != nil {
			log.Error("EXPIRE error(%v)", err)
			continue
		}
	}

	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2*len(mcases); i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	log.Info("SetWeight end len(%d) time(%v)", len(mcases), time.Now())
	return
}

//GetWeight 获取实时任务的权重配置
func (d *Dao) GetWeight(c context.Context, ids []int64) (mcases map[int64]*tmod.WeightParams, err error) {
	conn := d.secondary.Get(c)
	defer conn.Close()
	mcases = make(map[int64]*tmod.WeightParams)
	for _, id := range ids {
		var bs []byte
		key := key(int64(id))
		if bs, err = redis.Bytes(conn.Do("GET", key)); err != nil {
			if err == redis.ErrNil {
				err = nil
			} else {
				log.Error("conn.Do(GET, %v) error(%v)", key, err)
			}
			continue
		}
		p := &tmod.WeightParams{}
		if err = json.Unmarshal(bs, p); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", string(bs), err)
			err = nil
			continue
		}
		mcases[int64(id)] = p
	}
	return
}
