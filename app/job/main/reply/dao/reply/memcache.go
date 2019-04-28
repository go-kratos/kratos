package reply

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go-common/app/job/main/reply/conf"
	model "go-common/app/job/main/reply/model/reply"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_prefixSub      = "s_"
	_prefixRp       = "r_"
	_prefixAdminTop = "at_"
	_prefixUpperTop = "ut_"
)

// MemcacheDao define memcache info
type MemcacheDao struct {
	mc        *memcache.Pool
	expire    int32
	topExpire int32
}

// NewMemcacheDao return a new mc dao
func NewMemcacheDao(c *conf.Memcache) *MemcacheDao {
	return &MemcacheDao{
		mc:        memcache.NewPool(c.Config),
		expire:    int32(time.Duration(c.Expire) / time.Second),
		topExpire: int32(time.Duration(c.TopExpire) / time.Second),
	}
}

func keyAdminTop(oid int64, tp int8) string {
	if oid > _oidOverflow {
		return fmt.Sprintf("%s_%d_%d", _prefixAdminTop, oid, tp)
	}
	return _prefixAdminTop + strconv.FormatInt((oid<<8)|int64(tp), 10)
}

func keyUpperTop(oid int64, tp int8) string {
	if oid > _oidOverflow {
		return fmt.Sprintf("%s_%d_%d", _prefixUpperTop, oid, tp)
	}
	return _prefixUpperTop + strconv.FormatInt((oid<<8)|int64(tp), 10)
}

func keySub(oid int64, tp int8) string {
	if oid > _oidOverflow {
		return fmt.Sprintf("%s_%d_%d", _prefixSub, oid, tp)
	}
	return _prefixSub + strconv.FormatInt((oid<<8)|int64(tp), 10)
}

func keyRp(rpID int64) string {
	return _prefixRp + strconv.FormatInt(rpID, 10)
}

// Ping check connection success.
func (dao *MemcacheDao) Ping(c context.Context) (err error) {
	conn := dao.mc.Get(c)
	item := memcache.Item{Key: "ping", Value: []byte{1}, Expiration: dao.expire}
	err = conn.Set(&item)
	conn.Close()
	return
}

// AddSubject add subject into memcache.
func (dao *MemcacheDao) AddSubject(c context.Context, subs ...*model.Subject) (err error) {
	if len(subs) == 0 {
		return
	}
	conn := dao.mc.Get(c)
	for _, sub := range subs {
		key := keySub(sub.Oid, sub.Type)
		item := &memcache.Item{Key: key, Object: sub, Expiration: dao.expire, Flags: memcache.FlagJSON}
		if err = conn.Set(item); err != nil {
			log.Error("conn.Set(%s,%v) error(%v)", key, sub, err)
		}
	}
	conn.Close()
	return
}

// GetSubject get subject from memcache.
func (dao *MemcacheDao) GetSubject(c context.Context, oid int64, tp int8) (sub *model.Subject, err error) {
	key := keySub(oid, tp)
	conn := dao.mc.Get(c)
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

// AddReply add reply into memcache.
func (dao *MemcacheDao) AddReply(c context.Context, rs ...*model.Reply) (err error) {
	if len(rs) == 0 {
		return
	}
	conn := dao.mc.Get(c)
	for _, r := range rs {
		key := keyRp(r.RpID)
		item := &memcache.Item{Key: key, Object: r, Expiration: dao.expire, Flags: memcache.FlagJSON}
		if err = conn.Set(item); err != nil {
			log.Error("conn.Set(%s,%v) error(%v)", key, r, err)
		}
	}
	conn.Close()
	return
}

// GetTop get subject top reply from memcache
func (dao *MemcacheDao) GetTop(c context.Context, oid int64, tp int8, top uint32) (rp *model.Reply, err error) {
	var key string
	if top == model.ReplyAttrUpperTop {
		key = keyUpperTop(oid, tp)
	} else if top == model.ReplyAttrAdminTop {
		key = keyAdminTop(oid, tp)
	} else {
		return
	}
	conn := dao.mc.Get(c)
	defer conn.Close()
	item, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		}
		return
	}
	rp = new(model.Reply)
	if err = conn.Scan(item, &rp); err != nil {
		log.Error("conn.Scan(%s) error(%v)", item.Value, err)
		rp = nil
	}
	return
}

// AddTop add top reply into memcache.
func (dao *MemcacheDao) AddTop(c context.Context, rp *model.Reply) (err error) {
	if rp == nil {
		return
	}
	var key string
	if rp.AttrVal(model.ReplyAttrAdminTop) == 1 {
		key = keyAdminTop(rp.Oid, rp.Type)
	} else if rp.AttrVal(model.ReplyAttrUpperTop) == 1 {
		key = keyUpperTop(rp.Oid, rp.Type)
	} else {
		return
	}
	conn := dao.mc.Get(c)
	defer conn.Close()
	item := &memcache.Item{Key: key, Object: rp, Expiration: dao.topExpire, Flags: memcache.FlagJSON}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%s,%v) error(%v)", key, rp, err)
	}
	return
}

// DeleteTop delete topreply from memcache.
func (dao *MemcacheDao) DeleteTop(c context.Context, rp *model.Reply, tp uint32) (err error) {
	var key string
	if tp == model.SubAttrAdminTop {
		key = keyAdminTop(rp.Oid, rp.Type)
	} else if tp == model.SubAttrUpperTop {
		key = keyUpperTop(rp.Oid, rp.Type)
	}
	conn := dao.mc.Get(c)
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

// DeleteSub delete sub from memcache.
func (dao *MemcacheDao) DeleteSub(c context.Context, oid int64, tp int8) (err error) {
	conn := dao.mc.Get(c)
	if err = conn.Delete(keySub(oid, tp)); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Delete error(%v)", err)
		}
	}
	conn.Close()
	return
}

// DeleteReply delete reply from memcache.
func (dao *MemcacheDao) DeleteReply(c context.Context, rpID int64) (err error) {
	conn := dao.mc.Get(c)
	if err = conn.Delete(keyRp(rpID)); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Delete error(%v)", err)
		}
	}
	conn.Close()
	return
}

// GetReply get reply from memcache.
func (dao *MemcacheDao) GetReply(c context.Context, rpID int64) (rp *model.Reply, err error) {
	key := keyRp(rpID)
	conn := dao.mc.Get(c)
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
		log.Error("conn.Scan(%s) error(%v)", item.Value, err)
		rp = nil
	}
	return
}
