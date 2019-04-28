package dao

import (
	"context"
	"encoding/json"
	"fmt"

	"go-common/app/interface/main/web/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_keyRkFmt           = "r_v2_%d_%d_%d_%d"
	_keyRkIndexFmt      = "ri_%d"
	_keyRkRegionFmt     = "rc_%d_%d_%d"
	_keyRkRecommendFmt  = "rr_%d"
	_keyRkTagFmt        = "rt_%d_%d"
	_keyRegionCustom    = "krc"
	_keyRegionCustomBak = _keyBakPrefix + _keyRegionCustom
	_keyBakPrefix       = "b_"
)

func keyRkList(rid int16, rankType, day, arcType int) string {
	return fmt.Sprintf(_keyRkFmt, rid, rankType, day, arcType)
}

func keyRkListBak(rid int16, rankType, day, arcType int) string {
	return _keyBakPrefix + keyRkList(rid, rankType, day, arcType)
}

func keyRkIndex(day int) string {
	return fmt.Sprintf(_keyRkIndexFmt, day)
}

func keyRkIndexBak(day int) string {
	return _keyBakPrefix + keyRkIndex(day)
}

func keyRkRegionList(rid int16, day, original int) string {
	return fmt.Sprintf(_keyRkRegionFmt, rid, day, original)
}

func keyRkRegionListBak(rid int16, day, original int) string {
	return _keyBakPrefix + keyRkRegionList(rid, day, original)
}

func keyRkRecommendList(rid int16) string {
	return fmt.Sprintf(_keyRkRecommendFmt, rid)
}

func keyRkRecommendListBak(rid int16) string {
	return _keyBakPrefix + fmt.Sprintf(_keyRkRecommendFmt, rid)
}

func keyRkTagList(rid int16, tagID int64) string {
	return fmt.Sprintf(_keyRkTagFmt, rid, tagID)
}

func keyRkTagListBak(rid int16, tagID int64) string {
	return _keyBakPrefix + keyRkTagList(rid, tagID)
}

// RankingCache get rank list from cache.
func (d *Dao) RankingCache(c context.Context, rid int16, rankType, day, arcType int) (data *model.RankData, err error) {
	key := keyRkList(rid, rankType, day, arcType)
	conn := d.redis.Get(c)
	defer conn.Close()
	data, err = d.rankingCache(conn, key)
	return
}

// RankingBakCache get rank list from bak cache.
func (d *Dao) RankingBakCache(c context.Context, rid int16, rankType, day, arcType int) (data *model.RankData, err error) {
	d.cacheProm.Incr("ranking_remote_cache")
	key := keyRkListBak(rid, rankType, day, arcType)
	conn := d.redisBak.Get(c)
	defer conn.Close()
	data, err = d.rankingCache(conn, key)
	if data == nil || len(data.List) == 0 {
		log.Error("RankingBakCache(%s) is nil", key)
	}
	return
}

// RankingIndexCache get rank index from cache.
func (d *Dao) RankingIndexCache(c context.Context, day int) (arcs []*model.IndexArchive, err error) {
	key := keyRkIndex(day)
	conn := d.redis.Get(c)
	defer conn.Close()
	arcs, err = d.rankingIndexCache(conn, key)
	return
}

// RankingIndexBakCache get rank index from bak cache.
func (d *Dao) RankingIndexBakCache(c context.Context, day int) (arcs []*model.IndexArchive, err error) {
	d.cacheProm.Incr("ranking_index_remote_cache")
	key := keyRkIndexBak(day)
	conn := d.redisBak.Get(c)
	defer conn.Close()
	arcs, err = d.rankingIndexCache(conn, key)
	if len(arcs) == 0 {
		log.Error("RankingIndexBakCache(%s) is nil", key)
	}
	return
}

