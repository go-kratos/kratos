package service

import (
	"context"
	"encoding/json"
	"fmt"
	"go-common/app/service/bbq/common"
	"time"

	v1 "go-common/app/interface/bbq/app-bbq/api/http/v1"
	"go-common/app/interface/bbq/app-bbq/model"

	video "go-common/app/service/bbq/video/api/grpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/trace"
	xtime "go-common/library/time"
)

// FeedUpdateNum 关注页红点
func (s *Service) FeedUpdateNum(c context.Context, mid int64) (res v1.FeedUpdateNumResponse, err error) {

	// 0. 获取关注链
	followedMid, err := s.dao.FetchFollowList(c, mid)
	if err != nil {
		log.Errorv(c, log.KV("log", "fetch follow fail"), log.KV("mid", mid), log.KV("err", err))
		return
	}
	if len(followedMid) == 0 {
		res.Num = 0
		log.V(1).Infov(c, log.KV("log", "no_follow"), log.KV("mid", mid))
		return
	}

	// 1. 获取mid上次浏览点
	pubTime, _ := s.dao.GetMIDLastPubtime(c, mid)

	// 2. 获取新的视频
	newlySvID, err := s.dao.FetchAvailableOutboxList(c, common.FeedStates, followedMid, false, 0, xtime.Time(pubTime), 1)
	if err != nil {
		log.Errorv(c, log.KV("log", "fetch available outbox list fail"), log.KV("mid", mid), log.KV("pubTime", pubTime))
		return
	}

	// 3. form rsp
	if len(newlySvID) == 0 {
		res.Num = 0
	} else {
		res.Num = 1
	}

	return
}

