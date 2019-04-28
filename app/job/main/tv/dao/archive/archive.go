package archive

import (
	"context"
	"strconv"

	arccli "go-common/app/service/main/archive/api"
	"go-common/library/cache/memcache"

	"github.com/pkg/errors"
)

const (
	_prefixArc  = "a3p_"
	_prefixView = "avp_"
)

func keyArc(aid int64) string {
	return _prefixArc + strconv.FormatInt(aid, 10)
}

func keyView(aid int64) string {
	return _prefixView + strconv.FormatInt(aid, 10)
}

// UpArcCache update archive cache
func (d *Dao) UpArcCache(c context.Context, a *arccli.Arc) (err error) {
	conn := d.mc.Get(c)
	item := &memcache.Item{Key: keyArc(a.Aid), Object: a, Flags: memcache.FlagJSON, Expiration: 0}
	if err = conn.Set(item); err != nil {
		err = errors.Wrapf(err, "conn.Set(%v)", item)
	}
	conn.Close()
	return
}

// UpViewCache up all app cache .
func (d *Dao) UpViewCache(c context.Context, v *arccli.ViewReply) (err error) {
	conn := d.mc.Get(c)
	item := &memcache.Item{Key: keyView(v.Arc.Aid), Object: v, Flags: memcache.FlagJSON, Expiration: 0}
	if err = conn.Set(item); err != nil {
		err = errors.Wrapf(err, "conn.Set(%v)", item)
	}
	conn.Close()
	return
}
