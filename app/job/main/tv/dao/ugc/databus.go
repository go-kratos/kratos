package ugc

import (
	"context"

	ugcmdl "go-common/app/job/main/tv/model/ugc"
	arccli "go-common/app/service/main/archive/api"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_updateArc = "UPDATE ugc_archive SET title = ?, cover = ?, content = ?, pubtime = ?, " +
		"typeid = ?, submit = ?, state = ? WHERE aid = ? AND deleted = 0"
	_updateVideo = "UPDATE ugc_video SET eptitle = ?, index_order = ?, submit = ? WHERE cid = ? AND deleted = 0"
	_needSubmit  = 1
	_pickVideos  = "SELECT id, eptitle, cid, index_order FROM ugc_video WHERE aid = ? AND deleted = 0"
	_upInList    = "SELECT mid FROM ugc_uploader WHERE mid = ? AND deleted = 0"
	_setUploader = "REPLACE INTO ugc_uploader (mid, state) VALUES (?,?)"
)

// UpInList checks whether the upper is in our list
func (d *Dao) UpInList(c context.Context, mid int64) (realID int64, err error) {
	if err = d.DB.QueryRow(c, _upInList, mid).Scan(&realID); err != nil { // get the qualified aid to sync
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("d.UpInList.Query error(%v)", err)
	}
	return
}

// UpdateArc updates the key fields of an archive, used for databus monitoring
func (d *Dao) UpdateArc(c context.Context, arc *ugcmdl.ArchDatabus) (err error) {
	if _, err = d.DB.Exec(c, _updateArc,
		arc.Title, arc.Cover, arc.Content, arc.PubTime, arc.TypeID, _needSubmit, arc.State, arc.Aid); err != nil {
		log.Error("UpdateArc, failed to update: (%v), Error: %v", arc, err)
	}
	return
}

// TxUpdateVideo updates the ugc video's status, for databus update
func (d *Dao) TxUpdateVideo(tx *sql.Tx, video *arccli.Page) (err error) {
	if _, err = tx.Exec(_updateVideo,
		video.Part, video.Page, _needSubmit, video.Cid); err != nil {
		log.Error("TxUpdateVideo, failed to update: (%v), Error: %v", video, err)
	}
	return
}

// PickVideos picks the videos of an archive in one shot
func (d *Dao) PickVideos(c context.Context, aid int64) (res map[int64]*ugcmdl.SimpleVideo, err error) {
	var rows *sql.Rows
	res = make(map[int64]*ugcmdl.SimpleVideo)
	if rows, err = d.DB.Query(c, _pickVideos, aid); err != nil {
		log.Error("d.Import.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r = &ugcmdl.SimpleVideo{}
		if err = rows.Scan(&r.ID, &r.Eptitle, &r.CID, &r.IndexOrder); err != nil {
			log.Error("PickVideos row.Scan() error(%v)", err)
			return
		}
		res[r.CID] = r
	}
	if err = rows.Err(); err != nil {
		log.Error("d.PickVideos.Query error(%v)", err)
	}
	return
}

// TxAddVideos add into the db the new videos
func (d *Dao) TxAddVideos(tx *sql.Tx, pages []*arccli.Page, aid int64) (err error) {
	for _, v := range pages {
		if _, err = tx.Exec(_importVideo, aid, v.Cid, v.Part, v.Page, v.Duration, v.Desc); err != nil {
			log.Error("_importArc, failed to insert: (%v), Error: %v", v, err)
			return
		}
	}
	return
}

// TxUpAdd adds the upper
func (d *Dao) TxUpAdd(tx *sql.Tx, mid int64) (err error) {
	if _, err = tx.Exec(_setUploader, mid, 1); err != nil {
		log.Error("UpAdd Error %v", err)
	}
	return
}
