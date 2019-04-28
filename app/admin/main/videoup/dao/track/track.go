package track

import (
	"context"

	"go-common/app/admin/main/videoup/model/track"
	"go-common/library/log"
)

const (
	_arcTrackSQL = "SELECT state,round,remark,ctime,attribute FROM archive_track WHERE aid=? ORDER BY id DESC"
	_vdoTrackSQL = "SELECT aid,xcode_state,status,remark,ctime FROM archive_video_track WHERE filename=? AND aid=? ORDER BY id DESC"
)

// ArchiveTrack get archive track info.
func (d *Dao) ArchiveTrack(c context.Context, aid int64) (res []*track.Archive, err error) {
	rows, err := d.db.Query(c, _arcTrackSQL, aid)
	if err != nil {
		log.Error("ArchiveTrack d.arcTrackStmt.Query(%d) error(%v)", aid, err)
		return
	}
	defer rows.Close()
	res = make([]*track.Archive, 0)
	for rows.Next() {
		a := &track.Archive{}
		if err = rows.Scan(&a.State, &a.Round, &a.Remark, &a.Timestamp, &a.Attribute); err != nil {
			log.Error("ArchiveTrack rows.Scan() error(%v)", err)
			return
		}
		res = append(res, a)
	}
	return
}

// VideoTrack get video track info.
func (d *Dao) VideoTrack(c context.Context, filename string, aid int64) (res []*track.Video, er error) {
	rows, err := d.db.Query(c, _vdoTrackSQL, filename, aid)
	if err != nil {
		log.Error("VideoTrack d.vdoTrackStmt.Query(%s, %d) error(%v)", filename, aid, err)
		return
	}
	defer rows.Close()
	res = make([]*track.Video, 0)
	for rows.Next() {
		v := &track.Video{}
		if err = rows.Scan(&v.AID, &v.XCodeState, &v.Status, &v.Remark, &v.Timestamp); err != nil {
			log.Error("VideoTrack rows.Scan() error(%v)", err)
			return
		}
		res = append(res, v)
	}
	return
}
