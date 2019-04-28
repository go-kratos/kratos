package dao

import (
	"context"
	"strconv"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_prefixRecent = "dm_rct_"
)

func keyRecent(mid int64) string {
	return _prefixRecent + strconv.FormatInt(mid, 10)
}

// RecentDM get recent dm of up.
func (d *Dao) RecentDM(c context.Context, mid, start, end int64) (dms []*model.DM, oids []int64, total int64, err error) {
	var (
		conn   = d.dmRctRds.Get(c)
		key    = keyRecent(mid)
		oidMap = make(map[int64]struct{})
	)
	defer conn.Close()
	if err = conn.Send("ZREVRANGE", key, start, end); err != nil {
		log.Error("conn.Send(ZREVRANGE %s) error(%s)", key, err)
		return
	}
	if err = conn.Send("ZCARD", key); err != nil {
		log.Error("conn.Send(ZCARD %s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	values, err := redis.ByteSlices(conn.Receive())
	if err != nil {
		log.Error("conn.Receive(%s) error(%v)", key, err)
		return
	}
	if total, err = redis.Int64(conn.Receive()); err != nil {
		log.Error("conn.Receive(%s) error(%v)", key, err)
		return
	}
	for _, value := range values {
		dm := &model.DM{}
		if err = dm.Unmarshal(value); err != nil {
			log.Error("dm.Unmarshal(%s) error(%v)", value, err)
			return
		}
		dms = append(dms, dm)
		if _, ok := oidMap[dm.Oid]; !ok {
			oidMap[dm.Oid] = struct{}{}
			oids = append(oids, dm.Oid)
		}
	}
	return
}

// TrimUpRecent zrange remove recent dm of up.
func (d *Dao) TrimUpRecent(c context.Context, mid, count int64) (err error) {
	var (
		conn = d.dmRctRds.Get(c)
		key  = keyRecent(mid)
	)
	defer conn.Close()
	if _, err = conn.Do("ZREMRANGEBYRANK", key, 0, count-1); err != nil {
		log.Error("conn.Do(ZREMRANGEBYRANK %s) error(%v)", key, err)
	}
	return
}