// FeedList 关注页短视屏列表
func (s *Service) FeedList(c context.Context, req *v1.FeedListRequest) (res *v1.FeedListResponse, err error) {
	var (
		mid     = req.MID
		markStr = req.Mark
	)
	res = new(v1.FeedListResponse)
	res.List = make([]*v1.SvDetail, 0)
	res.RecList = make([]*v1.SvDetail, 0)
	// 0.前期校验
	// 解析mark，获取last_svid
	var mark model.FeedMark
	isFirstPage := false
	// 关注up的总视频列表是否空
	needRecList := false
	if len(markStr) != 0 {
		var markData = []byte(markStr)
		err = json.Unmarshal(markData, &mark)
		if err != nil {
			err = ecode.ReqParamErr
			log.Errorv(c, log.KV("log", "mark_unmarshal"), log.KV("mark", markStr))
			return
		}
		log.V(1).Infov(c, log.KV("mark", markData))
	}
	if mark.IsRec {
		needRecList = true
	} else if mark.LastSvID == 0 {
		isFirstPage = true
		mark.LastSvID = model.MaxInt64
		mark.LastPubtime = xtime.Time(time.Now().Unix())
		// 更新最新浏览svid
		s.dao.SetMIDLastPubtime(c, mid, int64(mark.LastPubtime))
	}

	// 1. 需要获取详情的svid列表
	var svids []int64
	if !needRecList {
		// 1.获取关注链
		var followedMid []int64
		followedMid, err = s.dao.FetchFollowList(c, mid)
		if err != nil {
			log.Errorv(c, log.KV("log", "fetch_follow"), log.KV("mid", mid))
			return
		}
		// 无关注人，直接返回
		if len(followedMid) == 0 {
			res.HasMore = false
			needRecList = true
			log.V(1).Infov(c, log.KV("log", "no_follow"), log.KV("mid", mid))
		} else {
			// 2.获取svid列表
			svids, err = s.dao.FetchAvailableOutboxList(c, common.FeedStates, followedMid, true, mark.LastSvID, mark.LastPubtime, model.FeedListLen)
			if err != nil {
				log.Warnw(c, "log", "fetch available outbox list fail")
				return
			}
			// 为了保护列表，所以/2，后面切换获取逻辑就可以去掉了
			if len(svids) < model.FeedListLen/2 {
				res.HasMore = false
				// 关注人了，但是这些人没有发布过视频
				if len(svids) == 0 && isFirstPage {
					needRecList = true
				}
			} else {
				res.HasMore = true
			}
		}
	}
	// 如果第一页就是empty或者mark携带了is_rec，需要为用户推荐一些视频
	if needRecList {
		svids, err = s.dao.AttentionRecList(c, model.FeedListLen, mid, req.BUVID)
		if err != nil {
			log.Warnw(c, "log", "get attention feed rec fail")
			return
		}
	}

	// 2.获取sv信息列表
	var list []*v1.SvDetail
	// 2.0 获取sv详情
	detailMap, err := s.getVideoDetail(c, req.MID, req.Qn, req.Device, svids, true)
	if err != nil {
		log.Warnv(c, log.KV("log", "get video detail fail"))
		return
	} else if len(detailMap) == 0 {
		log.Warnv(c, log.KV("log", "feed list empty"), log.KV("svid_num", len(svids)))
		return
	}
	// 2.1 获取热评
	hots, hotErr := s.dao.ReplyHot(c, mid, svids)
	if hotErr != nil {
		log.Warnv(c, log.KV("log", "get hot reply fail"))
	}
	log.V(1).Infov(c, log.KV("log", "get_video_detail"), log.KV("req_size", len(svids)),
		log.KV("rsp_size", len(detailMap)))
	for _, svID := range svids {
		v, exists := detailMap[svID]
		if exists {
			if hots, ok := hots[svID]; ok {
				v.HotReply.Hots = hots
			}
			list = append(list, v)
		} else {
			log.Warnv(c, log.KV("log", "sv_not_found"), log.KV("mid", mid), log.KV("svid", svID))
		}
	}

	// 3. 组装回包，判断往哪个list塞数据
	var nextMark model.FeedMark
	if needRecList {
		res.RecList = list
		nextMark.IsRec = true
	} else {
		res.List = list
		if len(list) > 0 && res.HasMore {
			nextMark.LastSvID = list[len(list)-1].SVID
			nextMark.LastPubtime = list[len(list)-1].Pubtime
		}
	}
	jsonStr, _ := json.Marshal(nextMark) // marshal的时候相信库函数，不做err判断
	res.Mark = string(jsonStr)

	return
}

