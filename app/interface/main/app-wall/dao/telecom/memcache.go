package telecom

import (
	"context"
	"fmt"

	"go-common/app/interface/main/app-wall/model/telecom"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_prefix     = "telecom_%d"
	_preorderid = "telecom_orderid_%d"
)

func keyTelecom(userphone int) string {
	return fmt.Sprintf(_prefix, userphone)
}

func keyTelecomOrderID(orderID int64) string {
	return fmt.Sprintf(_preorderid, orderID)
}

// AddTelecomCache
func (d *Dao) AddTelecomCache(c context.Context, userphone int, u *telecom.OrderInfo) (err error) {
	var (
		key  = keyTelecom(userphone)
		conn = d.mc.Get(c)
	)
	if err = conn.Set(&memcache.Item{Key: key, Object: u, Flags: memcache.FlagJSON, Expiration: d.expire}); err != nil {
		log.Error("addTelecomCache d.mc.Set(%s,%v) error(%v)", key, u, err)
	}
	conn.Close()
	return
}

// TelecomCache
func (d *Dao) TelecomCache(c context.Context, userphone int) (u *telecom.OrderInfo, err error) {
	var (
		key  = keyTelecom(userphone)
		conn = d.mc.Get(c)
		r    *memcache.Item
	)
	defer conn.Close()
	if r, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("telecomCache MemchDB.Get(%s) error(%v)", key, err)
		return
	}
	if err = conn.Scan(r, &u); err != nil {
		log.Error("r.Scan(%s) error(%v)", r.Value, err)
	}
	return
}

// AddTelecomOrderIDCache
func (d *Dao) AddTelecomOrderIDCache(c context.Context, orderID int64, u *telecom.OrderInfo) (err error) {
	var (
		key  = keyTelecomOrderID(orderID)
		conn = d.mc.Get(c)
	)
	if err = conn.Set(&memcache.Item{Key: key, Object: u, Flags: memcache.FlagJSON, Expiration: d.expire}); err != nil {
		log.Error("addTelecomCache d.mc.Set(%s,%v) error(%v)", key, u, err)
	}
	conn.Close()
	return
}

// TelecomOrderIDCache
func (d *Dao) TelecomOrderIDCache(c context.Context, orderID int64) (u *telecom.OrderInfo, err error) {
	var (
		key  = keyTelecomOrderID(orderID)
		conn = d.mc.Get(c)
		r    *memcache.Item
	)
	defer conn.Close()
	if r, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("telecomCache MemchDB.Get(%s) error(%v)", key, err)
		return
	}
	if err = conn.Scan(r, &u); err != nil {
		log.Error("r.Scan(%s) error(%v)", r.Value, err)
	}
	return
}
