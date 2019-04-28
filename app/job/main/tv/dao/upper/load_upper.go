package upper

import (
	"context"
	"database/sql"
	"fmt"

	ugcMdl "go-common/app/job/main/tv/model/ugc"
	"go-common/library/cache/memcache"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_loadUpper = "SELECT mid,toinit,submit,ori_name,cms_name,ori_face,cms_face,valid,deleted FROM ugc_uploader WHERE mid = ?"
)

func upperMetaKey(MID int64) string {
	return fmt.Sprintf("up_cms_%d", MID)
}

// LoadUpMeta loads the upper meta cms data from cache, for missed ones, pick them from the DB
func (d *Dao) LoadUpMeta(ctx context.Context, mid int64) (upper *ugcMdl.Upper, err error) {
	if upper, err = d.upMetaCache(ctx, mid); err != nil { // mc error
		log.Error("LoadUpMeta Get Mid [%d] from CMS Error (%v)", mid, err)
		return
	}
	if upper != nil { // mc found
		return
	}
	if upper, err = d.upMetaDB(ctx, mid); err != nil { // db error
		log.Error("LoadUpMeta Get Mid ERROR (%d) (%v)", mid, err)
		return
	}
	if upper == nil { // db not found
		err = ecode.NothingFound
		return
	}
	d.addUpMetaCache(ctx, upper) // db found, re-fill the cache
	return
}

// upMetaCache get upper meta cache.
func (d *Dao) upMetaCache(c context.Context, mid int64) (upper *ugcMdl.Upper, err error) {
	var (
		key  = upperMetaKey(mid)
		conn = d.mc.Get(c)
		item *memcache.Item
	)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			missedCount.Add("tv-meta", 1)
		} else {
			log.Error("conn.Get(%s) error(%v)", key, err)
		}
		return
	}
	if err = conn.Scan(item, &upper); err != nil {
		log.Error("conn.Get(%s) error(%v)", key, err)
	}
	cachedCount.Add("tv-meta", 1)
	return
}

// upMetaDB gets upper meta info from DB
func (d *Dao) upMetaDB(c context.Context, mid int64) (upper *ugcMdl.Upper, err error) {
	var row *xsql.Row
	if row = d.DB.QueryRow(c, _loadUpper, mid); err != nil {
		log.Error("d.db.QueryRow(%d) error(%v)", mid, err)
		return
	}
	upper = &ugcMdl.Upper{}
	// "SELECT id,mid,to_init,submit,ori_name,cms_name,ori_face,cms_face,valid,deleted FROM ugc_uploader WHERE mid = ?"
	if err = row.Scan(&upper.MID, &upper.Toinit, &upper.Submit, &upper.OriName,
		&upper.CMSName, &upper.OriFace, &upper.CMSFace, &upper.Valid, &upper.Deleted); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			upper = nil
		} else {
			log.Error("row.Scan(mid %d) error(%v)", mid, err)
		}
	}
	return
}

// setUpMetaCache save ugcMdl.Upper to memcache
func (d *Dao) addUpMetaCache(c context.Context, upper *ugcMdl.Upper) (err error) {
	var (
		key  = upperMetaKey(upper.MID)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: upper, Flags: memcache.FlagJSON, Expiration: d.mcExpire}); err != nil {
		log.Error("conn.Set error(%v)", err)
		return
	}
	return
}
