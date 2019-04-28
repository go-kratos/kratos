package unicom

import (
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/app-wall/model/unicom"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_prefix          = "unicoms_user_%v"
	_userbindkey     = "unicoms_user_bind_%d"
	_userbindreceive = "unicom_pack_receives_%d"
	_userflowkey     = "unicom_user_flow_%v"
	_userflowlistkey = "unicom_user_flow_list"
)

func keyUserBind(mid int64) string {
	return fmt.Sprintf(_userbindkey, mid)
}

func keyUserBindReceive(mid int64) string {
	return fmt.Sprintf(_userbindreceive, mid)
}

func keyUnicom(usermob string) string {
	return fmt.Sprintf(_prefix, usermob)
}

func keyUserFlow(key string) string {
	return fmt.Sprintf(_userflowkey, key)
}

// UserBindCache user bind cache
func (d *Dao) UserBindCache(c context.Context, mid int64) (ub *unicom.UserBind, err error) {
	var (
		key  = keyUserBind(mid)
		conn = d.mc.Get(c)
		r    *memcache.Item
	)
	defer conn.Close()
	if r, err = conn.Get(key); err != nil {
		if err != memcache.ErrNotFound {
			log.Error("UserBindCache error(%v) or mid(%v)", err, mid)
		}
		return
	}
	if err = conn.Scan(r, &ub); err != nil {
		log.Error("r.Scan(%s) error(%v)", r.Value, err)
	}
	return
}

// AddUserBindCache add user user bind cache
func (d *Dao) AddUserBindCache(c context.Context, mid int64, ub *unicom.UserBind) (err error) {
	var (
		key  = keyUserBind(mid)
		conn = d.mc.Get(c)
	)
	if err = conn.Set(&memcache.Item{Key: key, Object: ub, Flags: memcache.FlagJSON, Expiration: 0}); err != nil {
		log.Error("AddUserBindCache d.mc.Set(%s,%v) error(%v)", key, ub, err)
	}
	conn.Close()
	return
}

// DeleteUserBindCache delete user bind cache
func (d *Dao) DeleteUserBindCache(c context.Context, mid int64) (err error) {
	var (
		key  = keyUserBind(mid)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("DeleteUserBindCache MemchDB.Delete(%s) error(%v)", key, err)
		return
	}
	return
}

// UserPackReceiveCache user pack cache
func (d *Dao) UserPackReceiveCache(c context.Context, mid int64) (count int, err error) {
	var (
		key  = keyUserBindReceive(mid)
		conn = d.mc.Get(c)
		r    *memcache.Item
	)
	defer conn.Close()
	if r, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("UserBindCache MemchDB.Get(%s) error(%v)", key, err)
		return
	}
	if err = conn.Scan(r, &count); err != nil {
		log.Error("r.Scan(%s) error(%v)", r.Value, err)
	}
	return
}

// AddUserPackReceiveCache add user pack cache
func (d *Dao) AddUserPackReceiveCache(c context.Context, mid int64, count int, now time.Time) (err error) {
	var (
		key            = keyUserBindReceive(mid)
		conn           = d.mc.Get(c)
		currenttimeSec = int32((now.Hour() * 60 * 60) + (now.Minute() * 60) + now.Second())
		overtime       int32
	)
	if overtime = d.expire - currenttimeSec; overtime < 1 {
		overtime = d.expire
	}
	log.Info("AddUserPackReceiveCache mid(%d) overtime(%d)", mid, overtime)
	if err = conn.Set(&memcache.Item{Key: key, Object: count, Flags: memcache.FlagJSON, Expiration: overtime}); err != nil {
		log.Error("AddUserPackReceiveCache d.mc.Set(%s,%v) error(%v)", key, count, err)
	}
	conn.Close()
	return
}

// DeleteUserPackReceiveCache delete user pack cache
func (d *Dao) DeleteUserPackReceiveCache(c context.Context, mid int64) (err error) {
	var (
		key  = keyUserBindReceive(mid)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("DeleteUserPackReceiveCache MemchDB.Delete(%s) error(%v)", key, err)
		return
	}
	return
}

// UnicomCache
func (d *Dao) UnicomCache(c context.Context, usermob string) (u []*unicom.Unicom, err error) {
	var (
		key  = keyUnicom(usermob)
		conn = d.mc.Get(c)
		r    *memcache.Item
	)
	defer conn.Close()
	if r, err = conn.Get(key); err != nil {
		log.Error("unicomCache MemchDB.Get(%s) error(%v)", key, err)
		return
	}
	if err = conn.Scan(r, &u); err != nil {
		log.Error("r.Scan(%s) error(%v)", r.Value, err)
	}
	return
}

//UserFlowCache unicom flow cache
func (d *Dao) UserFlowCache(c context.Context, keyStr string) (err error) {
	var (
		key  = keyUserFlow(keyStr)
		conn = d.mc.Get(c)
		r    *memcache.Item
		res  struct{}
	)
	defer conn.Close()
	if r, err = conn.Get(key); err != nil {
		log.Error("UserFlowCache MemchDB.Get(%s) error(%v)", key, err)
		return
	}
	if err = conn.Scan(r, &res); err != nil {
		log.Error("r.Scan(%s) error(%v)", r.Value, err)
	}
	return
}

// AddUserFlowCache add user pack cache
func (d *Dao) AddUserFlowCache(c context.Context, keyStr string) (err error) {
	var (
		key  = keyUserFlow(keyStr)
		conn = d.mc.Get(c)
	)
	if err = conn.Set(&memcache.Item{Key: key, Object: struct{}{}, Flags: memcache.FlagJSON, Expiration: d.flowKeyExpired}); err != nil {
		log.Error("AddUserPackCache d.mc.Set(%s) error(%v)", key, err)
	}
	conn.Close()
	return
}

// DeleteUserPackCache delete user pack cache
func (d *Dao) DeleteUserFlowCache(c context.Context, keyStr string) (err error) {
	var (
		key  = keyUserFlow(keyStr)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("DeleteUserFlowCache MemchDB.Delete(%s) error(%v)", key, err)
		return
	}
	return
}

//UserFlowListCache unicom flow cache
func (d *Dao) UserFlowListCache(c context.Context) (res map[string]*unicom.UnicomUserFlow, err error) {
	var (
		key  = _userflowlistkey
		conn = d.mc.Get(c)
		r    *memcache.Item
	)
	defer conn.Close()
	if r, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("UserFlowListCache MemchDB.Get(%s) error(%v)", key, err)
		return
	}
	if err = conn.Scan(r, &res); err != nil {
		log.Error("r.Scan(%s) error(%v)", r.Value, err)
	}
	return
}

// AddUserFlowListCache add user pack cache
func (d *Dao) AddUserFlowListCache(c context.Context, list map[string]*unicom.UnicomUserFlow) (err error) {
	var (
		key  = _userflowlistkey
		conn = d.mc.Get(c)
	)
	if err = conn.Set(&memcache.Item{Key: key, Object: list, Flags: memcache.FlagJSON, Expiration: d.flowKeyExpired}); err != nil {
		log.Error("AddUserFlowListCache d.mc.Set(%s,%v) error(%v)", key, list, err)
	}
	conn.Close()
	return
}
