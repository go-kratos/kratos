package cms

import (
	"context"
	"database/sql"
	"fmt"

	"go-common/app/interface/main/tv/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_seasonCols  = "SELECT id, cover, title , upinfo, `desc`, category, area, play_time, role, staff, total_num, style,origin_name,alias,status FROM tv_ep_season "
	_epCols      = "SELECT epid, cover, title, subtitle, pay_status FROM tv_content "
	_seasonMetas = _seasonCols + " WHERE id IN (%s)"
	_arcMetas    = "SELECT title, aid, content, cover, typeid, pubtime, videos, valid, deleted, result FROM ugc_archive WHERE aid IN (%s)"
	_videoMeta   = "SELECT cid, eptitle, aid, index_order, valid, deleted, result FROM ugc_video WHERE cid = ?"
	_videoMetas  = "SELECT cid, eptitle, aid, index_order, valid, deleted, result FROM ugc_video WHERE cid IN (%s)"
	_epMetas     = _epCols + " WHERE epid IN (%s) AND is_deleted = 0 "
	_simpleEPC   = "SELECT `id`, `epid`, `state`, `is_deleted`, `valid`, `season_id` , `mark` FROM `tv_content` WHERE `epid` = ? LIMIT 1"
	_simpleEPCs  = "SELECT `id`, `epid`, `state`, `is_deleted`, `valid`, `season_id` , `mark` FROM `tv_content` WHERE `epid` IN (%s)"
	_simpleSea   = "SELECT `id`, `is_deleted`, `check`, `valid` FROM `tv_ep_season` WHERE `id` = ? LIMIT 1"
	_simpleSeas  = "SELECT `id`, `is_deleted`, `check`, `valid` FROM `tv_ep_season` WHERE `id` IN (%s)"
	_seasonCMS   = _seasonCols + "WHERE id = ? AND is_deleted = 0"
	_epCMS       = _epCols + " WHERE epid = ? AND is_deleted = 0 "
	_newestOrder = "SELECT a.epid,b.`order` FROM tv_content a LEFT JOIN tv_ep_content b ON a.epid=b.id " +
		"WHERE a.season_id=? AND a.state= ? AND a.valid= ? AND a.is_deleted=0 ORDER by `order` DESC LIMIT 1"
	epStatePass = 3
	_CMSValid   = 1
)

// VideoMetaDB picks video from DB
func (d *Dao) VideoMetaDB(c context.Context, cid int64) (meta *model.VideoCMS, err error) {
	rows, err := d.db.Query(c, _videoMeta, cid)
	if err != nil {
		log.Error("VideoMetaDB d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		li := &model.VideoCMS{}
		if err = rows.Scan(&li.CID, &li.Title, &li.AID, &li.IndexOrder, &li.Valid, &li.Deleted, &li.Result); err != nil {
			log.Error("VideoMetaDB row.Scan error(%v)", err)
			return
		}
		meta = li
	}
	return
}

// ArcMetaDB picks arc from DB
func (d *Dao) ArcMetaDB(c context.Context, aid int64) (meta *model.ArcCMS, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_arcMetas, fmt.Sprintf("%d", aid)))
	if err != nil {
		log.Error("ArcMetas d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		li := &model.ArcCMS{}
		if err = rows.Scan(&li.Title, &li.AID, &li.Content, &li.Cover, &li.TypeID, &li.Pubtime, &li.Videos, &li.Valid, &li.Deleted, &li.Result); err != nil {
			log.Error("ArcMetas row.Scan error(%v)", err)
			return
		}
		meta = li
	}
	return
}

