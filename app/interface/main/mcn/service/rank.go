package service

import (
	"context"
	"sort"

	"go-common/app/interface/main/mcn/dao/mcndao"
	"go-common/app/interface/main/mcn/model"
	"go-common/app/interface/main/mcn/model/mcnmodel"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	//SortFieldFans by fans
	SortFieldFans = "fans_count"
	//SortFieldMonthFans by month increase
	SortFieldMonthFans = "fans_count_increase_month"
	//SortFieldArchiveCount by archive count
	SortFieldArchiveCount = "archive_count"
)

var (
	sortFieldMap = map[string]mcndao.RecommendSortFunc{
		SortFieldFans:         mcndao.RecommendSortByFansDesc,
		SortFieldMonthFans:    mcndao.RecommendSortByMonthFansDesc,
		SortFieldArchiveCount: mcndao.RecommendSortByArchiveCountDesc,
	}
)

//McnGetRankUpFans get up rank fans
func (s *Service) McnGetRankUpFans(c context.Context, arg *mcnmodel.McnGetRankReq) (res *mcnmodel.McnGetRankUpFansReply, err error) {
	res, err = s.getRankResult(c, arg, s.mcndao.GetRankUpFans)
	return
}

//McnGetRankArchiveLikes get rank archive likes
func (s *Service) McnGetRankArchiveLikes(c context.Context, arg *mcnmodel.McnGetRankReq) (res *mcnmodel.McnGetRankUpFansReply, err error) {
	res, err = s.getRankResult(c, arg, s.mcndao.GetRankArchiveLikes)
	return
}

func (s *Service) getRankResult(c context.Context, arg *mcnmodel.McnGetRankReq, rankFunc mcndao.RankFunc) (res *mcnmodel.McnGetRankUpFansReply, err error) {
	mcnSign, err := s.getMcnWithState(c, arg.McnMid, model.MCNSignStateOnSign)
	if err != nil {
		if err != ecode.NothingFound {
			log.Error("error get state, err=%s", err)
		}
		return
	}

	v, err := rankFunc(mcnSign.ID)
	if err != nil || v == nil {
		log.Error("get rank fail, sign id=%d, err=%s", mcnSign.ID, err)
		return
	}

	res = new(mcnmodel.McnGetRankUpFansReply)
	res.Result = v.GetList(arg.Tid, arg.DataType)
	res.TypeList = v.GetTypeList(arg.DataType)
	return
}

//GetRecommendPool get recommend pool reply
func (s *Service) GetRecommendPool(c context.Context, arg *mcnmodel.McnGetRecommendPoolReq) (res *mcnmodel.McnGetRecommendPoolReply, err error) {
	mcnSign, err := s.getMcnWithState(c, arg.McnMid, model.MCNSignStateOnSign)
	if err != nil {
		if err != ecode.NothingFound {
			log.Error("error get state, err=%s", err)
		}
		return
	}

	var limit, offset = arg.CheckPageValidation()
	recommendCache, err := s.mcndao.GetRecommendPool()
	if err != nil {
		log.Error("get recommend pool fail, err=%s, mcn=%d", err, mcnSign.McnMid)
		return
	}

	if recommendCache == nil {
		log.Warn("recommend cache is nil, mcn=%d", mcnSign.McnMid)
		res.PageResult = arg.ToPageResult(0)
		return
	}

	res = new(mcnmodel.McnGetRecommendPoolReply)
	var upList = recommendCache.UpTidMap[arg.Tid]
	var listLen = len(upList)
	if offset >= listLen {
		return
	}

	res.PageResult = arg.ToPageResult(listLen)

	if upList == nil {
		return
	}

	var end = limit + offset
	if end >= listLen {
		end = listLen
	}

	var sortFunc mcndao.RecommendSortFunc
	switch arg.OrderField {
	case SortFieldMonthFans, SortFieldArchiveCount:
		sortFunc = sortFieldMap[arg.OrderField]
	}
	if sortFunc != nil {
		sort.Sort(&mcndao.RecommendDataSorter{Datas: upList, By: sortFunc})
	}

	var dest = make([]*mcnmodel.McnGetRecommendPoolInfo, listLen)

	// 如果是升序，那么把他们倒过来
	if arg.Sort == "asc" {
		copy(dest, upList)
		for left, right := 0, len(dest)-1; left < right; left, right = left+1, right-1 {
			dest[left], dest[right] = dest[right], dest[left]
		}
	} else {
		dest = upList
	}
	log.Info("offset, limit=%d,%d", offset, limit)
	dest = dest[offset:end]
	res.Result = dest
	return
}

//GetRecommendPoolTidList get tid list
func (s *Service) GetRecommendPoolTidList(c context.Context, arg *mcnmodel.McnGetRecommendPoolTidListReq) (res *mcnmodel.McnGetRecommendPoolTidListReply, err error) {
	mcnSign, err := s.getMcnWithState(c, arg.McnMid, model.MCNSignStateOnSign)
	if err != nil {
		if err != ecode.NothingFound {
			log.Error("error get state, err=%s", err)
		}
		return
	}

	recommendCache, err := s.mcndao.GetRecommendPool()
	if err != nil {
		log.Error("get recommend pool fail, err=%s, mcn=%d", err, mcnSign.McnMid)
		return
	}

	if recommendCache == nil {
		log.Warn("recommend cache is nil, mcn=%d", mcnSign.McnMid)
		return
	}

	res = new(mcnmodel.McnGetRecommendPoolTidListReply)
	res.Result = recommendCache.TidTypeList
	return
}
