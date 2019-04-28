package dao

import (
	"context"
	"fmt"
	"strconv"

	"go-common/app/service/main/relation/model"
	gmc "go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/log"
	"go-common/library/time"
)

const (
	_prefixFollowings = "at_"
	_prefixTags       = "tags_" // user tag info.
)

func tagsKey(mid int64) string {
	return _prefixTags + strconv.FormatInt(mid, 10)
}

func followingsKey(mid int64) string {
	return _prefixFollowings + strconv.FormatInt(mid, 10)
}

// ==== redis ===

// AddFollowingCache add following cache.
func (d *Dao) AddFollowingCache(c context.Context, mid int64, following *model.Following) (err error) {
	var (
		ok  bool
		key = followingsKey(mid)
	)
	conn := d.relRedis.Get(c)
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.relExpire)); err != nil {
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

// DelFollowing del following cache.
func (d *Dao) DelFollowing(c context.Context, mid int64, following *model.Following) (err error) {
	var (
		ok  bool
		key = followingsKey(mid)
	)
	conn := d.relRedis.Get(c)
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.relExpire)); err != nil {
		log.Error("redis.Bool(conn.Do(EXPIRE, %s)) error(%v)", key, err)
	} else if ok {
		if _, err = conn.Do("HDEL", key, following.Mid); err != nil {
			log.Error("conn.Do(HDEL, %s, %d) error(%v)", key, following.Mid, err)
		}
	}
	conn.Close()
	return
}

// encode
func (d *Dao) encode(attribute uint32, mtime time.Time, tagids []int64, special int32) (res []byte, err error) {
	ft := &model.FollowingTags{Attr: attribute, Ts: mtime, TagIds: tagids, Special: special}
	return ft.Marshal()
}

// ===== memcache =====
const (
	_prefixFollowing = "pb_a_"
	_prefixTagCount  = "rs_tmtc_%d" // key of relation tag by mid & tag's count
)

func followingKey(mid int64) string {
	return _prefixFollowing + strconv.FormatInt(mid, 10)
}

func tagCountKey(mid int64) string {
	return fmt.Sprintf(_prefixTagCount, mid)
}

// DelFollowingCache delete following cache.
func (d *Dao) DelFollowingCache(c context.Context, mid int64) (err error) {
	return d.delFollowingCache(c, followingKey(mid))
}

// delFollowingCache delete following cache.
func (d *Dao) delFollowingCache(c context.Context, key string) (err error) {
	conn := d.mc.Get(c)
	if err = conn.Delete(key); err != nil {
		if err == gmc.ErrNotFound {
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
		if err == gmc.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Delete(%s) error(%v)", tagCountKey(mid), err)
		}
	}
	conn.Close()
	return
}

// DelTagsCache is
func (d *Dao) DelTagsCache(c context.Context, mid int64) (err error) {
	conn := d.mc.Get(c)
	if err = conn.Delete(tagsKey(mid)); err != nil {
		if err == gmc.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Delete(%s) error(%v)", tagCountKey(mid), err)
		}
	}
	conn.Close()
	return
}
