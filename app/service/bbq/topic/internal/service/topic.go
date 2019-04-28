package service

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"go-common/app/service/bbq/topic/api"
	"go-common/app/service/bbq/topic/internal/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup.v2"
	"strings"
)

// UpdateVideoScore 更新视频的score分
func (s *Service) UpdateVideoScore(ctx context.Context, req *api.UpdateVideoScoreReq) (res *empty.Empty, err error) {
	res = new(empty.Empty)
	err = s.dao.UpdateVideoScore(ctx, req.Svid, req.Score)
	return
}

// UpdateVideoState 更新视频的状态，databus消费会调用
func (s *Service) UpdateVideoState(ctx context.Context, req *api.UpdateVideoStateReq) (res *empty.Empty, err error) {
	res = new(empty.Empty)
	err = s.dao.UpdateVideoState(ctx, req.Svid, req.State)
	return
}

// UpdateTopicDesc 更新话题信息，
func (s *Service) UpdateTopicDesc(ctx context.Context, in *api.TopicInfo) (res *empty.Empty, err error) {
	res = new(empty.Empty)

	// 0. check
	if in.TopicId == 0 {
		err = ecode.TopicIDErr
		return
	}
	if strings.Count(in.Desc, "")-1 > model.MaxTopicDescLen {
		err = ecode.TopicDescLenErr
		log.Warnw(ctx, "log", "topic desc too long", "desc", in.Desc)
		return
	}

	// 1. update
	err = s.dao.UpdateTopic(ctx, in.TopicId, "desc", in.Desc)
	if err != nil {
		log.Warnw(ctx, "log", "update topic desc fail", "topic_id", in.TopicId)
	}

	return
}

// UpdateTopicState 更新话题信息，
func (s *Service) UpdateTopicState(ctx context.Context, in *api.TopicInfo) (res *empty.Empty, err error) {
	res = new(empty.Empty)

	// 0. check
	if in.TopicId == 0 {
		err = ecode.TopicIDErr
		return
	}

	// 1. update
	state := model.TopicStateAvailable
	if in.State > 0 {
		state = model.TopicStateUnavailable
	}
	err = s.dao.UpdateTopic(ctx, in.TopicId, "state", state)
	if err != nil {
		log.Warnw(ctx, "log", "update topic state fail", "topic_id", in.TopicId)
	}

	return
}

// VideoTopic 返回视频关联的所有话题信息，cms使用，当前只有topic_id
func (s *Service) VideoTopic(ctx context.Context, in *api.VideoTopicReq) (res *api.VideoTopicReply, err error) {
	res = new(api.VideoTopicReply)
	// 0. check
	if in.Svid == 0 {
		log.Errorw(ctx, "log", "param svid error")
		err = ecode.ReqParamErr
		return
	}

	// 1. 获取话题id列表
	res.List, err = s.dao.GetVideoTopic(ctx, in.Svid)
	if err != nil {
		log.Warnw(ctx, "log", "get video topic fail")
		return
	}

	return
}

