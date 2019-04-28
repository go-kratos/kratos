package dao

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"strconv"

	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_subPrefix = "artsl_"
)

func midSub(mid int64, title string) string {
	ms := md5.Sum([]byte(title))
	return _subPrefix + strconv.FormatInt(mid, 10) + "_" + hex.EncodeToString(ms[:])
}

// SubmitCache get user submit cache.
func (d *Dao) SubmitCache(c context.Context, mid int64, title string) (exist bool, err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := midSub(mid, title)
	_, err = conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Get error(%+v) | key(%s) mid(%d) title(%s)", err, key, mid, title)
			PromError("creation:获取标题缓存")
		}
		return
	}
	exist = true
	return
}

// AddSubmitCache add submit cache into mc.
func (d *Dao) AddSubmitCache(c context.Context, mid int64, title string) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := midSub(mid, title)
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, 1)
	if err = conn.Set(&memcache.Item{Key: key, Object: bs, Flags: memcache.FlagJSON, Expiration: d.mcSubExp}); err != nil {
		log.Error("memcache.set error(%+v) | key(%s) mid(%d) title(%s)", err, key, mid, title)
		PromError("creation:设定标题缓存")
	}
	return
}

// DelSubmitCache del submit cache into mc.
func (d *Dao) DelSubmitCache(c context.Context, mid int64, title string) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(midSub(mid, title)); err == memcache.ErrNotFound {
		err = nil
	}
	if err != nil {
		PromError("creation:删除标题缓存")
		log.Error("creation: dao.DelSubmitCache(mid: %v, title: %v) err: %+v", mid, title, err)
	}
	return
}
