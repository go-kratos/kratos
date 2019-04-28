package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"go-common/app/interface/openplatform/article/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_prefixUpper    = "art_u_%d"           // upper's article list
	_prefixSorted   = "art_sort_%d_%d"     // sorted aids sort_category_field
	_prefixRank     = "art_ranks_%d"       // ranks by cid
	_prefixMaxLike  = "art_mlt_%d"         // like message number
	_readPingSet    = "art:readping"       // reading start set
	_prefixReadPing = "art:readping:%s:%d" // reading during on some device for some article
	_blank          = int64(-1)
)

func upperKey(mid int64) string {
	return fmt.Sprintf(_prefixUpper, mid)
}

func sortedKey(categoryID int64, field int) string {
	return fmt.Sprintf(_prefixSorted, categoryID, field)
}

func rankKey(cid int64) string {
	return fmt.Sprintf(_prefixRank, cid)
}

func hotspotKey(typ int8, id int64) string {
	return fmt.Sprintf("art_hotspot%d_%d", typ, id)
}

func authorCategoriesKey(mid int64) string {
	return fmt.Sprintf("author:categories:%d", mid)
}

func recommendsAuthorsKey(category int64) string {
	return fmt.Sprintf("recommends:authors:%d", category)
}

func readPingSetKey() string {
	return _readPingSet
}

func readPingKey(buvid string, aid int64) string {
	return fmt.Sprintf(_prefixReadPing, buvid, aid)
}

// pingRedis ping redis.
func (d *Dao) pingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	if _, err = conn.Do("SET", "PING", "PONG"); err != nil {
		PromError("redis: ping remote")
		log.Error("remote redis: conn.Do(SET,PING,PONG) error(%+v)", err)
	}
	conn.Close()
	return
}

// ExpireUpperCache expire the upper key.
func (d *Dao) ExpireUpperCache(c context.Context, mid int64) (ok bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("EXPIRE", upperKey(mid), d.redisUpperExpire)); err != nil {
		PromError("redis:up主设定过期")
		log.Error("conn.Send(EXPIRE, %s) error(%+v)", upperKey(mid), err)
	}
	return
}

