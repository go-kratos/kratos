package service

import (
	"context"
	"fmt"
	"go-common/app/service/bbq/common"
	"go-common/library/conf/env"
	"go-common/library/ecode"
	"time"

	"go-common/app/interface/bbq/app-bbq/api/http/v1"
	"go-common/app/interface/bbq/app-bbq/model"
	"go-common/app/interface/bbq/app-bbq/model/grpc"
	rec "go-common/app/service/bbq/recsys/api/grpc/v1"
	user "go-common/app/service/bbq/user/api"
	video "go-common/app/service/bbq/video/api/grpc/v1"
	"go-common/library/log"
	"go-common/library/net/trace"
)

// SvList 短视屏推荐列表
func (s *Service) SvList(c context.Context, pageSize int64, mid int64, base *v1.Base, deviceID string) (res []*v1.VideoResponse, err error) {
	var (
		svids []int64
		svRes map[int64]*v1.VideoResponse
	)
	res = make([]*v1.VideoResponse, 0)
	//推荐列表
	svids, err = s.dao.RawRecList(c, pageSize, mid, base.BUVID)
	if err != nil || len(svids) == 0 || env.DeployEnv == env.DeployEnvUat {
		log.Warnv(c, log.KV("log", fmt.Sprintf("s.dao.GetList err[%v]", err)))
		// 降级
		svids, _ = s.dao.GetList(c, pageSize)
	}
	svRes, err = s.svInfos(c, svids, mid, false)
	for _, id := range svids {
		if sv, ok := svRes[id]; ok {
			if common.IsRecommendSvStateAvailable(int64(sv.State)) {
				res = append(res, sv)
			} else {
				log.Warnw(c, "log", "get error svid in recommend list", "svid", id, "mid", mid, "sv", sv)
			}
		}
	}
	return
}

// svPlays 批量获取playurl(相对地址方法)
func (s *Service) svPlays(c context.Context, svids []int64) map[int64]*v1.VideoPlay {
	var (
		relAddr []string
		err     error
		bvcUrls map[string]*grpc.VideoKeyItem
		bvcKeys map[int64][]*model.SVBvcKey
	)
	playMap := make(map[int64]*v1.VideoPlay)
	bvcKeys, err = s.dao.RawSVBvcKey(c, svids)
	if err != nil {
		log.Error("s.dao.RawSVBvcKey err[%v]", err)
	}
	for id, keys := range bvcKeys {
		playMap[id] = &v1.VideoPlay{
			SVID: id,
		}
		for k, v := range keys {
			if k == 0 {
				playMap[id].Quality = int64(v.CodeRate)
			}
			fi := &v1.FileInfo{
				TimeLength: v.Duration,
				FileSize:   v.FileSize,
				Path:       v.Path,
			}
			playMap[id].FileInfo = append(playMap[id].FileInfo, fi)
			playMap[id].SupportQuality = append(playMap[id].SupportQuality, int64(v.CodeRate))
			relAddr = append(relAddr, v.Path)
		}
	}
	bvcUrls, err = s.dao.RelPlayURLs(c, relAddr)
	if err != nil {
		log.Error("s.dao.RelPlayURLs err[%v]", err)
	}
	//拼装playurl
	for _, svid := range svids {
		if play, ok := playMap[svid]; ok {
			for fk, f := range play.FileInfo {
				if urls, ok := bvcUrls[f.Path]; ok {
					playMap[svid].ExpireTime = int64(urls.Etime)
					playMap[svid].CurrentTime = time.Now().Unix()
					for _, u := range urls.URL {
						if playMap[svid].FileInfo[fk].URL == "" {
							playMap[svid].FileInfo[fk].URL = u
							if playMap[svid].URL == "" {
								playMap[svid].URL = u
							}
							continue
						}
						if playMap[svid].FileInfo[fk].URLBc == "" {
							playMap[svid].FileInfo[fk].URLBc = u
							break
						}
					}
				} else {
					delete(playMap, svid)
					break
				}
			}
			playMap[svid] = play
		}
	}
	return playMap
}

