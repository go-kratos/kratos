package video

import (
	"context"
	"encoding/json"
	"fmt"

	"go-common/app/service/main/favorite/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_covers = "fcs_"
)

func coversKey(mid, fid int64) string {
	return fmt.Sprintf("%s%d_%d", _covers, mid, fid)
}

// pingRedis check redis connection
func (d *Dao) pingRedis(c context.Context) (err error) {
	conn := d.redisPool.Get(c)
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}

// SetNewCoverCache set fav's cover to cache
func (d *Dao) SetNewCoverCache(c context.Context, mid, fid int64, covers []*model.Cover) (err error) {
	key := coversKey(mid, fid)
	conn := d.redisPool.Get(c)
	defer conn.Close()
	for _, cover := range covers {
		var bs []byte
		if bs, err = json.Marshal(cover); err != nil {
			log.Error("json.Marshal(%v) err(%v)", cover, err)
			return
		}
		if err = conn.Send("RPUSH", key, bs); err != nil {
			log.Error("conn.Send RPUSH error(%v)", err)
			return
		}
	}
	if err = conn.Send("EXPIRE", key, d.coverExpireRedis); err != nil {
		log.Error("conn.Send(EXPIRE) err(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush err(%v)", err)
		return
	}
	for i := 0; i < len(covers)+1; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive error(%v)", err)
		}
	}
	return
}

// NewCoversCache get multi cover of fids by pipeline
func (d *Dao) NewCoversCache(c context.Context, mid int64, fids []int64) (fcvs map[int64][]*model.Cover, mis []int64, err error) {
	conn := d.redisPool.Get(c)
	defer conn.Close()
	for _, fid := range fids {
		key := coversKey(mid, fid)
		if err = conn.Send("LRANGE", key, 0, 2); err != nil {
			log.Error("conn.Send(LRANGE) err(%v)", err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() err(%v)", err)
		return
	}
	fcvs = make(map[int64][]*model.Cover, len(fids))
	// receive lrange
	for i := 0; i < len(fids); i++ {
		var (
			bbs [][]byte
			cvs []*model.Cover
		)
		if bbs, err = redis.ByteSlices(conn.Receive()); err != nil {
			if err == redis.ErrNil {
				err = nil
				mis = append(mis, fids[i])
				continue
			}
			log.Error("redis.ByteSlices err(%v)", err)
			return
		}
		if len(bbs) == 0 {
			mis = append(mis, fids[i])
			continue
		}
		for _, bs := range bbs {
			cv := &model.Cover{}
			if err = json.Unmarshal(bs, cv); err != nil {
				log.Error("json.Unmarshal err(%v)", err)
				return
			}
			cvs = append(cvs, cv)
		}
		fcvs[fids[i]] = cvs
	}
	return
}

// DelCoverCache delete folder cover
func (d *Dao) DelCoverCache(c context.Context, mid, fid int64) (err error) {
	var (
		key  = coversKey(mid, fid)
		conn = d.redisPool.Get(c)
	)
	defer conn.Close()
	if _, err = conn.Do("DEL", key); err != nil {
		log.Error("conn.Do(DEL, %s) error(%v)", key, err)
	}
	return
}
