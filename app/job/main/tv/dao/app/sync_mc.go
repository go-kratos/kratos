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
	_seasonKey   = "sn_%d"
	_epKey       = "ep_%d"
	_countEP     = "SELECT COUNT(1) AS cnt FROM tv_content"
	_countSeason = "SELECT COUNT(1) AS cnt FROM tv_ep_season"
	_pickEPMC    = "SELECT is_deleted, state,valid, season_id, epid,id, mark, cover, title, subtitle, pay_status FROM tv_content " +
		"WHERE id > ? ORDER BY id LIMIT 0,"
	_pickSeasonMC = "SELECT is_deleted, `check`,valid, id, cover, `desc`, title, upinfo, category, area, play_time, role, staff, total_num, style, alias, origin_name, status FROM tv_ep_season " +
		"WHERE id > ? ORDER BY id LIMIT 0,"
	_singleSn = "SELECT id, cover, `desc`, title, upinfo, category, area, play_time, role, staff, total_num, style, alias, origin_name, status FROM tv_ep_season WHERE id = ?"
)

// EpCacheKey is used to generate the key of ep
func EpCacheKey(epid int) string {
	return fmt.Sprintf(_epKey, epid)
}

// SeasonCacheKey is used to generate the key of season
func SeasonCacheKey(sid int) string {
	return fmt.Sprintf(_seasonKey, sid)
}

// SetEP in MC
func (d *Dao) SetEP(ctx context.Context, res *model.SimpleEP) (err error) {
	var (
		key  = EpCacheKey(res.EPID)
		conn = d.mc.Get(ctx)
	)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: res, Expiration: d.mcExpire, Flags: memcache.FlagJSON}); err != nil {
		log.Error("conn.Set(%s,%v) error(%v)", key, res, err)
	}
	return
}

// SetSeason in MC
func (d *Dao) SetSeason(ctx context.Context, res *model.SimpleSeason) (err error) {
	var (
		key  = SeasonCacheKey(int(res.ID))
		conn = d.mc.Get(ctx)
	)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: res, Expiration: d.mcExpire, Flags: memcache.FlagJSON}); err != nil {
		log.Error("conn.Set(%s,%v) error(%v)", key, res, err)
	}
	return
}

// CountEP counts number of ep in DB
func (d *Dao) CountEP(ctx context.Context) (count int, err error) {
	row := d.DB.QueryRow(ctx, _countEP)
	err = row.Scan(&count)
	return
}

// CountSeason counts number of ep in DB
func (d *Dao) CountSeason(ctx context.Context) (count int, err error) {
	row := d.DB.QueryRow(ctx, _countSeason)
	err = row.Scan(&count)
	return
}

// RefreshEPMC picks data by Piece to sync in MC
func (d *Dao) RefreshEPMC(ctx context.Context, LastID int, nbData int) (myLast int, err error) {
	var rows *sql.Rows
	if rows, err = d.DB.Query(ctx, _pickEPMC+fmt.Sprintf("%d", nbData), LastID); err != nil {
		log.Error("d._pickEPMC.Query: %s error(%v)", _pickEPMC+fmt.Sprintf("%d", nbData), err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			r     = &model.SimpleEP{}
			rMeta = &model.EpCMS{}
		)
		if err = rows.Scan(&r.IsDeleted, &r.State, &r.Valid, &r.SeasonID, &r.EPID, &r.ID,
			&r.NoMark, &rMeta.Cover, &rMeta.Title, &rMeta.Subtitle, &rMeta.PayStatus); err != nil {
			log.Error("RefreshEPMC row.Scan() error(%v)", err)
			return
		}
		rMeta.EPID = r.EPID
		myLast = int(r.ID)
		if err = d.SetEP(ctx, r); err != nil {
			log.Warn("RefreshEPMC Auth Set EPID (%d), error (%v)", r.EPID, err)
		}
		if err = d.SetEpCMSCache(ctx, rMeta); err != nil {
			log.Warn("RefreshEPMC Meta Set EPID (%d), error (%v)", r.EPID, err)
		}
	}
	if err = rows.Err(); err != nil {
		log.Error("d.RefreshEpMC.Query error(%v)", err)
	}
	return
}

// RefreshSnMC picks data by Piece to sync into MC
func (d *Dao) RefreshSnMC(ctx context.Context, LastID int, nbData int) (myLast int, err error) {
	var rows *sql.Rows
	if rows, err = d.DB.Query(ctx, _pickSeasonMC+fmt.Sprintf("%d", nbData), LastID); err != nil {
		log.Error("d._pickSeasonMC.Query: %s error(%v)", _pickSeasonMC+fmt.Sprintf("%d", nbData), err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			auth  = &model.SimpleSeason{}
			media = &model.SeasonCMS{}
		)
		// SELECT is_deleted, `check`,valid, id, cover, `desc`, title, upinfo, category, area, play_time, role, staff, total_num, style FROM tv_ep_season
		if err = rows.Scan(&auth.IsDeleted, &auth.Check, &auth.Valid, &auth.ID, &media.Cover,
			&media.Desc, &media.Title, &media.UpInfo, &media.Category, &media.Area, &media.Playtime,
			&media.Role, &media.Staff, &media.TotalNum, &media.Style, &media.Alias, &media.OriginName,
			&media.PayStatus); err != nil { // refresh cache sn logic
			log.Error("RefreshSnMC row.Scan() error(%v)", err)
			return
		}
		media.SeasonID = int(auth.ID)                                        // season_id
		media.NewestEPID, media.NewestOrder, _ = d.NewestOrder(ctx, auth.ID) // newest info
		myLast = int(auth.ID)
		if err = d.SetSeason(ctx, auth); err != nil {
			log.Warn("RefreshSnMC Auth Set Sid (%d), error (%v)", auth.ID, err)
		}
		if err = d.SetSnCMSCache(ctx, media); err != nil {
			log.Warn("RefreshSnMC Meta Set Sid (%d), error (%v)", auth.ID, err)
		}
	}
	if err = rows.Err(); err != nil {
		log.Error("d.RefreshSnMC.Query error(%v)", err)
	}
	return
}

// PickSeason picks one season CMS struct
func (d *Dao) PickSeason(ctx context.Context, sid int) (media *model.SeasonCMS, err error) {
	media = &model.SeasonCMS{}
	if err = d.DB.QueryRow(ctx, _singleSn, sid).Scan(&media.SeasonID, &media.Cover,
		&media.Desc, &media.Title, &media.UpInfo, &media.Category, &media.Area, &media.Playtime,
		&media.Role, &media.Staff, &media.TotalNum, &media.Style, &media.Alias, &media.OriginName,
		&media.PayStatus); err != nil { // databus ep logic, with sn
		log.Error("d.PickSeason Sid %d Error %v", sid, err)
		if err == sql.ErrNoRows {
			err = nil
			media = nil
			return
		}
		return
	}
	media.NewestEPID, media.NewestOrder, _ = d.NewestOrder(ctx, int64(media.SeasonID)) // newest info
	return
}
