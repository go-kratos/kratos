package dao

import (
	"context"
	"go-common/app/service/bbq/common"
	video "go-common/app/service/bbq/video/api/grpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
)

// BatchVideoInfo .
func (d *Dao) BatchVideoInfo(c context.Context, svids []int64) (res map[int64]*video.VideoInfo, err error) {
	res = make(map[int64]*video.VideoInfo)

	videoReq := &video.ListVideoInfoRequest{SvIDs: svids}
	reply, err := d.videoClient.ListVideoInfo(c, videoReq)
	if err != nil {
		log.Errorw(c, "log", "call video service list vidoe info fail", "req", videoReq.String())
		return
	}

	for _, videoInfo := range reply.List {
		res[videoInfo.VideoBase.Svid] = videoInfo
	}
	log.V(10).Infow(c, "log", "batch video base", "video_info", res)

	return
}

// VideoBase 获取单个svid的VideoBase，不存在则会返回error
func (d *Dao) VideoBase(c context.Context, mid, svid int64) (res *video.VideoBase, err error) {
	videoInfos, err := d.BatchVideoInfo(c, []int64{svid})
	if err != nil {
		log.Warnw(c, "log", "batch fetch video info fail", "mid", mid, "svid", svid)
		return
	}
	if len(videoInfos) == 0 {
		err = ecode.VideoUnExists
		log.Warnw(c, "log", "get empty video base", "mid", mid, "svid", svid)
		return
	}
	videoInfo, exists := videoInfos[svid]
	if !exists {
		err = ecode.VideoUnExists
		log.Infow(c, "log", "get empty video base", "mid", mid, "svid", svid)
		return
	}

	res = videoInfo.VideoBase
	if res.State == common.VideoStPassReviewReject {
		log.Infow(c, "log", "video state in audit", "mid", mid, "svid", svid, "video_base", res)
		err = ecode.VideoInAudit
		return
	}
	if !common.IsSvStateAvailable(res.State) {
		err = ecode.VideoUnReachable
		log.Infow(c, "log", "video state not available", "mid", mid, "svid", svid, "video_base", res)
		return
	}

	if res.State == common.VideoStPassReviewReject && mid != res.Mid {
		err = ecode.VideoUnReachable
		log.Infow(c, "log", "video state only owner available", "mid", mid, "svid", svid, "video_base", res)
		return
	}

	return
}