// ListCmsTopics cms获取话题信息
func (s *Service) ListCmsTopics(ctx context.Context, in *api.ListCmsTopicsReq) (res *api.ListCmsTopicsReply, err error) {
	res = new(api.ListCmsTopicsReply)

	// 0. check

	// 1. 查找
	var topicIDs []int64
	res.HasMore = false
	stickTopic := make(map[int64]bool)
	if len(in.Name) > 0 {
		// 话题名搜索
		var topics map[string]int64
		if topics, err = s.dao.TopicID(ctx, []string{in.Name}); err != nil {
			log.Warnw(ctx, "log", "get topic id fail", "name", in.Name)
			return
		}
		topicID, exists := topics[in.Name]
		if !exists {
			log.Warnw(ctx, "log", "get topic id fail", "name", in.Name)
			return
		}
		topicIDs = append(topicIDs, topicID)
	} else if in.TopicId != 0 {
		// topic_id搜索
		topicIDs = append(topicIDs, in.TopicId)
	} else if in.State == model.TopicStateUnavailable {
		// 话题下架搜索
		topicIDs, res.HasMore, err = s.dao.ListUnAvailableTopics(ctx, in.Page, model.CmsTopicSize)
		if err != nil {
			log.Warnw(ctx, "log", "get topic ids fail")
			return
		}
	} else {
		// 话题详情页搜索
		topicIDs, res.HasMore, err = s.dao.ListRankTopics(ctx, in.Page, model.CmsTopicSize)
		if err != nil {
			log.Warnw(ctx, "log", "get topic ids fail")
			return
		}
	}

	// 2. 查找置顶数据
	stickList, _ := s.dao.GetStickTopic(ctx)
	for _, topicID := range stickList {
		stickTopic[topicID] = true
	}

	// 3. 获取TopicInfo
	topicInfos, err := s.dao.TopicInfo(ctx, topicIDs)
	if err != nil {
		log.Warnw(ctx, "log", "get topic info fail", "topic_ids", topicIDs)
		return
	}
	for _, topicID := range topicIDs {
		topicInfo, exists := topicInfos[topicID]
		if !exists {
			log.Errorw(ctx, "log", "get error topic id", "topic_id", topicID)
			continue
		}

		if _, exists2 := stickTopic[topicID]; exists2 {
			topicInfo.HotType = api.TopicHotTypeStick
		}
		res.List = append(res.List, topicInfo)
	}

	return
}

// ListDiscoveryTopics 用于发现页，这里采用的是page的样式
// 这里其实可以考虑只返回话题列表，由上游自行调用ListTopicVideo接口
func (s *Service) ListDiscoveryTopics(ctx context.Context, req *api.ListDiscoveryTopicReq) (res *api.ListDiscoveryTopicReply, err error) {
	res = new(api.ListDiscoveryTopicReply)
	page := req.Page
	// 0. check

	// 1. 获取话题列表
	//老逻辑，topicIDList, hasMore, err := s.dao.ListRankTopics(ctx, page, model.DiscoveryTopicSize)
	// 新逻辑，只取置顶的话题
	topicIDList, err := s.dao.GetStickTopic(ctx)
	if err != nil {
		log.Warnw(ctx, "log", "get recommend topic fail", "page", page)
		return
	}
	res.HasMore = false // 永远为false，只取置顶
	if len(topicIDList) == 0 {
		return
	}

	// 2. 获取话题信息
	topicInfos, err := s.getAvailableTopicInfo(ctx, topicIDList)
	if err != nil {
		log.Warnw(ctx, "log", "get topic info fail", "topic_id", topicIDList)
		return
	}

	// 3. 获取话题视频
	topicDetails := make(map[int64]*api.TopicDetail, len(topicIDList))
	topicDetailChan := make(chan *api.TopicDetail, len(topicIDList))
	group := errgroup.WithCancel(ctx)
	group.GOMAXPROCS(5)
	var groupTopicVideo = func(topicID int64) {
		group.Go(func(ctx context.Context) (err error) {
			log.V(1).Infow(ctx, "log", "get one topic videos", "topic_id", topicID)
			list, _, err := s.dao.ListTopicVideos(ctx, topicID, "", "", model.DiscoveryTopicVideoSize)
			if err != nil {
				log.Warnw(ctx, "log", "get topic videos fail", "topic_id", topicID)
				return
			}
			topicDetailChan <- &api.TopicDetail{TopicInfo: &api.TopicInfo{TopicId: topicID}, List: list}
			return
		})
	}
	for _, topicID := range topicIDList {
		groupTopicVideo(topicID)
	}
	err = group.Wait()
	close(topicDetailChan)
	if err != nil {
		log.Warnw(ctx, "log", "group Go occurs error", "err", err)
		return
	}
	for topicDetail := range topicDetailChan {
		topicDetails[topicDetail.TopicInfo.TopicId] = topicDetail
	}

	// 4. 获取置顶topic_id
	stickList, tmpErr := s.dao.GetStickTopic(ctx)
	if tmpErr != nil {
		log.Warnw(ctx, "log", "get stick topic fail")
	}
	stickMap := make(map[int64]bool)
	for _, topicID := range stickList {
		stickMap[topicID] = true
	}

	// 5. 组合回包
	// 只有TopicInfo和video都存在才会返回给前级
	for _, topicID := range topicIDList {
		// 获取话题视频
		topicDetail, exists := topicDetails[topicID]
		if !exists {
			log.Errorw(ctx, "log", "cannot find topic detail", "topic_id", topicID)
			continue
		}
		if len(topicDetail.List) == 0 {
			log.Warnw(ctx, "log", "video num is 0 in this topic", "topic", topicDetail)
			continue
		}

		// 获取话题info
		topicInfo, exists := topicInfos[topicID]
		if !exists {
			log.Errorw(ctx, "log", "cannot find topic info", "topic_id", topicID)
			continue
		}
		// 设置置顶标志
		if _, exists2 := stickMap[topicID]; exists2 {
			topicInfo.HotType = api.TopicHotTypeStick
		}

		topicDetail.TopicInfo = topicInfo

		res.List = append(res.List, topicDetail)
	}

	return
}

