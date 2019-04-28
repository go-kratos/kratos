package like

import (
	"context"
	"fmt"

	"go-common/app/interface/main/activity/model/like"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_prefixInfo = "m_"
)

func keyInfo(sid int64) string {
	return fmt.Sprintf("%s%d", _prefixInfo, sid)
}

// SetInfoCache Dao
func (dao *Dao) SetInfoCache(c context.Context, v *like.Subject, sid int64) (err error) {
	if v == nil {
		v = &like.Subject{}
	}
	var (
		conn  = dao.mc.Get(c)
		mckey = keyInfo(sid)
	)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: mckey, Object: v, Flags: memcache.FlagGOB, Expiration: dao.mcLikeExpire}); err != nil {
		log.Error("conn.Set error(%v)", err)
		return
	}
	return
}

// InfoCache Dao
func (dao *Dao) InfoCache(c context.Context, sid int64) (v *like.Subject, err error) {
	var (
		mckey = keyInfo(sid)
		conn  = dao.mc.Get(c)
		item  *memcache.Item
	)
	defer conn.Close()
	if item, err = conn.Get(mckey); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			v = nil
		} else {
			log.Error("conn.Get error(%v)", err)
		}
		return
	}
	if err = conn.Scan(item, &v); err != nil {
		log.Error("item.Scan error(%v)", err)
		return
	}
	return
}
