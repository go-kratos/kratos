package archive

import (
	"context"
	"go-common/library/log"
	"time"

	"database/sql"
	"go-common/app/job/main/videoup-report/model/archive"

	farm "github.com/dgryski/go-farm"
)

const (
	_updatedFilenamesByTime = "SELECT filename FROM video WHERE mtime >= ? AND mtime < ?"
	_videos2SQL             = `SELECT vr.id,v.filename,vr.cid,vr.aid,vr.title,vr.description,v.src_type,v.duration,v.filesize,v.resolutions,
	v.playurl,v.failcode,vr.index_order,v.attribute,v.xcode_state,v.status,vr.state,vr.ctime,vr.mtime FROM archive_video_relation AS vr JOIN video AS v ON vr.cid=v.id WHERE vr.aid=? ORDER BY vr.index_order`
	_newVideoByFnSQL = `SELECT avr.id,v.filename,avr.cid,avr.aid,avr.title,avr.description,v.src_type,v.duration,v.filesize,v.resolutions,v.playurl,v.failcode,
	avr.index_order,v.attribute,v.xcode_state,avr.state,v.status,avr.ctime,avr.mtime FROM archive_video_relation avr JOIN video v on avr.cid = v.id
	WHERE hash64=? AND filename=?`
)

// UpdatedFilenames Get updated video's filename between stime and etime.
func (d *Dao) UpdatedFilenames(c context.Context, stime, etime time.Time) (fns []string, err error) {
	rows, err := d.db.Query(c, _updatedFilenamesByTime, stime, etime)
	if err != nil {
		log.Error("d.UpdatedFilenames.Query(%v,%v) error(%v)", stime, etime, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		fn := ""
		if err = rows.Scan(&fn); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		fns = append(fns, fn)
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
			&v.Playurl, &v.FailCode, &v.Index, &v.Attribute, &v.XcodeState, &v.Status, &v.State, &v.CTime, &v.MTime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		vs = append(vs, v)
	}
	return
}

// NewVideo get video info by filename.
func (d *Dao) NewVideo(c context.Context, filename string) (v *archive.Video, err error) {
	hash64 := int64(farm.Hash64([]byte(filename)))
	row := d.db.QueryRow(c, _newVideoByFnSQL, hash64, filename)
	v = &archive.Video{}
	var avrState, vState int16
	if err = row.Scan(&v.ID, &v.Filename, &v.Cid, &v.Aid, &v.Title, &v.Desc, &v.SrcType, &v.Duration, &v.Filesize, &v.Resolutions,
		&v.Playurl, &v.FailCode, &v.Index, &v.Attribute, &v.XcodeState, &avrState, &vState, &v.CTime, &v.MTime); err != nil {
		if err == sql.ErrNoRows {
			v = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	// 2 state map to 1
	if avrState == archive.VideoStatusDelete {
		v.Status = archive.VideoStatusDelete
	} else {
		v.Status = vState
	}
	return
}