// RankingRegionCache get rank cate list from cache.
func (d *Dao) RankingRegionCache(c context.Context, rid int16, day, original int) (arcs []*model.RegionArchive, err error) {
	key := keyRkRegionList(rid, day, original)
	conn := d.redis.Get(c)
	defer conn.Close()
	arcs, err = d.rankingRegionCache(conn, key)
	return
}

// RankingRegionBakCache get rank cate list from bak cache.
func (d *Dao) RankingRegionBakCache(c context.Context, rid int16, day, original int) (arcs []*model.RegionArchive, err error) {
	d.cacheProm.Incr("ranking_region_remote_cache")
	key := keyRkRegionListBak(rid, day, original)
	conn := d.redisBak.Get(c)
	defer conn.Close()
	arcs, err = d.rankingRegionCache(conn, key)
	if len(arcs) == 0 {
		log.Error("RankingRegionBakCache(%s) is nil", key)
	}
	return
}

// RankingRecommendCache get rank recommend list from cache.
func (d *Dao) RankingRecommendCache(c context.Context, rid int16) (arcs []*model.IndexArchive, err error) {
	key := keyRkRecommendList(rid)
	conn := d.redis.Get(c)
	defer conn.Close()
	arcs, err = d.rankingIndexCache(conn, key)
	return
}

// RankingRecommendBakCache get rank recommend list from bak cache.
func (d *Dao) RankingRecommendBakCache(c context.Context, rid int16) (arcs []*model.IndexArchive, err error) {
	d.cacheProm.Incr("ranking_rec_remote_cache")
	key := keyRkRecommendListBak(rid)
	conn := d.redisBak.Get(c)
	defer conn.Close()
	arcs, err = d.rankingIndexCache(conn, key)
	if len(arcs) == 0 {
		log.Error("RankingRecommendBakCache(%s) is nil", key)
	}
	return
}

// RankingTagCache get ranking tag from cache.
func (d *Dao) RankingTagCache(c context.Context, rid int16, tagID int64) (arcs []*model.TagArchive, err error) {
	key := keyRkTagList(rid, tagID)
	conn := d.redis.Get(c)
	defer conn.Close()
	arcs, err = d.rankingTagCache(conn, key)
	return
}

// RankingTagBakCache get ranking tag from bak cache.
func (d *Dao) RankingTagBakCache(c context.Context, rid int16, tagID int64) (arcs []*model.TagArchive, err error) {
	d.cacheProm.Incr("ranking_tag_remote_cache")
	key := keyRkTagListBak(rid, tagID)
	conn := d.redisBak.Get(c)
	defer conn.Close()
	arcs, err = d.rankingTagCache(conn, key)
	if len(arcs) == 0 {
		log.Error("RankingTagBakCache(%s) is nil", key)
	}
	return
}

// RegionCustomCache get region custom data from cache
func (d *Dao) RegionCustomCache(c context.Context) (res []*model.Custom, err error) {
	key := _keyRegionCustom
	conn := d.redis.Get(c)
	defer conn.Close()
	res, err = regionCustomCache(conn, key)
	return
}

// RegionCustomBakCache get region custom data from cache
func (d *Dao) RegionCustomBakCache(c context.Context) (res []*model.Custom, err error) {
	key := _keyRegionCustomBak
	conn := d.redis.Get(c)
	defer conn.Close()
	res, err = regionCustomCache(conn, key)
	return
}

func regionCustomCache(conn redis.Conn, key string) (res []*model.Custom, err error) {
	var value []byte
	if value, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(GET, %s) error(%v)", key, err)
		}
		return
	}
	res = []*model.Custom{}
	if err = json.Unmarshal(value, &res); err != nil {
		log.Error("json.Unmarshal(%v) error(%v)", value, err)
	}
	return
}

// SetRegionCustomCache set region custom data cache
func (d *Dao) SetRegionCustomCache(c context.Context, data []*model.Custom) (err error) {
	key := _keyRegionCustom
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = d.setRegionCustomCache(conn, key, d.redisRcExpire, data); err != nil {
		return
	}
	key = _keyRegionCustomBak
	connBak := d.redisBak.Get(c)
	err = d.setRegionCustomCache(connBak, key, d.redisRcBakExpire, data)
	connBak.Close()
	return
}

