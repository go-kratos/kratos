package dao

import (
	"context"
	"database/sql"
	"fmt"

	"go-common/app/interface/main/space/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_indexOrderKeyFmt = "spc_io_%d"
	_indexOrderSQL    = `SELECT index_order FROM dede_member_up_settings%d WHERE mid = ?`
	_indexOrderAddSQL = `INSERT INTO dede_member_up_settings%d (mid,index_order) VALUES (?,?) ON DUPLICATE KEY UPDATE index_order = ?`
)

func indexOrderHit(mid int64) int64 {
	return mid % 10
}

func indexOrderKey(mid int64) string {
	return fmt.Sprintf(_indexOrderKeyFmt, mid)
}

// IndexOrder get index order info.
func (d *Dao) IndexOrder(c context.Context, mid int64) (indexOrder string, err error) {
	var row = d.db.QueryRow(c, fmt.Sprintf(_indexOrderSQL, indexOrderHit(mid)), mid)
	if err = row.Scan(&indexOrder); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("IndexOrder row.Scan() error(%v)", err)
		}
	}
	return
}

// IndexOrderModify index order modify.
func (d *Dao) IndexOrderModify(c context.Context, mid int64, orderStr string) (err error) {
	if _, err = d.db.Exec(c, fmt.Sprintf(_indexOrderAddSQL, indexOrderHit(mid)), mid, orderStr, orderStr); err != nil {
		log.Error("IndexOrderModify error d.db.Exec(%d,%s) error(%v)", mid, orderStr, err)
	}
	return
}

// IndexOrderCache get index order cache.
func (d *Dao) IndexOrderCache(c context.Context, mid int64) (data []*model.IndexOrder, err error) {
	var (
		conn = d.mc.Get(c)
		key  = indexOrderKey(mid)
	)
	defer conn.Close()
	reply, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("IndexOrderCache conn.Get(%v) error(%v)", key, err)
		return
	}
	if err = conn.Scan(reply, &data); err != nil {
		log.Error("IndexOrderCache reply.Scan(%s) error(%v)", reply.Value, err)
	}
	return
}

// SetIndexOrderCache set index order cache.
func (d *Dao) SetIndexOrderCache(c context.Context, mid int64, data []*model.IndexOrder) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &memcache.Item{Key: indexOrderKey(mid), Object: data, Flags: memcache.FlagJSON, Expiration: d.mcSettingExpire}
	if err = conn.Set(item); err != nil {
		log.Error("SetIndexOrderCache conn.Set(%s) error(%v)", indexOrderKey(mid), err)
	}
	return
}

// DelIndexOrderCache delete index order cache.
func (d *Dao) DelIndexOrderCache(c context.Context, mid int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := indexOrderKey(mid)
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("DelIndexOrderCache conn.Delete(%s) error(%v)", indexOrderKey(mid), err)
	}
	return
}