// VideoMetas picks video from DB
func (d *Dao) VideoMetas(c context.Context, cids []int64) (meta map[int64]*model.VideoCMS, err error) {
	meta = make(map[int64]*model.VideoCMS)
	rows, err := d.db.Query(c, fmt.Sprintf(_videoMetas, xstr.JoinInts(cids)))
	if err != nil {
		log.Error("VideoMetaDB d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		li := &model.VideoCMS{}
		if err = rows.Scan(&li.CID, &li.Title, &li.AID, &li.IndexOrder, &li.Valid, &li.Deleted, &li.Result); err != nil {
			log.Error("VideoMetaDB row.Scan error(%v)", err)
			return
		}
		meta[int64(li.CID)] = li
	}
	return
}

// ArcMetas picks seasons from DB
func (d *Dao) ArcMetas(c context.Context, aids []int64) (metas map[int64]*model.ArcCMS, err error) {
	metas = make(map[int64]*model.ArcCMS)
	rows, err := d.db.Query(c, fmt.Sprintf(_arcMetas, xstr.JoinInts(aids)))
	if err != nil {
		log.Error("ArcMetas d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		li := &model.ArcCMS{}
		if err = rows.Scan(&li.Title, &li.AID, &li.Content, &li.Cover, &li.TypeID, &li.Pubtime, &li.Videos, &li.Valid, &li.Deleted, &li.Result); err != nil {
			log.Error("ArcMetas row.Scan error(%v)", err)
			return
		}
		metas[li.AID] = li
	}
	return
}

// SeasonMetas picks seasons from DB
func (d *Dao) SeasonMetas(c context.Context, sids []int64) (metas map[int64]*model.SeasonCMS, err error) {
	metas = make(map[int64]*model.SeasonCMS)
	rows, err := d.db.Query(c, fmt.Sprintf(_seasonMetas, xstr.JoinInts(sids)))
	if err != nil {
		log.Error("SeasonMetas d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		li := &model.SeasonCMS{}
		// SELECT id, cover, title , upinfo, `desc`, category, area, play_time, role, staff, total_num, style, origin_name, status
		if err = rows.Scan(&li.SeasonID, &li.Cover, &li.Title, &li.UpInfo, &li.Desc, &li.Category, &li.Area,
			&li.Playtime, &li.Role, &li.Staff, &li.TotalNum, &li.Style, &li.OriginName, &li.Alias, &li.PayStatus); err != nil {
			log.Error("SeasonMetas row.Scan error(%v)", err)
			return
		}
		metas[int64(li.SeasonID)] = li
	}
	for _, v := range metas {
		v.NewestEPID, v.NewestOrder, _ = d.NewestOrder(c, v.SeasonID)
	}
	return
}

// NewestOrder picks one season's newest passed ep's order column value
func (d *Dao) NewestOrder(c context.Context, sid int64) (epid int64, newestOrder int, err error) {
	if err = d.db.QueryRow(c, _newestOrder, sid, epStatePass, _CMSValid).Scan(&epid, &newestOrder); err != nil { // get the qualified aid to sync
		log.Warn("d.NewestOrder(sid %d).Query error(%v)", sid, err)
	}
	return
}

// EpMetas picks ep info from DB
func (d *Dao) EpMetas(c context.Context, epids []int64) (metas map[int64]*model.EpCMS, err error) {
	metas = make(map[int64]*model.EpCMS)
	rows, err := d.db.Query(c, fmt.Sprintf(_epMetas, xstr.JoinInts(epids)))
	if err != nil {
		log.Error("EpMetas d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		li := &model.EpCMS{}
		if err = rows.Scan(&li.EPID, &li.Cover, &li.Title, &li.Subtitle, &li.PayStatus); err != nil {
			log.Error("EpMetas row.Scan error(%v)", err)
			return
		}
		metas[li.EPID] = li
	}
	return
}

// EpAuthDB pick ep data from DB for Cache missing case
func (d *Dao) EpAuthDB(c context.Context, cid int64) (ep *model.EpAuth, err error) {
	var row *xsql.Row
	if row = d.db.QueryRow(c, _simpleEPC, cid); err != nil {
		log.Error("d.db.QueryRow(%d) error(%v)", cid, err)
		return
	}
	ep = &model.EpAuth{}
	if err = row.Scan(&ep.ID, &ep.EPID, &ep.State, &ep.IsDeleted, &ep.Valid, &ep.SeasonID, &ep.NoMark); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			ep = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
	}
	return
}

// SnAuthDB .
func (d *Dao) SnAuthDB(c context.Context, cid int64) (s *model.SnAuth, err error) {
	var row *xsql.Row
	if row = d.db.QueryRow(c, _simpleSea, cid); err != nil {
		log.Error("d.db.QueryRow(%d) error(%v)", cid, err)
		return
	}
	s = &model.SnAuth{}
	if err = row.Scan(&s.ID, &s.IsDeleted, &s.Check, &s.Valid); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			s = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
	}
	return
}

// SnsAuthDB .
func (d *Dao) SnsAuthDB(c context.Context, sids []int64) (snsAuth map[int64]*model.SnAuth, err error) {
	snsAuth = make(map[int64]*model.SnAuth)
	rows, err := d.db.Query(c, fmt.Sprintf(_simpleSeas, xstr.JoinInts(sids)))
	if err != nil {
		log.Error("ArcMetas d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		s := &model.SnAuth{}
		if err = rows.Scan(&s.ID, &s.IsDeleted, &s.Check, &s.Valid); err != nil {
			log.Error("SnsAuthDB row.Scan error(%v)", err)
			return
		}
		snsAuth[s.ID] = s
	}
	return
}

// EpsAuthDB def.
func (d *Dao) EpsAuthDB(c context.Context, epids []int64) (epsAuth map[int64]*model.EpAuth, err error) {
	epsAuth = make(map[int64]*model.EpAuth)
	rows, err := d.db.Query(c, fmt.Sprintf(_simpleEPCs, xstr.JoinInts(epids)))
	if err != nil {
		log.Error("ArcMetas d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		ep := &model.EpAuth{}
		if err = rows.Scan(&ep.ID, &ep.EPID, &ep.State, &ep.IsDeleted, &ep.Valid, &ep.SeasonID, &ep.NoMark); err != nil {
			log.Error("SnsAuthDB row.Scan error(%v)", err)
			return
		}
		epsAuth[ep.EPID] = ep
	}
	return
}

// SeasonCMS gets the fields that can be changed from tv-cms side to offer the TV APP
func (d *Dao) SeasonCMS(c context.Context, sid int64) (season *model.SeasonCMS, err error) {
	var row *xsql.Row
	if row = d.db.QueryRow(c, _seasonCMS, sid); err != nil {
		log.Error("d.db.QueryRow(%d) error(%v)", sid, err)
		return
	}
	season = &model.SeasonCMS{}
	// select id, cover, `desc`, title , upinfo, category, area, play_time, role, staff, total_num, style, status
	if err = row.Scan(&season.SeasonID, &season.Cover, &season.Desc, &season.Title, &season.UpInfo,
		&season.Category, &season.Area, &season.Playtime, &season.Role, &season.Staff, &season.TotalNum,
		&season.Style, &season.OriginName, &season.Alias, &season.PayStatus); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			season = nil
		} else {
			log.Error("row.Scan(sid %d) error(%v)", sid, err)
		}
		return
	}
	// add newest info
	season.NewestEPID, season.NewestOrder, _ = d.NewestOrder(c, sid)
	return
}

// EpCMS gets the fields that can be changed from tv-cms side to offer the TV APP
func (d *Dao) EpCMS(c context.Context, epid int64) (ep *model.EpCMS, err error) {
	var row *xsql.Row
	if row = d.db.QueryRow(c, _epCMS, epid); err != nil {
		log.Error("d.db.QueryRow(%d) error(%v)", epid, err)
		return
	}
	ep = &model.EpCMS{}
	// select id, cover, `desc`, title , upinfo, category, area, play_time, role, staff, total_num, style
	if err = row.Scan(&ep.EPID, &ep.Cover, &ep.Title, &ep.Subtitle, &ep.PayStatus); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			ep = nil
		} else {
			log.Error("row.Scan(sid %d) error(%v)", epid, err)
		}
		return
	}
	return
}
