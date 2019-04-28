package dao

import (
	"context"
	"fmt"
	"strconv"

	relation "go-common/app/service/main/relation/model"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/log"
	xtime "go-common/library/time"
)

const (
	_prefixFollowings = "at_"
	_prefixTags       = "tags_" // user tag info.
	_prefixFollowing  = "pb_a_"
	_prefixStat       = "c_"         // key of stat
	_prefixTagCount   = "rs_tmtc_%d" // key of relation tag by mid & tag's count
)

func statKey(mid int64) string {
	return _prefixStat + strconv.FormatInt(mid, 10)
}

func tagsKey(mid int64) string {
	return _prefixTags + strconv.FormatInt(mid, 10)
}

func followingsKey(mid int64) string {
	return _prefixFollowings + strconv.FormatInt(mid, 10)
}

func followingKey(mid int64) string {
	return _prefixFollowing + strconv.FormatInt(mid, 10)
}

func tagCountKey(mid int64) string {
	return fmt.Sprintf(_prefixTagCount, mid)
}

// DelStatCache is
func (d *Dao) DelStatCache(ctx context.Context, mid int64) error {
	conn := d.mc.Get(ctx)
	defer conn.Close()
	if err := conn.Delete(statKey(mid)); err != nil {
		if err == memcache.ErrNotFound {
			return nil
		}
		log.Error("Failed to delete stat cache: mid: %d: %+v", mid, err)
		return err
	}
	return nil
}

// DelFollowerCache del follower cache
func (d *Dao) DelFollowerCache(ctx context.Context, fid int64) error {
	key := followingKey(fid)
	conn := d.mc.Get(ctx)
	defer conn.Close()
	if err := conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Delete(%s) error(%v)", key, err)
			return err
		}
	}
	return nil
}

// DelFollowing del following cache.
func (d *Dao) DelFollowing(c context.Context, mid int64, following *relation.Following) (err error) {
	var (
		ok  bool
		key = followingsKey(mid)
	)
	conn := d.redis.Get(c)
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.RelationTTL)); err != nil {
		log.Error("redis.Bool(conn.Do(EXPIRE, %s)) error(%v)", key, err)
	} else if ok {
		if _, err = conn.Do("HDEL", key, following.Mid); err != nil {
			log.Error("conn.Do(HDEL, %s, %d) error(%v)", key, following.Mid, err)
		}
	}
	conn.Close()
	return
}

// DelTagsCache is
func (d *Dao) DelTagsCache(ctx context.Context, mid int64) (err error) {
	conn := d.mc.Get(ctx)
	if err = conn.Delete(tagsKey(mid)); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Delete(%s) error(%v)", tagCountKey(mid), err)
		}
	}
	conn.Close()
	return
}

// AddFollowingCache is
func (d *Dao) AddFollowingCache(c context.Context, mid int64, following *relation.Following) (err error) {
	var (
		ok  bool
		key = followingsKey(mid)
	)
	conn := d.redis.Get(c)
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.RelationTTL)); err != nil {
		log.Error("redis.Bool(conn.Do(EXPIRE, %s)) error(%v)", key, err)
	} else if ok {
		var ef []byte
		if ef, err = d.encode(following.Attribute, following.MTime, following.Tag, following.Special); err != nil {
			return
		}
		if _, err = conn.Do("HSET", key, following.Mid, ef); err != nil {
			log.Error("conn.Do(HSET, %s, %d) error(%v)", key, following.Mid, err)
		}
	}
	conn.Close()
	return
}

// encode
func (d *Dao) encode(attribute uint32, mtime xtime.Time, tagids []int64, special int32) (res []byte, err error) {
	ft := &relation.FollowingTags{Attr: attribute, Ts: mtime, TagIds: tagids, Special: special}
	return ft.Marshal()
}

// DelFollowingCache delete following cache.
func (d *Dao) DelFollowingCache(c context.Context, mid int64) (err error) {
	return d.delFollowingCache(c, followingKey(mid))
}

// delFollowingCache delete following cache.
func (d *Dao) delFollowingCache(c context.Context, key string) (err error) {
	conn := d.mc.Get(c)
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Delete(%s) error(%v)", key, err)
		}
	}
	conn.Close()
	return
}

// DelTagCountCache del tag count cache.
func (d *Dao) DelTagCountCache(c context.Context, mid int64) (err error) {
	conn := d.mc.Get(c)
	if err = conn.Delete(tagCountKey(mid)); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Delete(%s) error(%v)", tagCountKey(mid), err)
		}
	}
	conn.Close()
	return
}
