package result

import (
	"context"

	"go-common/app/job/main/archive/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_inVideoSQL = `INSERT INTO archive_video (aid,eptitle,description,filename,src_type,cid,duration,index_order,attribute,weblink,dimensions) VALUES(?,?,?,?,?,?,?,?,?,?,?)
				ON DUPLICATE KEY UPDATE eptitle=?,description=?,filename=?,src_type=?,duration=?,index_order=?,attribute=?,weblink=?,dimensions=?`
	_delVideoByCidSQL = "DELETE FROM archive_video where aid=? and cid=?"
	_delVideosSQL     = "DELETE FROM archive_video WHERE aid=?"
)

// TxAddVideo add videos result
func (d *Dao) TxAddVideo(c context.Context, tx *sql.Tx, v *archive.Video) (rows int64, err error) {
	res, err := tx.Exec(_inVideoSQL, v.Aid, v.Title, v.Desc, v.Filename, v.SrcType, v.Cid, v.Duration, v.Index, v.Attribute, v.WebLink, v.Dimensions,
		v.Title, v.Desc, v.Filename, v.SrcType, v.Duration, v.Index, v.Attribute, v.WebLink, v.Dimensions)
	if err != nil {
		log.Error("tx.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// TxDelVideoByCid del videos by aid and cid
func (d *Dao) TxDelVideoByCid(c context.Context, tx *sql.Tx, aid, cid int64) (rows int64, err error) {
	res, err := tx.Exec(_delVideoByCidSQL, aid, cid)
	if err != nil {
		log.Error("tx.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// TxDelVideos del videos
func (d *Dao) TxDelVideos(c context.Context, tx *sql.Tx, aid int64) (rows int64, err error) {
	res, err := tx.Exec(_delVideosSQL, aid)
	if err != nil {
		log.Error("tx.Exec(%s, %d) error(%v)", _delVideosSQL, aid, err)
		return
	}
	return res.RowsAffected()
}
