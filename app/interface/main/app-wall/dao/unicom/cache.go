package unicom

import (
	"context"
	"fmt"

	"go-common/app/interface/main/app-wall/model/unicom"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_prefix          = "unicoms_user_%v"
	_userbindkey     = "unicoms_user_bind_%d"
	_userpackkey     = "unicom_user_pack_%d"
	_userflowkey     = "unicom_user_flow_%v"
	_userflowlistkey = "unicom_user_flow_list"
	_userflowWait    = "u_flowwait_%d"
)

func keyUnicom(usermob string) string {
	return fmt.Sprintf(_prefix, usermob)
}

func keyUserBind(mid int64) string {
	return fmt.Sprintf(_userbindkey, mid)
}

func keyUserPack(id int64) string {
	return fmt.Sprintf(_userpackkey, id)
}

func keyUserFlow(key string) string {
	return fmt.Sprintf(_userflowkey, key)
}

func keyUserflowWait(phone int) string {
	return fmt.Sprintf(_userflowWait, phone)
}

// AddUnicomCache
func (d *Dao) AddUnicomCache(c context.Context, usermob string, u []*unicom.Unicom) (err error) {
	var (
		key  = keyUnicom(usermob)
		conn = d.mc.Get(c)
	)
	if err = conn.Set(&memcache.Item{Key: key, Object: u, Flags: memcache.FlagJSON, Expiration: d.expire}); err != nil {
		log.Error("addUnicomCache d.mc.Set(%s,%v) error(%v)", key, u, err)
	}
	conn.Close()
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

// UpdateUnicomCache
func (d *Dao) UpdateUnicomCache(c context.Context, usermob string, u *unicom.Unicom) (err error) {
	var (
		us      []*unicom.Unicom
		unicoms []*unicom.Unicom
		uspid   = map[int]struct{}{}
	)
	if u.Spid == 979 && u.TypeInt == 1 {
		return d.DeleteUnicomCache(c, usermob)
	}
	if us, err = d.UnicomCache(c, usermob); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("d.UnicomCache error(%v)", err)
		return
	}
	if len(us) > 0 {
		for _, um := range us {
			if um.Spid == 979 && um.TypeInt == 1 {
				return d.DeleteUnicomCache(c, usermob)
			}
			tmp := &unicom.Unicom{}
			*tmp = *um
			if tmp.Spid == u.Spid {
				tmp = u
				uspid[u.Spid] = struct{}{}
			}
			unicoms = append(unicoms, tmp)
		}
		if _, ok := uspid[u.Spid]; !ok {
			unicoms = append(unicoms, u)
		}
		if err = d.AddUnicomCache(c, usermob, unicoms); err != nil {
			log.Error("d.AddUnicomCache error(%v)", err)
			return
		}
	}
	return
}

// DeleteUnicomCache
func (d *Dao) DeleteUnicomCache(c context.Context, usermob string) (err error) {
	var (
		key  = keyUnicom(usermob)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("unicomCache MemchDB.Delete(%s) error(%v)", key, err)
		return
	}
	return
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
		log.Error("UserBindCache MemchDB.Get(%s) error(%v)", key, err)
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

// UserPackCache user packs
func (d *Dao) UserPackCache(c context.Context, id int64) (res *unicom.UserPack, err error) {
	var (
		key  = keyUserPack(id)
		conn = d.mc.Get(c)
		r    *memcache.Item
	)
	defer conn.Close()
	if r, err = conn.Get(key); err != nil {
		log.Error("UserBindCache MemchDB.Get(%s) error(%v)", key, err)
		return
	}
	if err = conn.Scan(r, &res); err != nil {
		log.Error("r.Scan(%s) error(%v)", r.Value, err)
	}
	return
}

// AddUserPackCache add user pack cache
func (d *Dao) AddUserPackCache(c context.Context, id int64, u *unicom.UserPack) (err error) {
	var (
		key  = keyUserPack(id)
		conn = d.mc.Get(c)
	)
	if err = conn.Set(&memcache.Item{Key: key, Object: u, Flags: memcache.FlagJSON, Expiration: d.flowKeyExpired}); err != nil {
		log.Error("AddUserPackCache d.mc.Set(%s,%v) error(%v)", key, u, err)
	}
	conn.Close()
	return
}

// DeleteUserPackCache delete user pack cache
func (d *Dao) DeleteUserPackCache(c context.Context, id int64) (err error) {
	var (
		key  = keyUserPack(id)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("DeleteUserPackCache MemchDB.Delete(%s) error(%v)", key, err)
		return
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

//UserFlowWaitCache unicom flow wait
func (d *Dao) UserFlowWaitCache(c context.Context, phone int) (err error) {
	var (
		key  = keyUserflowWait(phone)
		conn = d.mc.Get(c)
		r    *memcache.Item
		res  struct{}
	)
	defer conn.Close()
	if r, err = conn.Get(key); err != nil {
		log.Error("UserFlowWaitCache MemchDB.Get(%s) error(%v)", key, err)
		return
	}
	if err = conn.Scan(r, &res); err != nil {
		log.Error("r.Scan(%s) error(%v)", r.Value, err)
	}
	return
}

// AddUserFlowWaitCache add user flow wait
func (d *Dao) AddUserFlowWaitCache(c context.Context, phone int) (err error) {
	var (
		key  = keyUserflowWait(phone)
		conn = d.mc.Get(c)
	)
	if err = conn.Set(&memcache.Item{Key: key, Object: struct{}{}, Flags: memcache.FlagJSON, Expiration: d.flowWait}); err != nil {
		log.Error("AddUserFlowWaitCache d.mc.Set(%s) error(%v)", key, err)
	}
	conn.Close()
	return
}
