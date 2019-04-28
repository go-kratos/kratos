package dao

import (
	"context"
	"fmt"
	"strconv"

	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/app/service/main/feed/model"
	feedmdl "go-common/app/service/main/feed/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
	"go-common/library/time"
	"go-common/library/xstr"
)

const (
	_prefixUpper          = "ua_"   // upper's archive list
	_prefixAppFeed        = "af_"   // user's app feed list
	_prefixWebFeed        = "wf_"   // user's web feed list
	_prefixAppLast        = "al_"   // user's last access
	_prefixWebLast        = "wl_"   // user's last access
	_prefixArtLast        = "tl_"   // user's last access of article
	_prefixBangumiFeed    = "banf_" // user's bangumi feed list
	_prefixArchiveFeed    = "arcf_" // user's archive feed list
	_prefixArticleFeed    = "artf_" // user's article feed list
	_prefixAppUnreadCount = "ac_"   // user's app unread count
	_prefixWebUnreadCount = "wc_"   // user's web unread count
	_prefixArtUnreadCount = "tc_"   // user's article unread count
)

func upperKey(mid int64) string {
	return _prefixUpper + strconv.FormatInt(mid, 10)
}

func bangumiFeedKey(mid int64) string {
	return _prefixBangumiFeed + strconv.FormatInt(mid, 10)
}

func archiveFeedKey(mid int64) string {
	return _prefixArchiveFeed + strconv.FormatInt(mid, 10)
}

func from(i int64) (time.Time, int8) {
	return time.Time((i >> 8)), int8(int64(i) & 0xff)
}

func combine(t time.Time, copyright int8) int64 {
	return int64(t)<<8 | int64(copyright)
}

func feedKey(ft int, mid int64) string {
	midStr := strconv.FormatInt(mid, 10)
	if ft == model.TypeApp {
		return _prefixAppFeed + midStr
	} else if ft == model.TypeWeb {
		return _prefixWebFeed + midStr
	} else {
		return _prefixArticleFeed + midStr
	}
}

func unreadCountKey(ft int, mid int64) string {
	midStr := strconv.FormatInt(mid%100000, 10)
	if ft == model.TypeApp {
		return _prefixAppUnreadCount + midStr
	} else if ft == model.TypeWeb {
		return _prefixWebUnreadCount + midStr
	} else {
		return _prefixArtUnreadCount + midStr
	}
}

func lastKey(ft int, mid int64) string {
	midStr := strconv.FormatInt(mid%100000, 10)
	if ft == model.TypeApp {
		return _prefixAppLast + midStr
	} else if ft == model.TypeWeb {
		return _prefixWebLast + midStr
	} else {
		return _prefixArtLast + midStr
	}
}

// appFeedValue convert Feed to string, format: "type,id,fold,fold,fold..."
func appFeedValue(f *feedmdl.Feed) string {
	ids := []int64{f.Type, f.ID}
	for _, arc := range f.Fold {
		ids = append(ids, arc.Aid)
	}
	return xstr.JoinInts(ids)
}

func recoverFeed(idsStr string) (fe *feedmdl.Feed, err error) {
	var (
		aid int64
		ids []int64
	)
	if ids, err = xstr.SplitInts(idsStr); err != nil {
		return
	}
	if len(ids) < 2 {
		err = fmt.Errorf("recoverFeed failed idsStr(%v)", idsStr)
		return
	}
	fe = &feedmdl.Feed{Type: ids[0], ID: ids[1]}
	for _, aid = range ids[2:] {
		fe.Fold = append(fe.Fold, &api.Arc{Aid: aid})
	}
	return
}

// pingRedis ping redis.
func (d *Dao) pingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	if _, err = conn.Do("SET", "PING", "PONG"); err != nil {
		PromError("redis: ping remote", "remote redis: conn.Do(SET,PING,PONG) error(%v)", err)
	}
	conn.Close()
	return
}

// LastAccessCache get last access time of user.
func (d *Dao) LastAccessCache(c context.Context, ft int, mid int64) (t int64, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := lastKey(ft, mid)
	if t, err = redis.Int64(conn.Do("HGET", key, mid)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			PromError("redis:获取上次访问时间", "conn.Do(HGET, %s, %s) error(%v)", key, mid, err)
		}
	}
	return
}

