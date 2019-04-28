package archive

import (
	"context"

	"go-common/app/job/main/archive/model/archive"
	"go-common/library/log"
)

const (
	_videosSQL = `SELECT id,filename,cid,aid,eptitle,description,src_type,duration,filesize,resolutions,playurl,failinfo,
	index_order,attribute,xcode_state,status,ctime,mtime FROM archive_video WHERE aid=? ORDER BY index_order`
	_videos2SQL = `SELECT vr.id,v.filename,vr.cid,vr.aid,vr.title,vr.description,v.src_type,v.duration,v.filesize,v.resolutions,
	v.playurl,v.failcode,vr.index_order,v.attribute,v.xcode_state,v.status,vr.state,vr.ctime,vr.mtime,v.dimensions FROM archive_video_relation AS vr JOIN video AS v ON vr.cid=v.id WHERE vr.aid=? ORDER BY vr.index_order`
	_aidsSQL = "SELECT aid FROM archive_video_relation WHERE cid=? AND state=0"
)

// Videos get videos by aid
func (d *Dao) Videos(c context.Context, aid int64) (vs []*archive.Video, err error) {
	rows, err := d.db.Query(c, _videosSQL, aid)
	if err != nil {
		log.Error("d.db.Query(%s, %d) error(%v)", _videosSQL, aid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		v := &archive.Video{}
		if err = rows.Scan(&v.ID, &v.Filename, &v.Cid, &v.Aid, &v.Title, &v.Desc, &v.SrcType, &v.Duration, &v.Filesize, &v.Resolutions,
			&v.Playurl, &v.FailCode, &v.Index, &v.Attribute, &v.XcodeState, &v.Status, &v.CTime, &v.MTime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		vs = append(vs, v)
	}
	return
}

// Videos2 get videos by 2 table em.......
func (d *Dao) Videos2(c context.Context, aid int64) (vs []*archive.Video, err error) {
	rows, err := d.db.Query(c, _videos2SQL, aid)
	if err != nil {
		log.Error("d.db.Query(%s, %d) error(%v)", _videos2SQL, aid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		v := &archive.Video{}
		if err = rows.Scan(&v.ID, &v.Filename, &v.Cid, &v.Aid, &v.Title, &v.Desc, &v.SrcType, &v.Duration, &v.Filesize, &v.Resolutions,
			&v.Playurl, &v.FailCode, &v.Index, &v.Attribute, &v.XcodeState, &v.Status, &v.State, &v.CTime, &v.MTime, &v.Dimensions); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		vs = append(vs, v)
	}
	return
}

// Aids get aids by cid
func (d *Dao) Aids(c context.Context, cid int64) (aids []int64, err error) {
	rows, err := d.db.Query(c, _aidsSQL, cid)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var aid int64
		if err = rows.Scan(&aid); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		aids = append(aids, aid)
	}
	return
}
