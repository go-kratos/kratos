package fav

import (
	"context"
	"fmt"
	favmdl "go-common/app/service/main/favorite/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
	"strconv"
)

const (
	_folder       = "f_%d_%d"      // key:f_{mid%100}_{fid},value:{*Fodler}.pb
	_relationFids = "rof_%d_%d_%d" // key:rof_{type}_{mid}_{oid},value:{[]int64}.pb
	_oidCount     = "oc_%d_%d"     // key:oc_{type}_{oid},value:int64
	_batchOids    = "bo_%d_%d"     // key:oc_{type}_{mid},value:{[]int64}.pb
	_recentOids   = "rcto_%d_%d"   // key:rcto_{type}_{mid},value:{[]int64}.pb
	_recentRes    = "rctr_%d_%d"   // key:rcto_{type}_{mid},value:{[]*Resource}.pb
)

// folderMcKey
func folderMcKey(mid, fid int64) string {
	return fmt.Sprintf(_folder, mid%100, fid)
}

// relationFidsKey
func relationFidsKey(typ int8, mid, oid int64) string {
	return fmt.Sprintf(_relationFids, typ, mid, oid)
}

func oidCountKey(typ int8, oid int64) string {
	return fmt.Sprintf(_oidCount, typ, oid)
}

func batchOidsKey(typ int8, mid int64) string {
	return fmt.Sprintf(_batchOids, typ, mid)
}

func recentOidsKey(typ int8, mid int64) string {
	return fmt.Sprintf(_recentOids, typ, mid)
}

func recentResKey(typ int8, mid int64) string {
	return fmt.Sprintf(_recentRes, typ, mid)
}

// pingMC ping mc is ok.
func (d *Dao) pingMC(c context.Context) error {
	conn := d.mc.Get(c)
	defer conn.Close()
	item := memcache.Item{Key: "ping", Value: []byte{1}, Expiration: d.mcExpire}
	return conn.Set(&item)
}

// SetFoldersMc add folders mc cache.
func (d *Dao) SetFoldersMc(c context.Context, vs ...*favmdl.Folder) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	for _, v := range vs {
		if v == nil {
			continue
		}
		item := &memcache.Item{Key: folderMcKey(v.Mid, v.ID), Object: v, Flags: memcache.FlagProtobuf, Expiration: d.mcExpire}
		if err = conn.Set(item); err != nil {
			log.Error("conn.Set(%s) error(%v)", folderMcKey(v.Mid, v.ID), err)
			return
		}
	}
	return
}

// FolderMc return one folder from mc.
func (d *Dao) FolderMc(c context.Context, typ int8, mid, fid int64) (f *favmdl.Folder, err error) {
	var (
		key  = folderMcKey(mid, fid)
		item *memcache.Item
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("memcache.Get(%s) error(%v)", key, err)
		}
		return
	}
	f = new(favmdl.Folder)
	if err = conn.Scan(item, f); err != nil {
		log.Error("conn.Scan(%s) error(%v)", item.Value, err)
		f = nil
	}
	return
}

// SetRelaitonFidsMc set fids binary data to mc.
func (d *Dao) SetRelaitonFidsMc(c context.Context, typ int8, mid, oid int64, fids []int64) (err error) {
	key := relationFidsKey(typ, mid, oid)
	conn := d.mc.Get(c)
	defer conn.Close()
	bytes := favmdl.ToBytes(fids)
	item := &memcache.Item{Key: key, Value: bytes, Flags: memcache.FlagRAW, Expiration: d.mcExpire}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%s) error(%v)", key, err)
		return
	}
	return
}

// RelaitonFidsMc return fids from mc.
func (d *Dao) RelaitonFidsMc(c context.Context, typ int8, mid, oid int64) (fids []int64, err error) {
	var (
		key  = relationFidsKey(typ, mid, oid)
		item *memcache.Item
		conn = d.mc.Get(c)
		b    []byte
	)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("memcache.Get(%s) error(%v)", key, err)
		}
		return
	}
	if err = conn.Scan(item, &b); err != nil {
		log.Error("conn.Scan(%s) error(%v)", item.Value, err)
		fids = nil
		return
	}
	if fids, err = favmdl.ToInt64s(b); err != nil {
		log.Error("fs.SetIndex(%v) error(%v)", b, err)
		err = nil
		fids = nil
	}
	return
}

// DelRelationFidsMc delete oid's fid mc cache.
func (d *Dao) DelRelationFidsMc(c context.Context, typ int8, mid, oid int64) (err error) {
	var (
		key  = relationFidsKey(typ, mid, oid)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Delete(%s) error(%v)", key, err)
		}
	}
	return
}

// SetOidCountMc return oid count from mc.
func (d *Dao) SetOidCountMc(c context.Context, typ int8, oid, count int64) (err error) {
	var (
		key  = oidCountKey(typ, oid)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	bs := []byte(strconv.FormatInt(int64(count), 10))
	item := &memcache.Item{Key: key, Value: bs, Flags: memcache.FlagRAW, Expiration: d.mcExpire}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%s) error(%v)", key, err)
		return
	}
	return
}

// DelBatchOidsMc delete oids mc cache.
func (d *Dao) DelBatchOidsMc(c context.Context, typ int8, mid int64) (err error) {
	key := batchOidsKey(typ, mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Delete(%s) error(%v)", key, err)
		}
	}
	return
}

// DelRecentOidsMc delete recent oids mc cache.
func (d *Dao) DelRecentOidsMc(c context.Context, typ int8, mid int64) (err error) {
	key := recentOidsKey(typ, mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Delete(%s) error(%v)", key, err)
		}
	}
	return
}

// DelRecentResMc delete recent oids mc cache.
func (d *Dao) DelRecentResMc(c context.Context, typ int8, mid int64) (err error) {
	key := recentResKey(typ, mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Delete(%s) error(%v)", key, err)
		}
	}
	return
}
