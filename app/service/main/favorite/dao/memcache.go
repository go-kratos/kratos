package dao

import (
	"context"
	"fmt"
	"go-common/app/service/main/favorite/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
	"go-common/library/xstr"
	"strconv"
)

const (
	_oldFolder    = "f_%d_%d_%d"   // key:f_{type}_{mid}_{fid},value:{*Fodler}.pb
	_folder       = "f_%d_%d"      // key:f_{mid%100}_{fid},value:{*Fodler}.pb
	_folderStat   = "ft_%d_%d"     // key:f_{mid%100}_{fid},value:{*FodlerStat}.pb
	_folderSort   = "fst_%d_%d"    // key:f_{type}_{mid},value:{fid,fid,fid...}.binary
	_relationFids = "rof_%d_%d_%d" // key:rof_{type}_{mid}_{oid},value:{[]int64}.pb
	_oidCount     = "oc_%d_%d"     // key:oc_{type}_{oid},value:int64
	_batchOids    = "bo_%d_%d"     // key:oc_{type}_{mid},value:{[]int64}.pb
	_recentOids   = "rcto_%d_%d"   // key:rcto_{type}_{mid},value:{[]int64}.pb
	_recentRes    = "rctr_%d_%d"   // key:rcto_{type}_{mid},value:{[]*Resource}.pb
)

func folderMcKey(mid, fid int64) string {
	return fmt.Sprintf(_folder, mid%100, fid)
}

func folderStatMcKey(mid, fid int64) string {
	return fmt.Sprintf(_folderStat, mid%100, fid)
}

func fsortMcKey(typ int8, mid int64) string {
	return fmt.Sprintf(_folderSort, typ, mid)
}

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

