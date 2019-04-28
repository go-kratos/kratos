package reply

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go-common/app/interface/main/reply/conf"
	model "go-common/app/interface/main/reply/model/reply"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_prefixSub      = "s_"
	_prefixRp       = "r_"
	_prefixAdminTop = "at_"
	_prefixUpperTop = "ut_"
	_prefixConfig   = "c_%d_%d_%d"
	_prefixCaptcha  = "pc_%d"
)

// MemcacheDao memcache dao.
type MemcacheDao struct {
	mc *memcache.Pool

	expire      int32
	emptyExpire int32
}

// NewMemcacheDao new a memcache dao and return.
func NewMemcacheDao(c *conf.Memcache) *MemcacheDao {
	m := &MemcacheDao{
		mc:          memcache.NewPool(c.Config),
		expire:      int32(time.Duration(c.Expire) / time.Second),
		emptyExpire: int32(time.Duration(c.EmptyExpire) / time.Second),
	}
	return m
}

func keyCaptcha(mid int64) string { return fmt.Sprintf(_prefixCaptcha, mid) }

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

func keyConfig(oid int64, typ, category int8) string {
	return fmt.Sprintf(_prefixConfig, oid, typ, category)
}

// Ping check connection success.
func (dao *MemcacheDao) Ping(c context.Context) (err error) {
	conn := dao.mc.Get(c)
	item := memcache.Item{Key: "ping", Value: []byte{1}, Expiration: dao.expire}
	err = conn.Set(&item)
	conn.Close()
	return
}

// CaptchaToken CaptchaToken
func (dao *MemcacheDao) CaptchaToken(c context.Context, mid int64) (string, error) {
	conn := dao.mc.Get(c)
	defer conn.Close()

	item, err := conn.Get(keyCaptcha(mid))
	if err == memcache.ErrNotFound {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	var token string
	if err = conn.Scan(item, &token); err != nil {
		log.Error("conn.Scan(%s) error(%v)", item.Value, err)
		return "", err
	}
	return token, nil
}

// SetCaptchaToken SetCaptchaToken
func (dao *MemcacheDao) SetCaptchaToken(c context.Context, mid int64, token string) error {
	conn := dao.mc.Get(c)
	defer conn.Close()

	return conn.Set(&memcache.Item{
		Key:        keyCaptcha(mid),
		Value:      []byte(token),
		Expiration: int32(time.Minute * 5 / time.Second),
	})
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

// GetMultiSubject get subject from memcache.
func (dao *MemcacheDao) GetMultiSubject(c context.Context, oids []int64, tp int8) (res map[int64]*model.Subject, missed []int64, err error) {
	var (
		keys     = make([]string, len(oids))
		missKeys = make(map[string]int64, len(oids))
	)
	for i, oid := range oids {
		key := keySub(oid, tp)
		keys[i] = key
		missKeys[key] = oid
	}
	conn := dao.mc.Get(c)
	defer conn.Close()
	items, err := conn.GetMulti(keys)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		}
		missed = oids
		return
	}
	res = make(map[int64]*model.Subject, len(items))
	for _, item := range items {
		sub := new(model.Subject)
		if err = conn.Scan(item, sub); err != nil {
			log.Error("conn.Scan(%s) error(%v)", item.Value, err)
			continue
		}
		res[sub.Oid] = sub
		delete(missKeys, item.Key)
	}
	missed = make([]int64, 0, len(missKeys))
	for _, oid := range missKeys {
		missed = append(missed, oid)
	}
	return
}

// DeleteSubject delete subject from memcache.
func (dao *MemcacheDao) DeleteSubject(c context.Context, oid int64, tp int8) (err error) {
	key := keySub(oid, tp)
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

// AddSubject add subject into memcache.
func (dao *MemcacheDao) AddSubject(c context.Context, subs ...*model.Subject) (err error) {
	if len(subs) == 0 {
		return
	}
	conn := dao.mc.Get(c)
	for _, sub := range subs {
		exp := dao.expire
		if sub.ID == -1 {
			exp = dao.emptyExpire
		}
		key := keySub(sub.Oid, sub.Type)
		item := &memcache.Item{Key: key, Object: sub, Expiration: exp, Flags: memcache.FlagJSON}
		if err = conn.Set(item); err != nil {
			log.Error("conn.Set(%s,%v) error(%v)", key, sub, err)
		}
	}
	conn.Close()
	return
}

// AddReply add reply into memcache.
func (dao *MemcacheDao) AddReply(c context.Context, rs ...*model.Reply) (err error) {
	if len(rs) == 0 {
		return
	}
	conn := dao.mc.Get(c)
	for _, r := range rs {
		if r == nil {
			continue
		}
		key := keyRp(r.RpID)
		item := &memcache.Item{Key: key, Object: r, Expiration: dao.expire, Flags: memcache.FlagJSON}
		if err = conn.Set(item); err != nil {
			log.Error("conn.Set(%s,%v) error(%v)", key, r, err)
		}
	}
	conn.Close()
	return
}

// AddTop add top reply into memcache.
func (dao *MemcacheDao) AddTop(c context.Context, oid int64, tp int8, rp *model.Reply) (err error) {
	if rp == nil {
		return
	}
	var key string
	if rp.AttrVal(model.ReplyAttrAdminTop) == 1 {
		key = keyAdminTop(oid, tp)
	} else if rp.AttrVal(model.ReplyAttrUpperTop) == 1 {
		key = keyUpperTop(oid, tp)
	} else {
		return
	}
	conn := dao.mc.Get(c)
	defer conn.Close()
	item := &memcache.Item{Key: key, Object: rp, Expiration: dao.expire, Flags: memcache.FlagJSON}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%s,%v) error(%v)", key, rp, err)
	}
	return
}

// DeleteReply delete reply from memcache.
func (dao *MemcacheDao) DeleteReply(c context.Context, rpID int64) (err error) {
	key := keyRp(rpID)
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

// GetMultiReply multi get replies from memcache.
func (dao *MemcacheDao) GetMultiReply(c context.Context, rpIDs []int64) (rpMap map[int64]*model.Reply, missed []int64, err error) {
	if len(rpIDs) == 0 {
		return
	}
	rpMap = make(map[int64]*model.Reply, len(rpIDs))
	keys := make([]string, len(rpIDs))
	mm := make(map[string]int64, len(rpIDs))
	for i, rpID := range rpIDs {
		key := keyRp(rpID)
		keys[i] = key
		mm[key] = rpID
	}
	conn := dao.mc.Get(c)
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

// GetReplyConfig get reply configuration from memocache by oid and type value
func (dao *MemcacheDao) GetReplyConfig(c context.Context, oid int64, typ, category int8) (config *model.Config, err error) {
	key := keyConfig(oid, typ, 1)
	conn := dao.mc.Get(c)
	defer conn.Close()
	item, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		}
		return
	}
	config = new(model.Config)
	if err = conn.Scan(item, config); err != nil {
		log.Error("conn.Scan(%s) error(%v)", item.Value, err)
		config = nil
	}
	return
}

// AddReplyConfigCache add/update reply configuration cache from memcache
func (dao *MemcacheDao) AddReplyConfigCache(c context.Context, m *model.Config) (err error) {
	key := keyConfig(m.Oid, m.Type, m.Category)
	conn := dao.mc.Get(c)
	item := &memcache.Item{Key: key, Object: m, Expiration: dao.expire, Flags: memcache.FlagJSON}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%s,%v) error(%v)", key, m, err)
	}
	conn.Close()
	return
}
