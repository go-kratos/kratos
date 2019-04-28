package dao

import (
	"context"
	"fmt"

	"go-common/app/interface/main/space/model"
	"go-common/library/cache/memcache"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_privacyKey       = "spc_pcy_%d"
	_privacySQL       = `SELECT privacy,status FROM member_privacy%d WHERE mid = ?`
	_privacyModifySQL = `UPDATE member_privacy%d SET status=? WHERE mid = ? AND privacy = ?`
	_privacyAddSQL    = `INSERT INTO member_privacy%d (mid,privacy,status) VALUES (?,?,?)`
)

func privacyHit(mid int64) int64 {
	return mid % 10
}

func privacyKey(mid int64) string {
	return fmt.Sprintf(_privacyKey, mid)
}

// Privacy get privacy data.
func (d *Dao) Privacy(c context.Context, mid int64) (data map[string]int, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_privacySQL, privacyHit(mid)), mid); err != nil {
		log.Error("d.Privacy.Query(%d) error(%v)", mid, err)
		return
	}
	defer rows.Close()
	data = make(map[string]int)
	for rows.Next() {
		r := new(model.Privacy)
		if err = rows.Scan(&r.Privacy, &r.Status); err != nil {
			log.Error("row.Scan() error(%v)", err)
			return
		}
		data[r.Privacy] = r.Status
	}
	return
}

// PrivacyModify modify privacy.
func (d *Dao) PrivacyModify(c context.Context, mid int64, field string, value int) (err error) {
	var privacy map[string]int
	if privacy, err = d.Privacy(c, mid); err != nil {
		return
	}
	if _, ok := privacy[field]; ok {
		_, err = d.db.Exec(c, fmt.Sprintf(_privacyModifySQL, privacyHit(mid)), value, mid, field)
		if err != nil {
			log.Error("Privacy Update error d.db.Exec(%d,%s,%d) error(%v)", mid, field, value, err)
		}
	} else {
		_, err = d.db.Exec(c, fmt.Sprintf(_privacyAddSQL, privacyHit(mid)), mid, field, value)
		if err != nil {
			log.Error("Privacy Add error d.db.Exec(%d,%s,%d) error(%v)", mid, field, value, err)
		}
	}
	return
}

// PrivacyCache get privacy cache.
func (d *Dao) PrivacyCache(c context.Context, mid int64) (data map[string]int, err error) {
	var (
		conn = d.mc.Get(c)
		key  = privacyKey(mid)
	)
	defer conn.Close()
	reply, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("conn.Get(%v) error(%v)", key, err)
		return
	}
	if err = conn.Scan(reply, &data); err != nil {
		log.Error("reply.Scan(%s) error(%v)", reply.Value, err)
	}
	return
}

// SetPrivacyCache set privacy cache.
func (d *Dao) SetPrivacyCache(c context.Context, mid int64, data map[string]int) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &memcache.Item{Key: privacyKey(mid), Object: data, Flags: memcache.FlagJSON, Expiration: d.mcSettingExpire}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Store(%s) error(%v)", privacyKey(mid), err)
		return
	}
	return
}

// DelPrivacyCache delete privacy cache.
func (d *Dao) DelPrivacyCache(c context.Context, mid int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := privacyKey(mid)
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("DelPrivacyCache conn.Delete(%s) error(%v)", privacyKey(mid), err)
	}
	return
}
