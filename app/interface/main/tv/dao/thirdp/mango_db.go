package thirdp

import (
	"context"
	"fmt"
	"go-common/app/interface/main/tv/model/thirdp"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_mangoPGCPage      = "SELECT id,mtime FROM tv_ep_season %s ORDER BY id"
	_mangoUGCPage      = "SELECT aid,mtime FROM ugc_archive %s ORDER BY aid"
	_mangoPGCOffset    = "SELECT id FROM tv_ep_season %s ORDER BY id"
	_mangoUGCOffset    = "SELECT aid FROM ugc_archive %s ORDER BY aid"
	_mangoPGCCount     = "SELECT count(1) FROM tv_ep_season"
	_mangoUGCCount     = "SELECT count(1) FROM ugc_archive"
	_mangoPGCSnCount   = "SELECT count(1) FROM tv_content WHERE season_id = ?"
	_mangoUGCArcCount  = "SELECT count(1) FROM ugc_video WHERE aid = ?"
	_mangoPGCSnOffset  = "SELECT epid,mtime FROM tv_content WHERE season_id = ?"
	_mangoUGCArcOffset = "SELECT cid,mtime FROM ugc_video WHERE aid = ?"
	_mangoPGCKey       = "TV_Mango_PGC_Page"
	_mangoUGCKey       = "TV_Mango_UGC_Page"
	// MangoPGC is mango pgc typeC
	MangoPGC = "mango_pgc"
	// MangoUGC is mango ugc typeC
	MangoUGC = "mango_ugc"
)

// MangoPages picks a page for dangbei api, lastID is the last page's biggest ID
func (d *Dao) MangoPages(ctx context.Context, req *thirdp.ReqDBeiPages) (sids []*thirdp.RespSid, myLast int64, err error) {
	var (
		rows  *sql.Rows
		query string
	)
	switch req.TypeC {
	case MangoPGC:
		query = fmt.Sprintf(_mangoPGCPage+" LIMIT %d", " WHERE id> ?", req.Ps)
	case MangoUGC:
		query = fmt.Sprintf(_mangoUGCPage+" LIMIT %d", " WHERE aid> ?", req.Ps)
	default:
		err = fmt.Errorf("MangoPages Wrong Type %s", req.TypeC)
		return
	}
	if rows, err = d.db.Query(ctx, query, req.LastID); err != nil {
		log.Error("MangoPages, lastID %d, Err %v", req.LastID, err)
		return
	}
	defer rows.Close()
	if sids, myLast, err = mangoRowsTreat(rows); err != nil {
		log.Error("MangoPages lastID %d, Err %v", req.LastID, err)
	}
	return
}

func mangoRowsTreat(rows *sql.Rows) (sids []*thirdp.RespSid, myLast int64, err error) {
	for rows.Next() {
		var r = thirdp.RespSid{}
		if err = rows.Scan(&r.Sid, &r.Mtime); err != nil {
			return
		}
		sids = append(sids, &r)
	}
	if len(sids) > 0 {
		myLast = sids[len(sids)-1].Sid
	}
	return
}

// MangoSnCnt counts ep/video number from DB
func (d *Dao) MangoSnCnt(ctx context.Context, isPGC bool, sid int64) (cnt int, err error) {
	var querySQL string
	if isPGC {
		querySQL = _mangoPGCSnCount
	} else {
		querySQL = _mangoUGCArcCount
	}
	if err = d.db.QueryRow(ctx, querySQL, sid).Scan(&cnt); err != nil {
		log.Error("MangoSnCnt isPGC %v, Sid %d, Err %v", isPGC, sid, err)
	}
	return
}

// MangoSnOffset picks season or arc's detail info by page ( limit + offset )
func (d *Dao) MangoSnOffset(ctx context.Context, isPGC bool, sid int64, pageN, pagesize int) (epids []*thirdp.RespSid, err error) {
	var (
		querySQL = fmt.Sprintf(" LIMIT %d, %d", (pageN-1)*pagesize, pagesize)
		rows     *sql.Rows
	)
	if isPGC {
		querySQL = _mangoPGCSnOffset + querySQL
	} else {
		querySQL = _mangoUGCArcOffset + querySQL
	}
	if rows, err = d.db.Query(ctx, querySQL, sid); err != nil {
		log.Error("MangoSnOffset, Sid %d, IsPGC %v, Err %v", sid, isPGC, err)
		return
	}
	defer rows.Close()
	if epids, _, err = mangoRowsTreat(rows); err != nil {
		log.Error("MangoSnOffset, Sid %d, IsPGC %v, Err %v", sid, isPGC, err)
	}
	return
}
