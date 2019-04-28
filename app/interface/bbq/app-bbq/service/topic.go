package service

import (
	"context"
	"encoding/json"
	"go-common/app/interface/bbq/app-bbq/api/http/v1"
	topic "go-common/app/service/bbq/topic/api"
	"go-common/library/ecode"
	"go-common/library/log"
	"net/http"
	"net/url"
)

func (s *Service) getExtension(ctx context.Context, svids []int64) (res map[int64]*topic.VideoExtension, err error) {
	res = make(map[int64]*topic.VideoExtension, len(svids))
	// 0. check

	// 1. get extension
	req := &topic.ListExtensionReq{Svids: svids}
	reply, err := s.topicClient.ListExtension(ctx, req)
	if err != nil {
		log.Warnw(ctx, "log", "get extension fail")
		return
	}

	// 2. form extension
	for _, extension := range reply.List {
		res[extension.Svid] = extension
	}
	return
}

// TopicDetail 获取话题详情
func (s *Service) TopicDetail(ctx context.Context, mid int64, req *topic.TopicVideosReq) (res *v1.TopicDetail, err error) {
	res = new(v1.TopicDetail)
	// 0. check
	if req.TopicId == 0 {
		err = ecode.TopicReqParamErr
		log.Warnw(ctx, "log", "topic id is 0")
		return
	}

	// 1. 获取话题信息及话题内视频
	topicDetails, err := s.topicClient.ListTopicVideos(ctx, req)
	if err != nil {
		log.Errorw(ctx, "log", "get list topic videos fail")
		return
	}

	// 2. 获取视频详情
	// 2.0 获取视频id
	var svids []int64
	topicVideoMap := make(map[int64]*topic.VideoItem)
	for _, item := range topicDetails.List {
		if _, exists := topicVideoMap[item.Svid]; !exists {
			topicVideoMap[item.Svid] = item
			svids = append(svids, item.Svid)
		}
	}
	// 2.1 获取详情
	svInfos, err := s.svInfos(ctx, svids, mid, false)
	if err != nil {
		log.Warnw(ctx, "log", "get sv infos fail", "svid", svids)
		return
	}

	// 3. 组装回包
	res.TopicInfo = topicDetails.TopicInfo
	res.HasMore = topicDetails.HasMore
	for _, item := range topicDetails.List {
		var topicVideo *v1.TopicVideo
		if svInfo, exists := svInfos[item.Svid]; !exists {
			log.Errorw(ctx, "log", "cannot find topicVideo response in topic detail", "svid", item.Svid)
			continue
		} else {
			topicVideo = new(v1.TopicVideo)
			topicVideo.VideoResponse = svInfo
		}
		topicVideo.CursorValue = item.CursorValue
		topicVideo.HotType = item.HotType
		res.List = append(res.List, topicVideo)
	}

	return
}

func (s *Service) getDiscoveryData(ctx context.Context, uri string) (data []byte, err error) {
	var ret struct {
		Code int             `json:"code"`
		Msg  string          `json:"message"`
		Data json.RawMessage `json:"data"`
	}

	req, err := s.httpClient.NewRequest(http.MethodGet, uri, "", url.Values{})
	if err != nil {
		log.Errorw(ctx, "log", "http.NewRequest error", "err", err)
		return
	}
	if err = s.httpClient.Do(ctx, req, &ret); err != nil {
		log.Errorw(ctx, "log", "client Do error", "err", err)
		return
	}
	if ret.Code != 0 {
		log.Errorw(ctx, "log", "return code error", "code", ret.Code)
		return
	}
	data = ret.Data
	return
}

func (s *Service) getHotWords(ctx context.Context) (list []string, err error) {
	list = make([]string, 0, 1)

	data, err := s.getDiscoveryData(ctx, "http://bbq-mng.bilibili.co/bbq/cms/hotword/api")
	if err != nil {
		log.Warnw(ctx, "log", "get discovery data fail")
		return
	}

	var hotWordResponse struct {
		OnshelfList []string `json:"onshelf_list"`
	}
	err = json.Unmarshal(data, &hotWordResponse)
	if err != nil {
		log.Errorw(ctx, "log", "unmarshal hot word response fail", "data", string(data))
		return
	}
	list = hotWordResponse.OnshelfList
	if list == nil {
		list = make([]string, 0, 1)
	}
	return
}

