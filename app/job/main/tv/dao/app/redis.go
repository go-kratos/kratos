package app

import (
	"context"
	"fmt"
	"time"

	commonMdl "go-common/app/job/main/tv/model/common"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

// keyZone gets the key of the zone in Redis
func keyZone(category int) string {
	return fmt.Sprintf("zone_idx_%d", category)
}

// Flush it flushes the list of one zone
func (d *Dao) Flush(c context.Context, category int, idxRanks []*commonMdl.IdxRank) (err error) {
	var (
		ctime  int64
		length int64
		conn   = d.redis.Get(c)
		key    = keyZone(category)
	)
	// remove the previous list
	if err = conn.Send("EXPIRE", key, 0); err != nil {
		log.Error("conn.Send(EXPIRE %s) error(%v)", key, err)
		return
	}
	// add new ids inside
	for _, v := range idxRanks {
		ctime = int64(v.Ctime)
		if err = conn.Send("ZADD", key, ctime, v.ID); err != nil {
			log.Error("conn.Send(ZADD %s %v) error(%v)", key, v.ID, ctime)
			return
		}
	}
	// set expiration
	if err = conn.Send("EXPIRE", key, d.redisExpire); err != nil {
		log.Error("conn.Send(EXPIRE %s) error(%v)", key, err)
		return
	}
	// check result
	if err = conn.Send("ZCARD", key); err != nil {
		log.Error("conn.Send(ZCARD %s) error(%v)", key, err)
		return
	}
	// flush result
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	for i := 0; i < len(idxRanks)+2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	if length, err = redis.Int64(conn.Receive()); err != nil {
		log.Error("conn.Receive() error(%v)", err)
	}
	log.Info("Flush Succ! Zone %d, Items: %d", category, length)
	conn.Close()
	return
}

// TimeTrans transform the time format to Unix Timestamp
func TimeTrans(stimeStr string) (stime int64, err error) {
	local, _ := time.LoadLocation("Local")
	var (
		timeValue time.Time
	)
	timeValue, err = time.ParseInLocation("2006-01-02 15:04:05", stimeStr, local)
	if err != nil {
		log.Warn("TimeTrans %s, Error %v", stimeStr, err)
		return
	}
	if stime = timeValue.Unix(); stime < 1 {
		err = fmt.Errorf("time %s transform %d error", stimeStr, stime)
		return
	}
	return
}

// ZAddIdx adds one valid season into the zone list
func (d *Dao) ZAddIdx(c context.Context, category int, ctimeStr string, id int64) (err error) {
	var (
		conn  = d.redis.Get(c)
		key   = keyZone(category)
		ctime int64
	)
	defer conn.Close()
	if ctime, err = TimeTrans(ctimeStr); err != nil {
		log.Warn("ZAddIdx Ctime %s Error %v", ctimeStr, err)
		return
	}
	if err = conn.Send("ZADD", key, ctime, id); err != nil {
		log.Error("conn.Send(ZADD %s - %v) error(%v)", key, id, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		log.Error("conn.Receive() error(%v)", err)
		return
	}
	return
}

// ZRemIdx ZREM trim from trim queue.
func (d *Dao) ZRemIdx(c context.Context, category int, id int64) (err error) {
	var (
		conn = d.redis.Get(c)
		key  = keyZone(category)
	)
	if _, err = conn.Do("ZREM", key, id); err != nil {
		log.Error("conn.Send(ZADD %s - %v) error(%v)", key, id, err)
	}
	conn.Close()
	return
}