// AddLastAccessCache add user's last access time.
func (d *Dao) AddLastAccessCache(c context.Context, ft int, mid int64, t int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := lastKey(ft, mid)
	if _, err = conn.Do("HSET", key, mid, t); err != nil {
		PromError("redis:增加上次访问时间", "conn.DO(HSET, %s, %d, %d) error(%v)", key, mid, t, err)
	}
	return
}

// ExpireFeedCache expire the user feed key.
func (d *Dao) ExpireFeedCache(c context.Context, ft int, mid int64) (ok bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := feedKey(ft, mid)
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.redisExpireFeed)); err != nil {
		PromError("redis:Feed缓存设定过期", "conn.Do(EXPIRE, %s, %d) error(%v)", key, d.redisExpireFeed, err)
	}
	return
}

// PurgeFeedCache purge the user feed key.
func (d *Dao) PurgeFeedCache(c context.Context, ft int, mid int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := feedKey(ft, mid)
	if _, err = redis.Bool(conn.Do("DEL", key)); err != nil {
		PromError("redis:删除feed", "conn.Do(DEL, %s, %d) error(%v)", key, err)
	}
	return
}

// FeedCache get upper feed by cache.
func (d *Dao) FeedCache(c context.Context, ft int, mid int64, start, end int) (as []*feedmdl.Feed, bids []int64, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := feedKey(ft, mid)
	vs, err := redis.Values(conn.Do("ZREVRANGE", key, start, end, "WITHSCORES"))
	if err != nil {
		PromError("redis:获取feed", "conn.Do(ZREVRANGE,%s,%d,%d) error(%v)", key, start, end, err)
		return
	}
	as = make([]*feedmdl.Feed, 0, len(vs))
	for len(vs) > 0 {
		var (
			ts     int64
			idsStr string
			fe     *feedmdl.Feed
		)
		if vs, err = redis.Scan(vs, &idsStr, &ts); err != nil {
			PromError("redis:获取feed", "redis.Scan(%v) error(%v)", vs, err)
			return
		}
		if idsStr != "" {
			fe, err = recoverFeed(idsStr)
			fe.PubDate = time.Time(ts)
			if err != nil {
				PromError("恢复feed", "redis.recoverFeed(%v) error(%v)", idsStr, err)
				err = nil
				continue
			}
			as = append(as, fe)
			switch fe.Type {
			case feedmdl.BangumiType:
				bids = append(bids, fe.ID)
			}
		}
	}
	return
}

// AddFeedCache add upper feed cache.
func (d *Dao) AddFeedCache(c context.Context, ft int, mid int64, as []*feedmdl.Feed) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := feedKey(ft, mid)
	if err = conn.Send("DEL", key); err != nil {
		PromError("redis:删除feed缓存", "conn.Send(DEL, %s) error(%v)", key, err)
		return
	}
	commondLen := 1
	if len(as) > 0 {
		var feedLen int
		if ft == model.TypeApp {
			feedLen = d.appFeedLength
		} else {
			feedLen = d.webFeedLength
		}
		if len(as) > feedLen {
			as = as[:feedLen]
		}
		commonds := []interface{}{key}
		for _, appFeed := range as {
			ts := appFeed.PubDate.Time().Unix()
			feedValue := appFeedValue(appFeed)
			commonds = append(commonds, ts, feedValue)
		}
		if err = conn.Send("ZADD", commonds...); err != nil {
			PromError("redis:增加feed缓存", "conn.Send(ZADD, %v, %v) error(%v)", key, commonds, err)
			return
		}
		commondLen++
		if err = conn.Send("EXPIRE", key, d.redisExpireFeed); err != nil {
			PromError("redis:expire-feed缓存", "conn.Send(expire, %s, %v) error(%v)", key, d.redisExpireFeed, err)
			return
		}
		commondLen++
	}
	if err = conn.Flush(); err != nil {
		PromError("redis:feed缓存flush", "conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < commondLen; i++ {
		if _, err = conn.Receive(); err != nil {
			PromError("redis:feed缓存receive", "conn.Receive error(%v)", err)
			return
		}
	}
	return
}