func (s *Service) getBanner(ctx context.Context) (list []*v1.Banner, err error) {
	list = make([]*v1.Banner, 0, 1)

	data, err := s.getDiscoveryData(ctx, "http://bbq-mng.bilibili.co/bbq/cms/banner/api")
	if err != nil {
		log.Warnw(ctx, "log", "get discovery data fail")
		return
	}

	type httpBanner struct {
		Title   string `json:"title"`
		ImgUrl  string `json:"img_url"`
		JumpUrl string `json:"jump_url"`
	}
	var bannerResponse struct {
		BannerList []*httpBanner `json:"banner_list"`
	}
	bannerResponse.BannerList = make([]*httpBanner, 0)
	err = json.Unmarshal(data, &bannerResponse)
	if err != nil {
		log.Errorw(ctx, "log", "unmarshal hot word response fail", "data", string(data))
		return
	}

	for _, item := range bannerResponse.BannerList {
		banner := new(v1.Banner)
		banner.Name = item.Title
		banner.PIC = item.ImgUrl
		banner.Scheme = item.JumpUrl
		list = append(list, banner)
	}

	if list == nil {
		list = make([]*v1.Banner, 0, 1)
	}
	return
}

// Discovery 发现页
func (s *Service) Discovery(ctx context.Context, mid int64, req *v1.DiscoveryReq) (res *v1.DiscoveryRes, err error) {
	res = new(v1.DiscoveryRes)
	res.BannerList = make([]*v1.Banner, 0, 10)
	res.HotWords = make([]string, 0, 10)
	res.TopicList = make([]*v1.TopicDetail, 0, 10)
	res.HasMore = false
	// check

	// 1. 条件判断
	if req.Page == 1 {
		// 请求热词
		if res.HotWords, err = s.getHotWords(ctx); err != nil {
			log.Warnw(ctx, "log", "get hot words fail", "err", err)
		}
		// 请求banner
		if res.BannerList, err = s.getBanner(ctx); err != nil {
			log.Warnw(ctx, "log", "get banner fail", "err", err)
		}
	}

	// 2. 请求话题详情
	reply, err := s.topicClient.ListDiscoveryTopics(ctx, &topic.ListDiscoveryTopicReq{Page: req.Page})
	if err != nil {
		log.Errorw(ctx, "log", "get discovery topics list fail")
		return
	}
	// 2.1 收集svid
	var svids []int64
	for _, topicDetail := range reply.List {
		for _, videoItem := range topicDetail.List {
			svids = append(svids, videoItem.Svid)
		}
	}
	// 2.2 获取视频详情
	svInfos, err := s.svInfos(ctx, svids, mid, false)
	if err != nil {
		log.Warnw(ctx, "log", "get sv infos fail", "svids", svids)
		return
	}

	// 3. 组装回包
	res.HasMore = reply.HasMore
	for _, topicDetail := range reply.List {
		newTopicDetail := new(v1.TopicDetail)
		newTopicDetail.TopicInfo = topicDetail.TopicInfo
		newTopicDetail.HasMore = topicDetail.HasMore
		for _, videoItem := range topicDetail.List {
			if val, exists := svInfos[videoItem.Svid]; exists {
				topicVideo := &v1.TopicVideo{CursorValue: videoItem.CursorValue, VideoResponse: val}
				newTopicDetail.List = append(newTopicDetail.List, topicVideo)
			} else {
				log.Warnw(ctx, "log", "get sv info fail", "svid", videoItem.Svid)
			}
		}
		if len(newTopicDetail.List) == 0 {
			log.Warnw(ctx, "log", "topic has nothing", "topic", topicDetail)
			continue
		}
		res.TopicList = append(res.TopicList, newTopicDetail)
	}

	return
}

// TopicSearch 话题搜索
func (s *Service) TopicSearch(ctx context.Context, req *v1.TopicSearchReq) (res *v1.TopicSearchResponse, err error) {
	res = new(v1.TopicSearchResponse)

	if len(req.Keyword) == 0 {
		var reply *topic.ListTopicsReply
		reply, err = s.topicClient.ListTopics(ctx, &topic.ListTopicsReq{Page: 1})
		if err != nil {
			log.Warnw(ctx, "log", "get topic list fail", "err", err)
			return
		}
		res.List = reply.List
		res.HasMore = reply.HasMore
	}
	res.HasMore = false

	return
}
