package dao

import (
	"context"
	"fmt"
	"strings"

	"go-common/app/service/main/tag/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_prefixTag   = "t_%d"
	_prefixName  = "n_%s"
	_prefixNames = "ns_%d_%d"
	_prefixCount = "tc_%d"

	_spaceReplace = "_^_"
)

func keyTag(tid int64) string {
	return fmt.Sprintf(_prefixTag, tid)
}

func keyName(name string) string {
	return fmt.Sprintf(_prefixName, strings.Replace(name, " ", _spaceReplace, -1))
}

func keyNames(oid int64, typ int8) string {
	return fmt.Sprintf(_prefixNames, oid, typ)
}

func keyCount(tid int64) string {
	return fmt.Sprintf(_prefixCount, tid)
}

// TagCache return tag by tid from cache.
func (d *Dao) TagCache(c context.Context, tid int64) (res *model.Tag, err error) {
	var (
		key  = keyTag(tid)
		conn = d.mc.Get(c)
		item *memcache.Item
	)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Get(%s) error(%v)", key, err)
		}
		return
	}
	res = &model.Tag{}
	if err = conn.Scan(item, res); err != nil {
		log.Error("mc.Scan(%s) error(%v)", item.Value, err)
	}
	return
}

// TagsCaches return tags by tids from cache.
func (d *Dao) TagsCaches(c context.Context, tids []int64) (res []*model.Tag, missed []int64, err error) {
	var (
		keys    = make([]string, 0, len(tids))
		keysMap = make(map[string]int64, len(tids))
	)
	for _, tid := range tids {
		key := keyTag(tid)
		keys = append(keys, key)
		keysMap[key] = tid
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	items, err := conn.GetMulti(keys)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.GetMulti(%v) error(%v)", keys, err)
		}
		return
	}
	for _, item := range items {
		t := &model.Tag{}
		if err = conn.Scan(item, t); err != nil {
			log.Error("conn.Scan(%s) error(%v)", item.Value, err)
			return
		}
		if t != nil {
			delete(keysMap, item.Key)
			res = append(res, t)
		}
	}
	for _, tid := range keysMap {
		missed = append(missed, tid)
	}
	return
}

// TagMapCaches return tag map by tids from caches.
func (d *Dao) TagMapCaches(c context.Context, tids []int64) (res map[int64]*model.Tag, missed []int64, err error) {
	var (
		keys    = make([]string, 0, len(tids))
		keysMap = make(map[string]int64, len(tids))
	)
	for _, tid := range tids {
		key := keyTag(tid)
		keys = append(keys, key)
		keysMap[key] = tid
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	items, err := conn.GetMulti(keys)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.GetMulti(%v) error(%v)", keys, err)
		}
		return
	}
	res = make(map[int64]*model.Tag, len(items))
	for _, item := range items {
		t := &model.Tag{}
		if err = conn.Scan(item, t); err != nil {
			log.Error("conn.Scan(%s) error(%v)", item.Value, err)
			return
		}
		if t != nil {
			delete(keysMap, item.Key)
			res[t.ID] = t
		}
	}
	for _, tid := range keysMap {
		missed = append(missed, tid)
	}
	return
}

// TagCacheByName return tag by name from cache.
func (d *Dao) TagCacheByName(c context.Context, name string) (res *model.Tag, err error) {
	var (
		key  = keyName(name)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	item, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Get(%v) error(%v)", key, err)
		}
		return
	}
	res = &model.Tag{}
	if err = conn.Scan(item, res); err != nil {
		log.Error("mc.Scan(%s) error(%v)", item.Value, err)
	}
	return
}

// TagCachesByNames return tag caches by names from cache.
func (d *Dao) TagCachesByNames(c context.Context, names []string) (res []*model.Tag, missed []string, err error) {
	var (
		keys    []string
		keysMap = make(map[string]string)
	)
	for _, name := range names {
		key := keyName(name)
		if _, ok := keysMap[key]; !ok {
			keys = append(keys, key)
			keysMap[key] = name
		}
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	items, err := conn.GetMulti(keys)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.GetMulti(%v) error(%v)", keys, err)
		}
		return
	}
	for _, item := range items {
		t := &model.Tag{}
		if err = conn.Scan(item, t); err != nil {
			log.Error("mc.Scan(%s) error(%v)", item.Value, err)
			return
		}
		if t != nil {
			res = append(res, t)
			delete(keysMap, item.Key)
		}
	}
	for _, name := range keysMap {
		missed = append(missed, strings.Replace(name, _spaceReplace, " ", -1))
	}
	return
}

// AddTagCache add a tag to cache.
func (d *Dao) AddTagCache(c context.Context, tag *model.Tag) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &memcache.Item{
		Key:        keyTag(tag.ID),
		Object:     tag,
		Flags:      memcache.FlagProtobuf,
		Expiration: d.tagExpire,
	}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%v) error(%v)", item, err)
	}
	item = &memcache.Item{
		Key:        keyName(tag.Name),
		Object:     tag,
		Flags:      memcache.FlagProtobuf,
		Expiration: d.tagExpire,
	}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%v) error(%v)", item, err)
	}
	return
}

