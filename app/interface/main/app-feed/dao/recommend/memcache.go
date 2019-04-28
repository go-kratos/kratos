package recommend

import (
	"context"

	"go-common/app/interface/main/app-card/model/card/ai"
	"go-common/library/cache/memcache"

	"github.com/pkg/errors"
)

const (
	_prefixRcmdAids       = "rc"
	_prefixRcmd           = "rc2"
	_prefixFollowModeList = "fml"
)

func keyRcmdAids() string {
	return _prefixRcmdAids
}

func keyRcmd() string {
	return _prefixRcmd
}

func keyFollowModeList() string {
	return _prefixFollowModeList
}

// AddRcmdCache add ai into cahce.
func (d *Dao) AddRcmdAidsCache(c context.Context, aids []int64) (err error) {
	conn := d.mc.Get(c)
	key := keyRcmdAids()
	item := &memcache.Item{Key: key, Object: aids, Flags: memcache.FlagJSON, Expiration: d.expireMc}
	if err = conn.Set(item); err != nil {
		err = errors.Wrapf(err, "%v", aids)
	}
	conn.Close()
	return
}

// RcmdCache get ai cache data from cache
func (d *Dao) RcmdAidsCache(c context.Context) (aids []int64, err error) {
	var r *memcache.Item
	conn := d.mc.Get(c)
	key := keyRcmdAids()
	defer conn.Close()
	if r, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrap(err, key)
		return
	}
	if err = conn.Scan(r, &aids); err != nil {
		err = errors.Wrapf(err, "%s", r.Value)
	}
	return
}

// AddRcmdCache add ai into cahce.
func (d *Dao) AddRcmdCache(c context.Context, is []*ai.Item) (err error) {
	conn := d.mc.Get(c)
	key := keyRcmd()
	item := &memcache.Item{Key: key, Object: is, Flags: memcache.FlagJSON, Expiration: d.expireMc}
	if err = conn.Set(item); err != nil {
		err = errors.Wrap(err, key)
	}
	conn.Close()
	return
}

// RcmdCache get ai cache data from cache
func (d *Dao) RcmdCache(c context.Context) (is []*ai.Item, err error) {
	var r *memcache.Item
	conn := d.mc.Get(c)
	key := keyRcmd()
	defer conn.Close()
	if r, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrap(err, key)
		return
	}
	if err = conn.Scan(r, &is); err != nil {
		err = errors.Wrapf(err, "%s", r.Value)
	}
	return
}

// AddFollowModeListCache is.
func (d *Dao) AddFollowModeListCache(c context.Context, list map[int64]struct{}) (err error) {
	conn := d.mc.Get(c)
	key := keyFollowModeList()
	item := &memcache.Item{Key: key, Object: list, Flags: memcache.FlagJSON, Expiration: d.expireMc}
	if err = conn.Set(item); err != nil {
		err = errors.Wrap(err, key)
	}
	conn.Close()
	return
}

// FollowModeListCache is.
func (d *Dao) FollowModeListCache(c context.Context) (list map[int64]struct{}, err error) {
	var r *memcache.Item
	conn := d.mc.Get(c)
	key := keyFollowModeList()
	defer conn.Close()
	if r, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrap(err, key)
		return
	}
	if err = conn.Scan(r, &list); err != nil {
		err = errors.Wrapf(err, "%s", r.Value)
	}
	return
}