// SvStatistics 视频统计服务
func (s *Service) SvStatistics(c context.Context, mid int64, svids []int64) (res []*v1.SvStatRes, err error) {
	var (
		stMap map[int64]*model.SvStInfo
		ulike map[int64]bool
		upIDs []int64
	)
	svInfos, _ := s.dao.RawVideos(c, svids)
	for _, sv := range svInfos {
		upIDs = append(upIDs, sv.MID)
	}
	stMap, err = s.dao.RawVideoStatistic(c, svids)
	if err != nil {
		log.Error("s.dao.RawVideoStatistic err[%v]", err)
	}
	//点赞状态
	ulike, err = s.dao.CheckUserLike(c, mid, svids)
	if err != nil {
		log.Error("s.dao.CheckUserLike err[%v]", err)
	}
	uflw, _ := s.dao.BatchUserInfo(c, mid, upIDs, false, false, true)
	if err != nil {
		log.Error("s.dao.IsFollow err[%v]", err)
	}
	for _, id := range svids {
		rp := &v1.SvStatRes{}
		rp.SVID = id
		if st, ok := stMap[id]; ok {
			rp.Like = st.Like
			rp.Share = st.Share
			rp.Play = st.Play
			rp.Subtitles = st.Subtitles
			rp.Reply = st.Reply
		}
		if l, ok := ulike[id]; ok {
			rp.IsLike = l
		}
		if sv, ok := svInfos[id]; ok {
			if f, ok2 := uflw[sv.MID]; ok2 {
				rp.FollowState = f.FollowState
			}
		}
		res = append(res, rp)
	}
	return
}

// SvCPlays 批量拉取playurl
func (s *Service) SvCPlays(c context.Context, svids []int64, mid int64) (res []*v1.VideoPlay, err error) {
	res = make([]*v1.VideoPlay, 0)

	//视频列表
	svRes, err := s.dao.RawVideos(c, svids)
	if err != nil {
		log.Error("s.dao.RawVideos err[%v]", err)
		return
	}

	avaliableSvids := make([]int64, 0)
	for _, v := range svRes {
		if (mid == 0 || v.MID != mid) && common.IsSvStateGuestAvailable(int64(v.State)) {
			avaliableSvids = append(avaliableSvids, v.SVID)
		} else if v.MID == mid && common.IsSvStateOwnerAvailable(int64(v.State)) {
			avaliableSvids = append(avaliableSvids, v.SVID)
		}
	}

	playMap := s.svPlays(c, avaliableSvids)
	for _, svid := range avaliableSvids {
		var play *v1.VideoPlay
		if p, ok := playMap[svid]; !ok {
			log.Warn("play不存在 svid[%d]", svid)
			continue
		} else {
			play = p
		}
		res = append(res, play)
	}
	return
}

// SvDetail 单条sv的视频详情，暂时只用于评论中转页
func (s *Service) SvDetail(c context.Context, svid int64, mid int64) (res *v1.VideoResponse, err error) {
	_, err = s.dao.VideoBase(c, mid, svid)
	if err != nil {
		return
	}

	svInfos, err := s.svInfos(c, []int64{svid}, mid, true)
	if err != nil {
		return
	}

	if val, exists := svInfos[svid]; exists {
		res = val
	} else {
		err = ecode.VideoUnExists
		log.Infow(c, "log", "not sv info", "svid", svid)
	}

	return
}