// AddTagsCache add tags to cache.
func (d *Dao) AddTagsCache(c context.Context, tags []*model.Tag) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	for _, tag := range tags {
		item := &memcache.Item{
			Key:        keyTag(tag.ID),
			Object:     tag,
			Flags:      memcache.FlagProtobuf,
			Expiration: d.tagExpire,
		}
		if err = conn.Set(item); err != nil {
			log.Error("conn.Set(%s) error(%v)", item.Key, err)
			return
		}
		item = &memcache.Item{
			Key:        keyName(tag.Name),
			Object:     tag,
			Flags:      memcache.FlagProtobuf,
			Expiration: d.tagExpire,
		}
		if err = conn.Set(item); err != nil {
			log.Error("conn.Set(%s) (byteKey:%v) error(%v)", item.Key, []byte(item.Key), err)
			return
		}
	}
	return
}

// DelTagCache .
func (d *Dao) DelTagCache(c context.Context, tid int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	var tag *model.Tag
	if tag, err = d.TagCache(c, tid); err != nil {
		return
	}
	if err = conn.Delete(keyTag(tid)); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Set(%d) error(%v)", tid, err)
		}
	}
	if tag != nil {
		if err = conn.Delete(keyName(tag.Name)); err != nil {
			if err == memcache.ErrNotFound {
				err = nil
			} else {
				log.Error("conn.Set(%s) error(%v)", keyName(tag.Name), err)
			}
		}
	}
	return
}

// AddCountCache .
func (d *Dao) AddCountCache(c context.Context, count *model.Count) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &memcache.Item{
		Key:        keyCount(count.Tid),
		Object:     count,
		Flags:      memcache.FlagProtobuf,
		Expiration: d.resExpire,
	}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%v) error(%v)", item, err)
	}
	return
}

// AddCountsCache .
func (d *Dao) AddCountsCache(c context.Context, counts []*model.Count) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	for _, count := range counts {
		item := &memcache.Item{
			Key:        keyCount(count.Tid),
			Object:     count,
			Flags:      memcache.FlagProtobuf,
			Expiration: d.resExpire,
		}
		if err = conn.Set(item); err != nil {
			log.Error("conn.Set(%v) error(%v)", item, err)
			return
		}
	}
	return
}

// CountCache return tag caches by names from cache.
func (d *Dao) CountCache(c context.Context, tid int64) (res *model.Count, err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err := conn.Get(keyCount(tid))
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Get(%v) error(%v)", tid, err)
		}
		return
	}
	res = &model.Count{}
	if err = conn.Scan(item, res); err != nil {
		log.Error("mc.Scan(%s) error(%v)", item.Value, err)
	}
	return
}

// CountMapCache return tag caches by names from cache.
func (d *Dao) CountMapCache(c context.Context, tids []int64) (res map[int64]*model.Count, missed []int64, err error) {
	var (
		keys    []string
		keysMap = make(map[int64]string)
	)
	for _, tid := range tids {
		key := keyCount(tid)
		if _, ok := keysMap[tid]; !ok {
			keys = append(keys, key)
			keysMap[tid] = key
		}
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	items, err := conn.GetMulti(keys)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.GetMulti(%v) error(%v)", keys, err)
		}
		return
	}
	res = make(map[int64]*model.Count)
	for _, item := range items {
		tc := &model.Count{}
		if err = conn.Scan(item, tc); err != nil {
			log.Error("mc.Scan(%s) error(%v)", item.Value, err)
			return
		}
		if c != nil {
			res[tc.Tid] = tc
			delete(keysMap, tc.Tid)
		}
	}
	for k := range keysMap {
		missed = append(missed, k)
	}
	return
}

// DelCountCache .
func (d *Dao) DelCountCache(c context.Context, tid int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(keyCount(tid)); err != nil {
		if err == memcache.ErrNotFound {
			return nil
		}
		log.Error("conn.Set(%d) error(%v)", tid, err)
	}
	return
}

// DelCountsCache .
func (d *Dao) DelCountsCache(c context.Context, tids []int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	for _, tid := range tids {
		if err = conn.Delete(keyCount(tid)); err != nil {
			if err == memcache.ErrNotFound {
				return nil
			}
			log.Error("conn.Set(%d) error(%v)", tid, err)
		}
	}
	return
}

// TagNamesCache .
func (d *Dao) TagNamesCache(c context.Context, oid int64, typ int8) (names []string, err error) {
	key := keyNames(oid, typ)
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.GeGettMulti(%v) error(%v)", key, err)
		}
		return
	}
	if err = conn.Scan(item, &names); err != nil {
		log.Error("conn.Scan(%s) error(%v)", item.Value, err)
	}
	return
}
