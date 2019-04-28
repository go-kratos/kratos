package dao

import (
	"context"
	"encoding/json"

	"go-common/library/cache/memcache"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	//appkeys represents key name
	_asspkeys = "sappkeys"
)

// pingMC ping memcache .
func (d *Dao) pingMC(c context.Context) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: "ping", Value: []byte{1}, Expiration: 0}); err != nil {
		log.Error("conn.Store(set, ping, 1) error (%v)", err)
	}
	return
}

// AppkeyCache .
func (d *Dao) AppkeyCache(c context.Context) (res map[string]string, err error) {
	var (
		conn memcache.Conn
		item *memcache.Item
	)
	conn = d.mc.Get(c)
	defer conn.Close()
	if item, err = conn.Get(_asspkeys); err != nil {
		if err == memcache.ErrNotFound {
			err = ecode.NothingFound
		}
		return
	}
	res = make(map[string]string)
	err = json.Unmarshal([]byte(item.Value), &res)
	return
}

// SetAppkeyCache .
func (d *Dao) SetAppkeyCache(c context.Context, newData map[string]string) (err error) {
	var (
		conn         memcache.Conn
		item         *memcache.Item
		jsonAppCache []byte
	)
	if jsonAppCache, err = json.Marshal(newData); err != nil {
		return
	}
	conn = d.mc.Get(c)
	defer conn.Close()
	item = &memcache.Item{
		Key:   _asspkeys,
		Value: jsonAppCache,
	}
	err = conn.Set(item)
	return
}
