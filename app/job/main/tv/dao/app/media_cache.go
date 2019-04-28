package app

import (
	"context"
	"fmt"

	model "go-common/app/job/main/tv/model/pgc"
	"go-common/library/cache/memcache"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_mcSnCMSKey  = "sn_cms_%d"
	_mcEPCMSKey  = "ep_cms_%d"
	_newestOrder = "SELECT a.epid,b.`order` FROM tv_content a LEFT JOIN tv_ep_content b ON a.epid=b.id " +
		"WHERE a.season_id=? AND a.state= ? AND a.valid= ? AND a.is_deleted=0 ORDER BY b.`order` DESC LIMIT 1"
	_AllEPs       = "SELECT subtitle, epid FROM tv_content WHERE season_id = ?"
	_lessStrategy = 1
)

// SnCMSCacheKey .
func (d *Dao) SnCMSCacheKey(sid int) string {
	return fmt.Sprintf(_mcSnCMSKey, sid)
}

// EpCMSCacheKey .
func (d *Dao) EpCMSCacheKey(epid int) string {
	return fmt.Sprintf(_mcEPCMSKey, epid)
}

// SetSnCMSCache save model.SeasonCMS to memcache
func (d *Dao) SetSnCMSCache(c context.Context, s *model.SeasonCMS) (err error) {
	var (
		key  = d.SnCMSCacheKey(s.SeasonID)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: s, Flags: memcache.FlagJSON, Expiration: d.mcMediaExpire}); err != nil {
		log.Error("conn.Set error(%v)", err)
		return
	}
	return
}

// SetEpCMSCache save model.EpCMS to memcache
func (d *Dao) SetEpCMSCache(c context.Context, s *model.EpCMS) (err error) {
	var (
		key  = d.EpCMSCacheKey(s.EPID)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: s, Flags: memcache.FlagJSON, Expiration: d.mcMediaExpire}); err != nil {
		log.Error("conn.Set error(%v)", err)
		return
	}
	return
}

// NewestOrder picks one season's newest passed ep's order column value
func (d *Dao) NewestOrder(c context.Context, sid int64) (epid, newestOrder int, err error) {
	if err = d.DB.QueryRow(c, _newestOrder, sid, EPPassed, _CMSValid).Scan(&epid, &newestOrder); err != nil { // get the qualified aid to sync
		log.Info("d.NewestOrder(sid %d).Query error(%v)", sid, err)
	}
	return
}

// AllEP picks all the not deleted ep of a season
func (d *Dao) AllEP(c context.Context, sid int, strategy int) (eps []*model.EpCMS, err error) {
	var (
		rows  *sql.Rows
		query = _AllEPs
	)
	if strategy == _lessStrategy {
		query = _AllEPs + " AND is_deleted = 0" // less strategy
	}
	if rows, err = d.DB.Query(c, query, sid); err != nil {
		log.Error("d.AllEP.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r = &model.EpCMS{}
		if err = rows.Scan(&r.Title, &r.EPID); err != nil {
			log.Error("AllEP row.Scan() error(%v)", err)
			return
		}
		eps = append(eps, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("d.AllEp.Query error(%v)", err)
	}
	return
}
