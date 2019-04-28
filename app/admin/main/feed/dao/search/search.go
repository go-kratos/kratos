package search

import (
	"context"
	"time"

	searchModel "go-common/app/admin/main/feed/model/search"
	"go-common/library/cache/memcache"
)

//SetSearchAuditStat set hot publish state to MC
func (d *Dao) SetSearchAuditStat(c context.Context, key string, state bool) (err error) {
	var (
		conn memcache.Conn
		p    searchModel.PublishState
	)
	p.Date = time.Now().Format("2006-01-02")
	p.State = state
	conn = d.MC.Get(c)
	defer conn.Close()
	itemJSON := &memcache.Item{
		Key:        key,
		Flags:      memcache.FlagJSON,
		Object:     p,
		Expiration: 0,
	}
	if err = conn.Set(itemJSON); err != nil {
		return
	}
	return
}

//GetSearchAuditStat get hot publish state from MC
func (d *Dao) GetSearchAuditStat(c context.Context, key string) (f bool, date string, err error) {
	var (
		conn memcache.Conn
		item *memcache.Item
		p    searchModel.PublishState
	)
	conn = d.MC.Get(c)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			return false, "", nil
		}
		return
	}
	if err = conn.Scan(item, &p); err != nil {
		return
	}
	return p.State, p.Date, nil
}
