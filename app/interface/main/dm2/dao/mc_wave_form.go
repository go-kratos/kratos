package dao

import (
	"context"
	"fmt"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_waveFormFmt = "wf_%d_%d"
)

func (d *Dao) waveFormKey(oid int64, tp int32) string {
	return fmt.Sprintf(_waveFormFmt, oid, tp)
}

// SetWaveFormCache .
func (d *Dao) SetWaveFormCache(c context.Context, waveForm *model.WaveForm) (err error) {
	var (
		key  = d.waveFormKey(waveForm.Oid, waveForm.Type)
		conn = d.subtitleMc.Get(c)
		item *memcache.Item
	)
	defer conn.Close()
	item = &memcache.Item{
		Key:        key,
		Object:     waveForm,
		Flags:      memcache.FlagJSON | memcache.FlagGzip,
		Expiration: d.subtitleMcExpire,
	}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%v) error(%v)", item, err)
	}
	return
}

// WaveFormCache .
func (d *Dao) WaveFormCache(c context.Context, oid int64, tp int32) (waveForm *model.WaveForm, err error) {
	var (
		item *memcache.Item
		conn = d.subtitleMc.Get(c)
		key  = d.waveFormKey(oid, tp)
	)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("memcache.Get(%s) error(%v)", key, err)
		}
		return
	}
	if err = conn.Scan(item, &waveForm); err != nil {
		log.Error("mc.Scan() error(%v)", err)
	}
	return
}

// DelWaveFormCache .
func (d *Dao) DelWaveFormCache(c context.Context, oid int64, tp int32) (err error) {
	var (
		key  = d.waveFormKey(oid, tp)
		conn = d.subtitleMc.Get(c)
	)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("memcache.Delete(%s) error(%v)", key, err)
		}
	}
	return
}
