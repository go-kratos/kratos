package mc

import (
	"context"
	"encoding/json"
	"fmt"

	"go-common/app/admin/main/aegis/model/common"
	gmc "go-common/library/cache/memcache"
	"go-common/library/log"
)

// ConsumerOn 登入
func (d *Dao) ConsumerOn(c context.Context, opt *common.BaseOptions) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := mcKey(opt)

	if err = conn.Set(&gmc.Item{Key: key, Value: []byte{1}, Expiration: d.c.Consumer.OnExp}); err != nil {
		log.Error("conn.Set error(%v)", err)
	}

	return
}

// ConsumerOff .
func (d *Dao) ConsumerOff(c context.Context, opt *common.BaseOptions) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := mcKey(opt)
	if err = conn.Delete(key); err != nil {
		if err != gmc.ErrNotFound {
			log.Error("conn.Delete error(%v)", err)
		}
		err = nil
	}
	return
}

// IsConsumerOn .
func (d *Dao) IsConsumerOn(c context.Context, opt *common.BaseOptions) (isOn bool, err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := mcKey(opt)
	if _, err = conn.Get(key); err != nil {
		if err == gmc.ErrNotFound {
			err = nil
		} else {
			log.Error("IsConsumerOn error(%v)", err)
		}
		return
	}
	isOn = true
	if err = conn.Set(&gmc.Item{Key: key, Value: []byte{1}, Expiration: d.c.Consumer.OnExp}); err != nil {
		log.Error("conn.Set error(%v)", err)
	}

	return
}

func mcKey(opt *common.BaseOptions) string {
	return fmt.Sprintf("aegis%d_%d_%d", opt.BusinessID, opt.FlowID, opt.UID)
}

func roleKey(opt *common.BaseOptions) string {
	return fmt.Sprintf("aegis_role%d_%d_%d", opt.BusinessID, opt.FlowID, opt.UID)
}

// Role .
type Role struct {
	Role  int8   `json:"role"`
	Uname string `json:"uname"`
}

// GetRole TODO 目前缓存组长组员,以后扩展到存管理员
func (d *Dao) GetRole(c context.Context, opt *common.BaseOptions) (role int8, uname string, err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := roleKey(opt)
	var item *gmc.Item
	if item, err = conn.Get(key); err != nil {
		if err == gmc.ErrNotFound {
			log.Info("GetRole opt(%+v) miss", opt)
			err = nil
		} else {
			log.Error("GetRole opt(%+v) error(%v)", opt, err)
		}
		return
	}

	rt := new(Role)
	if err = json.Unmarshal(item.Value, rt); err != nil {
		log.Error("GetRole value(%s) error(%v)", string(item.Value), err)
		return
	}
	role = rt.Role
	uname = rt.Uname
	return
}

// SetRole .
func (d *Dao) SetRole(c context.Context, opt *common.BaseOptions, role int8) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := roleKey(opt)
	val := &Role{
		Role:  role,
		Uname: opt.Uname,
	}

	var roleb []byte
	if roleb, err = json.Marshal(val); err != nil {
		log.Error("SetRole error(%v)", err)
		return
	}
	if err = conn.Set(&gmc.Item{Key: key, Value: roleb, Expiration: d.c.Consumer.RoleExp}); err != nil {
		log.Error("conn.Set error(%v)", err)
	}
	return
}

//DelRole .
func (d *Dao) DelRole(c context.Context, bizid, flowid int64, uids []int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()

	for _, uid := range uids {
		key := roleKey(&common.BaseOptions{
			BusinessID: bizid,
			FlowID:     flowid,
			UID:        uid})
		conn.Delete(key)
	}
	return
}
