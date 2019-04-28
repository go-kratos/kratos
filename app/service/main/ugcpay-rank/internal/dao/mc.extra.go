package dao

import (
	"context"

	"go-common/app/service/main/ugcpay-rank/internal/conf"
	"go-common/library/cache/memcache"
)

// CASCacheElecPrepRank cas的方式存储预备榜单，防止分布式脏写
func (d *Dao) CASCacheElecPrepRank(c context.Context, val interface{}, rawItem *memcache.Item) (ok bool, err error) {
	var (
		conn = d.mc.Get(c)
	)
	defer conn.Close()

	rawItem.Object = val
	rawItem.Expiration = conf.Conf.CacheTTL.ElecPrepAVRankTTL
	rawItem.Flags = memcache.FlagProtobuf

	if err = conn.CompareAndSwap(rawItem); err != nil {
		if err == memcache.ErrCASConflict { // CAS冲突, 则返回ok == false, 准备重试
			err = nil
			return
		}
		if err == memcache.ErrNotStored { // 如果CAS中恰好失效，尝试Add
			if err = conn.Add(rawItem); err != nil {
				if err == memcache.ErrNotStored { // 在Add时恰好又被其他实例Add过, 则返回ok == false, 准备重试
					err = nil
					return
				}
				return // 在Add时发生未知错误
			}
		}
		return // 在CAS时发生未知错误
	}
	ok = true
	return
}