// SpaceSvList 个人空间视频列表
func (s *Service) SpaceSvList(c context.Context, req *v1.SpaceSvListRequest) (res *v1.SpaceSvListResponse, err error) {
	// 0.前期校验
	// 这里就不校验up主是否存在
	res = new(v1.SpaceSvListResponse)
	res.List = make([]*v1.SvDetail, 0)
	upMid := req.UpMid
	if upMid == 0 {
		err = ecode.ReqParamErr
		log.Errorv(c, log.KV("log", "up mid is 0"), log.KV("up_mid", 0))
		return
	}
	// parseCursor
	cursor, cursorNext, err := parseCursor(req.CursorPrev, req.CursorNext)
	if err != nil {
		return
	}

	// 1. 如果是主人态，其第一页，则进行额外prepare_list
	if req.MID == req.UpMid && len(req.CursorNext) == 0 && len(req.CursorPrev) == 0 {
		prepareRes, tmpErr := s.videoClient.ListPrepareVideo(c, &video.PrepareVideoRequest{Mid: req.MID})
		if tmpErr != nil {
			log.Warnw(c, "log", "get prepare video fail", "mid", req.MID)
		} else {
			res.PrepareList = prepareRes.List
		}
	}

	// 2. 获取svid列表
	states := common.SpaceFanStates
	if req.MID == req.UpMid {
		states = common.SpaceOwnerStates
	}
	svids, err := s.dao.FetchAvailableOutboxList(c, states, []int64{upMid}, cursorNext, cursor.CursorID, cursor.CursorTime, req.Size)
	if err != nil {
		log.Infov(c, log.KV("log", "fetch_outbox_list"), log.KV("error", err))
		return
	}
	if len(svids) < req.Size/2 {
		res.HasMore = false
		if len(svids) == 0 {
			return
		}
	} else {
		res.HasMore = true
	}

	// 3.获取sv详情
	detailMap, err := s.svInfos(c, svids, req.MID, true)
	if err != nil {
		log.Errorv(c, log.KV("log", "get video detail fail"))
		return
	} else if len(detailMap) == 0 {
		log.Warnv(c, log.KV("log", "feed list empty"), log.KV("svid_num", len(svids)))
		return
	}
	for _, svID := range svids {
		item, exists := detailMap[svID]
		if exists {
			sv := new(v1.SvDetail)
			sv.VideoResponse = *item
			res.List = append(res.List, sv)
		} else {
			log.Warnv(c, log.KV("log", "sv_not_found"), log.KV("svid", svID))
		}
	}

	// query id
	tracer, _ := trace.FromContext(c)
	queryID := fmt.Sprintf("%s", tracer)

	// 4. 后处理，为每个item添加cursor值
	var itemCursor model.CursorValue
	for _, item := range res.List {
		itemCursor.CursorID = item.SVID
		itemCursor.CursorTime = item.Pubtime
		jsonStr, _ := json.Marshal(itemCursor) // marshal的时候相信库函数，不做err判断
		item.CursorValue = string(jsonStr)
		item.QueryID = queryID
	}

	return
}

// getVideoDetail 返回SvDetail的列表，返回的list顺序和svids顺序一致，但不保证svid都能出现在list中
func (s *Service) getVideoDetail(c context.Context, mid int64, qn int64, device *bm.Device, svids []int64, fullVersion bool) (res map[int64]*v1.SvDetail, err error) {
	res = make(map[int64]*v1.SvDetail)
	if len(svids) == 0 {
		return
	}
	// 拉取视频详情
	svRes, err := s.svInfos(c, svids, mid, true)
	if err != nil {
		log.Errorv(c, log.KV("log", "get video detail from dao fail"))
		return
	}

	log.V(1).Infov(c, log.KV("log", "get_video_detail"), log.KV("req_size", len(svids)), log.KV("rsp_size", len(svRes)))
	if len(svRes) == 0 {
		log.Warnv(c, log.KV("log", "get_video_detail_empty"), log.KV("req_size", len(svids)), log.KV("rsp_size", len(svRes)))
		return
	}

	// 开始组装回包
	currentTs := time.Now().Unix()
	for svid, svInfo := range svRes {
		sv := new(v1.SvDetail)
		// 组装video基础信息
		sv.VideoResponse = *svInfo
		if currentTs > int64(sv.Pubtime) {
			sv.ElapsedTime = currentTs - int64(sv.Pubtime)
		}
		res[svid] = sv
	}

	return
}

// parseCursor从cursor_prev和cursor_next，判断请求的方向，以及生成cursor
func parseCursor(cursorPrev string, cursorNext string) (cursor model.CursorValue, directionNext bool, err error) {
	// 判断是向前还是向后查询
	directionNext = true
	cursorStr := cursorNext
	if len(cursorNext) == 0 && len(cursorPrev) > 0 {
		directionNext = false
		cursorStr = cursorPrev
	}
	// 解析cursor中的cursor_id
	if len(cursorStr) != 0 {
		var cursorData = []byte(cursorStr)
		err = json.Unmarshal(cursorData, &cursor)
		if err != nil {
			err = ecode.ReqParamErr
			return
		}
	}
	// 第一次请求的时候，携带的svid=0，需要转成max传给dao层
	if directionNext && cursor.CursorID == 0 {
		cursor.CursorID = model.MaxInt64
		cursor.CursorTime = xtime.Time(time.Now().Unix())
	}
	return
}
