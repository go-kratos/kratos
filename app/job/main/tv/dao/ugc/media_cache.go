package ugc

import (
	"context"
	"fmt"

	ugcmdl "go-common/app/job/main/tv/model/ugc"
	"go-common/library/cache/memcache"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_mcArcCMSKey   = "arc_cms_%d"
	_mcVideoCMSKey = "video_cms_%d"
	_totalArcs     = "SELECT count(1) FROM ugc_archive WHERE mid = ?"
	_totalVideos   = "SELECT count(1) FROM ugc_video"
	_activeVideos  = "SELECT count(1) FROM ugc_video WHERE aid = ? AND deleted = 0"
	_cntArcVideo   = "SELECT count(1) FROM ugc_video WHERE aid = ?"
	_pickArcMC     = "SELECT title, aid, content, cover, typeid, pubtime, videos, valid, deleted, result,copyright, state, mid,duration FROM ugc_archive " +
		"WHERE mid = %d LIMIT %d,%d"
	_pickArcVideoMC = "SELECT cid, eptitle, aid, index_order, valid, deleted, result FROM ugc_video " +
		"WHERE aid = ? AND cid > ? ORDER BY cid LIMIT 0,%d"
	_transFailVideos = "SELECT cid FROM ugc_video WHERE aid = ? AND cid > %d AND transcoded = 2 and deleted = 0"
)

func arcCMSCacheKey(sid int64) string {
	return fmt.Sprintf(_mcArcCMSKey, sid)
}

func videoCMSCacheKey(sid int) string {
	return fmt.Sprintf(_mcVideoCMSKey, sid)
}

// SetArcCMS in MC
func (d *Dao) SetArcCMS(ctx context.Context, res *ugcmdl.ArcCMS) (err error) {
	var (
		key  = arcCMSCacheKey(res.AID)
		conn = d.mc.Get(ctx)
	)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: res, Expiration: d.mcExpire, Flags: memcache.FlagJSON}); err != nil {
		log.Error("conn.Set(%s,%v) error(%v)", key, res, err)
	}
	return
}

// SetVideoCMS in MC
func (d *Dao) SetVideoCMS(ctx context.Context, res *ugcmdl.VideoCMS) (err error) {
	var (
		key  = videoCMSCacheKey(res.CID)
		conn = d.mc.Get(ctx)
	)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: res, Expiration: d.mcExpire, Flags: memcache.FlagJSON}); err != nil {
		log.Error("conn.Set(%s,%v) error(%v)", key, res, err)
	}
	return
}

// UpArcsCnt counts the total number of arcs, including the deleted ones
func (d *Dao) UpArcsCnt(c context.Context, mid int64) (count int, err error) {
	if err = d.DB.QueryRow(c, _totalArcs, mid).Scan(&count); err != nil {
		log.Error("d.UpArcsCnt.Query error(%v)", err)
	}
	return
}

// TotalVideos counts the total number of arcs, including the deleted ones
func (d *Dao) TotalVideos(c context.Context) (count int, err error) {
	if err = d.DB.QueryRow(c, _totalVideos).Scan(&count); err != nil {
		log.Error("d.TotalVideos.Query error(%v)", err)
	}
	return
}

// ActVideos checks whether there is still some active videos
func (d *Dao) ActVideos(c context.Context, aid int64) (has bool, err error) {
	var count int
	if err = d.DB.QueryRow(c, _activeVideos, aid).Scan(&count); err != nil {
		log.Error("d.TotalVideos.Query error(%v)", err)
		return
	}
	if count == 0 {
		has = false
	} else {
		has = true
	}
	return
}

// ArcVideoCnt counts one archive's video
func (d *Dao) ArcVideoCnt(c context.Context, aid int64) (count int, err error) {
	if err = d.DB.QueryRow(c, _cntArcVideo, aid).Scan(&count); err != nil {
		log.Error("d.TotalVideos.Query error(%v)", err)
	}
	return
}

// PickUpArcs picks data by Piece to sync in MC, attention: nbPiece begins from zero
func (d *Dao) PickUpArcs(ctx context.Context, mid, nbPiece, nbData int) (res []*ugcmdl.ArcFull, err error) {
	var (
		rows  *sql.Rows
		query = fmt.Sprintf(_pickArcMC, mid, nbPiece*nbData, nbData)
	)
	if rows, err = d.DB.Query(ctx, query); err != nil {
		log.Error("d.PickUpArcs.Query: %s error(%v)", query, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var li = &ugcmdl.ArcFull{}
		// SELECT title, aid, content, cover, typeid, pubtime, videos, valid, deleted, result,copyright, state, mid,duration
		if err = rows.Scan(&li.Title, &li.AID, &li.Content, &li.Cover, &li.TypeID, &li.Pubtime,
			&li.Videos, &li.Valid, &li.Deleted, &li.Result, &li.Copyright, &li.State, &li.MID, &li.Duration); err != nil {
			log.Error("PickUpArcs row.Scan() error(%v)", err)
			return
		}
		res = append(res, li)
	}
	if err = rows.Err(); err != nil {
		log.Error("d.PickUpArcs.Query error(%v)", err)
	}
	return
}

// PickArcVideo picks data by Piece to sync in MC
func (d *Dao) PickArcVideo(ctx context.Context, aid int64, LastID int, nbData int) (res []*ugcmdl.VideoCMS, myLast int, err error) {
	var (
		rows  *sql.Rows
		query = fmt.Sprintf(_pickArcVideoMC, nbData)
	)
	if rows, err = d.DB.Query(ctx, query, aid, LastID); err != nil {
		log.Error("d._pickArcVideoMC.Query: %s error(%v)", query, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var li = new(ugcmdl.VideoCMS)
		if err = rows.Scan(&li.CID, &li.Title, &li.AID, &li.IndexOrder, &li.Valid, &li.Deleted, &li.Result); err != nil {
			log.Error("PickVideoMC row.Scan() error(%v)", err)
			return
		}
		res = append(res, li)
		myLast = li.CID
	}
	if err = rows.Err(); err != nil {
		log.Error("d.PickArcVideo.Query error(%v)", err)
	}
	return
}

// TransFailVideos picks transcoding failure videos
func (d *Dao) TransFailVideos(ctx context.Context, aid int64) (cids []int64, err error) {
	var (
		rows  *sql.Rows
		query = fmt.Sprintf(_transFailVideos, d.criCID)
	)
	if rows, err = d.DB.Query(ctx, query, aid); err != nil {
		log.Error("d.TransFailVideos.Query: Cid %d error(%v)", aid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var li int64
		if err = rows.Scan(&li); err != nil {
			log.Error("PickVideoMC row.Scan() error(%v)", err)
			return
		}
		cids = append(cids, li)
	}
	if err = rows.Err(); err != nil {
		log.Error("d.TransFailVideos.Query error(%v)", err)
	}
	return
}