// ListTopics 话题列表
func (s *Service) ListTopics(ctx context.Context, req *api.ListTopicsReq) (res *api.ListTopicsReply, err error) {
	res = new(api.ListTopicsReply)
	page := &req.Page

	// 0. check

	// 1. 获取话题列表
	// 新逻辑，只取置顶的话题
	topicIDList, err := s.dao.GetStickTopic(ctx)
	if err != nil {
		log.Warnw(ctx, "log", "get recommend topic fail", "page", page)
		return
	}
	res.HasMore = false // 永远为false，只取置顶
	if len(topicIDList) == 0 {
		return
	}
	stickMap := make(map[int64]bool)
	for _, topicID := range topicIDList {
		stickMap[topicID] = true
	}

	// 2. 获取话题信息
	topicInfos, err := s.getAvailableTopicInfo(ctx, topicIDList)
	if err != nil {
		log.Warnw(ctx, "log", "get topic info fail", "topic_id", topicIDList)
		return
	}

	// 3. 组合回包
	for _, topicID := range topicIDList {
		// 获取话题info
		topicInfo, exists := topicInfos[topicID]
		if !exists {
			log.Errorw(ctx, "log", "cannot find topic info", "topic_id", topicID)
			continue
		}
		// 设置置顶标志
		if _, exists2 := stickMap[topicID]; exists2 {
			topicInfo.HotType = api.TopicHotTypeStick
		}
		res.List = append(res.List, topicInfo)
	}
	return
}

// ListTopicVideos 获取话题下的视频
func (s *Service) ListTopicVideos(ctx context.Context, req *api.TopicVideosReq) (res *api.TopicDetail, err error) {
	res = new(api.TopicDetail)
	// 0.check
	if req.TopicId == 0 {
		err = ecode.TopicIDErr
		return
	}

	// 1. 获取话题信息
	topicInfos, err := s.dao.TopicInfo(ctx, []int64{req.TopicId})
	if err != nil {
		log.Warnw(ctx, "log", "get topic info fail")
		return
	}
	topicInfo, exists := topicInfos[req.TopicId]
	if !exists {
		err = ecode.TopicIDNotFound
		return
	}
	res.TopicInfo = topicInfo

	// 2. 获取话题内的视频
	list, hasMore, err := s.dao.ListTopicVideos(ctx, req.TopicId, req.CursorPrev, req.CursorNext, model.TopicVideoSize)
	if err != nil {
		log.Warnw(ctx, "log", "get topic video list fail", "err", err, "req", req)
		return
	}

	res.HasMore = hasMore
	res.List = list

	return
}