// FolderStatsMc return folders stat by mid & fid from mc.
func (d *Dao) FolderStatsMc(c context.Context, fvmids []*model.ArgFVmid) (fs map[int64]*model.Folder, missFvmids []*model.ArgFVmid, err error) {
	var (
		keys    = make([]string, 0, len(fvmids))
		keysMap = make(map[string]*model.ArgFVmid, len(fvmids))
	)
	for _, v := range fvmids {
		key := folderStatMcKey(v.Vmid, v.Fid)
		keys = append(keys, key)
		keysMap[key] = v
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	items, err := conn.GetMulti(keys)
	if err != nil {
		log.Error("conn.GetMulti(%v) error(%v)", keys, err)
		return
	}
	fs = make(map[int64]*model.Folder, len(items))
	for _, item := range items {
		stat := new(model.Folder)
		if err = conn.Scan(item, stat); err != nil {
			log.Error("conn.Scan(%s) error(%v)", item.Value, err)
			return
		}
		fs[keysMap[item.Key].MediaID()] = stat
		delete(keysMap, item.Key)
	}
	for _, v := range keysMap {
		missFvmids = append(missFvmids, v)
	}
	return
}

// SetFoldersMc add folders mc cache.
func (d *Dao) SetFoldersMc(c context.Context, vs ...*model.Folder) (err error) {
	if len(vs) == 0 {
		return
	}
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

// SetFolderStatsMc add folders mc cache.
func (d *Dao) SetFolderStatsMc(c context.Context, stats map[int64]*model.Folder) (err error) {
	if len(stats) == 0 {
		return
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	for k, v := range stats {
		if v == nil {
			continue
		}
		item := &memcache.Item{Key: folderStatMcKey(k%100, k/100), Object: v, Flags: memcache.FlagProtobuf, Expiration: d.mcExpire}
		if err = conn.Set(item); err != nil {
			log.Error("conn.Set(%s) error(%v)", folderStatMcKey(k%100, k/100), err)
			return
		}
	}
	return
}

// FoldersMc return folders by mid & fid from mc.
func (d *Dao) FoldersMc(c context.Context, fvmids []*model.ArgFVmid) (fs map[string]*model.Folder, missFvmids []*model.ArgFVmid, err error) {
	var (
		keys    = make([]string, 0, len(fvmids))
		keysMap = make(map[string]*model.ArgFVmid, len(fvmids))
	)
	for _, v := range fvmids {
		key := folderMcKey(v.Vmid, v.Fid)
		keys = append(keys, key)
		keysMap[key] = v
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	items, err := conn.GetMulti(keys)
	if err != nil {
		log.Error("conn.GetMulti(%v) error(%v)", keys, err)
		return
	}
	fs = make(map[string]*model.Folder, len(items))
	for _, item := range items {
		folder := new(model.Folder)
		if err = conn.Scan(item, folder); err != nil {
			log.Error("conn.Scan(%s) error(%v)", item.Value, err)
			return
		}
		delete(keysMap, item.Key)
		fvmidStr := xstr.JoinInts([]int64{folder.ID, folder.Mid})
		fs[fvmidStr] = folder
	}
	for _, v := range keysMap {
		missFvmids = append(missFvmids, v)
	}
	return
}

// FolderMc return folder pb from mc.
func (d *Dao) FolderMc(c context.Context, typ int8, mid, fid int64) (f *model.Folder, err error) {
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
	f = new(model.Folder)
	if err = conn.Scan(item, f); err != nil {
		log.Error("conn.Scan(%s) error(%v)", item.Value, err)
		f = nil
	}
	return
}

// DelFolderMc delete folder mc cache.
func (d *Dao) DelFolderMc(c context.Context, typ int8, mid, fid int64) (err error) {
	var (
		key  = folderMcKey(mid, fid)
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

// SetFolderSortMc set folder's sort binary data to mc.
func (d *Dao) SetFolderSortMc(c context.Context, fst *model.FolderSort) (err error) {
	key := fsortMcKey(fst.Type, fst.Mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	sortBin := fst.Index()
	item := &memcache.Item{
		Key:        key,
		Value:      sortBin,
		Expiration: d.mcExpire,
	}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%s) error(%v)", key, err)
		return
	}
	return
}

// FolderSortMc return folder sort binary from mc.
func (d *Dao) FolderSortMc(c context.Context, typ int8, mid int64) (fst *model.FolderSort, err error) {
	var (
		key  = fsortMcKey(typ, mid)
		item *memcache.Item
		b    []byte
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
	if err = conn.Scan(item, &b); err != nil {
		log.Error("conn.Scan(%s) error(%v)", item.Value, err)
		fst = nil
		return
	}
	fst = &model.FolderSort{Type: typ, Mid: mid}
	if err = fst.SetIndex(b); err != nil {
		log.Error("fs.SetIndex(%v) error(%v)", b, err)
		err = nil
		fst = nil
	}
	return
}

// DelFolderSortMc delete folder's sort mc cache.
func (d *Dao) DelFolderSortMc(c context.Context, typ int8, mid int64) (err error) {
	var (
		key  = fsortMcKey(typ, mid)
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

// SetRelaitonFidsMc set fids binary data to mc.
func (d *Dao) SetRelaitonFidsMc(c context.Context, typ int8, mid, oid int64, fids []int64) (err error) {
	key := relationFidsKey(typ, mid, oid)
	conn := d.mc.Get(c)
	defer conn.Close()
	bytes := model.ToBytes(fids)
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
	if fids, err = model.ToInt64s(b); err != nil {
		log.Error("fs.SetIndex(%v) error(%v)", b, err)
		err = nil
		fids = nil
	}
	return
}

// DelRelationFidsMc delete oid's fid mc cache.
func (d *Dao) DelRelationFidsMc(c context.Context, typ int8, mid int64, oids ...int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	for _, oid := range oids {
		key := relationFidsKey(typ, mid, oid)
		if err = conn.Delete(key); err != nil {
			if err == memcache.ErrNotFound {
				err = nil
			} else {
				log.Error("conn.Delete(%s) error(%v)", key, err)
			}
		}
	}
	return
}

// SetOidCountMc set oid count to mc.
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

// OidCountMc return oid count from mc.
func (d *Dao) OidCountMc(c context.Context, typ int8, oid int64) (count int64, err error) {
	var (
		key  = oidCountKey(typ, oid)
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
	var v string
	if err = conn.Scan(item, &v); err != nil {
		log.Error("conn.Scan(%s) error(%v)", item.Value, err)
		return
	}
	count, err = strconv.ParseInt(v, 10, 64)
	return
}

// OidsCountMc return oids's count from mc.
func (d *Dao) OidsCountMc(c context.Context, typ int8, oids []int64) (counts map[int64]int64, misOids []int64, err error) {
	var (
		keys    = make([]string, 0, len(oids))
		keysMap = make(map[string]int64, len(oids))
	)
	for _, oid := range oids {
		key := oidCountKey(typ, oid)
		keys = append(keys, key)
		keysMap[key] = oid
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	items, err := conn.GetMulti(keys)
	if err != nil {
		log.Error("conn.GetMulti(%v) error(%v)", keys, err)
		return
	}
	counts = make(map[int64]int64, len(items))
	for _, item := range items {
		var v string
		if err = conn.Scan(item, &v); err != nil {
			log.Error("conn.Scan(%s) error(%v)", item.Value, err)
			return
		}
		var cnt int64
		if cnt, err = strconv.ParseInt(v, 10, 64); err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", v, err)
			return
		}
		counts[keysMap[item.Key]] = cnt
		delete(keysMap, item.Key)
	}
	for _, v := range keysMap {
		misOids = append(misOids, v)
	}
	return
}

// SetOidsCountMc set oids count to mc.
func (d *Dao) SetOidsCountMc(c context.Context, typ int8, cnts map[int64]int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	for oid, cnt := range cnts {
		bs := []byte(strconv.FormatInt(int64(cnt), 10))
		item := &memcache.Item{Key: oidCountKey(typ, oid), Object: bs, Flags: memcache.FlagRAW, Expiration: d.mcExpire}
		if err = conn.Set(item); err != nil {
			log.Error("conn.Set(%s) error(%v)", oidCountKey(typ, oid), err)
			return
		}
	}
	return
}

// BatchOidsMc return oids from mc.
func (d *Dao) BatchOidsMc(c context.Context, typ int8, mid int64) (oids []int64, err error) {
	var (
		key  = batchOidsKey(typ, mid)
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
		return
	}
	if oids, err = model.ToInt64s(b); err != nil {
		log.Error("fs.SetIndex(%v) error(%v)", b, err)
	}
	return
}

// SetBatchOidsMc set oids binary data to mc.
func (d *Dao) SetBatchOidsMc(c context.Context, typ int8, mid int64, oids []int64) (err error) {
	key := batchOidsKey(typ, mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	bytes := model.ToBytes(oids)
	item := &memcache.Item{Key: key, Value: bytes, Flags: memcache.FlagRAW, Expiration: d.mcExpire}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%s) error(%v)", key, err)
		return
	}
	return
}

// UserRecentOidsMc return oids from mc.
func (d *Dao) UserRecentOidsMc(c context.Context, typ int8, mid int64) (oids []int64, err error) {
	var (
		key  = recentOidsKey(typ, mid)
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
		return
	}
	if oids, err = model.ToInt64s(b); err != nil {
		log.Error("fs.SetIndex(%v) error(%v)", b, err)
	}
	return
}

// SetUserRecentOidsMc set oids binary data to mc.
func (d *Dao) SetUserRecentOidsMc(c context.Context, typ int8, mid int64, oids []int64) (err error) {
	key := recentOidsKey(typ, mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	bytes := model.ToBytes(oids)
	item := &memcache.Item{Key: key, Value: bytes, Flags: memcache.FlagRAW, Expiration: d.mcExpire}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%s) error(%v)", key, err)
		return
	}
	return
}

// SetUserRecentOidsMc set oids binary data to mc.
func (d *Dao) SetUserRecentResourcesMc(c context.Context, typ int8, mid int64, recents []*model.Resource) (err error) {
	key := recentResKey(typ, mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &memcache.Item{Key: key, Object: recents, Flags: memcache.FlagJSON, Expiration: d.mcExpire}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%s) error(%v)", key, err)
		return
	}
	return
}

// UserRecentOidsMc return oids from mc.
func (d *Dao) UserRecentResourcesMc(c context.Context, typ int8, mid int64) (recents []*model.Resource, err error) {
	var (
		key  = recentResKey(typ, mid)
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
	recents = make([]*model.Resource, 0)
	if err = conn.Scan(item, &recents); err != nil {
		log.Error("conn.Scan(%s) error(%v)", item.Value, err)
		return
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
