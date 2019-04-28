package stat

import (
	"context"
	"fmt"

	favmdl "go-common/app/service/main/favorite/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_folderStat = "ft_%s_%d" // key:ft_{mid%100}_{fid},value:{*FodlerStat}.pb
)

// folderStatMcKey
func folderStatMcKey(table string, fid int64) string {
	return fmt.Sprintf(_folderStat, table, fid)
}

// SetFolderStatMc add folder stat mc cache.
func (d *Dao) SetFolderStatMc(c context.Context, id int64, s *favmdl.Folder) (err error) {
	fid, table := hit(id)
	conn := d.mc.Get(c)
	defer conn.Close()

	item := &memcache.Item{Key: folderStatMcKey(table, fid), Object: s, Flags: memcache.FlagProtobuf, Expiration: d.mcExpire}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%s) error(%v)", folderStatMcKey(table, fid), err)
	}
	return
}

// FolderStatMc return one folder stat from mc.
func (d *Dao) FolderStatMc(c context.Context, id int64) (f *favmdl.Folder, err error) {
	fid, table := hit(id)
	var (
		key  = folderStatMcKey(table, fid)
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
