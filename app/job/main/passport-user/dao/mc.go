package dao

import (
	"context"
	"fmt"

	"go-common/app/job/main/passport-user/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

func ubKey(mid int64) string {
	return fmt.Sprintf("ub_%d", mid)
}

func utKey(mid int64) string {
	return fmt.Sprintf("ut_%d", mid)
}

func ueKey(mid int64) string {
	return fmt.Sprintf("ue_%d", mid)
}

func uroKey(mid int64) string {
	return fmt.Sprintf("uro_%d", mid)
}

func usqKey(mid int64) string {
	return fmt.Sprintf("usq_%d", mid)
}

func qqKey(mid int64) string {
	return fmt.Sprintf("utb_qq_%d", mid)
}

func sinaKey(mid int64) string {
	return fmt.Sprintf("utb_sina_%d", mid)
}

// SetUserBaseCache set user base to cache
func (d *Dao) SetUserBaseCache(c context.Context, ub *model.UserBase) (err error) {
	key := ubKey(ub.Mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &memcache.Item{Key: key, Object: ub.ConvertToProto(), Flags: memcache.FlagProtobuf, Expiration: d.mcExpire}
	if err = conn.Set(item); err != nil {
		log.Error("fail to set user base to mc, key(%s) expire(%d) error(%+v)", key, d.mcExpire, err)
	}
	return
}

// SetUserTelCache set user tel to cache
func (d *Dao) SetUserTelCache(c context.Context, ut *model.UserTel) (err error) {
	key := utKey(ut.Mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &memcache.Item{Key: key, Object: ut.ConvertToProto(), Flags: memcache.FlagProtobuf, Expiration: d.mcExpire}
	if err = conn.Set(item); err != nil {
		log.Error("fail to set user tel to mc, key(%s) expire(%d) error(%+v)", key, d.mcExpire, err)
	}
	return
}

// SetUserEmailCache set user email to cache
func (d *Dao) SetUserEmailCache(c context.Context, ue *model.UserEmail) (err error) {
	key := ueKey(ue.Mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &memcache.Item{Key: key, Object: ue.ConvertToProto(), Flags: memcache.FlagProtobuf, Expiration: d.mcExpire}
	if err = conn.Set(item); err != nil {
		log.Error("fail to set user email to mc, key(%s) expire(%d) error(%+v)", key, d.mcExpire, err)
	}
	return
}

// SetUserRegOriginCache set user reg origin to cache
func (d *Dao) SetUserRegOriginCache(c context.Context, uro *model.UserRegOrigin) (err error) {
	key := uroKey(uro.Mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &memcache.Item{Key: key, Object: uro.ConvertToProto(), Flags: memcache.FlagProtobuf, Expiration: d.mcExpire}
	if err = conn.Set(item); err != nil {
		log.Error("fail to set user reg origin to mc, key(%s) expire(%d) error(%+v)", key, d.mcExpire, err)
	}
	return
}

// SetUserSafeQuestionCache set user safe question to cache
func (d *Dao) SetUserSafeQuestionCache(c context.Context, usq *model.UserSafeQuestion) (err error) {
	key := usqKey(usq.Mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &memcache.Item{Key: key, Object: usq.ConvertToProto(), Flags: memcache.FlagProtobuf, Expiration: d.mcExpire}
	if err = conn.Set(item); err != nil {
		log.Error("fail to set user safe question to mc, key(%s) expire(%d) error(%+v)", key, d.mcExpire, err)
	}
	return
}

// SetUserThirdBindQQCache set user third bind qq to cache
func (d *Dao) SetUserThirdBindQQCache(c context.Context, utb *model.UserThirdBind) (err error) {
	key := qqKey(utb.Mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &memcache.Item{Key: key, Object: utb.ConvertToProto(), Flags: memcache.FlagProtobuf, Expiration: d.mcExpire}
	if err = conn.Set(item); err != nil {
		log.Error("fail to set user third bind qq to mc, key(%s) expire(%d) error(%+v)", key, d.mcExpire, err)
	}
	return
}

// SetUserThirdBindSinaCache set user third bind sina to cache
func (d *Dao) SetUserThirdBindSinaCache(c context.Context, utb *model.UserThirdBind) (err error) {
	key := sinaKey(utb.Mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &memcache.Item{Key: key, Object: utb.ConvertToProto(), Flags: memcache.FlagProtobuf, Expiration: d.mcExpire}
	if err = conn.Set(item); err != nil {
		log.Error("fail to set user third bind sina to mc, key(%s) expire(%d) error(%+v)", key, d.mcExpire, err)
	}
	return
}

// DelUserBaseCache del user base cache
func (d *Dao) DelUserBaseCache(c context.Context, mid int64) (err error) {
	key := ubKey(mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("fail to del user base cache, key(%s) error(%+v)", key, err)
	}
	return
}

// DelUserTelCache del user tel cache
func (d *Dao) DelUserTelCache(c context.Context, mid int64) (err error) {
	key := utKey(mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("fail to del user tel cache, key(%s) error(%+v)", key, err)
	}
	return
}

// DelUserEmailCache del user email cache
func (d *Dao) DelUserEmailCache(c context.Context, mid int64) (err error) {
	key := ueKey(mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("fail to del user email cache, key(%s) error(%+v)", key, err)
	}
	return
}

// DelUserThirdBindQQCache del user third bind qq cache
func (d *Dao) DelUserThirdBindQQCache(c context.Context, mid int64) (err error) {
	key := qqKey(mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("fail to del user third bind qq cache, key(%s) error(%+v)", key, err)
	}
	return
}

// DelUserThirdBindSinaCache del user third bind sina cache
func (d *Dao) DelUserThirdBindSinaCache(c context.Context, mid int64) (err error) {
	key := sinaKey(mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("fail to del user third bind sina cache, key(%s) error(%+v)", key, err)
	}
	return
}

// pingMC
func (d *Dao) pingMC(c context.Context) (err error) {
	conn := d.mc.Get(c)
	item := memcache.Item{Key: "ping", Value: []byte{1}, Expiration: d.mcExpire}
	err = conn.Set(&item)
	conn.Close()
	return
}
