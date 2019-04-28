package dao

import (
	"context"
	"go-common/app/job/bbq/video/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_updatevideostatus    = "update video set state = ? where svid = ?"
	_updateCmsvideoStatus = "insert into cms_video (svid,cms_status,sv_status,`from`,title,pubtime,mid) values (?,?,?,?,?,?,?) on duplicate key update cms_status = values(cms_status),sv_status=values(sv_status),cms_uname = '',`from` = values(`from`),title=values(title),pubtime=values(pubtime),mid = values(mid)"
	_selectVideoRows      = "select id,svid,title,mid,`from`,pubtime from video where state = ? and id > ? limit 1000"
	_updateVR             = "update video_repository set state = ? where svid=?"
)

//UpdateCms ..
func (d *Dao) UpdateCms(c context.Context, v *model.VideoRaw) (err error) {
	if _, err = d.dbCms.Exec(c,
		_updateCmsvideoStatus,
		v.SVID,
		v.State,
		v.State,
		v.From,
		v.Title,
		v.Pubtime,
		v.MID,
	); err != nil {
		log.Error("DeliveryNewVdieoToCms insert cms_video err,svid : %v,err :%v", v.SVID, err)
		return
	}
	return
}

//TransToCheckBack ..
func (d *Dao) TransToCheckBack() (err error) {
	var (
		rows  *xsql.Rows
		count int64
		id    int64
		c     = context.Background()
	)
	for {
		if rows, err = d.db.Query(c, _selectVideoRows, model.VideoStPassReview, count); err != nil {
			log.Error("DeliveryNewVdieoToCms select video failed ,err:%v", err)
			return
		}
		flag := false
		for rows.Next() {
			videoinfo := model.VideoInfo{}
			if err = rows.Scan(
				&id,
				&videoinfo.SVID,
				&videoinfo.Title,
				&videoinfo.MID,
				&videoinfo.From,
				&videoinfo.Pubtime,
			); err != nil {
				if err == xsql.ErrNoRows {
					return
				}
				continue
			}
			count = id
			//满足运营导入规则
			if d.CmsRule(videoinfo.SVID) {
				if _, err = d.dbCms.Exec(c,
					_updateCmsvideoStatus,
					videoinfo.SVID,
					model.VideoStCheckBack,
					model.VideoStCheckBack,
					videoinfo.From,
					videoinfo.Title,
					videoinfo.Pubtime,
					videoinfo.MID,
				); err != nil {
					log.Error("DeliveryNewVdieoToCms insert cms_video err,svid : %v,err :%v", videoinfo.SVID, err)
					continue
				}
				if _, err = d.db.Exec(c,
					_updatevideostatus,
					model.VideoStCheckBack,
					videoinfo.SVID,
				); err != nil {
					log.Error("DeliveryNewVdieoToCms update video status err : %v,svid : %v", err, videoinfo.SVID)
					continue
				}
				if _, err = d.dbCms.Exec(c, _updateVR, model.VideoStCheckBack, videoinfo.SVID); err != nil {
					log.Error("DeliveryNewVdieoToCms update vr err :%v,svid : %v", err, videoinfo.SVID)
					continue
				}
			}
			flag = true
		}
		rows.Close()
		if !flag {
			return
		}
	}
}

// CmsRule ...
func (d *Dao) CmsRule(svid int64) (flag bool) {
	return true
}

//TransToReview ...
func (d *Dao) TransToReview() (err error) {
	var (
		rows  *xsql.Rows
		count int64
		id    int64
		c     = context.Background()
	)
	for {
		if rows, err = d.db.Query(c, _selectVideoRows, model.VideoStPendingPassReview, count); err != nil {
			log.Error("TransToReview select video failed ,err:%v", err)
			continue
		}
		flag := false
		for rows.Next() {
			videoinfo := model.VideoInfo{}
			if err = rows.Scan(
				&id,
				&videoinfo.SVID,
				&videoinfo.Title,
				&videoinfo.MID,
				&videoinfo.From,
				&videoinfo.Pubtime,
			); err != nil {
				if err == xsql.ErrNoRows {
					return
				}
				continue
			}
			count = id
			//满足运营导入规则
			if d.CmsRule(videoinfo.SVID) {
				var st int
				if videoinfo.From == model.VideoFromBILI || videoinfo.From == model.VideoFromCMS {
					st = model.VideoStPassReview
				} else {
					st = model.VideoStPassReviewReject
				}

				if _, err = d.dbCms.Exec(c,
					_updateCmsvideoStatus,
					videoinfo.SVID,
					st,
					st,
					videoinfo.From,
					videoinfo.Title,
					videoinfo.Pubtime,
					videoinfo.MID,
				); err != nil {
					log.Error("TransToReview insert cms_video err,svid : %v,err :%v", videoinfo.SVID, err)
					continue
				}
				if _, err = d.db.Exec(c,
					_updatevideostatus,
					st,
					videoinfo.SVID,
				); err != nil {
					log.Error("TransToReview update video status err : %v,svid : %v", err, videoinfo.SVID)
					continue
				}
			}
			flag = true
		}
		rows.Close()
		if !flag {
			return
		}
	}
}
