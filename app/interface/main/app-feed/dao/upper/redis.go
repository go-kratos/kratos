package upper

import (
	"context"
	"strconv"

	"go-common/app/service/main/archive/api"
	feed "go-common/app/service/main/feed/model"
	"go-common/library/log"
	"go-common/library/xstr"

	"go-common/library/cache/redis"

	"github.com/pkg/errors"
)

const (
	// up items
	_prefixUpItems = "u3_"
	// unread count
	_prefixUnreadCount = "uc2_"
)

func keyUpItem(mid int64) string {
	return _prefixUpItems + strconv.FormatInt(mid, 10)
}

func keyUnreadCount(mid int64) string {
	return _prefixUnreadCount + strconv.FormatInt(mid%100000, 10)
}

func (d *Dao) UpItemCaches(c context.Context, mid int64, start, end int) (uis []*feed.Feed, aids []int64, seasonIDs []int64, err error) {
	var vs []interface{}
	conn := d.redis.Get(c)
	key := keyUpItem(mid)
	defer conn.Close()
	if vs, err = redis.Values(conn.Do("ZREVRANGE", key, start, end, "WITHSCORES")); err != nil {
		err = errors.Wrapf(err, "conn.Do(ZREVRANGE,%s,%d,%d)", key, start, end)
		return
	}
	uis = make([]*feed.Feed, 0, len(vs))
	aids = make([]int64, 0, len(vs))
	seasonIDs = make([]int64, 0, len(vs))
Loop:
	for len(vs) > 0 {
		var (
			i      int64
			value  string
			values []int64
		)
		if vs, err = redis.Scan(vs, &value, &i); err != nil {
			err = errors.Wrapf(err, "%v", vs)
			return
		}
		if values, err = xstr.SplitInts(value); err != nil {
			log.Error("xstr.SplitInts(%v) error(%v)", value, err)
			continue Loop
		}
		if len(values) >= 2 {
			ua := &feed.Feed{}
			urs := make([]*api.Arc, 0, len(values)-2)
			for k, v := range values {
				if k == 0 {
					ua.Type = v
				} else if k == 1 {
					ua.ID = v
					switch ua.Type {
					case feed.ArchiveType:
						aids = append(aids, v)
					case feed.BangumiType:
						seasonIDs = append(seasonIDs, v)
					}
				} else if k >= 2 {
					switch ua.Type {
					case feed.ArchiveType:
						aids = append(aids, v)
						urs = append(urs, &api.Arc{Aid: v})
					}
				}
			}
			ua.Fold = urs
			uis = append(uis, ua)
		}
	}
	return
}

func (d *Dao) AddUpItemCaches(c context.Context, mid int64, uis ...*feed.Feed) (err error) {
	var (
		ucKey = keyUnreadCount(mid)
		upKey = keyUpItem(mid)
		conn  = d.redis.Get(c)
	)
	defer conn.Close()
	if err = conn.Send("HSET", ucKey, mid, 0); err != nil {
		err = errors.Wrapf(err, "conn.Send(HSET,%s,%d,%d) error(%v)", ucKey, mid, 0)
		return
	}
	if err = conn.Send("ZREMRANGEBYRANK", upKey, 0, -1); err != nil {
		err = errors.Wrapf(err, "conn.Send(ZREMRANGEBYRANK,%s,%d,%d)", upKey, 0, -1)
		return
	}
	for _, ui := range uis {
		if ui.ID != 0 && (ui.Type == feed.ArchiveType || ui.Type == feed.BangumiType) {
			var vs = []int64{ui.Type, ui.ID}
			for _, r := range ui.Fold {
				if r.Aid != 0 {
					vs = append(vs, r.Aid)
				}
			}
			var (
				score = ui.PubDate.Time().Unix()
				value = xstr.JoinInts(vs)
			)
			if err = conn.Send("ZADD", upKey, score, value); err != nil {
				err = errors.Wrapf(err, "conn.Send(ZADD,%s,%d,%d)", upKey, score, value)
				return
			}
		}
	}
	if err = conn.Send("EXPIRE", upKey, d.expireRds); err != nil {
		err = errors.Wrapf(err, "conn.Send(EXPIRE,%s,%d)", upKey, d.expireRds)
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	for i := 0; i < len(uis)+3; i++ {
		if _, err = conn.Receive(); err != nil {
			return
		}
	}
	return
}

func (d *Dao) ExpireUpItem(c context.Context, mid int64) (ok bool, err error) {
	var (
		key  = keyUpItem(mid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.expireRds)); err != nil {
		err = errors.Wrapf(err, "conn.Do(EXPIRE,%s,%d)", key, d.expireRds)
	}
	return
}

func (d *Dao) UnreadCountCache(c context.Context, mid int64) (unread int, err error) {
	var (
		key  = keyUnreadCount(mid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if unread, err = redis.Int(conn.Do("HGET", key, mid)); err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		err = errors.Wrapf(err, "conn.Do(HGET,%s,%d)", key, mid)
	}
	return
}

func (d *Dao) AddUnreadCountCache(c context.Context, mid int64, unread int) (err error) {
	var (
		key  = keyUnreadCount(mid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if _, err = conn.Do("HSET", key, mid, unread); err != nil {
		err = errors.Wrapf(err, "conn.DO(HSET,%s,%d,%d)", key, mid, unread)
	}
	return
}
