package view

import (
	"context"

	"go-common/app/interface/main/tv/model"
	"go-common/library/log"
)

// ArcMsg returns the arc auth msg
func (s *Service) ArcMsg(aid int64) (arc *model.ArcCMS, ok bool, msg string, err error) {
	if arc, err = s.cmsDao.LoadArcMeta(context.TODO(), aid); err != nil {
		log.Error("ArcMsg loadArcMeta aid %d, Err %v", aid, err)
		return
	}
	ok, msg = s.cmsDao.UgcErrMsg(arc.Deleted, arc.Result, arc.Valid)
	return
}

// VideoMsg returns the arc auth msg
func (s *Service) VideoMsg(ctx context.Context, cid int64) (ok bool, msg string, err error) {
	var video *model.VideoCMS
	if video, err = s.cmsDao.LoadVideoMeta(context.TODO(), cid); err != nil {
		log.Error("VideoMsg LoadVideoMeta aid %d, Err %v", cid, err)
		return
	}
	ok, msg = s.cmsDao.UgcErrMsg(video.Deleted, video.Result, video.Valid)
	if ok { // if video is normal, we also need to check it's archive
		_, ok, msg, err = s.ArcMsg(int64(video.AID))
	}
	return
}
