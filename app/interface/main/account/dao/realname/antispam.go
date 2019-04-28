package realname

import (
	"context"
	"fmt"

	"go-common/app/interface/main/account/conf"
	"go-common/library/cache/memcache"

	"github.com/pkg/errors"
)

func alipayAntispamKey(mid int64) string {
	return fmt.Sprintf("raa_%d", mid)
}

//AlipayAntispamValue 最低位为pass flag ，高位为计数
//计数:判断用户的申请次数
//flag:是否通过本次防刷验证（极验是否通过）
type AlipayAntispamValue int

// IncreaseCount add antispam count
func (a *AlipayAntispamValue) IncreaseCount() {
	*a = AlipayAntispamValue((a.Count()+1)<<1 + a.Flag())
}

// SetPass set is antispam verified (such as when geetest passed)
func (a *AlipayAntispamValue) SetPass(pass bool) {
	var flag int
	if pass {
		flag = 1
	}
	*a = AlipayAntispamValue(a.Count()<<1 + flag)
}

// Count return antispam hit count
func (a *AlipayAntispamValue) Count() int {
	return int(*a) >> 1
}

// Flag return antispam pass flag
func (a *AlipayAntispamValue) Flag() int {
	return int(*a) & 0x1
}

// Pass return is antispam passed (such as when geetest passed)
func (a *AlipayAntispamValue) Pass() bool {
	return a.Flag() > 0
}

// AlipayAntispam get alipay antispam count by mid
func (d *Dao) AlipayAntispam(c context.Context, mid int64) (value *AlipayAntispamValue, err error) {
	var (
		key  = alipayAntispamKey(mid)
		conn = d.mc.Get(c)
		item *memcache.Item
	)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrapf(err, "conn.Get(%s)", key)
		return
	}
	value = new(AlipayAntispamValue)
	if err = conn.Scan(item, &value); err != nil {
		err = errors.Wrapf(err, "conn.Scan(%+v)", item)
		return
	}
	return
}

// SetAlipayAntispam set alipay antispam count by mid
func (d *Dao) SetAlipayAntispam(c context.Context, mid int64, value *AlipayAntispamValue) (err error) {
	var (
		key  = alipayAntispamKey(mid)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: value, Flags: memcache.FlagJSON, Expiration: conf.Conf.Realname.AlipayAntispamTTL}); err != nil {
		err = errors.Wrapf(err, "conn.Set(%s,%+v)", key, value)
		return
	}
	return
}

// DeleteAlipayAntispam delete alipay antispam count by mid
func (d *Dao) DeleteAlipayAntispam(c context.Context, mid int64) (err error) {
	var (
		key  = alipayAntispamKey(mid)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrapf(err, "conn.Delete(%s)", key)
		return
	}
	return
}