// ExpireUppersCache expire the upper key.
func (d *Dao) ExpireUppersCache(c context.Context, mids []int64) (res map[int64]bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	res = make(map[int64]bool, len(mids))
	for _, mid := range mids {
		if err = conn.Send("TTL", upperKey(mid)); err != nil {
			PromError("redis:up主ttl", "conn.Send(TTL, %s) error(%v)", upperKey(mid), err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		PromError("redis:up主flush", "conn.Flush error(%v)", err)
		return
	}
	var state int64
	for _, mid := range mids {
		if state, err = redis.Int64(conn.Receive()); err != nil {
			PromError("redis:up主receive", "conn.Receive() error(%v)", err)
			return
		}
		if int32(state) > (d.redisTTLUpper - d.redisExpireUpper) {
			res[mid] = true
		} else {
			res[mid] = false
		}
	}
	return
}

// UppersCaches batch get new archives of uppers by cache.
func (d *Dao) UppersCaches(c context.Context, mids []int64, start, end int) (res map[int64][]*archive.AidPubTime, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	res = make(map[int64][]*archive.AidPubTime, len(mids))
	for _, mid := range mids {
		if err = conn.Send("ZREVRANGE", upperKey(mid), start, end, "withscores"); err != nil {
			PromError("redis:获取up主", "conn.Send(%s) error(%v)", upperKey(mid), err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		PromError("redis:获取up主flush", "conn.Flush error(%v)", err)
		return
	}
	for _, mid := range mids {
		values, err := redis.Values(conn.Receive())
		if err != nil {
			PromError("redis:获取up主receive", "conn.Send(ZREVRANGE, %d) error(%v)", mid, err)
			err = nil
			continue
		}
		for len(values) > 0 {
			arc := archive.AidPubTime{}
			var score int64
			if values, err = redis.Scan(values, &arc.Aid, &score); err != nil {
				PromError("redis:scan UP主", "redis.Scan() error(%v)", err)
				err = nil
				continue
			}
			arc.PubDate, arc.Copyright = from(score)
			res[mid] = append(res[mid], &arc)
		}
	}
	CachedCount.Add("up", int64(len(res)))
	return
}

// AddUpperCaches batch add passed archive of upper.
// set max num of upper's passed list.
func (d *Dao) AddUpperCaches(c context.Context, mArcs map[int64][]*archive.AidPubTime) (err error) {
	var (
		mid   int64
		arcs  []*archive.AidPubTime
		conn  = d.redis.Get(c)
		count int
	)
	defer conn.Close()
	if len(mArcs) == 0 {
		return
	}
	for mid, arcs = range mArcs {
		if len(arcs) == 0 {
			continue
		}
		key := upperKey(mid)
		if err = conn.Send("DEL", key); err != nil {
			PromError("redis:删除up主缓存", "conn.Send(DEL, %s) error(%v)", key, err)
			return
		}
		count++
		for _, arc := range arcs {
			score := combine(arc.PubDate, arc.Copyright)
			if err = conn.Send("ZADD", key, "CH", score, arc.Aid); err != nil {
				PromError("redis:增加up主缓存", "conn.Send(ZADD, %s, %d, %d) error(%v)", key, arc.Aid, err)
				return
			}
			count++
		}
		if err = conn.Send("ZREMRANGEBYRANK", key, 0, -(d.c.MultiRedis.MaxArcsNum + 1)); err != nil {
			PromError("redis:清理up主缓存", "conn.Send(ZREMRANGEBYRANK, %s) error(%v)", key, err)
			return
		}
		count++
		if err = conn.Send("EXPIRE", key, d.redisTTLUpper); err != nil {
			PromError("redis:expireup主缓存", "conn.Send(EXPIRE, %s, %v) error(%v)", key, d.redisTTLUpper, err)
			return
		}
		count++
	}
	if err = conn.Flush(); err != nil {
		PromError("redis:增加up主flush", "conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			PromError("redis:增加up主receive", "conn.Receive error(%v)", err)
			return
		}
	}
	return
}

// AddUpperCache .
func (d *Dao) AddUpperCache(c context.Context, mid int64, arc *archive.AidPubTime) (err error) {
	var conn = d.redis.Get(c)
	defer conn.Close()
	key := upperKey(mid)
	score := combine(arc.PubDate, arc.Copyright)
	if err = conn.Send("ZADD", key, "CH", score, arc.Aid); err != nil {
		PromError("redis:增加up主缓存", "conn.Send(ZADD, %s, %d, %d) error(%v)", key, arc.Aid, err)
		return
	}
	if err = conn.Send("ZREMRANGEBYRANK", key, 0, -(d.c.MultiRedis.MaxArcsNum + 1)); err != nil {
		PromError("redis:清理up主缓存", "conn.Send(ZREMRANGEBYRANK, %s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		PromError("redis:增加up主flush", "conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			PromError("redis:增加up主receive", "conn.Receive error(%v)", err)
			return
		}
	}
	return
}

// DelUpperCache delete archive of upper cache.
func (d *Dao) DelUpperCache(c context.Context, mid int64, aid int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("ZREM", upperKey(mid), aid); err != nil {
		PromError("redis:删除up主", "conn.Do(ZERM, %s, %d) error(%v)", upperKey(mid), aid, err)
	}
	return
}

// AddArchiveFeedCache add archive feed cache.
func (d *Dao) AddArchiveFeedCache(c context.Context, mid int64, as []*feedmdl.Feed) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if len(as) == 0 {
		return
	}
	if len(as) > d.c.Feed.ArchiveFeedLength {
		as = as[:d.c.Feed.ArchiveFeedLength]
	}
	key := archiveFeedKey(mid)
	commonds := []interface{}{key}
	for _, f := range as {
		ts := f.PubDate.Time().Unix()
		value := appFeedValue(f)
		commonds = append(commonds, ts, value)
	}
	if err = conn.Send("ZADD", commonds...); err != nil {
		PromError("redis:增加archive-feed缓存", "conn.Send(ZADD, %v, %v) error(%v)", key, commonds, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.redisExpireArchiveFeed); err != nil {
		PromError("redis:expire-archive-feed缓存", "conn.Send(expire, %s, %v) error(%v)", key, d.redisExpireArchiveFeed, err)
		return
	}
	if err = conn.Flush(); err != nil {
		PromError("redis:archive-feed-flush", "conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			PromError("redis:archive-feed-receive", "conn.Receive error(%v)", err)
			return
		}
	}
	return
}

// AddBangumiFeedCache add bangumi feed cache.
func (d *Dao) AddBangumiFeedCache(c context.Context, mid int64, as []*feedmdl.Feed) (err error) {
	if len(as) == 0 {
		return
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	key := bangumiFeedKey(mid)
	commonds := []interface{}{key}
	for _, f := range as {
		ts := f.PubDate.Time().Unix()
		value := appFeedValue(f)
		commonds = append(commonds, ts, value)
	}
	if err = conn.Send("ZADD", commonds...); err != nil {
		PromError("redis:增加bangumi-feed缓存", "conn.Send(ZADD, %v, %v) error(%v)", key, commonds, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.redisExpireBangumiFeed); err != nil {
		PromError("redis:expire-bangumi-feed", "conn.Send(expire, %s, %v) error(%v)", key, d.redisExpireBangumiFeed, err)
		return
	}
	if err = conn.Flush(); err != nil {
		PromError("redis:bangumi-feed-flush", "conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			PromError("redis:bangumi-feed-receive", "conn.Receive error(%v)", err)
			return
		}
	}
	return
}

// ArchiveFeedCache get archive feed by cache.
func (d *Dao) ArchiveFeedCache(c context.Context, mid int64, start, end int) (as []*feedmdl.Feed, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := archiveFeedKey(mid)
	vs, err := redis.Values(conn.Do("ZREVRANGE", key, start, end, "WITHSCORES"))
	if err != nil {
		PromError("redis:获取archive-feed", "conn.Do(ZREVRANGE,%s,%d,%d) error(%v)", key, start, end, err)
		return
	}
	for len(vs) > 0 {
		var (
			ts     int64
			idsStr string
			fe     *feedmdl.Feed
		)
		if vs, err = redis.Scan(vs, &idsStr, &ts); err != nil {
			PromError("redis:获取archive-feed", "redis.Scan(%v) error(%v)", vs, err)
			return
		}
		if idsStr != "" {
			fe, err = recoverFeed(idsStr)
			fe.PubDate = time.Time(ts)
			if err != nil {
				PromError("恢复archive-feed", "redis.recoverFeed(%v) error(%v)", idsStr, err)
				err = nil
				continue
			}
			as = append(as, fe)
		}
	}
	return
}

// BangumiFeedCache get bangumi feed by cache.
func (d *Dao) BangumiFeedCache(c context.Context, mid int64, start, end int) (bids []int64, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := bangumiFeedKey(mid)
	vs, err := redis.Values(conn.Do("ZREVRANGE", key, start, end, "WITHSCORES"))
	if err != nil {
		PromError("redis:获取feed", "conn.Do(ZREVRANGE,%s,%d,%d) error(%v)", key, start, end, err)
		return
	}
	bids = make([]int64, 0, len(vs))
	for len(vs) > 0 {
		var (
			ts     int64
			idsStr string
			fe     *feedmdl.Feed
		)
		if vs, err = redis.Scan(vs, &idsStr, &ts); err != nil {
			PromError("redis:获取bangumi-feed", "redis.Scan(%v) error(%v)", vs, err)
			return
		}
		if idsStr != "" {
			fe, err = recoverFeed(idsStr)
			if err != nil {
				PromError("恢复bangumi-feed", "redis.recoverFeed(%v) error(%v)", idsStr, err)
				err = nil
				continue
			}
			fe.PubDate = time.Time(ts)
			bids = append(bids, fe.ID)
		}
	}

	return
}

// ArticleFeedCache get article feed by cache.
func (d *Dao) ArticleFeedCache(c context.Context, mid int64, start, end int) (aids []int64, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := feedKey(model.TypeArt, mid)
	vs, err := redis.Values(conn.Do("ZREVRANGE", key, start, end))
	if err != nil {
		log.Error("ArticleFeedCache conn.Do(ZREVRANGE,%s,%d,%d) error(%v)", key, start, end, err)
		return
	}
	for len(vs) > 0 {
		var aid int64
		if vs, err = redis.Scan(vs, &aid); err != nil {
			log.Error("ArticleFeedCache redis.Scan(%v) error(%v)", vs, err)
			return
		}
		aids = append(aids, aid)
	}
	return
}

// AddArticleFeedCache add article feed cache.
func (d *Dao) AddArticleFeedCache(c context.Context, mid int64, as []*artmdl.Meta) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if len(as) == 0 {
		return
	}
	if len(as) > d.c.Feed.ArticleFeedLength {
		as = as[:d.c.Feed.ArticleFeedLength]
	}
	key := feedKey(model.TypeArt, mid)
	commonds := []interface{}{key}
	for _, a := range as {
		ts := a.PublishTime.Time().Unix()
		commonds = append(commonds, ts, a.ID)
	}
	if err = conn.Send("ZADD", commonds...); err != nil {
		log.Error("AddArticleFeedCache conn.Send(ZADD, %v, %v) error(%v)", key, commonds, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.redisExpireArchiveFeed); err != nil {
		log.Error("AddArticleFeedCache conn.Send(expire, %s, %v) error(%v)", key, d.redisExpireArchiveFeed, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("AddArticleFeedCache conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("AddArticleFeedCache conn.Receive error(%v)", err)
			return
		}
	}
	return
}

// UnreadCountCache get unread count cache of user.
func (d *Dao) UnreadCountCache(c context.Context, ft int, mid int64) (count int, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := unreadCountKey(ft, mid)
	if count, err = redis.Int(conn.Do("HGET", key, mid)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			PromError("redis:获取未读数", "conn.Do(HGET, %s, %v) error(%v)", key, mid, err)
		}
	}
	return
}

// AddUnreadCountCache add user's unread count cache.
func (d *Dao) AddUnreadCountCache(c context.Context, ft int, mid int64, count int) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := unreadCountKey(ft, mid)
	if _, err = conn.Do("HSET", key, mid, count); err != nil {
		PromError("redis:增加未读数", "conn.DO(HSET, %s, %d, %d) error(%v)", key, mid, count, err)
	}
	return
}
