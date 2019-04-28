package dao

import (
	"encoding/json"

	"context"
	"go-common/app/admin/main/tv/model"
	"go-common/library/cache/memcache"
)

//SiMcOutKey is used for tv search intervene key with MC
const SiMcOutKey = "_tv_search"

//SiMcStateKey is used for tv search intervene publish status key with MC
const SiMcStateKey = "_tv_search_state"

// SetSearchInterv is used for setting search inter rank cache
func (d *Dao) SetSearchInterv(c context.Context, rank []*model.OutSearchInter) (err error) {
	var (
		conn memcache.Conn
		item *memcache.Item
		bs   []byte
	)
	bs, err = json.Marshal(rank)
	if err != nil {
		return
	}
	conn = d.mc.Get(c)
	defer conn.Close()
	item = &memcache.Item{
		Key:        SiMcOutKey,
		Value:      bs,
		Expiration: 0,
	}
	err = conn.Set(item)
	return
}

// GetSearchInterv is used for getting search inter rank cache
func (d *Dao) GetSearchInterv(c context.Context) (rank []*model.OutSearchInter, err error) {
	var (
		conn memcache.Conn
		item *memcache.Item
	)
	conn = d.mc.Get(c)
	defer conn.Close()
	if item, err = conn.Get(SiMcOutKey); err != nil {
		return
	}
	if err = json.Unmarshal(item.Value, &rank); err != nil {
		return
	}
	return
}

// SetPublishCache is used for setting publish status
func (d *Dao) SetPublishCache(c context.Context, state *model.PublishStatus) (err error) {
	var (
		conn memcache.Conn
		item *memcache.Item
		bs   []byte
	)
	bs, err = json.Marshal(state)
	if err != nil {
		return
	}
	conn = d.mc.Get(c)
	defer conn.Close()
	item = &memcache.Item{
		Key:        SiMcStateKey,
		Value:      bs,
		Expiration: 0,
	}
	err = conn.Set(item)
	return
}

// GetPublishCache is used for getting search inter rank cache
func (d *Dao) GetPublishCache(c context.Context) (state *model.PublishStatus, err error) {
	var (
		conn memcache.Conn
		item *memcache.Item
	)
	conn = d.mc.Get(c)
	defer conn.Close()
	if item, err = conn.Get(SiMcStateKey); err != nil {
		return
	}
	if err = json.Unmarshal(item.Value, &state); err != nil {
		return
	}
	return
}