// StickTopic 话题的置顶、取消置顶操作
func (s *Service) StickTopic(ctx context.Context, in *api.StickTopicReq) (res *empty.Empty, err error) {
	res = new(empty.Empty)
	err = s.dao.StickTopic(ctx, in.TopicId, in.Op)
	return
}

// StickTopicVideo 话题下视频的置顶、取消置顶操作
func (s *Service) StickTopicVideo(ctx context.Context, in *api.StickTopicVideoReq) (res *empty.Empty, err error) {
	res = new(empty.Empty)
	err = s.dao.StickTopicVideo(ctx, in.TopicId, in.Svid, in.Op)
	return
}

// SetStickTopicVideo 替换话题下的置顶视频
func (s *Service) SetStickTopicVideo(ctx context.Context, in *api.SetStickTopicVideoReq) (res *empty.Empty, err error) {
	res = new(empty.Empty)
	err = s.dao.SetStickTopicVideo(ctx, in.TopicId, in.Svids)
	if err != nil {
		log.Warnw(ctx, "log", "set stick topic video fail")
		return
	}
	return
}

func (s *Service) registerTopic(ctx context.Context, svid int64, list []*api.TitleExtraItem) (res []*api.TitleExtraItem, err error) {
	// 0. 校验请求
	if svid == 0 {
		log.Warnw(ctx, "log", "svid=0")
		return
	}
	if len(list) == 0 {
		return
	}
	if len(list) > model.MaxSvTopicNum {
		err = ecode.TopicTooManyInOneVideo
		return
	}

	// 1 话题操作
	// 1.0 插入新话题，同时更新老话题，这里dao层用on duplicate key就行
	topics := make(map[string]*api.TopicInfo)
	for _, item := range list {
		topics[item.Name] = &api.TopicInfo{Name: item.Name}
	}
	// 插入，同时获取话题ID
	newTopics, err := s.dao.InsertTopics(ctx, topics)
	if err != nil {
		log.Warnw(ctx, "log", "insert topic fail")
		return
	}

	// 2. 视频插入topic_video
	var topicIDs []int64
	for _, topicInfo := range newTopics {
		topicIDs = append(topicIDs, topicInfo.TopicId)
	}
	num, err := s.dao.InsertTopicVideo(ctx, svid, topicIDs)
	if err != nil {
		log.Warnw(ctx, "log", "insert topic video fail")
		return
	}
	// 打error日志，但是这个情况是可以接受的，保持幂等即可
	if int(num) != len(topicIDs) {
		log.Errorw(ctx, "log", "insert topic_video num not match", "rows_affected", num, "topic_num", len(topicIDs), "svid", svid)
	}

	// 3. 返回新的数组，结构补充完整
	for _, item := range list {
		if topicInfo, exists := newTopics[item.Name]; exists {
			item.Scheme = fmt.Sprintf("qing://topic?topic_id=%d", topicInfo.TopicId)
		} else {
			log.Errorw(ctx, "log", "get topic id fail", "name", item.Name)
		}
	}
	res = list

	return
}

// 用于获取状态可见的话题，为的是和审核分开
func (s *Service) getAvailableTopicInfo(c context.Context, keys []int64) (res map[int64]*api.TopicInfo, err error) {
	res, err = s.dao.TopicInfo(c, keys)
	if err != nil {
		log.Warnw(c, "log", "get topic info fail")
		return
	}
	var toDeletedTopic []int64
	for topicID, topicInfo := range res {
		if topicInfo.State == api.TopicStateUnAvailable {
			toDeletedTopic = append(toDeletedTopic, topicID)
		}
	}
	for _, topicID := range toDeletedTopic {
		log.Warnw(c, "log", "get one unavailable topic", "topic", res[topicID])
		delete(res, topicID)
	}
	return
}
