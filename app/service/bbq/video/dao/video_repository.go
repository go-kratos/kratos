package dao

import (
	"context"
	"go-common/app/service/bbq/video/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	insertVR           = "insert into video_repository (`cid`,`svid`,`mid`,`title`,`from`,`sync_status`) values (?,?,?,?,?,?)"
	updateVRSyncStatus = "update video_repository set sync_status = ?"
	queryVRBySvid      = "select `title`,`mid`,`home_img_url`,`home_img_width`,`home_img_height` from video_repository where svid = ?"
)

//InsertVR ..
func (d *Dao) InsertVR(c context.Context, vr *model.VideoRepository) (err error) {
	if vr == nil {
		err = ecode.BBQSystemErr
		log.Errorw(c, "event", "InsertVR req nil")
		return
	}
	if _, err = d.cmsdb.Exec(c, insertVR, vr.SVID, vr.SVID, vr.MID, vr.Title, vr.From,
		vr.SyncStatus); err != nil {
		log.Errorw(c, "event", "InsertVR err", "err", err, "param", vr)
		return
	}
	return
}

//UpdateVR ..
func (d *Dao) UpdateVR(c context.Context, vr *model.VideoRepository) (err error) {
	if vr == nil {
		err = ecode.BBQSystemErr
		log.Errorw(c, "event", "InsertVR req nil")
		return
	}
	if _, err = d.cmsdb.Exec(c, updateVRSyncStatus, vr.SyncStatus); err != nil {
		log.Errorw(c, "event", "UpdateVR err", "err", err, "param", vr)
		return
	}
	return
}

//QueryVR ..
func (d *Dao) QueryVR(c context.Context, vr *model.VideoRepository) (res *model.VideoRepository, err error) {
	if vr == nil {
		err = ecode.BBQSystemErr
		log.Errorw(c, "event", "InsertVR req nil")
		return
	}
	res = new(model.VideoRepository)
	if err = d.cmsdb.QueryRow(c, queryVRBySvid, vr.SVID).Scan(&res.Title, &res.MID, &res.HomeImgURL, &res.HomeImgWidth, &res.HomeImgHeight); err != nil {
		log.Errorw(c, "event", "queryVR scan err", "err", err, "param", vr)
		return
	}
	return
}

//HomeImgCreate ..
func (d *Dao) HomeImgCreate(c context.Context, vr *model.VideoRepository) (err error) {
	if vr == nil {
		err = ecode.BBQSystemErr
		log.Errorw(c, "event", "HomeImgCreate req nil")
		return
	}
	if _, err = d.cmsdb.Exec(c, "update video_repository set home_img_url = ?,home_img_width = ? ,home_img_height = ? where svid = ? and mid = ?",
		vr.HomeImgURL, vr.HomeImgWidth, vr.HomeImgHeight, vr.SVID, vr.MID); err != nil {
		log.Errorw(c, "update home_img err", "err", err, "param", vr)
		return
	}
	return
}