// ExpireUppersCache expire the upper key.
func (d *Dao) ExpireUppersCache(c context.Context, mids []int64) (res map[int64]bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	res = make(map[int64]bool, len(mids))
	for _, mid := range mids {
		if err = conn.Send("EXPIRE", upperKey(mid), d.redisUpperExpire); err != nil {
			PromError("redis:up主设定过期")
			log.Error("conn.Send(EXPIRE, %s) error(%+v)", upperKey(mid), err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		PromError("redis:up主flush")
		log.Error("conn.Flush error(%+v)", err)
		return
	}
	var ok bool
	for _, mid := range mids {
		if ok, err = redis.Bool(conn.Receive()); err != nil {
			PromError("redis:up主receive")
			log.Error("conn.Receive() error(%+v)", err)
			return
		}
		res[mid] = ok
	}
	return
}

// UppersCaches batch get new articles of uppers by cache.
func (d *Dao) UppersCaches(c context.Context, mids []int64, start, end int) (res map[int64][]int64, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	res = make(map[int64][]int64, len(mids))
	for _, mid := range mids {
		if err = conn.Send("ZREVRANGE", upperKey(mid), start, end); err != nil {
			PromError("redis:获取up主")
			log.Error("conn.Send(%s) error(%+v)", upperKey(mid), err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		PromError("redis:获取up主flush")
		log.Error("conn.Flush error(%+v)", err)
		return
	}
	for _, mid := range mids {
		aids, err := redis.Int64s(conn.Receive())
		if err != nil {
			PromError("redis:获取up主receive")
			log.Error("conn.Send(ZREVRANGE, %d) error(%+v)", mid, err)
		}
		l := len(aids)
		if l == 0 {
			continue
		}
		if aids[l-1] == _blank {
			aids = aids[:l-1]
		}
		res[mid] = aids
	}
	cachedCount.Add("up", int64(len(res)))
	return
}

// AddUpperCache adds passed article of upper.
func (d *Dao) AddUpperCache(c context.Context, mid, aid int64, ptime int64) (err error) {
	art := map[int64][][2]int64{mid: [][2]int64{[2]int64{aid, ptime}}}
	err = d.AddUpperCaches(c, art)
	return
}

// AddUpperCaches batch add passed article of upper.
func (d *Dao) AddUpperCaches(c context.Context, idsm map[int64][][2]int64) (err error) {
	var (
		mid, aid, ptime int64
		arts            [][2]int64
		conn            = d.redis.Get(c)
		count           int
	)
	defer conn.Close()
	for mid, arts = range idsm {
		key := upperKey(mid)
		if len(arts) == 0 {
			arts = [][2]int64{[2]int64{_blank, _blank}}
		}
		for _, art := range arts {
			aid = art[0]
			ptime = art[1]
			if err = conn.Send("ZADD", key, "CH", ptime, aid); err != nil {
				PromError("redis:增加up主缓存")
				log.Error("conn.Send(ZADD, %s, %d, %d) error(%+v)", key, aid, err)
				return
			}
			count++
		}
		if err = conn.Send("EXPIRE", key, d.redisUpperExpire); err != nil {
			PromError("redis:增加up主expire")
			log.Error("conn.Expire error(%+v)", err)
			return
		}
		count++
	}
	if err = conn.Flush(); err != nil {
		PromError("redis:增加up主flush")
		log.Error("conn.Flush error(%+v)", err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			PromError("redis:增加up主receive")
			log.Error("conn.Receive error(%+v)", err)
			return
		}
	}
	return
}

// DelUpperCache delete article of upper cache.
func (d *Dao) DelUpperCache(c context.Context, mid int64, aid int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("ZREM", upperKey(mid), aid); err != nil {
		PromError("redis:删除up主")
		log.Error("conn.Do(ZERM, %s, %d) error(%+v)", upperKey(mid), aid, err)
	}
	return
}

// UpperArtsCountCache get upper articles count
func (d *Dao) UpperArtsCountCache(c context.Context, mid int64) (res int, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if res, err = redis.Int(conn.Do("ZCOUNT", upperKey(mid), 0, "+inf")); err != nil {
		PromError("redis:up主文章计数")
		log.Error("conn.Do(ZCARD, %s) error(%+v)", upperKey(mid), err)
	}
	return
}

// MoreArtsCaches batch get early articles of upper by publish time.
func (d *Dao) MoreArtsCaches(c context.Context, mid, ptime int64, num int) (before []int64, after []int64, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("ZREVRANGEBYSCORE", upperKey(mid), fmt.Sprintf("(%d", ptime), "-inf", "LIMIT", 0, num); err != nil {
		PromError("redis:获取up主更早文章")
		log.Error("conn.Send(%s) error(%+v)", upperKey(mid), err)
		return
	}
	if err = conn.Send("ZRANGEBYSCORE", upperKey(mid), fmt.Sprintf("(%d", ptime), "+inf", "LIMIT", 0, num); err != nil {
		PromError("redis:获取up主更晚文章")
		log.Error("conn.Send(%s) error(%+v)", upperKey(mid), err)
		return
	}
	if err = conn.Flush(); err != nil {
		PromError("redis:获取up主更晚文章")
		log.Error("conn.Flush error(%+v)", err)
		return
	}
	if before, err = redis.Int64s(conn.Receive()); err != nil {
		PromError("redis:获取up主更早文章")
		log.Error("conn.Receive error(%+v)", err)
		return
	}
	if after, err = redis.Int64s(conn.Receive()); err != nil {
		PromError("redis:获取up主更晚文章")
		log.Error("conn.Receive error(%+v)", err)
		return
	}
	l := len(before)
	if l == 0 {
		return
	}
	if before[l-1] == _blank {
		before = before[:l-1]
	}
	return
}

// ExpireRankCache expire rank cache
func (d *Dao) ExpireRankCache(c context.Context, cid int64) (res bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	var ttl int64
	if ttl, err = redis.Int64(conn.Do("TTL", rankKey(cid))); err != nil {
		PromError("redis:排行榜expire")
		log.Error("ExpireRankCache(ttl %s) error(%+v)", rankKey(cid), err)
		return
	}
	if ttl > (d.redisRankTTL - d.redisRankExpire) {
		res = true
		return
	}
	return
}

// RankCache get rank cache
func (d *Dao) RankCache(c context.Context, cid int64) (res model.RankResp, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := rankKey(cid)
	var s string
	if s, err = redis.String(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		PromError("redis:获取排行榜")
		log.Error("dao.RankCache zrevrange(%s) err: %+v", key, err)
		return
	}
	err = json.Unmarshal([]byte(s), &res)
	return
}

// AddRankCache add rank cache
func (d *Dao) AddRankCache(c context.Context, cid int64, arts model.RankResp) (err error) {
	var (
		key   = rankKey(cid)
		conn  = d.redis.Get(c)
		count int
	)
	defer conn.Close()
	if len(arts.List) == 0 {
		return
	}
	if err = conn.Send("DEL", key); err != nil {
		PromError("redis:删除排行榜缓存")
		log.Error("conn.Send(DEL, %s) error(%+v)", key, err)
		return
	}
	count++
	value, _ := json.Marshal(arts)
	if err = conn.Send("SET", key, value); err != nil {
		PromError("redis:增加排行榜缓存")
		log.Error("conn.Send(SET, %s, %s) error(%+v)", key, value, err)
		return
	}
	count++
	if err = conn.Send("EXPIRE", key, d.redisRankTTL); err != nil {
		PromError("redis:expire排行榜")
		log.Error("conn.Send(EXPIRE, %s, %v) error(%+v)", key, d.redisRankTTL, err)
		return
	}
	count++
	if err = conn.Flush(); err != nil {
		PromError("redis:增加排行榜flush")
		log.Error("conn.Flush error(%+v)", err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			PromError("redis:增加排行榜主receive")
			log.Error("conn.Receive error(%+v)", err)
			return
		}
	}
	return
}

// AddCacheHotspotArts .
func (d *Dao) AddCacheHotspotArts(c context.Context, typ int8, id int64, arts [][2]int64, replace bool) (err error) {
	var (
		key   = hotspotKey(typ, id)
		conn  = d.redis.Get(c)
		count int
	)
	defer conn.Close()
	if len(arts) == 0 {
		return
	}
	if replace {
		if err = conn.Send("DEL", key); err != nil {
			PromError("redis:删除热点标签缓存")
			log.Error("conn.Send(DEL, %s) error(%+v)", key, err)
			return
		}
		count++
	}
	for _, art := range arts {
		id := art[0]
		score := art[1]
		if err = conn.Send("ZADD", key, "CH", score, id); err != nil {
			PromError("redis:增加热点标签缓存")
			log.Error("conn.Send(ZADD, %s, %d, %v) error(%+v)", key, score, id, err)
			return
		}
		count++
	}
	if err = conn.Send("EXPIRE", key, d.redisHotspotExpire); err != nil {
		PromError("redis:热点标签设定过期")
		log.Error("conn.Send(EXPIRE, %s, %d) error(%+v)", key, d.redisHotspotExpire, err)
		return
	}
	count++
	if err = conn.Flush(); err != nil {
		PromError("redis:增加热点标签缓存flush")
		log.Error("conn.Flush error(%+v)", err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			PromError("redis:增加热点标签缓存receive")
			log.Error("conn.Receive error(%+v)", err)
			return
		}
	}
	return
}

// HotspotArtsCache .
func (d *Dao) HotspotArtsCache(c context.Context, typ int8, id int64, start, end int) (res []int64, err error) {
	key := hotspotKey(typ, id)
	conn := d.redis.Get(c)
	defer conn.Close()
	res, err = redis.Int64s(conn.Do("ZREVRANGE", key, start, end))
	if err != nil {
		PromError("redis:获取热点标签列表receive")
		log.Error("conn.Send(ZREVRANGE, %s) error(%+v)", key, err)
	}
	return
}

// HotspotArtsCacheCount .
func (d *Dao) HotspotArtsCacheCount(c context.Context, typ int8, id int64) (res int64, err error) {
	key := hotspotKey(typ, id)
	conn := d.redis.Get(c)
	defer conn.Close()
	res, err = redis.Int64(conn.Do("ZCARD", key))
	if err != nil {
		PromError("redis:获取热点标签计数")
		log.Error("conn.Send(ZCARD, %s) error(%+v)", key, err)
	}
	return
}

// ExpireHotspotArtsCache .
func (d *Dao) ExpireHotspotArtsCache(c context.Context, typ int8, id int64) (ok bool, err error) {
	key := hotspotKey(typ, id)
	conn := d.redis.Get(c)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.redisHotspotExpire)); err != nil {
		PromError("redis:热点运营设定过期")
		log.Error("conn.Send(EXPIRE, %s) error(%+v)", key, err)
	}
	return
}

// DelHotspotArtsCache .
func (d *Dao) DelHotspotArtsCache(c context.Context, typ int8, hid int64, aid int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := hotspotKey(typ, hid)
	if _, err = conn.Do("ZREM", key, aid); err != nil {
		PromError("redis:删除热点运营文章")
		log.Error("conn.Do(ZERM, %s, %d) error(%+v)", key, aid, err)
	}
	return
}

// AuthorMostCategories .
func (d *Dao) AuthorMostCategories(c context.Context, mid int64) (categories []int64, err error) {
	var (
		categoriesInts []string
		category       int64
	)
	conn := d.redis.Get(c)
	defer conn.Close()
	key := authorCategoriesKey(mid)
	if categoriesInts, err = redis.Strings(conn.Do("SMEMBERS", key)); err != nil {
		PromError("redis:获取作者分区")
		log.Error("conn.Do(GET, %s) error(%+v)", key, err)
	}
	for _, categoryInt := range categoriesInts {
		if category, err = strconv.ParseInt(categoryInt, 10, 64); err != nil {
			PromError("redis:获取作者分区")
			log.Error("strconv.Atoi(%s) error(%+v)", categoryInt, err)
			return
		}
		categories = append(categories, category)
	}
	return
}

// CategoryAuthors .
func (d *Dao) CategoryAuthors(c context.Context, category int64, count int) (authors []int64, err error) {
	var (
		authorsInts []string
		author      int64
	)
	conn := d.redis.Get(c)
	defer conn.Close()
	key := recommendsAuthorsKey(category)
	if authorsInts, err = redis.Strings(conn.Do("SRANDMEMBER", key, count)); err != nil {
		PromError("redis:获取分区作者")
		log.Error("conn.Do(GET, %s) error(%+v)", key, err)
	}
	for _, authorInt := range authorsInts {
		if author, err = strconv.ParseInt(authorInt, 10, 64); err != nil {
			PromError("redis:获取作者分区")
			log.Error("strconv.Atoi(%s) error(%+v)", authorInt, err)
			return
		}
		authors = append(authors, author)
	}
	return
}
