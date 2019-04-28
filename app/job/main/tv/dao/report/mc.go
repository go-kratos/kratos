package report

import (
	"context"
	"time"

	mdlpgc "go-common/app/job/main/tv/model/pgc"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_report = "_report"
	_style  = "style_label"
	_label  = "label_data"
)

// SetReportCache set report cache .
func (d *Dao) SetReportCache(c context.Context, val map[string]interface{}) (err error) {
	var (
		conn = d.mc.Get(c)
		key  = _report
	)
	defer conn.Close()
	item := &memcache.Item{
		Key:        key,
		Object:     val,
		Flags:      memcache.FlagJSON,
		Expiration: int32(time.Duration(d.conf.Report.Expire) / time.Second),
	}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%v) error(%v)", item, err)
	}
	return
}

// GetReportCache get report all data .
func (d *Dao) GetReportCache(c context.Context) (res map[string]interface{}, err error) {
	var (
		conn = d.mc.Get(c)
		key  = _report
		rp   *memcache.Item
	)
	res = make(map[string]interface{})
	defer conn.Close()
	if rp, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("mc.Get(%s) error(%v)", key, err)
		}
		return
	}
	if err = conn.Scan(rp, &res); err != nil {
		log.Error("conn.Scan error(%v)", err)
	}
	return
}

// SetStyleCache style show .
func (d *Dao) SetStyleCache(c context.Context, val map[int][]*mdlpgc.ParamStyle) (err error) {
	var (
		conn = d.mc.Get(c)
		key  = _style
	)
	defer conn.Close()
	item := &memcache.Item{
		Key:        key,
		Object:     val,
		Flags:      memcache.FlagJSON,
		Expiration: 0,
	}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%v) error(%v)", item, err)
	}
	return
}

// SetLabelCache label  show .
func (d *Dao) SetLabelCache(c context.Context, val map[int]map[string]int) (err error) {
	var (
		conn = d.mc.Get(c)
		key  = _label
	)
	defer conn.Close()
	item := &memcache.Item{
		Key:        key,
		Object:     val,
		Flags:      memcache.FlagJSON,
		Expiration: 0,
	}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%v) error(%v)", item, err)
	}
	return
}

//
// GetLabelCache get label all data .
func (d *Dao) GetLabelCache(c context.Context) (res map[int]map[string]int, err error) {
	var (
		conn = d.mc.Get(c)
		key  = _label
		rp   *memcache.Item
	)
	defer conn.Close()
	if rp, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("mc.Get(%s) error(%v)", key, err)
		}
		return
	}
	if err = conn.Scan(rp, &res); err != nil {
		log.Error("conn.Scan error(%v)", err)
	}
	return
}