func (d *Dao) setRegionCustomCache(conn redis.Conn, key string, expire int32, data []*model.Custom) (err error) {
	var bs []byte
	if bs, err = json.Marshal(data); err != nil {
		log.Error("json.Marshal(%v) error (%v)", data, err)
		return
	}
	if err = conn.Send("SET", key, bs); err != nil {
		log.Error("conn.Send(SET, %s, %s) error(%v)", key, string(bs), err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.redisRkExpire); err != nil {
		log.Error("conn.Send(Expire, %s, %d) error(%v)", key, d.redisRkExpire, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

func (d *Dao) rankingCache(conn redis.Conn, key string) (arcs *model.RankData, err error) {
	var value []byte
	if value, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(GET, %s) error(%v)", key, err)
		}
		return
	}
	arcs = new(model.RankData)
	if err = json.Unmarshal(value, &arcs); err != nil {
		log.Error("json.Unmarshal(%v) error(%v)", value, err)
	}
	return
}

func (d *Dao) rankingIndexCache(conn redis.Conn, key string) (arcs []*model.IndexArchive, err error) {
	var value []byte
	if value, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(GET, %s) error(%v)", key, err)
		}
		return
	}
	arcs = []*model.IndexArchive{}
	if err = json.Unmarshal(value, &arcs); err != nil {
		log.Error("json.Unmarshal(%v) error(%v)", value, err)
	}
	return
}

func (d *Dao) rankingRegionCache(conn redis.Conn, key string) (arcs []*model.RegionArchive, err error) {
	var value []byte
	if value, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(GET, %s) error(%v)", key, err)
		}
		return
	}
	arcs = []*model.RegionArchive{}
	if err = json.Unmarshal(value, &arcs); err != nil {
		log.Error("json.Unmarshal(%v) error(%v)", value, err)
	}
	return
}

func (d *Dao) rankingTagCache(conn redis.Conn, key string) (arcs []*model.TagArchive, err error) {
	var value []byte
	if value, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(GET, %s) error(%v)", key, err)
		}
		return
	}
	arcs = []*model.TagArchive{}
	if err = json.Unmarshal(value, &arcs); err != nil {
		log.Error("json.Unmarshal(%v) error(%v)", value, err)
	}
	return
}

// SetRankingCache set ranking data to cache
func (d *Dao) SetRankingCache(c context.Context, rid int16, rankType, day, arcType int, data *model.RankData) (err error) {
	key := keyRkList(rid, rankType, day, arcType)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = d.setRkCache(c, conn, key, d.redisRkExpire, data); err != nil {
		return
	}
	key = keyRkListBak(rid, rankType, day, arcType)
	connBak := d.redisBak.Get(c)
	err = d.setRkCache(c, connBak, key, d.redisRkBakExpire, data)
	connBak.Close()
	return
}

// SetRankingIndexCache set ranking index data to cache
func (d *Dao) SetRankingIndexCache(c context.Context, day int, arcs []*model.IndexArchive) (err error) {
	key := keyRkIndex(day)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = d.setRkIndexCache(c, conn, key, d.redisRkExpire, arcs); err != nil {
		return
	}
	key = keyRkIndexBak(day)
	connBak := d.redisBak.Get(c)
	err = d.setRkIndexCache(c, connBak, key, d.redisRkBakExpire, arcs)
	connBak.Close()
	return
}

// SetRankingRegionCache set ranking data to cache
func (d *Dao) SetRankingRegionCache(c context.Context, rid int16, day, original int, arcs []*model.RegionArchive) (err error) {
	key := keyRkRegionList(rid, day, original)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = d.setRkRegionCache(c, conn, key, d.redisRkExpire, arcs); err != nil {
		return
	}
	key = keyRkRegionListBak(rid, day, original)
	connBak := d.redisBak.Get(c)
	err = d.setRkRegionCache(c, connBak, key, d.redisRkBakExpire, arcs)
	connBak.Close()
	return
}

