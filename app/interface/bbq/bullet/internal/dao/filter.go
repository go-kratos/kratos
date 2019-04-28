package dao

import (
	"context"
	user "go-common/app/service/bbq/user/api"
	video "go-common/app/service/bbq/video/api/grpc/v1"
	filter "go-common/app/service/main/filter/api/grpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
)

// 用于filter
const (
	FilterAreaAccount = "BBQ_account"
	FilterAreaVideo   = "BBQ_video"
	FilterAreaReply   = "BBQ_reply"
	FilterAreaSearch  = "BBQ_search"
	FilterAreaDanmu   = "BBQ_danmu"
	FilterLevel       = 20
)

// Filter .
func (d *Dao) Filter(ctx context.Context, content string, area string) (level int32, err error) {
	req := new(filter.FilterReq)
	req.Message = content
	req.Area = area
	reply, err := d.filterClient.Filter(ctx, req)
	if err != nil {
		log.Errorv(ctx, log.KV("log", "filter fail : req="+req.String()))
		return
	}
	level = reply.Level
	log.V(1).Infov(ctx, log.KV("log", "get filter reply="+reply.String()))
	return
}

// PhoneCheck .
func (d *Dao) PhoneCheck(c context.Context, mid int64) (telStatus int32, err error) {
	req := &user.PhoneCheckReq{Mid: mid}
	res, err := d.userClient.PhoneCheck(c, req)
	if err != nil {
		log.Errorw(c, "log", "call phone check fail", "err", err, "mid", mid)
		return
	}
	telStatus = res.TelStatus
	return
}

// batchVideoInfo .
func (d *Dao) batchVideoInfo(c context.Context, svids []int64) (res map[int64]*video.VideoInfo, err error) {
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
	log.V(1).Infow(c, "log", "batch video base", "video_info", res)

	return
}

// VideoBase 获取单个svid的VideoBase，不存在则会返回error
func (d *Dao) VideoBase(c context.Context, svid int64) (res *video.VideoBase, err error) {
	videoInfos, err := d.batchVideoInfo(c, []int64{svid})
	if err != nil {
		log.Warnw(c, "log", "batch fetch video info fail", "svid", svid)
		return
	}
	if len(videoInfos) == 0 {
		err = ecode.VideoUnExists
		return
	}
	videoInfo, exists := videoInfos[svid]
	if !exists {
		err = ecode.VideoUnExists
		return
	}
	res = videoInfo.VideoBase
	return
}
