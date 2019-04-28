package dao

import (
	"context"
	"encoding/json"
	"go-common/library/conf/paladin"
	"time"

	"go-common/app/admin/main/workflow/model/search"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_prefixChallPendingCount = "wkf_cpc_uid_%d"
)

// pingMC ping memcache.
func (d *Dao) pingMC(c context.Context) (err error) {
	conn := d.mc.Get(c)
	//if err = conn.Store("set", "ping", []byte{1}, 0, d.mcExpire, 0); err != nil {
	if err = conn.Set(&memcache.Item{Key: "ping", Value: []byte{1}, Expiration: 0}); err != nil {
		log.Error("conn.Store(set, ping, 1) error(%v)", err)
	}
	conn.Close()
	return
}

// ChallCountCache read pending chall count by uid from memcache
func (d *Dao) ChallCountCache(c context.Context, uid int64) (challCount *search.ChallCount, err error) {
	var (
		conn memcache.Conn
		item *memcache.Item
		key  string
	)
	conn = d.mc.Get(c)
	defer conn.Close()
	key = d.keyChallCount(uid)
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		return
	}
	err = json.Unmarshal(item.Value, &challCount)
	return
}

// UpChallCountCache will write chall count to cache
func (d *Dao) UpChallCountCache(c context.Context, challCount *search.ChallCount, uid int64) (err error) {
	var (
		conn           memcache.Conn
		item           *memcache.Item
		jsonChallCount []byte
	)
	jsonChallCount, err = json.Marshal(challCount)
	if err != nil {
		return
	}
	conn = d.mc.Get(c)
	defer conn.Close()
	item = &memcache.Item{
		Key:        d.keyChallCount(uid),
		Value:      jsonChallCount,
		Expiration: int32(paladin.Duration(d.c.Get("expireCount"), time.Duration(10*time.Second)) / time.Second),
	}
	err = conn.Set(item)
	return
}
