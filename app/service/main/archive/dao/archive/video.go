package archive

import (
	"context"
	"database/sql"
	"fmt"

	"go-common/app/service/main/archive/api"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_vdosSQL       = "SELECT cid,src_type,index_order,eptitle,duration,filename,weblink,dimensions FROM archive_video WHERE aid=? ORDER BY index_order"
	_vdosByAidsSQL = "SELECT aid,cid,src_type,index_order,eptitle,duration,filename,weblink,dimensions FROM archive_video WHERE aid in (%s) ORDER BY index_order"
	_vdosByCidsSQL = "SELECT aid,cid,src_type,index_order,eptitle,duration,filename,weblink FROM archive_video WHERE cid IN (%s)"
	_vdoSQL        = "SELECT cid,src_type,index_order,eptitle,duration,filename,weblink,description,dimensions FROM archive_video WHERE aid=? AND cid=?"
	_firstCidSQL   = "SELECT cid,src_type,dimensions FROM archive_video WHERE aid=? ORDER BY index_order LIMIT 1"
)

// firstCid get first cid
func (d *Dao) firstCid(c context.Context, aid int64) (cid int64, dimensions string, err error) {
	var srcType string
	row := d.resultDB.QueryRow(c, _firstCidSQL, aid)
	if err = row.Scan(&cid, &srcType, &dimensions); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	if srcType != "vupload" {
		cid = 0
	}
	return
}

// videos3 get videos by aid.
func (d *Dao) videos3(c context.Context, aid int64) (ps []*api.Page, err error) {
	d.infoProm.Incr("videos3")
	rows, err := d.vdosStmt.Query(c, aid)
	if err != nil {
		d.errProm.Incr("archive_db")
		log.Error("d.vdosStmt.Query(%d) error(%v)", aid, err)
		return
	}
	defer rows.Close()
	var page = int32(0)
	for rows.Next() {
		var (
			p          = &api.Page{}
			fn         string
			dimensions string
		)
		page++
		if err = rows.Scan(&p.Cid, &p.From, &p.Page, &p.Part, &p.Duration, &fn, &p.WebLink, &dimensions); err != nil {
			log.Error("rows.Scan error(%v)", err)
			d.errProm.Incr("archive_db")
			return
		}
		p.Page = page
		if p.From != "vupload" {
			p.Vid = fn
		}
		p.FillDimension(dimensions)
		ps = append(ps, p)
	}
	return
}

// videosByAids3 get videos by aids
func (d *Dao) videosByAids3(c context.Context, aids []int64) (vs map[int64][]*api.Page, err error) {
	d.infoProm.Incr("videosByAids3")
	rows, err := d.resultDB.Query(c, fmt.Sprintf(_vdosByAidsSQL, xstr.JoinInts(aids)))
	if err != nil {
		d.errProm.Incr("archive_db")
		log.Error("d.resultDB.Query(%s) error(%v)", fmt.Sprintf(_vdosByAidsSQL, xstr.JoinInts(aids)), err)
		return
	}
	vs = make(map[int64][]*api.Page, len(aids))
	var pages = make(map[int64]int32, len(aids))
	defer rows.Close()
	for rows.Next() {
		var (
			p          = &api.Page{}
			aid        int64
			fn         string
			dimensions string
		)
		if err = rows.Scan(&aid, &p.Cid, &p.From, &p.Page, &p.Part, &p.Duration, &fn, &p.WebLink, &dimensions); err != nil {
			d.errProm.Incr("archive_db")
			log.Error("rows.Scan error(%v)", err)
			return
		}
		pages[aid]++
		p.Page = pages[aid]
		if p.From != "vupload" {
			p.Vid = fn
		}
		p.FillDimension(dimensions)
		vs[aid] = append(vs[aid], p)
	}
	return
}

// VideosByCids get videos by cids.
func (d *Dao) VideosByCids(c context.Context, cids []int64) (vs map[int64]map[int64]*api.Page, err error) {
	if len(cids) == 0 {
		return
	}
	d.infoProm.Incr("videosByCids")
	rows, err := d.arcReadDB.Query(c, fmt.Sprintf(_vdosByCidsSQL, xstr.JoinInts(cids)))
	if err != nil {
		log.Error("d.arcReadDB.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	vs = make(map[int64]map[int64]*api.Page, len(cids))
	for rows.Next() {
		var (
			aid int64
			p   = &api.Page{}
			fn  string
		)
		if err = rows.Scan(&aid, &p.Cid, &p.From, &p.Page, &p.Part, &p.Duration, &fn, &p.WebLink); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		if p.From != "vupload" {
			p.Vid = fn
		}
		if _, ok := vs[aid]; !ok {
			vs[aid] = make(map[int64]*api.Page)
		}
		vs[aid][int64(p.Cid)] = p
	}
	return
}

// video3 get video by aid & cid.
func (d *Dao) video3(c context.Context, aid, cid int64) (p *api.Page, err error) {
	d.infoProm.Incr("video3")
	var fn, dimension string
	row := d.resultDB.QueryRow(c, _vdoSQL, aid, cid)
	p = &api.Page{}
	if err = row.Scan(&p.Cid, &p.From, &p.Page, &p.Part, &p.Duration, &fn, &p.WebLink, &p.Desc, &dimension); err != nil {
		if err == sql.ErrNoRows {
			p = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
			d.errProm.Incr("result_db")
		}
		return
	}
	if p.From != "vupload" {
		p.Vid = fn
	}
	p.FillDimension(dimension)
	return
}
