package dao

import (
	"context"
	"fmt"
	"strconv"

	"go-common/app/admin/main/reply/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_prefixSub      = "s_"         // sub_oid<<8|type
	_prefixReply    = "r_"         // r_rpID
	_prefixAdminTop = "at_"        // at_rpID
	_prefixUpperTop = "ut_"        // ut_rpID
	_prefixConfig   = "c_%d_%d_%d" // oid_type_category

	_oidOverflow = 1 << 48
)

func keyReply(rpID int64) string {
	return _prefixReply + strconv.FormatInt(rpID, 10)
}

func keySubject(oid int64, typ int32) string {
	if oid > _oidOverflow {
		return fmt.Sprintf("%s_%d_%d", _prefixSub, oid, typ)
	}
	return _prefixSub + strconv.FormatInt((oid<<8)|int64(typ), 10)
}

func keyConfig(oid int64, typ, category int32) string {
	return fmt.Sprintf(_prefixConfig, oid, typ, category)
}

func keyAdminTop(oid int64, attr uint32) string {
	if oid > _oidOverflow {
		return fmt.Sprintf("%s_%d_%d", _prefixAdminTop, oid, attr)
	}
	return _prefixAdminTop + strconv.FormatInt((oid<<8)|int64(attr), 10)
}

func keyUpperTop(oid int64, attr uint32) string {
	if oid > _oidOverflow {
		return fmt.Sprintf("%s_%d_%d", _prefixUpperTop, oid, attr)
	}
	return _prefixUpperTop + strconv.FormatInt((oid<<8)|int64(attr), 10)
}

// PingMC check connection success.
func (d *Dao) pingMC(c context.Context) (err error) {
	conn := d.mc.Get(c)
	item := memcache.Item{Key: "ping", Value: []byte{1}, Expiration: d.mcExpire}
	err = conn.Set(&item)
	conn.Close()
	return
}

// SubjectCache get subject from memcache.
func (d *Dao) SubjectCache(c context.Context, oid int64, typ int32) (sub *model.Subject, err error) {
	key := keySubject(oid, typ)
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		}
		return
	}
	sub = new(model.Subject)
	if err = conn.Scan(item, sub); err != nil {
		log.Error("conn.Scan(%s) error(%v)", item.Value, err)
		sub = nil
	}
	return
}

// AddSubjectCache add subject into memcache.
func (d *Dao) AddSubjectCache(c context.Context, subs ...*model.Subject) (err error) {
	conn := d.mc.Get(c)
	for _, sub := range subs {
		key := keySubject(sub.Oid, sub.Type)
		item := &memcache.Item{Key: key, Object: sub, Expiration: d.mcExpire, Flags: memcache.FlagJSON}
		if err = conn.Set(item); err != nil {
			log.Error("conn.Set(%s,%v) error(%v)", key, sub, err)
		}
	}
	conn.Close()
	return
}

// DelSubjectCache delete subject from memcache.
func (d *Dao) DelSubjectCache(c context.Context, oid int64, typ int32) (err error) {
	key := keySubject(oid, typ)
	conn := d.mc.Get(c)
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		}
	}
	conn.Close()
	return
}

// ReplyCache get a reply from memcache.
func (d *Dao) ReplyCache(c context.Context, rpID int64) (rp *model.Reply, err error) {
	key := keyReply(rpID)
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		}
		return
	}
	rp = new(model.Reply)
	if err = conn.Scan(item, rp); err != nil {
		rp = nil
	}
	return
}

// RepliesCache multi get replies from memcache.
func (d *Dao) RepliesCache(c context.Context, rpIDs []int64) (rpMap map[int64]*model.Reply, missed []int64, err error) {
	rpMap = make(map[int64]*model.Reply, len(rpIDs))
	keys := make([]string, len(rpIDs))
	mm := make(map[string]int64, len(rpIDs))
	for i, rpID := range rpIDs {
		key := keyReply(rpID)
		keys[i] = key
		mm[key] = rpID
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	items, err := conn.GetMulti(keys)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		}
		return
	}
	for _, item := range items {
		rp := new(model.Reply)
		if err = conn.Scan(item, rp); err != nil {
			log.Error("conn.Scan(%s) error(%v)", item.Value, err)
			continue
		}
		rpMap[mm[item.Key]] = rp
		delete(mm, item.Key)
	}
	missed = make([]int64, 0, len(mm))
	for _, valIn := range mm {
		missed = append(missed, valIn)
	}
	return
}

// AddReplyCache add reply into memcache.
func (d *Dao) AddReplyCache(c context.Context, rps ...*model.Reply) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	for _, rp := range rps {
		item := &memcache.Item{
			Key:        keyReply(rp.ID),
			Object:     rp,
			Expiration: d.mcExpire,
			Flags:      memcache.FlagJSON,
		}
		if err = conn.Set(item); err != nil {
			return
		}
	}
	return
}

// DelReplyCache delete reply from memcache.
func (d *Dao) DelReplyCache(c context.Context, rpID int64) (err error) {
	conn := d.mc.Get(c)
	if err = conn.Delete(keyReply(rpID)); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		}
	}
	conn.Close()
	return
}

// DelConfigCache delete reply config from memcache.
func (d *Dao) DelConfigCache(c context.Context, oid int64, typ, category int32) (err error) {
	key := keyConfig(oid, typ, category)
	conn := d.mc.Get(c)
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		}
	}
	conn.Close()
	return
}

// TopCache get a reply from memcache.
func (d *Dao) TopCache(c context.Context, oid int64, attr uint32) (rp *model.Reply, err error) {
	var key string
	if attr == model.SubAttrTopAdmin {
		key = keyAdminTop(oid, attr)
	} else if attr == model.SubAttrTopUpper {
		key = keyUpperTop(oid, attr)
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		}
		return
	}
	rp = new(model.Reply)
	if err = conn.Scan(item, rp); err != nil {
		rp = nil
	}
	return
}

// DelTopCache delete topreply from memcache.
func (d *Dao) DelTopCache(c context.Context, oid int64, attr uint32) (err error) {
	var key string
	if attr == model.SubAttrTopAdmin {
		key = keyAdminTop(oid, attr)
	} else if attr == model.SubAttrTopUpper {
		key = keyUpperTop(oid, attr)
	}
	conn := d.mc.Get(c)
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Delete(%s) error(%v)", key, err)
		}
	}
	conn.Close()
	return
}

// AddTopCache add top reply into memcache.
func (d *Dao) AddTopCache(c context.Context, rp *model.Reply) (err error) {
	var key string
	if rp.AttrVal(model.AttrTopAdmin) == model.AttrYes {
		key = keyAdminTop(rp.Oid, model.AttrTopAdmin)
	} else if rp.AttrVal(model.AttrTopUpper) == model.AttrYes {
		key = keyUpperTop(rp.Oid, model.AttrTopUpper)
	} else {
		return
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &memcache.Item{Key: key, Object: rp, Expiration: d.mcExpire, Flags: memcache.FlagJSON}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%s,%v) error(%v)", key, rp, err)
	}
	return
}
