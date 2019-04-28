package dao

import (
	"context"
	"fmt"

	"go-common/library/cache/memcache"

	"github.com/pkg/errors"
)

func userKey(mid int64) (key string) {
	key = fmt.Sprintf("u_%d", mid)
	return
}

// DeleteUserBlockCache is.
func (d *Dao) DeleteUserBlockCache(c context.Context, mid int64) (err error) {
	var (
		key  = userKey(mid)
		conn = d.blockMemcache.Get(c)
	)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			return nil
		}
		err = errors.WithStack(err)
		return
	}
	return
}
