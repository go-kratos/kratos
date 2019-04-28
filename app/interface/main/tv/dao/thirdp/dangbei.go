package thirdp

import (
	"context"
	"fmt"

	"go-common/app/interface/main/tv/model/thirdp"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_dangbeiPGCPage  = "SELECT id FROM tv_ep_season WHERE `check`= 1 AND valid = 1 AND is_deleted=0 %s ORDER BY id"
	_dangbeiUGCPage  = "SELECT aid FROM ugc_archive WHERE result = 1 AND valid = 1 AND deleted = 0 %s ORDER BY aid"
	_dangbeiPGCCount = "SELECT count(1) FROM tv_ep_season WHERE `check`= 1 AND valid= 1 AND is_deleted=0"
	_dangbeiUGCCount = "SELECT count(1) FROM ugc_archive WHERE result = 1 AND valid= 1 AND deleted =0"
	_dangbeiPGCKey   = "TV_Dangbei_PGC_Page"
	_dangbeiUGCKey   = "TV_Dangbei_UGC_Page"
	_countField      = "Dangbei_Count"
	// DBeiPGC is dangbei pgc typeC
	DBeiPGC = "pgc"
	// DBeiUGC is dangbei ugc typeC
	DBeiUGC = "ugc"
)

// KeyThirdp returns the key in Redis according to the type input
func KeyThirdp(typeC string) (key string, err error) {
	switch typeC {
	case DBeiPGC:
		key = _dangbeiPGCKey
	case DBeiUGC:
		key = _dangbeiUGCKey
	case MangoPGC:
		key = _mangoPGCKey
	case MangoUGC:
		key = _mangoUGCKey
	default:
		err = ecode.TvDangbeiWrongType
	}
	return
}

// ThirdpCnt counts the number of pgc/ugc data to display for dangbei
func (d *Dao) ThirdpCnt(ctx context.Context, typeC string) (count int, err error) {
	var query string
	switch typeC {
	case DBeiPGC:
		query = _dangbeiPGCCount
	case DBeiUGC:
		query = _dangbeiUGCCount
	case MangoPGC:
		query = _mangoPGCCount
	case MangoUGC:
		query = _mangoUGCCount
	default:
		err = ecode.TvDangbeiWrongType
		return
	}
	if err = d.db.QueryRow(ctx, query).Scan(&count); err != nil {
		log.Error("PickDBeiPage ThirdpCnt Error %v", err)
	}
	return
}

// DBeiPages picks a page for dangbei api, lastID is the last page's biggest ID
func (d *Dao) DBeiPages(ctx context.Context, req *thirdp.ReqDBeiPages) (sids []int64, myLast int64, err error) {
	var (
		rows   *sql.Rows
		query  string
		lastID = req.LastID
		ps     = req.Ps
	)
	switch req.TypeC {
	case DBeiPGC:
		query = fmt.Sprintf(_dangbeiPGCPage+" LIMIT %d", "AND id> ?", ps)
	case DBeiUGC:
		query = fmt.Sprintf(_dangbeiUGCPage+" LIMIT %d", "AND aid> ?", ps)
	default:
		err = ecode.TvDangbeiWrongType
		return
	}
	if rows, err = d.db.Query(ctx, query, lastID); err != nil {
		log.Error("DangbeiPage, lastID %d, Err %v", lastID, err)
		return
	}
	defer rows.Close()
	if sids, myLast, err = dbeiRowsTreat(rows); err != nil {
		log.Error("dbeiOffset lastID %d, Err %v", lastID, err)
	}
	return
}

