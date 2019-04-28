package dao

import (
	"context"
	"fmt"

	"go-common/library/cache/memcache"

	"github.com/pkg/errors"
)

const (
	_pointInfo = "pt:%d"
)

func pointKey(mid int64) string {
	return fmt.Sprintf(_pointInfo, mid)
}

//DelPointInfoCache .
func (d *Dao) DelPointInfoCache(c context.Context, mid int64) (err error) {
	return d.delCache(c, pointKey(mid))
}

func (d *Dao) delCache(c context.Context, key string) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			err = errors.Wrapf(err, "conn.Delete(%s)", key)
		}
	}
	return
}