// SetRankingRecommendCache set ranking data to bak cache
func (d *Dao) SetRankingRecommendCache(c context.Context, rid int16, arcs []*model.IndexArchive) (err error) {
	key := keyRkRecommendList(rid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = d.setRkIndexCache(c, conn, key, d.redisRkExpire, arcs); err != nil {
		return
	}
	key = keyRkRecommendListBak(rid)
	connBak := d.redisBak.Get(c)
	err = d.setRkIndexCache(c, connBak, key, d.redisRkBakExpire, arcs)
	connBak.Close()
	return
}

// SetRankingTagCache set ranking tag data to cache
func (d *Dao) SetRankingTagCache(c context.Context, rid int16, tagID int64, arcs []*model.TagArchive) (err error) {
	key := keyRkTagList(rid, tagID)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = d.setRkTagCache(c, conn, key, d.redisRkExpire, arcs); err != nil {
		return
	}
	key = keyRkTagListBak(rid, tagID)
	connBak := d.redisBak.Get(c)
	err = d.setRkTagCache(c, connBak, key, d.redisRkBakExpire, arcs)
	connBak.Close()
	return
}

func (d *Dao) setRkCache(c context.Context, conn redis.Conn, key string, expire int32, arcs *model.RankData) (err error) {
	var bs []byte
	if bs, err = json.Marshal(arcs); err != nil {
		log.Error("json.Marshal(%v) error (%v)", arcs, err)
		return
	}
	if err = conn.Send("SET", key, bs); err != nil {
		log.Error("conn.Send(SET, %s, %s) error(%v)", key, string(bs), err)
		return
	}
	if err = conn.Send("EXPIRE", key, expire); err != nil {
		log.Error("conn.Send(Expire, %s, %d) error(%v)", key, expire, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

func (d *Dao) setRkIndexCache(c context.Context, conn redis.Conn, key string, expire int32, arcs []*model.IndexArchive) (err error) {
	var bs []byte
	if bs, err = json.Marshal(arcs); err != nil {
		log.Error("json.Marshal(%v) error (%v)", arcs, err)
		return
	}
	if err = conn.Send("SET", key, bs); err != nil {
		log.Error("conn.Send(SET, %s, %s) error(%v)", key, string(bs), err)
		return
	}
	if err = conn.Send("EXPIRE", key, expire); err != nil {
		log.Error("conn.Send(Expire, %s, %d) error(%v)", key, expire, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

func (d *Dao) setRkRegionCache(c context.Context, conn redis.Conn, key string, expire int32, arcs []*model.RegionArchive) (err error) {
	var bs []byte
	if bs, err = json.Marshal(arcs); err != nil {
		log.Error("json.Marshal(%v) error (%v)", arcs, err)
		return
	}
	if err = conn.Send("SET", key, bs); err != nil {
		log.Error("conn.Send(SET, %s, %s) error(%v)", key, string(bs), err)
		return
	}
	if err = conn.Send("EXPIRE", key, expire); err != nil {
		log.Error("conn.Send(Expire, %s, %d) error(%v)", key, expire, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

func (d *Dao) setRkTagCache(c context.Context, conn redis.Conn, key string, expire int32, arcs []*model.TagArchive) (err error) {
	var bs []byte
	if bs, err = json.Marshal(arcs); err != nil {
		log.Error("json.Marshal(%v) error (%v)", arcs, err)
		return
	}
	if err = conn.Send("SET", key, bs); err != nil {
		log.Error("conn.Send(SET, %s, %s) error(%v)", key, string(bs), err)
		return
	}
	if err = conn.Send("EXPIRE", key, expire); err != nil {
		log.Error("conn.Send(Expire, %s, %d) error(%v)", key, expire, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}