// svInfos 批量获取视频信息
// @params allowState 可放出状态，传空为app整体可露出状态
// @params needStInfo 是否需要视频统计数据
func (s *Service) svInfos(c context.Context, ids []int64, mid int64, needStInfo bool) (res map[int64]*v1.VideoResponse, err error) {
	var (
		mids  []int64
		svRes map[int64]*model.SvInfo
		ulike map[int64]bool
		stMap map[int64]*model.SvStInfo
	)
	res = make(map[int64]*v1.VideoResponse)
	stMap = make(map[int64]*model.SvStInfo)

	//视频列表
	svRes, err = s.dao.RawVideos(c, ids)
	if err != nil {
		log.Error("s.dao.RawVideos err[%v]", err)
		return
	}
	for _, v := range svRes {
		mids = append(mids, v.MID)
	}
	if mid != 0 {
		ulike, err = s.dao.CheckUserLike(c, mid, ids)
		if err != nil {
			log.Error("s.dao.CheckUserLike err[%v]", err)
		}
	}
	// query id
	tracer, _ := trace.FromContext(c)
	queryID := fmt.Sprintf("%s", tracer)
	//账号
	var userMap map[int64]*user.UserBase
	userMap, err = s.dao.JustGetUserBase(c, mids)
	if err != nil {
		log.Error("s.dao.UserBase err[%v]", err)
	}
	// play信息
	playMap := s.svPlays(c, ids)
	if needStInfo {
		stMap, err = s.dao.RawVideoStatistic(c, ids)
		if err != nil {
			log.Error("s.dao.RawVideoStatistic err[%v]", err)
		}
	}
	// extension信息
	extensions, tmpErr := s.getExtension(c, ids)
	if tmpErr != nil {
		log.Warnw(c, "log", "get extension fail")
	}

	for _, v := range svRes {
		if common.IsSvStateAvailable(int64(v.State)) {
			sv := &v1.VideoResponse{}
			if acc, ok := userMap[v.MID]; ok {
				sv.UserInfo = *acc
			}
			if lk, ok := ulike[v.SVID]; ok {
				sv.IsLike = lk
			}
			sv.SVID = v.SVID
			sv.Title = v.Title
			sv.Content = v.Content
			sv.MID = v.MID
			sv.Duration = v.Duration
			sv.Pubtime = v.Pubtime
			sv.Ctime = v.Ctime
			sv.AVID = v.AVID
			sv.CID = v.CID
			sv.From = v.From
			sv.CoverURL = v.CoverURL
			sv.CoverHeight = v.CoverHeight
			sv.CoverWidth = v.CoverWidth
			sv.QueryID = queryID
			sv.State = v.State
			if play, ok := playMap[v.SVID]; ok {
				sv.Play = *play
				res[v.SVID] = sv
			} else {
				log.Warn("play不存在 svid[%d]，此条记录直接舍弃", v.SVID)
			}
			if st, ok := stMap[v.SVID]; ok {
				sv.SvStInfo = *st
			}
			if extension, exists := extensions[v.SVID]; exists {
				sv.Extension = extension.Extension
			}

			res[v.SVID] = sv
		}
	}

	return
}

// SvRelRec 相关推荐服务
func (s *Service) SvRelRec(c context.Context, data *v1.SvRelReq) (res map[string]interface{}, err error) {
	res = make(map[string]interface{})
	var svMap map[int64]*v1.VideoResponse
	list := make([]*v1.VideoResponse, 0)
	relReq := &rec.RecsysRequest{
		SVID:       data.SVID,
		Offset:     data.Offset,
		Limit:      data.Limit,
		QueryID:    data.QueryID,
		App:        data.APP,
		AppVersion: data.APPVersion,
		BUVID:      data.BUVID,
		MID:        data.MID,
	}
	IDList, err := s.dao.RelRecList(c, relReq)
	if err != nil {
		err = nil
		return
	}
	svMap, err = s.svInfos(c, IDList, data.MID, false)
	if err != nil {
		err = nil
		return
	}
	for _, id := range IDList {
		if sv, ok := svMap[id]; ok {
			list = append(list, sv)
		}

	}
	res["list"] = list
	return
}

// SvDel 视频删除
func (s *Service) SvDel(c context.Context, in *video.VideoDeleteRequest) (interface{}, error) {
	return s.dao.SvDel(c, in)
}
