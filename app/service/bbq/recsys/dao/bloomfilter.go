package dao

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go-common/app/service/bbq/recsys/dao/parallel"
	"go-common/library/log"

	"go-common/library/cache/redis"

	"github.com/Dai0522/go-hash/bloomfilter"
	"github.com/Dai0522/workpool"
)

const (
	_baseBfKey = "BBQ:BF:V1:%s:%s"
)

func userBFRedisKey(k string) string {
	d := time.Now().Format("20060102")
	return fmt.Sprintf(_baseBfKey, k, d)
}

func (d *Dao) loadBF(c context.Context, mid int64, buvid string) (bf *bloomfilter.BloomFilter, err error) {
	var tasks []workpool.Task
	if buvid != "" {
		buvidK := userBFRedisKey(buvid)
		t := parallel.NewRedisTask(&c, d.bfRedis, "GET", buvidK)
		tasks = append(tasks, t)
	}

	if mid != 0 {
		midK := userBFRedisKey(strconv.FormatInt(mid, 10))
		t := parallel.NewRedisTask(&c, d.bfRedis, "GET", midK)
		tasks = append(tasks, t)
	}
	ftTasks := d.parallelTask(tasks)

	for _, ft := range *ftTasks {
		raw, e := ft.Wait(100 * time.Millisecond)
		if e != nil && e != redis.ErrNil {
			log.Errorv(c, log.KV("BF_GET_ERROR", e), log.KV("TASK", ft.T.(*parallel.RedisTask)))
			continue
		}
		if raw == nil || len(*raw) == 0 {
			continue
		}
		tmp, e := bloomfilter.Load(raw)
		if e != nil || tmp == nil {
			log.Errorv(c, log.KV("BF_LOAD_ERROR", e), log.KV("TASK", ft.T.(*parallel.RedisTask)), log.KV("raw", *raw))
			continue
		}
		bf = bloomfilter.Merge(bf, tmp)
	}
	if bf == nil {
		bf, err = bloomfilter.New(1000, 0.0001)
	}

	return
}

// WriteBF .
func (d *Dao) WriteBF(c context.Context, mid int64, buvid string, svid []uint64) (bool, error) {
	if mid == int64(0) && buvid == "" {
		return false, errors.New("mid && buvid can't be empty")
	}
	// load bf from redis
	bf, err := d.loadBF(c, mid, buvid)
	if err != nil {
		return false, err
	}

	// put svid
	for _, v := range svid {
		bf.PutUint64(v)
	}

	// store bf into redis
	var tasks []workpool.Task
	b := bf.Serialized()
	if buvid != "" {
		buvidK := userBFRedisKey(buvid)
		t := parallel.NewRedisTask(&c, d.bfRedis, "SETEX", buvidK, 86400, *b)
		tasks = append(tasks, t)
	}
	if mid != int64(0) {
		midK := userBFRedisKey(strconv.FormatInt(mid, 10))
		t := parallel.NewRedisTask(&c, d.bfRedis, "SETEX", midK, 86400, *b)
		tasks = append(tasks, t)
	}

	ftTasks := d.parallelTask(tasks)
	for _, ft := range *ftTasks {
		_, err = ft.Wait(100 * time.Millisecond)
		if err != nil {
			log.Errorv(c, log.KV("BF_SET_ERROR", err))
		}
	}

	return true, err
}