func dbeiRowsTreat(rows *sql.Rows) (sids []int64, myLast int64, err error) {
	for rows.Next() {
		var r int64
		if err = rows.Scan(&r); err != nil {
			return
		}
		sids = append(sids, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
		return
	}
	// record my biggest id
	if len(sids) > 0 {
		myLast = sids[len(sids)-1]
	}
	return
}

// thirdpOffset is used in case of missing pageID record in Redis
func (d *Dao) thirdpOffset(ctx context.Context, page int64, ps int64, typeC string) (lastPageMax int64, err error) {
	if page <= 0 {
		return 0, nil
	}
	var query string
	switch typeC {
	case DBeiPGC:
		query = _dangbeiPGCPage
	case DBeiUGC:
		query = _dangbeiUGCPage
	case MangoPGC:
		query = _mangoPGCOffset
	case MangoUGC:
		query = _mangoUGCOffset
	default:
		err = ecode.TvDangbeiWrongType
		return
	}
	querySQL := fmt.Sprintf(query+" LIMIT %d,%d", "", page*ps-1, 1)
	if err = d.db.QueryRow(ctx, querySQL).Scan(&lastPageMax); err != nil {
		log.Error("DangbeiPage, page %d, Err %v", page, err)
	}
	return
}

// SetPageID is used to record each dangbei page's biggest ID, it's to ease the next page's pickup
func (d *Dao) SetPageID(c context.Context, req *thirdp.ReqPageID) (err error) {
	var key string
	if key, err = KeyThirdp(req.TypeC); err != nil {
		log.Error("PickDBeiPage Dangbei Key TypeC = %v, Error %v", req.TypeC, err)
		return
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("HSET", key, req.Page, req.ID); err != nil {
		log.Error("PickDBeiPage conn.Send(HSET Key %v, field %d, value %d) error(%v)", key, req.Page, req.ID, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.dbeiExpire); err != nil {
		log.Error("PickDBeiPage conn.Send error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("PickDBeiPage conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("PickDBeiPage conn.Receive error(%v)", err)
			return
		}
	}
	return
}

// SetThirdpCnt is used to record dangbei pgc data count
func (d *Dao) SetThirdpCnt(c context.Context, count int, typeC string) (err error) {
	var key string
	if key, err = KeyThirdp(typeC); err != nil {
		log.Error("PickDBeiPage Dangbei Key TypeC = %v, Error %v", typeC, err)
		return
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("HSET", key, _countField, count); err != nil {
		log.Error("PickDBeiPage conn.Send(HSET Key %v, field %d, value %d) error(%v)", key, _countField, count, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.dbeiExpire); err != nil {
		log.Error("PickDBeiPage conn.Send error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("PickDBeiPage conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("PickDBeiPage conn.Receive error(%v)", err)
			return
		}
	}
	return
}

// GetThirdpCnt get dangbei pgc data count
func (d *Dao) GetThirdpCnt(c context.Context, typeC string) (count int, err error) {
	var key string
	if key, err = KeyThirdp(typeC); err != nil {
		log.Error("PickDBeiPage Dangbei Key TypeC = %v, Error %v", typeC, err)
		return
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	if count, err = redis.Int(conn.Do("HGET", key, _countField)); err != nil {
		if err == redis.ErrNil {
			log.Info("PickDBeiPage conn.HGET(field:%s) not found", _countField)
		} else {
			log.Error("PickDBeiPage conn.HGET(field:%s) error(%v)", _countField, err)
		}
	}
	return
}

// getPageID get pageNumber's biggestID from redis hashmap
func (d *Dao) getPageID(c context.Context, pageNumber int64, typeC string) (biggestID int64, err error) {
	var key string
	if key, err = KeyThirdp(typeC); err != nil {
		log.Error("Dangbei Key TypeC = %v, Error %v", typeC, err)
		return
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	if biggestID, err = redis.Int64(conn.Do("HGET", key, pageNumber)); err != nil {
		if err == redis.ErrNil {
			log.Info("conn.HGET(page:%d) not found", pageNumber)
		} else {
			log.Error("conn.HGET(page:%d) error(%v)", pageNumber, err)
		}
	}
	return
}

// LoadPageID picks the last page's biggest ID
func (d *Dao) LoadPageID(c context.Context, req *thirdp.ReqDBeiPages) (biggestID int64, err error) {
	var (
		pageNum = req.Page - 1
		typeC   = req.TypeC
	)
	if pageNum <= 0 {
		return 0, nil
	}
	if biggestID, err = d.getPageID(c, pageNum, typeC); err == nil { // get directly from cache
		cachedCount.Add("thirdp-page", 1)
		return
	}
	missedCount.Add("thirdp-page", 1)
	if biggestID, err = d.thirdpOffset(c, pageNum, req.Ps, typeC); err != nil {
		log.Error("ThirdpOffSet TypeC %s, PageNum %d, Pagesize %d", typeC, pageNum, req.Ps)
		return
	}
	if err = d.SetPageID(c, &thirdp.ReqPageID{
		Page:  pageNum,
		ID:    biggestID,
		TypeC: typeC,
	}); err != nil {
		log.Error("ThirdpOffset TypeC %s, PageNum %d, SetPageID Err %v", typeC, pageNum, err)
	}
	return
}
