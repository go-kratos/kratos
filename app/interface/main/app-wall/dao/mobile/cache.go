package mobile

import (
	"context"
	"fmt"

	"go-common/app/interface/main/app-wall/model/mobile"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_prefix = "mobiles_user_%v"
)

func keyMobile(usermob string) string {
	return fmt.Sprintf(_prefix, usermob)
}

// AddMobileCache
func (d *Dao) AddMobileCache(c context.Context, usermob string, m []*mobile.Mobile) (err error) {
	var (
		key  = keyMobile(usermob)
		conn = d.mc.Get(c)
	)
	if err = conn.Set(&memcache.Item{Key: key, Object: m, Flags: memcache.FlagJSON, Expiration: d.expire}); err != nil {
		log.Error("addMobileCache d.mc.Set(%s,%v) error(%v)", key, m, err)
	}
	conn.Close()
	return
}

// MobileCache
func (d *Dao) MobileCache(c context.Context, usermob string) (m []*mobile.Mobile, err error) {
	var (
		key  = keyMobile(usermob)
		conn = d.mc.Get(c)
		r    *memcache.Item
	)
	defer conn.Close()
	if r, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("mobileCache MemchDB.Get(%s) error(%v)", key, err)
		return
	}
	if err = conn.Scan(r, &m); err != nil {
		log.Error("r.Scan(%s) error(%v)", r.Value, err)
	}
	return
}

// DeleteMobileCache
func (d *Dao) DeleteMobileCache(c context.Context, usermob string) (err error) {
	var (
		key  = keyMobile(usermob)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("mobileCache MemchDB.Delete(%s) error(%v)", key, err)
		return
	}
	return
}

// UpdateMobileCache
func (d *Dao) UpdateMobileCache(c context.Context, usermob string, m *mobile.Mobile) (err error) {
	var (
		ms         []*mobile.Mobile
		mobiles    []*mobile.Mobile
		uproductid = map[string]struct{}{}
	)
	if ms, err = d.MobileCache(c, usermob); err != nil && len(ms) > 0 {
		log.Error("d.MobileCache error(%v)", err)
		return
	}
	if len(ms) > 0 {
		for _, ml := range ms {
			tmp := &mobile.Mobile{}
			*tmp = *ml
			if tmp.Productid == m.Productid {
				tmp = m
				if m.Threshold == 0 {
					tmp.Threshold = ml.Threshold
				}
				uproductid[m.Productid] = struct{}{}
			}
			mobiles = append(mobiles, tmp)
		}
		if _, ok := uproductid[m.Productid]; !ok {
			mobiles = append(mobiles, m)
		}
		if err = d.AddMobileCache(c, usermob, mobiles); err != nil {
			log.Error("d.AddMobileCache error(%v)", err)
			return
		}
	}
	return
}

// UpdateMobileFlowCache
func (d *Dao) UpdateMobileFlowCache(c context.Context, usermob string, m *mobile.Mobile) (err error) {
	var (
		ms      []*mobile.Mobile
		mobiles []*mobile.Mobile
	)
	if ms, err = d.MobileCache(c, usermob); err != nil && len(ms) > 0 {
		log.Error("d.MobileCache error(%v)", err)
		return
	}
	if len(ms) > 0 {
		for _, ml := range ms {
			tmp := &mobile.Mobile{}
			*tmp = *ml
			if tmp.Productid == m.Productid {
				tmp.Threshold = m.Threshold
			}
			mobiles = append(mobiles, tmp)
		}
		if err = d.AddMobileCache(c, usermob, mobiles); err != nil {
			log.Error("d.AddMobileCache error(%v)", err)
			return
		}
	}
	return
}
