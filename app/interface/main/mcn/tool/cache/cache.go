package cache

import (
	"context"
	"fmt"
	"go-common/library/cache"
	"go-common/library/cache/memcache"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/stat/prom"
	"time"
)

//DataLoader cache interface
type DataLoader interface {
	Key() (key string)
	Value() (value interface{})

	// LoadValue return value need cache
	// if err, nothing will cache
	// if value == nil, and IsNullCached is true, empty will be cached
	LoadValue(c context.Context) (value interface{}, err error)
	Expire() time.Duration
	Desc() string
}

// Get
// Delete
// Add

//MCWrapper wrapper for mc
type MCWrapper struct {
	mc    *memcache.Pool
	cache *cache.Cache

	// 是否缓存空值，防止缓存穿透
	IsNullCached bool
}

// null definition
const (
	IsNull  = 1
	NotNull = 0
)

type cacheValue struct {
	Null  int8        `json:"n"` // not 0 means null
	Value interface{} `json:"v"`
}

//IsNull return true is it's null
func (s *cacheValue) IsNull() bool {
	return s.Null != NotNull
}

//New new memcache wrapper
func New(mc *memcache.Pool) *MCWrapper {
	return &MCWrapper{
		mc:    mc,
		cache: cache.New(10, 1024),
	}
}
func (m *MCWrapper) addRaw(c context.Context, data DataLoader, cacheV *cacheValue) (err error) {
	if data == nil {
		return
	}
	conn := m.mc.Get(c)
	defer conn.Close()
	key := data.Key()

	item := &memcache.Item{Key: key, Object: cacheV, Expiration: int32(data.Expire() / time.Second), Flags: memcache.FlagJSON}
	if err = conn.Set(item); err != nil {
		actionDesc := "Add" + data.Desc()
		prom.BusinessErrCount.Incr("mc:" + actionDesc)
		log.Errorv(c, log.KV(actionDesc, fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	log.Info("Add key ok, key=%s, null=%d", key, cacheV.Null)
	return
}

//Add add cache data
func (m *MCWrapper) Add(c context.Context, data DataLoader) (err error) {
	var cacheV = &cacheValue{
		Value: data.Value(),
	}
	return m.addRaw(c, data, cacheV)
}

//Delete delete cache data
func (m *MCWrapper) Delete(c context.Context, data DataLoader) (err error) {
	conn := m.mc.Get(c)
	defer conn.Close()
	key := data.Key()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		actionDesc := "Del" + data.Desc()
		prom.BusinessErrCount.Incr("mc:" + actionDesc)
		log.Errorv(c, log.KV(actionDesc, fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}

//Get get data
func (m *MCWrapper) Get(c context.Context, data DataLoader) (err error) {
	_, err = m.getRaw(c, data)
	return
}

func (m *MCWrapper) getRaw(c context.Context, data DataLoader) (v *cacheValue, err error) {
	conn := m.mc.Get(c)
	defer conn.Close()
	key := data.Key()
	value, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		actionDesc := "Cache" + data.Desc()
		prom.BusinessErrCount.Incr("mc:" + actionDesc)
		log.Errorv(c, log.KV(actionDesc, fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	var cacheV = cacheValue{
		Value: data.Value(),
	}
	err = conn.Scan(value, &cacheV)
	if err != nil {
		actionDesc := "Cache" + data.Desc()
		prom.BusinessErrCount.Incr("mc:" + actionDesc)
		log.Errorv(c, log.KV(actionDesc, fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	v = &cacheV
	return
}

//GetOrLoad get from cache, if not found, then call data.LoadValue to load
func (m *MCWrapper) GetOrLoad(c context.Context, data DataLoader) (err error) {
	var v *cacheValue
	v, err = m.getRaw(c, data)
	if err != nil {
		return
	}

	if v != nil && !v.IsNull() {
		prom.CacheHit.Incr(data.Desc())
		return
	}

	// 没有找到对应的缓存，需求去拉取
	prom.CacheMiss.Incr(data.Desc())
	res, err := data.LoadValue(c)
	if err != nil {
		return
	}

	// 没有查到值，并且不缓存空值
	if res == nil && !m.IsNullCached {
		return
	}

	var cacheV = &cacheValue{
		Value: res,
	}

	if res == nil {
		cacheV.Null = IsNull
	}

	m.cache.Save(func() {
		m.addRaw(metadata.WithContext(c), data, cacheV)
	})
	return
}
