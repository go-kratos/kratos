package block

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

func userDetailKey(mid int64) (key string) {
	key = fmt.Sprintf("ud_%d", mid)
	return
}

// DeleteUserCache .
func (d *Dao) DeleteUserCache(c context.Context, mid int64) (err error) {
	var (
		key  = userKey(mid)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrapf(err, "key : %s", key)
		return
	}
	return
}

// DeleteUserDetailCache delete user detail cache
func (d *Dao) DeleteUserDetailCache(c context.Context, mid int64) (err error) {
	var (
		key  = userDetailKey(mid)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	err = conn.Delete(key)
	return
}
