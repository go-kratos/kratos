package dao

import (
	"context"
	"fmt"
	"strconv"

	"go-common/app/interface/main/credit/model"
	gmc "go-common/library/cache/memcache"
)

const (
	_prefixBlockedUserList = "bul_%d"
	_prefixBlockInfo       = "blo_%d"
)

func userBlockedListKey(mid int64) string {
	return fmt.Sprintf(_prefixBlockedUserList, mid)
}

func blockedInfoKey(id int64) string {
	return fmt.Sprintf(_prefixBlockInfo, id)
}

// BlockedUserListCache get user blocked list.
func (d *Dao) BlockedUserListCache(c context.Context, mid int64) (ls []*model.BlockedInfo, err error) {
	var (
		reply *gmc.Item
		conn  = d.mc.Get(c)
	)
	defer conn.Close()
	reply, err = conn.Get(userBlockedListKey(mid))
	if err != nil {
		if err == gmc.ErrNotFound {
			err = nil
		}
		return
	}
	ls = make([]*model.BlockedInfo, 0)
	err = conn.Scan(reply, &ls)
	return
}

// SetBlockedUserListCache set user blocked list cache.
func (d *Dao) SetBlockedUserListCache(c context.Context, mid int64, ls []*model.BlockedInfo) (err error) {
	var (
		item = &gmc.Item{Key: userBlockedListKey(mid), Object: ls, Expiration: d.userExpire, Flags: gmc.FlagJSON}
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	err = conn.Set(item)
	return
}

// BlockedInfoCache get blocked info by blocked id.
func (d *Dao) BlockedInfoCache(c context.Context, id int64) (info *model.BlockedInfo, err error) {
	var (
		reply *gmc.Item
		conn  = d.mc.Get(c)
		key   = blockedInfoKey(id)
	)
	defer conn.Close()
	reply, err = conn.Get(key)
	if err != nil {
		if err == gmc.ErrNotFound {
			err = nil
		}
		return
	}
	info = &model.BlockedInfo{}
	err = conn.Scan(reply, &info)
	return
}

// SetBlockedInfoCache set user blocked list cache.
func (d *Dao) SetBlockedInfoCache(c context.Context, id int64, info *model.BlockedInfo) (err error) {
	var (
		item = &gmc.Item{Key: blockedInfoKey(id), Object: info, Expiration: d.minCommonExpire, Flags: gmc.FlagJSON}
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	err = conn.Set(item)
	return
}

// BlockedInfosCache get blocked infos by ids.
func (d *Dao) BlockedInfosCache(c context.Context, ids []int64) (infos []*model.BlockedInfo, miss []int64, err error) {
	var (
		rs   map[string]*gmc.Item
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	keys := make([]string, len(ids))
	for _, id := range ids {
		keys = append(keys, blockedInfoKey(id))
	}
	rs, err = conn.GetMulti(keys)
	if err != nil {
		if err == gmc.ErrNotFound {
			err = nil
		}
		return
	}
	for _, id := range ids {
		if r, ok := rs[strconv.FormatInt(id, 10)]; ok {
			info := &model.BlockedInfo{}
			conn.Scan(r, &info)
			infos = append(infos, info)
		} else {
			miss = append(miss, id)
		}
	}
	return
}

// SetBlockedInfosCache set user blocked list cache.
func (d *Dao) SetBlockedInfosCache(c context.Context, infos []*model.BlockedInfo) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	for _, info := range infos {
		item := &gmc.Item{Key: blockedInfoKey(info.ID), Object: info, Expiration: d.minCommonExpire, Flags: gmc.FlagJSON}
		if err = conn.Set(item); err != nil {
			return
		}
	}
	return
}
