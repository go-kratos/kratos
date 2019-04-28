package dao

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

func userKey(mid int64) (key string) {
	key = fmt.Sprintf("u_%d", mid)
	return
}

// UserCache .
// func (d *Dao) UserCache(c context.Context, mid int64) (info *model.MCBlockInfo, err error) {
// 	var (
// 		key  = userKey(mid)
// 		conn = d.mc.Get(c)
// 		item *memcache.Item
// 	)
// 	defer conn.Close()
// 	if item, err = conn.Get(key); err != nil {
// 		if err == memcache.ErrNotFound {
// 			err = nil
// 			return
// 		}
// 		err = errors.WithStack(err)
// 		return
// 	}
// 	info = &model.MCBlockInfo{}
// 	if err = conn.Scan(item, info); err != nil {
// 		if err == memcache.ErrNotFound {
// 			info = nil
// 			err = nil
// 			return
// 		}
// 		err = errors.WithStack(err)
// 		return
// 	}
// 	return
// }

// DeleteUserCache .
func (d *Dao) DeleteUserCache(c context.Context, mid int64) (err error) {
	var (
		key  = userKey(mid)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		err = errors.Wrapf(err, "key : %s", key)
		return
	}
	return
}
