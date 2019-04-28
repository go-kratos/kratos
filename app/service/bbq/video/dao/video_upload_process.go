package dao

import (
	"context"
	"fmt"
	"go-common/app/service/bbq/video/api/grpc/v1"
	"go-common/app/service/bbq/video/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	insertOrUpdateVUP = "insert into video_upload_process (`svid`,`title`,`mid`,`upload_status`,`retry_times`,`home_img_url`,`home_img_width`,`home_img_height`) values(?,?,?,?,?,?,?,?) on duplicate key update `title` = values(`title`),`mid` = values(`mid`),`upload_status` = values(`upload_status`),`retry_times`= values(`retry_times`)"
	selectPrepareVUP  = "select `svid`,`title`,`upload_status`,`home_img_url`,`home_img_height`,`home_img_width` from video_upload_process where mid=%d and is_deleted=0 and upload_status != 1 order by ctime desc limit 20"
)

//InsertOrUpdateVUP ..
func (d *Dao) InsertOrUpdateVUP(c context.Context, vup *model.VideoUploadProcess) (err error) {
	if vup == nil {
		err = ecode.BBQSystemErr
		log.Errorw(c, "event", "InsertVUP req nil")
		return
	}
	if _, err = d.db.Exec(c, insertOrUpdateVUP, vup.SVID, vup.Title, vup.Mid, vup.UploadStatus, vup.RetryTimes, vup.HomeImgURL, vup.HomeImgWidth, vup.HomeImgHeight); err != nil {
		log.Errorw(c, "event", "InsertVR err", "err", err, "param", vup)
		return
	}
	return
}

// GetPrepareVUP 获取数据
func (d *Dao) GetPrepareVUP(c context.Context, mid int64) (vups []*v1.UploadingVideo, err error) {
	querySQL := fmt.Sprintf(selectPrepareVUP, mid)
	rows, err := d.db.Query(c, querySQL)
	if err != nil {
		log.Errorw(c, "log", "get prepare vup fail", "mid", mid)
		return
	}
	defer rows.Close()

	for rows.Next() {
		vup := new(v1.UploadingVideo)
		if err = rows.Scan(&vup.Svid, &vup.Title, &vup.UploadStatus, &vup.HomeImgUrl, &vup.HomeImgHeight, &vup.HomeImgWidth); err != nil {
			log.Errorw(c, "log", "scan vup fail")
			return
		}
		vups = append(vups, vup)
	}

	return
}
